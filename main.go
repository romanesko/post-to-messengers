package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/xmppo/go-xmpp"
)

type MessageRequest struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

var (
	xmppClient *xmpp.Client
)

func getEnv(name string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		log.Fatalf("Environment variable '%s' is not set", name)
	}
	return value
}
func main() {
	// Configuration variables
	server := getEnv("XMPP_SERVER")
	username := getEnv("XMPP_USERNAME")
	password := getEnv("XMPP_PASSWORD")

	// XMPP client options
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

	// Connect to the XMPP server
	var err error
	xmppClient, err = options.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to XMPP server: %v", err)
	}
	defer xmppClient.Close()

	// Start HTTP server
	http.HandleFunc("/send", sendMessageHandler)

	port := getEnv("SERVICE_PORT")

	log.Printf("Starting HTTP server on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var msgReq MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&msgReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if msgReq.Recipient == "" || msgReq.Message == "" {
		http.Error(w, "Recipient and message fields are required", http.StatusBadRequest)
		return
	}

	// Send the message
	_, err := xmppClient.Send(xmpp.Chat{
		Remote: msgReq.Recipient,
		Type:   "chat", // Type of message (e.g., "chat", "groupchat")
		Text:   msgReq.Message,
	})
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	log.Printf("Message sent to %s: %s", msgReq.Recipient, msgReq.Message)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Message sent successfully!"))
}
