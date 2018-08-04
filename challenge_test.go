package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dtsang7/ASAPP/config"
	"github.com/dtsang7/ASAPP/controllers"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	waitTime    = time.Second * 30
	contentType = "application/json"
)

var baseUrl string

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func assertNotEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Fatalf("%s != %s", a, b)
	}
}

func loginHelper(username, password string) (string, error) {
	payload := []byte(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
	resp, err := http.Post(baseUrl+"/login", contentType, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Fail to login")
	}
	var loginResp controllers.LoginResponse
	json.NewDecoder(resp.Body).Decode(&loginResp)
	return loginResp.Token, nil
}

func createUserHelper(username, password string) (int, error) {
	payload := []byte(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
	resp, err := http.Post(baseUrl+"/users", contentType, bytes.NewBuffer(payload))
	if err != nil {
		return 0, err
	}
	var cUser controllers.CreateUserResponse
	json.NewDecoder(resp.Body).Decode(&cUser)
	return cUser.Id, nil
}

func waitOnServerStart() {
	mult := time.Second * 1
	for mult < waitTime {
		resp, cErr := http.Post(baseUrl+"/check", contentType, nil)
		if cErr == nil {
			var health controllers.Health
			json.NewDecoder(resp.Body).Decode(&health)
			if health.Health == "ok" {
				return
			}
		}
		time.Sleep(mult)
		mult = mult * 2
	}
	log.Fatal("unable to connect")
}

func TestMain(m *testing.M) {
	// clean up test database before tests
	if err := os.Remove("challenge_test.db"); err != nil {
		log.Println("unable to remove file", err.Error())
	}
	// for server to load test config
	os.Setenv("ASAPP_ENV", "test")
	// starting app in a goroutine
	go main()

	// load test config
	config, _ := config.GetConfig("test")
	baseUrl = "http://" + config.Host + ":" + config.Port

	// wait for server to start
	waitOnServerStart()

	// run tests
	code := m.Run()

	// clean up
	os.Unsetenv("ASAPP_ENV")
	os.Exit(code)
}

func TestLoginUser(t *testing.T) {
	username := "test_login"
	password := "test_password"
	testLoginUserId, err := createUserHelper(username, password)
	if err != nil {
		t.Fatal(err)
	}
	// Test login with wrong password
	{
		payload := []byte(fmt.Sprintf(`{"username": "%s", "password": "bad_password"}`, username))
		resp, err := http.Post(baseUrl+"/login", contentType, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, resp.StatusCode, http.StatusBadRequest)
	}

	// Test login successfully
	{
		payload := []byte(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
		resp, err := http.Post(baseUrl+"/login", contentType, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		var loginResp controllers.LoginResponse
		json.NewDecoder(resp.Body).Decode(&loginResp)
		assertEqual(t, resp.StatusCode, http.StatusOK)
		assertEqual(t, loginResp.Id, testLoginUserId)
		assertNotEqual(t, loginResp.Token, "")
	}
}

func TestCreateUser(t *testing.T) {
	// Test create user with password missing
	{
		payload := []byte(`{"username": "testuser"}`)
		resp, err := http.Post(baseUrl+"/users", contentType, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, resp.StatusCode, http.StatusBadRequest)
	}

	// Test create user with username missing
	{
		payload := []byte(`{"password": "testpassword"}`)
		resp, err := http.Post(baseUrl+"/users", contentType, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, resp.StatusCode, http.StatusBadRequest)
	}

	// Test create user successfully
	{
		payload := []byte(`{"username": "testuser", "password": "testpassword"}`)
		resp, err := http.Post(baseUrl+"/users", contentType, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}

		var cUser controllers.CreateUserResponse
		json.NewDecoder(resp.Body).Decode(&cUser)
		assertEqual(t, resp.StatusCode, http.StatusOK)
		assertNotEqual(t, cUser.Id, 0)
	}
}

/*
Test Scenario:
1. Create two users(user1 and user2) and login
2. User1 send 3 messages to user2 and user2 send 1 message to user1
3. Check user1 can get back 3 messages
4. Check user2 can get back 1 message
*/
func TestMessages(t *testing.T) {
	// Users and credentials
	username1 := "testuser1"
	username2 := "testuser2"
	password1 := "testpassword1"
	password2 := "testpassword2"
	// Create two users
	user1id, uErr1 := createUserHelper(username1, password1)
	if uErr1 != nil {
		t.Fatal(uErr1)
	}
	user2id, uErr2 := createUserHelper(username2, password2)
	if uErr2 != nil {
		t.Fatal(uErr2)
	}
	// Login users
	token1, tErr1 := loginHelper(username1, password1)
	if tErr1 != nil {
		t.Fatal(tErr1)
	}
	token2, tErr2 := loginHelper(username2, password2)
	if tErr2 != nil {
		t.Fatal(tErr2)
	}
	bearer1 := "Bearer " + token1
	bearer2 := "Bearer " + token2
	client := &http.Client{}

	// Test send message missing recipient id
	{
		payload := []byte(fmt.Sprintf(`{"sender": %d, "content":{"type": "text", "text": "test message"}}`, user1id))
		req, _ := http.NewRequest("POST", baseUrl+"/messages", bytes.NewBuffer(payload))
		req.Header.Set("Authorization", bearer1)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, resp.StatusCode, http.StatusBadRequest)
	}

	// Test send message missing type
	{
		payload := []byte(fmt.Sprintf(`{"sender": %d, "recipient": %d, "content":{"text: "test message"}}`, user1id, user2id))
		req, _ := http.NewRequest("POST", baseUrl+"/messages", bytes.NewBuffer(payload))
		req.Header.Set("Authorization", bearer1)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, resp.StatusCode, http.StatusBadRequest)
	}

	// Test send text message successfully
	{
		payload := []byte(fmt.Sprintf(`{"sender": %d, "recipient": %d, "content":{"type": "text", "text": "test message 1"}}`, user1id, user2id))
		req, err := http.NewRequest("POST", baseUrl+"/messages", bytes.NewBuffer(payload))

		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", bearer1)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, resp.StatusCode, http.StatusOK)
		var sr controllers.SendMessageResponse
		json.NewDecoder(resp.Body).Decode(&sr)
		assertNotEqual(t, sr.Id, 0)
		assertNotEqual(t, sr.Timestamp, "")
	}
	// Test send text message successfully
	{
		payload := []byte(fmt.Sprintf(`{"sender": %d, "recipient": %d, "content":{"type": "text", "text": "test message 2"}}`, user2id, user1id))
		req, err := http.NewRequest("POST", baseUrl+"/messages", bytes.NewBuffer(payload))

		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", bearer2)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, resp.StatusCode, http.StatusOK)
		var sr controllers.SendMessageResponse
		json.NewDecoder(resp.Body).Decode(&sr)
		assertNotEqual(t, sr.Id, 0)
		assertNotEqual(t, sr.Timestamp, "")
	}
	// Test send image message successfully
	{
		payload := []byte(fmt.Sprintf(`{"sender": %d, "recipient": %d, "content":{"type": "image", "width": 10, "height": 10, "url": "http://some_image_url"}}`, user1id, user2id))
		req, _ := http.NewRequest("POST", baseUrl+"/messages", bytes.NewBuffer(payload))
		req.Header.Set("Authorization", bearer1)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, resp.StatusCode, http.StatusOK)
		var sr controllers.SendMessageResponse
		json.NewDecoder(resp.Body).Decode(&sr)
		assertNotEqual(t, sr.Id, 0)
		assertNotEqual(t, sr.Timestamp, "")
	}
	// Test send video message successfully
	{
		payload := []byte(fmt.Sprintf(`{"sender": %d, "recipient": %d, "content":{"type": "video", "source": "youtube", "url": "http://some_video_url"}}`, user1id, user2id))
		req, _ := http.NewRequest("POST", baseUrl+"/messages", bytes.NewBuffer(payload))
		req.Header.Set("Authorization", bearer1)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, resp.StatusCode, http.StatusOK)
		var sr controllers.SendMessageResponse
		json.NewDecoder(resp.Body).Decode(&sr)
		assertNotEqual(t, sr.Id, 0)
		assertNotEqual(t, sr.Timestamp, "")
	}
	// Test get messages for user1 successfully
	{
		req, _ := http.NewRequest("GET", fmt.Sprintf(baseUrl+"/messages?recipient=%d&start=1", user2id), nil)
		req.Header.Set("Authorization", bearer1)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, resp.StatusCode, http.StatusOK)
		var gr controllers.GetMessagesResponse
		json.NewDecoder(resp.Body).Decode(&gr)
		assertEqual(t, len(gr.Messages), 3)
		// Check first message
		message0 := gr.Messages[0]
		assertEqual(t, message0.SenderID, user1id)
		assertEqual(t, message0.RecipientID, user2id)
		assertEqual(t, message0.Content.Type, "text")
		assertEqual(t, message0.Content.Text, "test message 1")

		// Check second message
		message1 := gr.Messages[1]
		assertEqual(t, message1.SenderID, user1id)
		assertEqual(t, message1.RecipientID, user2id)
		assertEqual(t, message1.Content.Type, "image")
		assertEqual(t, message1.Content.Width, 10)
		assertEqual(t, message1.Content.Height, 10)
		assertEqual(t, message1.Content.Url, "http://some_image_url")
		// Check third message
		message2 := gr.Messages[2]
		assertEqual(t, message2.SenderID, user1id)
		assertEqual(t, message2.RecipientID, user2id)
		assertEqual(t, message2.Content.Type, "video")
		assertEqual(t, message2.Content.Source, "youtube")
		assertEqual(t, message2.Content.Url, "http://some_video_url")
	}
	// Test get messages for user1 successfully
	{
		req, _ := http.NewRequest("GET", fmt.Sprintf(baseUrl+"/messages?recipient=%d&start=1", user1id), nil)
		req.Header.Set("Authorization", bearer2)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, resp.StatusCode, http.StatusOK)
		var gr controllers.GetMessagesResponse
		json.NewDecoder(resp.Body).Decode(&gr)
		assertEqual(t, len(gr.Messages), 1)
		// Check first message
		message0 := gr.Messages[0]
		assertEqual(t, message0.SenderID, user2id)
		assertEqual(t, message0.RecipientID, user1id)
		assertEqual(t, message0.Content.Type, "text")
		assertEqual(t, message0.Content.Text, "test message 2")
	}

}
