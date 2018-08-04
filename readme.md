
# PROJECT: ASAPP chat backend

This project is a basic chat backend that implements the following:
	* Create user
	* Login user
	* Send message
	* Fetch messages

#Install
=========================

##On OS X using homebrew:
	$ brew install go

	# after Go is installed run the following to get SQLite3 driver:
	$ go get -u github.com/mattn/go-sqlite3
	# install dependencies
	$ brew install dep
	$ brew upgrade dep
	$ dep ensure

#Run server
=========================

	$ go run challenge.go

	#to run test:
	$ go test -v

#Examples
=========================

##Check system
	$ curl -s -d '' -XPOST http://localhost:8080/check
##Response:
{"health":"ok"}

##Create user
#Required: username, password
	$ curl -XPOST -d '{"username": "testuser", "password": "testpassword"}' http://localhost:8080/users
##Rsponse:
{"id":1}

##Login user
#Required: username, password
	$ curl -XPOST -d '{"username": "testuser", "password": "testpassword"}' http://localhost:8080/login
##Response:
{"Id":1,"Token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzMzNzY3MzQsImlkIjoxfQ.nzl2yeZMFZdPz9XE26yoNJlfpoUjIvUOmsaraclsKw4"}

##Send message
#Required: token, senderID, recipientID, type = ('text': text), ('image': width, height, url), ('video': source, url)
	$ export TKN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzMzNzY3MzQsImlkIjoxfQ.nzl2yeZMFZdPz9XE26yoNJlfpoUjIvUOmsaraclsKw4"
	$ curl -XPOST -H "Authorization: Bearer $TKN" -d '{"sender": 1, "recipient": 2, "content":{"type": "text", "text": "Test Message"}}' http://localhost:8080/messages
##Resonse:
{"Id":1,"Timestamp":"2018-08-04T05:06:22Z"}

##Fetch messages
#Required: token, recipientID 
#Optional: limit (default is 100)
	$ curl -XGET -H "Authorization: Bearer $TKN" "http://localhost:8080/messages?recipient=2&start=1&limit=1"
##Response:
{"messages":[{"id":1,"timestamp":"2018-08-04T05:06:22Z","sender":1,"recipient":2,"content":{"type":"text","text":"Test Message"}}]}

