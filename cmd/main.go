package main

import (
	"log"

	"github.com/xinliangnote/go-gin-api/cmd/command"
)

func main() {

	err := command.Execute()
	if err != nil {
		log.Fatalf("cmd.Execute error: %v", err)
	}
}
