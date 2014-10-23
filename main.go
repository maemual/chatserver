package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/codegangsta/cli"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const APP_VER = "0.1.0"

var (
	db     *sql.DB
	logger *log.Logger
)

func main() {
	app := cli.NewApp()
	app.Name = "chatserver"
	app.Usage = "A server of chat system"
	app.Author = "maemual"
	app.Version = APP_VER
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port",
			Value: "9999",
			Usage: "The port which server listen",
		},
	}
	app.Action = func(c *cli.Context) {
		db, _ = sql.Open("mysql", "root:xxxxxxx@/chat")
		defer db.Close()

		chatServer := NewChatServer()

		addr, _ := net.ResolveTCPAddr("tcp", ":"+c.String("port"))
		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			panic(err)
		}
		defer listener.Close()
		for {
			conn, _ := listener.AcceptTCP()
			conn.SetNoDelay(true)
			conn.SetKeepAlive(true)
			conn.SetKeepAlivePeriod(45 * time.Second)
			chatServer.joins <- conn
		}
	}
	app.Run(os.Args)
}
