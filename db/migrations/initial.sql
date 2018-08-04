-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS 'messageType' (
	mtype TEXT PRIMARY KEY NOT NULL
);

REPLACE INTO messageType (mtype) VALUES ('text');
REPLACE INTO messageType (mtype) VALUES ('image');
REPLACE INTO messageType (mtype) VALUES ('video');

CREATE TABLE IF NOT EXISTS 'sourceType' (
	stype TEXT PRIMARY KEY NOT NULL
);

REPLACE INTO sourceType (stype) VALUES ('youtube');
REPLACE INTO sourceType (stype) VALUES ('vimeo');


CREATE TABLE IF NOT EXISTS 'users' (
	uid INTEGER PRIMARY KEY,
	username VARCHAR(50) NOT NULL, 
	password VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS 'messages' (
	msg_id INTEGER PRIMARY KEY,
	sender_id INTEGER,
	recipient_id INTEGER,
	type TEXT,
	created_on DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(sender_id) REFERENCES users(uid),
	FOREIGN KEY(recipient_id) REFERENCES users(uid),
	FOREIGN KEY(type) REFERENCES messageType(mtype)
);

CREATE TABLE IF NOT EXISTS 'texts' (
	msg_id INTEGER,
	msg TEXT NOT NULL,
	FOREIGN KEY(msg_id) REFERENCES messages(msg_id)
);

CREATE TABLE IF NOT EXISTS 'images' (
	msg_id INTEGER,
	width INTEGER NOT NULL,
	height INTEGER NOT NULL,
	i_url TEXT NOT NULL,
	FOREIGN KEY(msg_id) REFERENCES messages(msg_id)
);

CREATE TABLE IF NOT EXISTS 'videos' (
	msg_id INTEGER,
	source TEXT,
	v_url TEXT NOT NULL,
	FOREIGN KEY(msg_id) REFERENCES messages(msg_id)
	FOREIGN KEY(source) REFERENCES sourceType(stype)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS 'videos';
DROP TABLE IF EXISTS 'images';
DROP TABLE IF EXISTS 'texts';
DROP TABLE IF EXISTS 'messages';
DROP TABLE IF EXISTS 'users';
DROP TABLE IF EXISTS 'sourceType';
DROP TABLE IF EXISTS 'messageType';



