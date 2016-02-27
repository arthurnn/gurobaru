package main

import (
	"io"
	"bufio"
	"strconv"
	"log"
	"net"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func FetchId() int64 {

	res, err := db.Exec("UPDATE sequence SET id=LAST_INSERT_ID(id+1)")
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()

	log.Println("POP", id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func handleRequest(conn net.Conn) {
	client := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	cmd, err := client.ReadString('\n')
	// handle http requests
	for err == nil && cmd != "\r\n" {
		log.Print(strconv.QuoteToASCII(cmd))
		cmd, err = client.ReadString('\n')
	}

	if err != nil && err != io.EOF {
		log.Fatal("Error reading buffer:", err.Error())
	}

	id := FetchId()

	client.WriteString(strconv.FormatInt(id, 10))
	client.Flush()
	return
}

func main() {
	var err error

	db, err = sql.Open("mysql", "root@/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	socket, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleRequest(conn)
	}
}
