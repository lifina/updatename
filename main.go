package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const projectID = "positive-apex-202905"

func publish(ctx context.Context, topic Topic) error {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Errorf(ctx, "Failed to create the pubsub client %v", err)
		return err
	}

	baseDate := time.Date(2019, 2, 1, 0, 0, 0, 0, time.Local)
	now := time.Now()
	h := math.Trunc(now.Sub(baseDate).Hours())
	d := int64(math.Ceil(h / 24))
	message := fmt.Sprintf("無職%d日目", d)

	t := client.TopicInProject(string(topic), projectID)
	_, err = t.Publish(ctx, &pubsub.Message{
		Data: []byte(message),
	}).Get(ctx)
	if err != nil {
		log.Errorf(ctx, "Failed to publish a message %v", err)
		return err
	}
	log.Debugf(ctx, "Successed to publish a message %v", message)
	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	if topic, found := strToTopic[r.URL.Query().Get("topic")]; found {
		ctx := appengine.NewContext(r)
		if err := publish(ctx, topic); err != nil {
			http.Error(w, "internal", http.StatusInternalServerError)
		} else {
			fmt.Fprintf(w, "published")
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
