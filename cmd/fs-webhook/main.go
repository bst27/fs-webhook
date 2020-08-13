package main

import (
	"bytes"
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"log"
	"net/http"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				go sendWebhook(event.Name, event.Op.String())
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add("/tmp/foo")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func sendWebhook(path string, action string) error {
	body := make(map[string]string)
	body["path"] = path
	body["action"] = action

	jsonBody, err := json.MarshalIndent(body, "", "   ")
	if err != nil {
		return err
	}

	_, err = http.Post("https://webhook.site/5405b754-19a0-4048-8d2e-d5f1de72b037", "application/json", bytes.NewReader(jsonBody))
	return err
}
