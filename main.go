package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-telegram/bot"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"xmpp/internal/common"
	"xmpp/internal/matrix"
	"xmpp/internal/telegram"
	"xmpp/internal/xmpp"
)

type MessageRequest struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) WithMessage(message string) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: message,
	}
}

var (
	InvalidRequestBody       = &AppError{Code: 1, Message: "Invalid request body"}
	FieldsAreRequired        = &AppError{Code: 2, Message: "Missing required fields"}
	WrongRecipientFormat     = &AppError{Code: 3, Message: "Wrong Recipient format"}
	UnexpectedMessengerError = &AppError{Code: 4, Message: "Unexpected messenger error"}
	TelegramForbidden        = &AppError{Code: 5, Message: "User blocked telegram bot"}
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if common.GetEnvBool("USE_XMPP") {
		xmpp.Init(ctx, common.GetEnv("XMPP_SERVER"), common.GetEnv("XMPP_USERNAME"), common.GetEnv("XMPP_PASSWORD"))
	}

	if common.GetEnvBool("USE_TELEGRAM") {
		telegram.Init(ctx, common.GetEnv("TELEGRAM_BOT_TOKEN"))
	}

	if common.GetEnvBool("USE_MATRIX") {
		matrix.Init(ctx)
	}

	http.HandleFunc("/xmpp", withValidation(xmppSendMessageHandler))
	http.HandleFunc("/telegram", withValidation(telegramSendMessageHandler))
	http.HandleFunc("/matrix", withValidation(matrixSendMessageHandler))

	port := common.GetEnv("SERVICE_PORT", "8080")

	log.Printf("Starting HTTP server on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func withValidation(next MessageHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var msgReq MessageRequest
		if err := json.NewDecoder(r.Body).Decode(&msgReq); err != nil {
			sendErrorResponse(w, *InvalidRequestBody)
			return
		}

		if msgReq.Message == "" {
			sendErrorResponse(w, *FieldsAreRequired, "field 'message is required")
			return
		}

		if msgReq.Recipient == "" {
			sendErrorResponse(w, *FieldsAreRequired, "fields 'recipient' is required")
			return
		}

		err := next(msgReq)
		if err != nil {
			sendErrorResponse(w, *err)
			return
		}

		w.WriteHeader(http.StatusOK)
		var msgResponse = map[string]string{"message": msgReq.Message, "status": "sent"}
		_ = json.NewEncoder(w).Encode(msgResponse)
	}
}

func sendErrorResponse(w http.ResponseWriter, appError AppError, customMessage ...string) {
	if len(customMessage) > 0 {
		appError.Message = customMessage[0]
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(appError)
}

type MessageHandler func(req MessageRequest) *AppError

func xmppSendMessageHandler(req MessageRequest) *AppError {
	err := xmpp.SendMessage(req.Recipient, req.Message)
	if err != nil {
		return UnexpectedMessengerError.WithMessage(err.Error())
	}
	return nil
}

func telegramSendMessageHandler(req MessageRequest) *AppError {

	chatId, err := strconv.ParseInt(req.Recipient, 10, 64)
	if err != nil {
		return WrongRecipientFormat.WithMessage("Expected chat id (number), got: " + req.Recipient)
	}

	err = telegram.SendMessage(chatId, req.Message)
	if err != nil {

		if errors.Is(err, bot.ErrorForbidden) {
			return TelegramForbidden
		}

		return UnexpectedMessengerError.WithMessage(err.Error())
	}
	return nil
}

func matrixSendMessageHandler(req MessageRequest) *AppError {

	if !strings.HasPrefix(req.Recipient, "@") {
		return WrongRecipientFormat.WithMessage("Expected user id (@username), got: " + req.Recipient)
	}

	err := matrix.SendMessage(req.Recipient, req.Message)
	if err != nil {
		return UnexpectedMessengerError.WithMessage(err.Error())
	}

	return nil
}
