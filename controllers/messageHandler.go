package controllers

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/dtsang7/ASAPP/models"
	"net/http"
)

const (
	ErrorMismatchIDMessage = "ID in token doesn't match sender_id. Stop pretending to be someone else :("
)

var errorMismatchIDMessage = errors.New(ErrorMismatchIDMessage)

type MessagesResponse struct {
	Messages []models.GetMessage `json:"messages"`
}

type Response struct {
	Id        int
	Timestamp string
}

type Retriever struct {
	RecipientId int
	MsgID       int
	Limit       int
}

func verifyTokenID(r *http.Request, senderID int) bool {
	user := r.Context().Value("user")
	tokenID, found := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"]
	if !found {
		return false
	}
	if id, ok := tokenID.(float64); !ok || id != float64(senderID) {
		return false
	}
	return true
}

func (h Handler) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var msg models.Message
	json.NewDecoder(r.Body).Decode(&msg)

	//validate
	err := ValidateSendMessage(msg)
	if err != nil {
		WriteJsonError(err, w)
		return
	}

	if !verifyTokenID(r, msg.SenderID) {
		WriteJsonError(errorMismatchIDMessage, w)
		return
	}

	msgID, timeStamp, err := h.DB.SendMessage(msg)

	if err != nil {
		WriteJsonError(err, w)
		return
	}

	resp := Response{msgID, timeStamp}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Write error", http.StatusInternalServerError)
	}
}

func (h Handler) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	var retriever Retriever
	json.NewDecoder(r.Body).Decode(&retriever)

	err := ValidateGetMessage(retriever)
	if err != nil {
		WriteJsonError(err, w)
		return
	}

	if retriever.Limit == 0 {
		retriever.Limit = 100 //default
	}

	messages, err := h.DB.GetMessages(retriever.RecipientId, retriever.MsgID, retriever.Limit)

	if err != nil {
		WriteJsonError(err, w)
		return
	}

	jsonErr := json.NewEncoder(w).Encode(MessagesResponse{messages})
	if jsonErr != nil {
		http.Error(w, "Write error", http.StatusInternalServerError)
	}
}
