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

	randomWait()
}

func randomWait() {
	// Generate a random percentile (0-99)
	p := rand.Intn(100)
	var sleepDuration time.Duration
	switch {
	case p < 90:
		sleepDuration = time.Duration(rand.Intn(451)+0) * time.Millisecond
	case p < 95:
		sleepDuration = time.Duration(rand.Intn(450)+451) * time.Millisecond
	default:
		// Top 1%: 1001ms to 1200ms (simulate rare slow requests)
		sleepDuration = time.Duration(rand.Intn(501)+901) * time.Millisecond
	}
	time.Sleep(sleepDuration)
}
