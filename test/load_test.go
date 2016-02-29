package main

import (
	"os/exec"
	"log"
)

func main() {


	for i := 0; i < 100; i++ {
		log.Println("process", i)

		cmd := exec.Command("./client")
		cmd.Start()
	}
}
