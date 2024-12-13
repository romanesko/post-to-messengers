package matrix

import (
	"context"
	"errors"
	"log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"xmpp/internal/common"
)

var matrixClient *mautrix.Client
var ctx context.Context

var userChats = make(map[id.UserID]id.RoomID)

func Init(context context.Context) {

	ctx = context

	userId := id.UserID(common.GetEnv("MATRIX_USER_ID", ""))
	accessToken := common.GetEnv("MATRIX_ACCESS_TOKEN", "")

	var err error
	matrixClient, err = mautrix.NewClient("https://"+common.GetEnv("MATRIX_SERVER"), userId, accessToken)
	if err != nil {
		panic(err)
	}

	if userId == "" || accessToken == "" {
		log.Printf("Matrix logging in by username and password")
		_, err = matrixClient.Login(ctx, &mautrix.ReqLogin{
			Type: mautrix.AuthTypePassword,
			Identifier: mautrix.UserIdentifier{
				Type: mautrix.IdentifierTypeUser,
				User: common.GetEnv("MATRIX_USERNAME"),
			},
			Password:                 common.GetEnv("MATRIX_PASSWORD"),
			InitialDeviceDisplayName: "bot",
			RefreshToken:             false,
			StoreCredentials:         true,
			StoreHomeserverURL:       false,
		})
		log.Printf("Matrix access token: %s", matrixClient.AccessToken)
	} else {
		log.Printf("Matrix logging in by auth token with user ID: %s", userId)
	}

	if err != nil {
		log.Fatal(err)
	}

	err = getRooms()
	if err != nil {
		log.Fatal(err)
	}

}

func getRooms() error {

	resp, err := matrixClient.JoinedRooms(ctx)

	if err != nil {
		log.Printf("Failed to get rooms: %v", err)
		return err
	}

	for _, room := range resp.JoinedRooms {
		resp, err := matrixClient.JoinedMembers(ctx, room)
		if err != nil {
			log.Printf("Failed to get members: %v", err)
			continue
		}

		if len(resp.Joined) != 2 {
			log.Printf("Unexpected number of members in room %s: %d, removing", room, len(resp.Joined))
			leaveRoom(room)
			continue
		}

		for key, _ := range resp.Joined {
			if key != matrixClient.UserID {
				userChats[key] = room
			}
		}
	}

	return nil
}

func leaveRoom(roomID id.RoomID) {
	_, _ = matrixClient.LeaveRoom(ctx, roomID)
}

func getOrCreateDMRoom(userID id.UserID) (id.RoomID, error) {

	room := userChats[userID]

	if room != "" {
		return room, nil
	}

	var userIDs []id.UserID
	userIDs = append(userIDs, userID)

	creationResult, err := matrixClient.CreateRoom(ctx, &mautrix.ReqCreateRoom{
		Visibility: "private",
		Invite:     userIDs,
	})
	if err != nil {
		log.Printf("Failed to create DM room: %v", err)
		return "", err
	}

	log.Printf("Direct message room created with ID: %s", creationResult.RoomID.String())

	userChats[userID] = creationResult.RoomID
	return creationResult.RoomID, nil
}

func sendMessageMatrix(roomID id.RoomID, message string) error {

	content := event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    message,
	}
	_, err := matrixClient.SendMessageEvent(ctx, roomID, event.EventMessage, content)
	return err
}

func SendMessage(recipient string, message string) error {

	if matrixClient == nil {
		return errors.New("Matrix client is not initialized")
	}

	room, err := getOrCreateDMRoom(id.UserID(recipient))
	if err != nil {
		log.Printf("Failed to get or create DM room: %v", err)
		return err
	}
	return sendMessageMatrix(room, message)
}
