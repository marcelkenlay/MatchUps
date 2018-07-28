package handlers

import (
  "github.com/pusher/pusher-http-go"
)

var PusherClient pusher.Client

func SetUpPusher() {
	PusherClient = pusher.Client{
	   AppId: "544934",
	   Key: "112f8743d26528cd9b7e",
	   Secret: "fd6bb26b65d969c08f93",
	   Cluster: "eu",
	}
}
