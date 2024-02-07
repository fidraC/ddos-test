package main

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

var (
	Client = &http.Client{}
)

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func sendRequest(e Endpoint, host string) error {
	var body io.Reader = nil
	var params string = ""
	if e.Method != "GET" && e.Payload != nil {
		body = strings.NewReader(e.Payload.Data)
	}
	if e.Payload != nil {
		params = "?" + e.Payload.Params
		if strings.Contains(e.Payload.Params, "${!RANDOM}") {
			params = strings.Replace(e.Payload.Params, "${!RANDOM}", randomString(5), -1)
		}
	}
	req, _ := http.NewRequest(e.Method, host+e.Path+params, body)
	if e.Headers != nil {
		for k, v := range e.Headers {
			req.Header.Set(k, v)
		}
	}
	_, err := Client.Do(req)
	return err
}

func Punch(e Endpoint, threads int, host string, cancel chan struct{}) {
	ch := make(chan struct{}, threads)
	var wg sync.WaitGroup

	// Loop forever until something is sent to the cancel channel
	for {
		select {
		case <-cancel:
			return
		default:
			ch <- struct{}{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				sendRequest(e, host)
				<-ch
			}()
		}
	}

}
