package main

import (
	"log"

	"gopkg.in/toast.v1"
)

const (
	APP_ID = "Satisfactory launcher"
	TITLE  = "Copy saves in synch directory"
)

func pushToast(title, message string) {
	notification := toast.Notification{
		AppID:   APP_ID,
		Title:   title,
		Message: message,
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}
