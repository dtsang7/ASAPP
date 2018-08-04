package models

import (
	"database/sql"
	"errors"
	"log"
)

type GetMessage struct {
	MsgID       int    `json:"msg_id"`
	SenderID    int    `json:"sender_id"`
	RecipientID int    `json:"recipient_id"`
	Type        string `json:"type"`
	Message     string `json:"message, omitempty"`
	Width       int    `json:"width, omitempty"`
	Height      int    `json:"height, omitempty"`
	Source      string `json:"source, omitempty"`
	Url         string `json:"url, omitempty"`
	TimeStamp   string `json:"created_on"`
}

type Message struct {
	MsgID       int            `db:"msg_id" json:"msg_id"`
	SenderID    int            `db:"sender_id" json:"sender_id"`
	RecipientID int            `db:"recipient_id" json:"recipient_id"`
	Type        string         `db:"type"`
	Message     sql.NullString `db:"msg"`
	Width       sql.NullInt64  `db:"width"`
	Height      sql.NullInt64  `db:"height"`
	Url         sql.NullString `db:"i_url, v_url"`
	Source      sql.NullString `db:"source"`
	TimeStamp   string         `db:"created_on`
}

const (
	ErrorCreatingMessage         = "Error creating message"
	ErrorMessageTypeNotSupported = "Message type not supported"
)

var errorCreateMessage = errors.New(ErrorCreatingMessage)
var errorMessageTypeNotSupported = errors.New(ErrorMessageTypeNotSupported)

func (dao *DAO) SendMessage(msg Message) (int, string, error) {
	var timeStamp string
	mtype := msg.Type

	tx, err := dao.db.Begin()
	if err != nil {
		log.Print("error starting Tx", err.Error())
		return 0, timeStamp, err
	}

	//store message info
	query := "INSERT INTO messages (sender_id, recipient_id, type) VALUES (?, ?, ?)"
	res, err := tx.Exec(query, msg.SenderID, msg.RecipientID, mtype)
	if err != nil {
		tx.Rollback()
		log.Println("error inserting message into messages table", err.Error())
		return 0, timeStamp, errorCreateMessage
	}
	//retrieve message id
	msgID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Println("error retrieving last inserted message id", err.Error())
		return 0, timeStamp, err
	}
	//retrieve timestamp
	query = "SELECT created_on FROM messages WHERE msg_id = ?"
	err = tx.QueryRow(query, msgID).Scan(&timeStamp)
	if err != nil {
		tx.Rollback()
		log.Println("error retrieving timestamp", err.Error())
		return 0, timeStamp, err
	}

	// store messages based on type
	switch mtype {
	case "text":
		query := "INSERT INTO texts (msg_id, msg) VALUES (?, ?)"

		_, err = tx.Exec(query, msgID, msg.Message)
		if err != nil {
			tx.Rollback()
			log.Println("error inserting message into texts table", err.Error())
			return 0, timeStamp, err
		}
	case "image":
		query := "INSERT INTO images (msg_id, width, height, i_url) VALUES (?, ?, ?, ?)"

		_, err = tx.Exec(query, msgID, msg.Width, msg.Height, msg.Url)
		if err != nil {
			tx.Rollback()
			log.Println("error inserting image into images table", err.Error())
			return 0, timeStamp, err
		}
	case "video":
		query := "INSERT INTO videos (msg_id, source, v_url) VALUES (?, ?, ?)"

		_, err = tx.Exec(query, msgID, msg.Source, msg.Url)
		if err != nil {
			tx.Rollback()
			log.Println("error inserting video into videos table", err.Error())
			return 0, timeStamp, err
		}
	default:
		tx.Rollback()
		log.Println(errorMessageTypeNotSupported.Error())
		return 0, timeStamp, errorMessageTypeNotSupported
	}

	tx.Commit()
	return int(msgID), timeStamp, nil
}

func (dao *DAO) GetMessages(recipient_id int, msg_id int, limit int) ([]Message, error) {

	tx, err := dao.db.Begin()
	if err != nil {
		log.Println("error starting Tx", err.Error())
		return nil, err
	}

	query := `SELECT messages.msg_id, sender_id, recipient_id, type, msg, width, height, i_url, v_url, source, created_on
			  FROM messages
			  LEFT JOIN texts ON messages.msg_id = texts.msg_id
			  LEFT JOIN images ON messages.msg_id = images.msg_id
			  LEFT JOIN videos ON messages.msg_id = videos.msg_id
			  WHERE recipient_id = ? AND messages.msg_id >= ?
			  ORDER BY messages.msg_id
			  LIMIT ?`

	res, err := tx.Query(query, recipient_id, msg_id, limit)
	if err != nil {
		tx.Rollback()
		log.Println("error retrieving messages", err.Error())
		return nil, err
	}
	defer res.Close()

	msgs := []Message{}

	for res.Next() {
		var msg Message
		var imageUrl sql.NullString
		var videoUrl sql.NullString
		err := res.Scan(&msg.MsgID, &msg.SenderID, &msg.RecipientID, &msg.Type, &msg.Message, &msg.Width, &msg.Height, &imageUrl, &videoUrl, &msg.Source, &msg.TimeStamp)
		if err != nil {
			tx.Rollback()
			log.Println("error scanning messages", err.Error())
			return nil, err
		}
		if imageUrl.Valid {
			msg.Url = imageUrl
		} else if videoUrl.Valid {
			msg.Url = videoUrl
		}

		msgs = append(msgs, msg)
	}
	err = res.Err()
	if err != nil {
		tx.Rollback()
		log.Println("error occured during iteration", err.Error())
		return nil, err
	}
	tx.Commit()
	return msgs, nil
}
