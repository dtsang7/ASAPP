package controllers

import (
	"errors"
	"github.com/dtsang7/ASAPP/models"
	"log"
	"net/http"
	"strconv"
)

const (
	ErrorMissingArgument    = "error missing argument"
	ErrorUsernameExceedSize = "error username exceed size limit"
	ErrorPasswordExceedSize = "error password exceed size limit"
	ErrorSourceNotSupported = "error video source not supported"
	ErrorTypeNotSupported   = "error type of message not supported"
)

var errorMissingArgument = errors.New(ErrorMissingArgument)
var errorUsernameExceedSize = errors.New(ErrorUsernameExceedSize)
var errorPasswordExceedSize = errors.New(ErrorPasswordExceedSize)
var errorSourceNotSupported = errors.New(ErrorSourceNotSupported)
var errorTypeNotSupported = errors.New(ErrorTypeNotSupported)

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

func ValidateSendMessage(req Message) error {
	if req.SenderID <= 0 || req.RecipientID <= 0 || req.Content.Type == "" {
		log.Println(errorMissingArgument)
		return errorMissingArgument
	}

	switch req.Content.Type {
	case "text":
		if req.Content.Text == "" {
			log.Println(errorMissingArgument)
			return errorMissingArgument
		}
	case "image":
		if req.Content.Width <= 0 || req.Content.Height <= 0 || req.Content.Url == "" {
			log.Println(errorMissingArgument)
			return errorMissingArgument
		}
	case "video":
		if req.Content.Source == "" || req.Content.Url == "" {
			log.Println(errorMissingArgument)
			return errorMissingArgument
		}
		if req.Content.Source != "youtube" && req.Content.Source != "vimeo" {
			log.Println(errorSourceNotSupported)
			return errorSourceNotSupported
		}
	default:
		log.Println(errorTypeNotSupported)
		return errorTypeNotSupported
	}
	return nil
}

// Parse int from string, expect greater than zero
func parsePositiveInt(str string) (int, error) {
	intVal, parseErr := strconv.ParseInt(str, 10, 64)
	if parseErr != nil {
		return int(intVal), parseErr
	}
	if intVal <= 0 {
		return 0, errorMissingArgument
	}
	return int(intVal), nil
}

// Parse query parameters and validate them
func ParseAndValidateGetMessageRequest(r *http.Request) (req GetMessagesRequest, err error) {
	params := r.URL.Query()
	// parse recipient, required
	if val, parseErr := parsePositiveInt(params.Get("recipient")); parseErr == nil {
		req.RecipientID = val
	} else {
		err = parseErr
		return
	}
	// parse start, required
	if val, parseErr := parsePositiveInt(params.Get("start")); parseErr == nil {
		req.StartMsgID = val
	} else {
		err = parseErr
		return
	}
	// parse limit, optional
	if val, parseErr := parsePositiveInt(params.Get("limit")); parseErr == nil {
		req.Limit = val
	} else {
		req.Limit = 100
	}
	return
}
