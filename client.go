package main

import (
	"bufio"
	"net"
)

type Client struct {
	incoming chan string
	outgoing chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func NewClient(connetction net.Conn) *Client {
	writer := bufio.NewWriter(connetction)
	reader := bufio.NewReader(connetction)

	client := &Client{
		incoming: make(chan string),
		outgoing: make(chan string),
		reader:   reader,
		writer:   writer,
	}

	client.Listen()
	return client
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func (client *Client) Read() {
	for {
		line, _ := client.reader.ReadString('\n')
		client.incoming <- line
	}
}

func (client *Client) Write() {
	for data := range client.outgoing {
		client.writer.WriteString(data)
		client.writer.Flush()
	}
}
