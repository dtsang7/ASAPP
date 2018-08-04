package controllers

import (
	"database/sql"
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

type GetMessagesRequest struct {
	RecipientID int `json:"recipient"`
	StartMsgID  int `json:"start"`
	Limit       int
}
type GetMessagesResponse struct {
	Messages []Message `json:"messages"`
}

type MessageContent struct {
	Type   string `json:"type"`
	Text   string `json:"text,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	Source string `json:"source,omitempty"`
	Url    string `json:"url,omitempty"`
}

type Message struct {
	MsgID       int            `json:"id"`
	TimeStamp   string         `json:"timestamp"`
	SenderID    int            `json:"sender"`
	RecipientID int            `json:"recipient"`
	Content     MessageContent `json:"content"`
}

type SendMessageResponse struct {
	Id        int
	Timestamp string
}

func (h Handler) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var req Message
	json.NewDecoder(r.Body).Decode(&req)

	//validate
	err := ValidateSendMessage(req)
	if err != nil {
		WriteHttpError(err, w)
		return
	}

	if !verifyTokenID(r, req.SenderID) {
		WriteHttpError(errorMismatchIDMessage, w)
		return
	}
	dbMsg := models.Message{
		SenderID:    req.SenderID,
		RecipientID: req.RecipientID,
	}
	switch req.Content.Type {
	case "text":
		dbMsg.Type = "text"
		dbMsg.Message = sql.NullString{String: req.Content.Text, Valid: true}
	case "image":
		dbMsg.Type = "image"
		dbMsg.Url = sql.NullString{String: req.Content.Url, Valid: true}
		dbMsg.Width = sql.NullInt64{Int64: int64(req.Content.Width), Valid: true}
		dbMsg.Height = sql.NullInt64{Int64: int64(req.Content.Height), Valid: true}
	case "video":
		dbMsg.Type = "video"
		dbMsg.Url = sql.NullString{String: req.Content.Url, Valid: true}
		dbMsg.Source = sql.NullString{String: req.Content.Source, Valid: true}
	default:
		WriteHttpError(errors.New("Unsupported Message content type"), w)
		return
	}
	msgID, timeStamp, err := h.DB.SendMessage(dbMsg)
	if err != nil {
		WriteHttpError(err, w)
		return
	}

	resp := SendMessageResponse{msgID, timeStamp}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Write error", http.StatusInternalServerError)
	}
}

func (h Handler) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// Get query paramters
	req, err := ParseAndValidateGetMessageRequest(r)
	if err != nil {
		WriteHttpError(err, w)
		return
	}

	if !verifyTokenID(r, req.RecipientID) {
		WriteHttpError(errorMismatchIDMessage, w)
		return
	}

	if req.Limit == 0 {
		req.Limit = 100 //default
	}

	dbMsgs, err := h.DB.GetMessages(req.RecipientID, req.StartMsgID, req.Limit)
	var messages []Message
	for _, dbMsg := range dbMsgs {
		msg := Message{
			MsgID:       dbMsg.MsgID,
			TimeStamp:   dbMsg.TimeStamp,
			SenderID:    dbMsg.SenderID,
			RecipientID: dbMsg.RecipientID,
		}
		switch dbMsg.Type {
		case "text":
			msg.Content = MessageContent{
				Type: "text",
				Text: dbMsg.Message.String,
			}
		case "image":
			msg.Content = MessageContent{
				Type:   "image",
				Width:  int(dbMsg.Width.Int64),
				Height: int(dbMsg.Height.Int64),
				Url:    dbMsg.Url.String,
			}
		case "video":
			msg.Content = MessageContent{
				Type:   "video",
				Url:    dbMsg.Url.String,
				Source: dbMsg.Source.String,
			}
		}
		messages = append(messages, msg)
	}

	if err != nil {
		WriteHttpError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonErr := json.NewEncoder(w).Encode(GetMessagesResponse{messages})
	if jsonErr != nil {
		http.Error(w, "Write error", http.StatusInternalServerError)
	}
}

func verifyTokenID(r *http.Request, id int) bool {
	user := r.Context().Value("user")
	tokenID, found := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"]
	if !found {
		return false
	}
	if id, ok := tokenID.(float64); !ok || id != float64(id) {
		return false
	}
	return true
}
