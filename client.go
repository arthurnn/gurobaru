package main
import (
	"net"
	"log"
	"bufio"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	defer conn.Close()

	if err != nil {
		log.Fatal(err)
	}

	client := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	for {
		client.WriteString("FETCH ID\n")
		client.Flush()

		cmd, err := client.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		log.Print(cmd)
	}
}
