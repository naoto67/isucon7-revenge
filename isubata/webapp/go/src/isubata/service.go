package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

func getUser(userID int64) (*User, error) {
	u := User{}
	if err := db.Get(&u, "SELECT * FROM user WHERE id = ?", userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func addMessage(channelID, userID int64, content string) (int64, error) {
	res, err := db.Exec(
		"INSERT INTO message (channel_id, user_id, content, created_at) VALUES (?, ?, ?, NOW())",
		channelID, userID, content)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func queryMessages(chanID, lastID int64) ([]Message, error) {
	msgs := []Message{}
	err := db.Select(&msgs, "SELECT * FROM message WHERE id > ? AND channel_id = ? ORDER BY id DESC LIMIT 100",
		lastID, chanID)
	return msgs, err
}

func sessUserID(c echo.Context) int64 {
	sess, _ := session.Get("session", c)
	var userID int64
	if x, ok := sess.Values["user_id"]; ok {
		userID, _ = x.(int64)
	}
	return userID
}

func sessSetUserID(c echo.Context, id int64) {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		HttpOnly: true,
		MaxAge:   360000,
	}
	sess.Values["user_id"] = id
	sess.Save(c.Request(), c.Response())
}

const LettersAndDigits = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(n int) string {
	b := make([]byte, n)
	z := len(LettersAndDigits)

	for i := 0; i < n; i++ {
		b[i] = LettersAndDigits[rand.Intn(z)]
	}
	return string(b)
}

func register(name, password string) (int64, error) {
	salt := randomString(20)
	digest := fmt.Sprintf("%x", sha1.Sum([]byte(salt+password)))

	res, err := db.Exec(
		"INSERT INTO user (name, salt, password, display_name, avatar_icon, created_at)"+
			" VALUES (?, ?, ?, ?, ?, NOW())",
		name, salt, digest, name, "default.png")
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func jsonifyMessage(m Message) (map[string]interface{}, error) {
	u := User{}
	err := db.Get(&u, "SELECT name, display_name, avatar_icon FROM user WHERE id = ?",
		m.UserID)
	if err != nil {
		return nil, err
	}

	r := make(map[string]interface{})
	r["id"] = m.ID
	r["user"] = u
	r["date"] = m.CreatedAt.Format("2006/01/02 15:04:05")
	r["content"] = m.Content
	return r, nil
}

func createResponse(chanID, lastID int64) ([]map[string]interface{}, error) {
	response := make([]map[string]interface{}, 0)
	rows, err := db.Query("select m.*, u.name, u.display_name, u.avatar_icon from message m inner join user u on u.id = m.user_id where channel_id = ? and id > ? order by m.id desc limit 100", chanID, lastID)
	if err != nil {
		return response, err
	}
	for rows.Next() {
		var m Message
		var u User
		if err = rows.Scan(&m.ID, &m.ChannelID, &m.UserID, &m.Content, &m.CreatedAt, &u.Name, &u.DisplayName, &u.AvatarIcon); err != nil {
			return response, err
		}
		r := make(map[string]interface{})
		r["id"] = m.ID
		r["user"] = u
		r["date"] = m.CreatedAt.Format("2006/01/02 15:04:05")
		r["content"] = m.Content
		response = append(response, r)
	}
	return response, err
}
