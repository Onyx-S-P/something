package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xedflix/auto-approval-system/policy"
)

func RunServer() {

	r := gin.Default()

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	r.POST("/event", handleEvent)

	log.Println("Listening on port: ", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Println(err.Error())
	}

}

func handleEvent(ctx *gin.Context) {

	event := ctx.GetHeader("X-GitHub-Event")
	log.Println(event)
	payload, _ := ioutil.ReadAll(ctx.Request.Body)

	//TODO: improve error handeling
	if policy.VerifySignature(payload, ctx.GetHeader("X-Hub-Signature-256")) {
		fmt.Println("Signature verified")
		switch event {
		case "pull_request":
			policy.PullRequestEvent(payload)
		default:
		}
	}
}
