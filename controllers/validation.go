package controllers

import (
	"errors"
	"github.com/dtsang7/ASAPP/models"
	"log"
)

const (
	ErrorMissingArgument    = "error missing argument"
	ErrorUsernameExceedSize = "error username exceed size limit"
	ErrorPasswordExceedSize = "error password exceed size limit"
	ErrorSourceNotSupported = "error video source not supported"
)

var errorMissingArgument = errors.New(ErrorMissingArgument)
var errorUsernameExceedSize = errors.New(ErrorUsernameExceedSize)
var errorPasswordExceedSize = errors.New(ErrorPasswordExceedSize)
var errorSourceNotSupported = errors.New(ErrorSourceNotSupported)

func ValidateUser(usr models.User) error {

	if usr.Username == "" || usr.Password == "" {
		log.Println(errorMissingArgument.Error())
		return errorMissingArgument
	}

	if len(usr.Username) > 50 {
		log.Println(errorUsernameExceedSize.Error())
		return errorUsernameExceedSize
	}

	if len(usr.Password) > 100 {
		log.Println(errorPasswordExceedSize.Error())
		return errorPasswordExceedSize
	}
	return nil
}

func ValidateSendMessage(msg models.Message) error {

	if msg.SenderID == 0 || msg.RecipientID == 0 || msg.Type == "" {
		log.Println(errorMissingArgument)
		return errorMissingArgument
	}

	switch msg.Type {
	case "text":
		if msg.Message == "" {
			log.Println(errorMissingArgument)
			return errorMissingArgument
		}
	case "image":
		if int(msg.Width.Int64) == 0 || int(msg.Height.Int64) == 0 || msg.Url.String == "" {
			log.Println(errorMissingArgument)
			return errorMissingArgument
		}
	case "video":
		if msg.Source.String == "" || msg.Url.String == "" {
			log.Println(errorMissingArgument)
			return errorMissingArgument
		}
		if msg.Source.String != "youtube" || msg.Source.String != "vimeo" {
			log.Println(errorSourceNotSupported)
			return errorSourceNotSupported
		}
	}
	return nil
}

func ValidateGetMessage(r Retriever) error {
	if r.RecipientId == 0 || r.MsgID == 0 {
		log.Println(errorMissingArgument)
		return errorMissingArgument
	}
	return nil
}
