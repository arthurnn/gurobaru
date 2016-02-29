package main

import (
	"flag"
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

	// log.Println("POP", id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func handleRequest(conn net.Conn) {
	client := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	cmd, err := client.ReadString('\n')

	for err == nil || err == io.EOF {
		// log.Print(strconv.QuoteToASCII(cmd))
		switch cmd {
		case "GET / HTTP/1.1\r\n", "FETCH ID\r\n", "FETCH ID\n", "FETCH ID" :
			id := FetchId()
			client.WriteString(strconv.FormatInt(id, 10))
			client.WriteString("\n")
			client.Flush()
		case "\r\n":
			return
		}
		cmd, err = client.ReadString('\n')
		if err == io.EOF {
			return
		}
	}

	if err != nil && err != io.EOF {
		log.Fatal("Error reading buffer:", err.Error())
	}


	return
}

func main() {
	var err error
	user := flag.String("user", "root", "MySQL user")
	password := flag.String("password", "", "MySQL password")
	db_name := flag.String("database", "test", "MySQL database name")
	host := flag.String("host", "127.0.0.1", "MySQL host")
	port := flag.String("port", "3306", "MySQL port")
	flag.Parse()
	hostport := net.JoinHostPort(*host, *port)

	db, err = sql.Open("mysql", *user+":"+*password+"@tcp("+hostport+")/"+*db_name)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE sequence (id INT NOT NULL)")
	if err == nil {
		db.Exec("INSERT INTO sequence VALUES (0)")
	}

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
