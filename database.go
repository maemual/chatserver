package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func CreateUser(nickname string, password string) (int64, error) {
	pd := PasswordHash(password)
	stmt, err := db.Prepare("insert into chat.user (nickname, password) values (?, ?)")
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(nickname, pd)
	if err != nil {
		log.Println(err)
	}
	return result.LastInsertId()
}

func CheckLogin(user_id int, password string) bool {
	var pd string
	err := db.QueryRow("select password from chat.user where id = ?", user_id).Scan(&pd)
	if err != nil {
		panic(err)
		return false
	}
	p := PasswordHash(password)
	if p == pd {
		return true
	} else {
		return false
	}
}

func GetUserUUID(user_id int) string {
	var uuid string
	err := db.QueryRow("SELECT last_token FROM chat.user WHERE id = ?", user_id).Scan(&uuid)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		return uuid
	}
	return ""
}

func GetUserName(user_id int) string {
	var name string
	err := db.QueryRow("SELECT nickname FROM chat.user WHERE id = ?", user_id).Scan(&name)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		return name
	}
	return ""
}

func UpdateUserUUID(uuid string, user_id int) {
	stmt, _ := db.Prepare("update chat.user set last_token = ? where id = ?")
	defer stmt.Close()
	stmt.Exec(uuid, user_id)
}

func InsertMessage(send_id, receive_id int, target_type string, message string) {
	stmt, _ := db.Prepare("insert into chat.message (send_id, receiver_id, type, message, time) values(?, ?, ?, ?, NOW())")
	defer stmt.Close()
	stmt.Exec(send_id, receive_id, target_type, message)
}

func GetBuddyList(user_id int) (result []map[string]interface{}) {
	rows, err := db.Query("select buddy_id from chat.buddy where user_id = ?", user_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var x int
		rows.Scan(&x)
		tmp := map[string]interface{}{
			"id":   x,
			"name": GetUserName(x),
		}
		result = append(result, tmp)
	}
	return
}

func GetTalkMessage(send_id, recv_id int) (result []map[string]interface{}) {
	rows, err := db.Query("select message, time from chat.message where send_id = ? and receiver_id = ? order by time DESC limit 10", send_id, recv_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var msg, t string
		rows.Scan(&msg, &t)
		tmp := map[string]interface{}{
			"message": msg,
			"time":    t,
		}
		result = append(result, tmp)
	}
	return
}

func GetGroupList(user_id int) {

}

func AddBuddy(user_id, friend_id int) {
	var tmp int
	err := db.QueryRow("select user_id from chat.buddy where user_id = ? and buddy_id = ?", user_id, friend_id).Scan(&tmp)
	if err == sql.ErrNoRows {
		err = db.QueryRow("select id from chat.user where id = ?", user_id).Scan(&tmp)
		if err == nil {
			err = db.QueryRow("select id from chat.user where id = ?", friend_id).Scan(&tmp)
			if err == nil {
				stmt, _ := db.Prepare("insert into chat.buddy (user_id, buddy_id) values(?, ?)")
				defer stmt.Close()
				stmt.Exec(user_id, friend_id)
			}
		}
	}
}

func DeleteBuddy(user_id, friend_id int) {
	stmt, _ := db.Prepare("delete from chat.buddy where user_id = ? and buddy_id = ?")
	defer stmt.Close()
	stmt.Exec(user_id, friend_id)
}
