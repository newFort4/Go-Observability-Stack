package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	logger.Info("Request from user", "name", name)

	fmt.Fprintf(w, "Hello, %s!\n", name)
	r.Response = new(http.Response)
	if rand.Int()%2 == 0 {
		r.Response.StatusCode = http.StatusInternalServerError
	} else {
		r.Response.StatusCode = http.StatusOK
	}

	r.Response.ContentLength = int64(len(name) + 8)

	// Generate a random duration between 100ms and 500ms
	sleepDuration := time.Duration(rand.Intn(401)+100) * time.Millisecond
	time.Sleep(sleepDuration)
}
