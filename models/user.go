package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username`
	Password string `json:"password"`
}

const (
	ErrorUserExist        = "Username taken"
	ErrorWrongPassword    = "Wrong password"
	ErrorUserDoesNotExist = "User does not exist"
)

var errUserExist = errors.New(ErrorUserExist)
var errWrongPassword = errors.New(ErrorWrongPassword)
var errUserDoesNotExist = errors.New(ErrorUserDoesNotExist)

//insert new user into the database
func (dao *DAO) CreateUser(usr User) (int, error) {
	var exist bool
	username := usr.Username

	tx, err := dao.db.Begin()
	if err != nil {
		log.Println("error starting Tx", err.Error())
		return 0, err
	}

	query := "SELECT EXISTS (SELECT username FROM users WHERE username = ?)"

	err = tx.QueryRow(query, username).Scan(&exist)
	if err != nil {
		tx.Rollback()
		log.Println("error checking if username exist", err.Error())
		return 0, err
	}

	if exist {
		tx.Rollback()
		log.Println(errUserExist.Error())
		return 0, errUserExist
	}

	query = "INSERT INTO users (username, password) VALUES(?, ?)"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		log.Println("error hashing password", err.Error())
		return 0, err
	}

	res, err := tx.Exec(query, username, string(hashedPassword))
	if err != nil {
		tx.Rollback()
		log.Println("error inserting user", err.Error())
		return 0, err
	}

	//get user id
	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Println("error getting last inserted user id", err.Error())
		return 0, err
	}

	tx.Commit()
	return int(id), nil
}

//login user
func (dao *DAO) LoginUser(existingUser User) (int, error) {
	var uid int
	var dbPassword string

	tx, err := dao.db.Begin()
	if err != nil {
		log.Println("error starting Tx", err.Error())
		return 0, err
	}

	query := "SELECT uid, password FROM users WHERE username = ?"

	err = tx.QueryRow(query, existingUser.Username).Scan(&uid, &dbPassword)
	if err != nil {
		tx.Rollback()
		log.Println("error finding username", err.Error())
		return 0, errUserDoesNotExist
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(existingUser.Password))
	if err != nil {
		tx.Rollback()
		log.Println("error comparing passwords", err.Error())
		return 0, errWrongPassword
	}
	tx.Commit()
	return uid, nil
}
