package xmpp

import (
	"context"
	"errors"
	"fmt"
	"github.com/xmppo/go-xmpp"
	"log"
)

var (
	xmppClient *xmpp.Client
)

func Init(ctx context.Context, server string, username string, password string) {

	options := xmpp.Options{
		Host:          fmt.Sprintf("%s:%s", server, "5222"),
		User:          fmt.Sprintf("%s@%s", username, server),
		Password:      password,
		NoTLS:         true,
		Debug:         false,
		Session:       true,
		Status:        "chat",
		StatusMessage: "I'm here!",
	}

	var err error
	xmppClient, err = options.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to XMPP server: %v", err)
	}

	go func() {
		<-ctx.Done()
		fmt.Println("Received interrupt signal, closing XMPP client...")
		xmppClient.Close()
	}()

	log.Printf("XMPP client started with user ID %s", username)
}

func SendMessage(recipient string, message string) error {

	if xmppClient == nil {
		return errors.New("XMPP client is not initialized")
	}

	_, err := xmppClient.Send(xmpp.Chat{
		Remote: recipient,
		Type:   "chat",
		Text:   message,
	})
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}
	return nil

}
