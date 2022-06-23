package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/xedflix/auto-approval-system/dataops"
	"github.com/xedflix/auto-approval-system/server"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Println(err.Error())
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dataops.SetupBase("resources/config.yaml", "tmp/")

}

func main() {
	server.RunServer()
}
