package activitymonitor

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	activitylogger "github.com/dhanushs3366/activity-logger"
)

type LoggedActivity struct {
	Key          int `json:"all_keys"`
	MiddleClicks int `json:"middle_clicks"`
	RightClicks  int `json:"right_clicks"`
	LeftClicks   int `json:"left_clicks"`
	ExtraClicks  int `json:"extra_clicks"`
}

func SetupLoggers(devPaths []string) ([]activitylogger.Keylogger, error) {
	var keyloggers []activitylogger.Keylogger
	for _, devPath := range devPaths {
		keylogger, err := activitylogger.New(devPath)
		if err != nil {
			continue
		}
		keyloggers = append(keyloggers, *keylogger)
	}

	return keyloggers, nil
}

func sendRequest(method, url string, buff *bytes.Buffer) error {
	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Printf("Response: %s\n", res.Status)
	return nil
}

func SendDataFromLogger(keylogger activitylogger.Keylogger, delay time.Duration) {
	events := keylogger.Read()
	storedApiData := LoggedActivity{}
	var mu sync.Mutex

	go func() {
		for event := range events {
			if IsKeyInputValid(event) {
				log.Printf("Key pressed: %s keyCode:%d", event.ToString(), event.Code)
				mu.Lock()
				storedApiData.Key++
				mu.Unlock()
			}
		}
	}()

	for {
		time.Sleep(delay)
		mu.Lock()
		currentData := storedApiData
		storedApiData = LoggedActivity{} // Reset the stored data
		mu.Unlock()

		go func(data LoggedActivity) {
			jsonBody, err := json.Marshal(data)
			if err != nil {
				log.Printf("Error marshaling JSON: %v", err)
				return
			}
			url := "http://localhost:8000/log"
			err = sendRequest("POST", url, bytes.NewBuffer(jsonBody))
			if err != nil {
				log.Printf("Error sending request: %v", err)
			}
		}(currentData)
	}
}
func SendDataFromLoggers(keyloggers []activitylogger.Keylogger, delay time.Duration) {
	var mu sync.Mutex
	storedApiData := LoggedActivity{}

	var wg sync.WaitGroup

	for _, keylogger := range keyloggers {
		wg.Add(1)
		go func(k activitylogger.Keylogger) {
			defer wg.Done()
			events := k.Read()
			for event := range events {
				if IsKeyInputValid(event) {
					log.Printf("Key pressed: %s keyCode:%d", event.ToString(), event.Code)
					mu.Lock()
					eventType := CategorizeEvent(event)
					UpdateLogFromEventType(eventType, &storedApiData)
					mu.Unlock()
				}
			}
		}(keylogger)
	}

	go func() {
		for {
			time.Sleep(delay)
			mu.Lock()
			currentData := storedApiData
			storedApiData = LoggedActivity{} // Reset the stored data
			mu.Unlock()

			go func(data LoggedActivity) {
				jsonBody, err := json.Marshal(data)
				if err != nil {
					log.Printf("Error marshaling JSON: %v", err)
					return
				}
				url := "http://localhost:8000/log"
				err = sendRequest("POST", url, bytes.NewBuffer(jsonBody))
				if err != nil {
					log.Printf("Error sending request: %v", err)
				} else {
					log.Printf("Data sent successfully: %v", data)
				}
			}(currentData)
		}
	}()
	wg.Wait()
}

func ReadDataFromLogger(keylogger activitylogger.Keylogger, activityChannel chan LoggedActivity, wg *sync.WaitGroup, delay time.Duration) {
	defer wg.Done()
	var mu sync.Mutex
	readData := LoggedActivity{}
	events := keylogger.Read()

	go func() {
		for event := range events {
			if IsKeyInputValid(event) {
				evtType := CategorizeEvent(event)
				// log.Printf("KEY CLICKED: %s CODE: %d", event.ToString(), event.Code)
				mu.Lock()
				UpdateLogFromEventType(evtType, &readData)
				mu.Unlock()
			}
		}
	}()

	time.Sleep(delay)
	mu.Lock()
	activityChannel <- readData
	readData = LoggedActivity{}
	mu.Unlock()
}
