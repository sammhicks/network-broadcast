package main

import (
	"net/http"

	"github.com/sammhicks/network-broadcast/messages"
)

type messageHandler struct {
	mc chan<- messages.Message
}

func (eh *messageHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	eh.mc <- messages.Message{
		Sender:  req.PostFormValue("name"),
		Content: req.PostFormValue("message")}

	http.Redirect(rw, req, "/", http.StatusTemporaryRedirect)
}
