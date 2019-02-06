package main

type Topic string

const (
	UpdateName Topic = "update-name"
)

var strToTopic = map[string]Topic{
	"update-name": UpdateName,
}
