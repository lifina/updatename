package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const projectID = "positive-apex-202905"

func publish(ctx context.Context, topic Topic) (string, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Errorf(ctx, "Failed to create the pubsub client %v", err)
		return "", err
	}

	baseDate := time.Date(2019, 2, 1, 0, 0, 0, 0, time.Local)
	now := time.Now()
	d := now.Sub(baseDate).Hours() / 24

	t := client.TopicInProject(string(topic), projectID)
	defer t.Stop()
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(fmt.Sprintf("無職%d日目", d)),
	})

	id, err := result.Get(ctx)
	if err != nil {
		log.Errorf(ctx, "Failed to publish a message %v", err)
		return "", err
	}
	return id, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	if topic, found := strToTopic[r.URL.Query().Get("topic")]; found {
		ctx := appengine.NewContext(r)
		if id, err := publish(ctx, topic); err != nil {
			http.Error(w, "internal", http.StatusInternalServerError)
		} else {
			fmt.Fprintf(w, "published: %v", id)
		}
	} else {
		http.Error(w, "bad topic", http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/publish", publishHandler)
	appengine.Main()
}
