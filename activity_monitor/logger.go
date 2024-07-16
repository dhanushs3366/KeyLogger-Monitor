package activitymonitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

func SetupLoggers(devEvents []uint) ([]activitylogger.Keylogger, error) {
	var keyloggers []activitylogger.Keylogger
	for _, devEvent := range devEvents {
		devPath := fmt.Sprintf("%s/event%d", DEV_PATH, devEvent)
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
func GetDevPaths() []uint {
	dir, err := os.Open(SYS_PATH)
	if err != nil {
		log.Printf("Error:%v\n", err)
	}
	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		log.Printf("Error:%v\n", err)
	}

	var events []int
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "event") {
			event, err := strconv.Atoi(strings.ReplaceAll(file.Name(), "event", ""))
			if err != nil {
				log.Printf("No event number found for %s", file.Name())
				continue
			}
			events = append(events, event)
		}
	}
	sort.Ints(events[:])
	inputDevices := getInputDevices(events)
	allowedDevEvents := getAllowedDevEvents(inputDevices)

	return allowedDevEvents
}

func getAllowedDevEvents(devMap map[int]string) []uint {
	var devEvents []uint

	for _, allowedDevice := range AllowedDevices {
		for event, device := range devMap {
			if normalizeString(device) == normalizeString(allowedDevice) {
				devEvents = append(devEvents, uint(event))
			}
		}
	}

	return devEvents
}

func getInputDeviceName(event uint) (string, error) {
	buff, err := os.ReadFile(fmt.Sprintf("%s/event%d/device/name", SYS_PATH, event))
	if err != nil {
		return "", err
	}
	return normalizeString(string(buff)), nil
}

func getInputDevices(events []int) map[int]string {
	inputDevices := make(map[int]string)
	for _, event := range events {
		nameFile := fmt.Sprintf("%s/event%d/device/name", SYS_PATH, event)
		inputDevice, err := os.ReadFile(nameFile)
		if err != nil {
			log.SetPrefix("ERROR: ")
			log.Printf("%v", err)
			continue
		}
		inputDevices[event] = strings.TrimSpace(string(inputDevice))

	}
	return inputDevices
}

func normalizeString(s string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.ToLower(re.ReplaceAllString(s, ""))
}
