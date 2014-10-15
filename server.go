package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/bitly/go-simplejson"
)

type ChatServer struct {
	clients  map[string]*Client
	joins    chan net.Conn
	incoming chan string
	outgoing chan string
}

func NewChatServer() *ChatServer {
	chatServer := &ChatServer{
		clients:  make(map[string]*Client, 0),
		joins:    make(chan net.Conn),
		incoming: make(chan string),
		outgoing: make(chan string),
	}

	chatServer.Listen()

	return chatServer
}

func (chatServer *ChatServer) Listen() {
	go func() {
		for {
			select {
			case data := <-chatServer.incoming:
				chatServer.DispatchMessage(data)
			case data := <-chatServer.joins:
				chatServer.Join(data)
			}
		}
	}()
}

func (chatServer *ChatServer) DispatchMessage(message string) {
	js, _ := simplejson.NewJson([]byte(message))
	messageType := js.Get("type").MustString()
	if messageType == "request" {
		action := js.Get("action").MustString()
		switch action {
		case "send":
			go chatServer.DealSendMessage(message)
		case "register":
			go chatServer.DealRegister(message)
		case "login":
			go chatServer.DealLogin(message)
		case "logout":
			go chatServer.DealLogout(message)
		case "addbuddy":
			go chatServer.DealAddBuddy(message)
		case "deletebuddy":
			go chatServer.DealDeleteBuddy(message)
		case "getbuddylist":
			go chatServer.DealGetBuddyList(message)
		}
	} else if messageType == "response" {

	}
}

func (chatServer *ChatServer) Join(conn net.Conn) {
	client := NewClient(conn)
	uuid := NewUUID()
	chatServer.clients[uuid] = client
	go chatServer.DealConnect(client, uuid)
	go func() {
		for {
			chatServer.incoming <- <-client.incoming
		}
	}()
}

func (chatServer *ChatServer) SendMessage(message []byte, uuid string) {
	chatServer.clients[uuid].outgoing <- string(message)
}

func (chatServer *ChatServer) DealConnect(client *Client, uuid string) {
	var resp ConnectResponse
	resp.Type = "response"
	resp.Action = "connect"
	resp.Ok = "ok"
	resp.Token = uuid
	response, _ := json.Marshal(resp)
	chatServer.SendMessage(response, uuid)
}

func (chatServer *ChatServer) DealRegister(message string) {
	var r RegisterRequest
	json.Unmarshal([]byte(message), &r)
	id, err := CreateUser(r.Nickname, r.Password)
	if err != nil {
		fmt.Println(err)
	}
	var rsp RegisterResponse
	rsp.Action = r.Action
	rsp.Type = "response"
	rsp.Token = r.Token
	rsp.Ok = "ok"
	rsp.Message = strconv.FormatInt(id, 10)
	response, _ := json.Marshal(rsp)
	go chatServer.SendMessage(response, r.Token)
}

func (chatServer *ChatServer) DealLogin(message string) {
	var r LoginRequest
	json.Unmarshal([]byte(message), &r)
	var resp LoginResponse
	resp.Action = r.Action
	resp.Type = "response"
	resp.Token = r.Token

	if CheckLogin(r.SendID, r.Password) {
		UpdateUserUUID(r.Token, r.SendID)
		resp.Ok = "ok"
	} else {
		resp.Ok = "no"
	}
	response, _ := json.Marshal(resp)
	go chatServer.SendMessage(response, resp.Token)
}

func (chatServer *ChatServer) DealLogout(message string) {
	var r LogoutRequest
	json.Unmarshal([]byte(message), &r)
	var rsp LogoutResponse
	rsp.Action = r.Action
	rsp.Type = "response"
	rsp.Token = r.Token
	rsp.Ok = "ok"
	response, _ := json.Marshal(rsp)
	delete(chatServer.clients, r.Token)
	go chatServer.SendMessage(response, rsp.Token)
}

func (chatServer *ChatServer) DealSendMessage(message string) {
	var r SendMessageRequest
	json.Unmarshal([]byte(message), &r)
	if r.MsgType == "single" {
		uuid := GetUserUUID(r.RecvID)
		InsertMessage(r.SendID, r.RecvID, r.MsgType, r.Message)
		go chatServer.SendMessage([]byte(message), uuid)
	} else if r.MsgType == "group" {

	}
}

func (chatServer *ChatServer) DealAddBuddy(message string) {
	var r AddBuddyRequest
	json.Unmarshal([]byte(message), &r)
	AddBuddy(r.SendID, r.BuddyID)
	AddBuddy(r.BuddyID, r.SendID)
	var rsp AddBuddyResponse
	rsp.Action = r.Action
	rsp.Type = "response"
	rsp.Ok = "ok"
	rsp.Token = r.Token
	response, _ := json.Marshal(rsp)
	go chatServer.SendMessage(response, r.Token)
}

func (chatServer *ChatServer) DealDeleteBuddy(message string) {
	var r DeleteBuddyRequest
	json.Unmarshal([]byte(message), &r)
	DeleteBuddy(r.SendID, r.BuddyID)
	DeleteBuddy(r.BuddyID, r.SendID)
	var rsp DeleteBuddyResponse
	rsp.Action = r.Action
	rsp.Type = "response"
	rsp.Ok = "ok"
	rsp.Token = r.Token
	response, _ := json.Marshal(rsp)
	go chatServer.SendMessage(response, r.Token)
}

func (chatServer *ChatServer) DealGetBuddyList(message string) {
	var r GetBuddyListRequest
	json.Unmarshal([]byte(message), &r)
	list := GetBuddyList(r.SendID)
	var rsp GetBuddyListResponse
	rsp.Action = r.Action
	rsp.Type = "response"
	rsp.Token = r.Token
	rsp.List = list
	response, _ := json.Marshal(rsp)
	go chatServer.SendMessage(response, r.Token)
}

func (chatServer *ChatServer) DealGetGroupList(message string) {

}
