package main

type base struct {
	Type   string `json:"type"`
	Action string `json:"action"`
	Token  string `json:"token"`
}

type RegisterRequest struct {
	base
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type response struct {
	base
	Ok      string `json:"ok"`
	Message string `json:"message"`
}

type ConnectResponse struct {
	response
}

type RegisterResponse struct {
	response
}

type LoginRequest struct {
	base
	SendID   int    `json:"sendid"`
	Password string `json:"password"`
}

type LoginResponse struct {
	response
}

type LogoutRequest struct {
	base
	SendID int `json:"sendid"`
}

type LogoutResponse struct {
	response
}

type SendMessageRequest struct {
	base
	SendID  int    `json:"sendid"`
	RecvID  int    `json:"recvid"`
	MsgType string `json:"msgtype"`
	Message string `json:"message"`
}
