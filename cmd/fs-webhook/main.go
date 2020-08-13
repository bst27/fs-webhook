package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/fsnotify/fsnotify"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	url := flag.String("url", "", "URL to receive http post requests (webhook target)")
	path := flag.String("path", "", "Filesystem path to be monitored for changes (file or folder)")
	flag.Parse()

	if *url == "" || *path == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

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

				go sendWebhook(*url, event.Name, strings.ToLower(event.Op.String()))
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(*path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func sendWebhook(url string, path string, action string) error {
	body := make(map[string]string)
	body["path"] = path
	body["action"] = action

	jsonBody, err := json.MarshalIndent(body, "", "   ")
	if err != nil {
		return err
	}

	_, err = http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		log.Println(err)
	}

	return err
}
