package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func CreateUser(nickname string, password string) (int64, error) {
	fmt.Println(APP_VER)
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

func UpdateUserUUID(uuid string, user_id int) {
	stmt, _ := db.Prepare("update chat.user set last_token = ? where id = ?")
	defer stmt.Close()
	stmt.Exec(uuid, user_id)
}

func InsertMessage(send_id int, receive_id, target_type string, message string) {

}

func GetBuddyList(user_id int) {

}

func GetGroupList(user_id int) {

}
