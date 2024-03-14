package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func Alu(times int, results chan<- float64) {
	a := rand.Intn(91) + 10
	b := rand.Intn(91) + 10
	var temp float64
	for i := 0; i < times; i++ {
		if i%4 == 0 {
			temp = float64(a + b)
		} else if i%4 == 1 {
			temp = float64(a - b)
		} else if i%4 == 2 {
			temp = float64(a * b)
		} else if i%4 == 1 {
			temp = float64(a) / float64(b)
		}
	}
	results <- temp
}

func inFunctionHandler(w http.ResponseWriter, r *http.Request) {
	numProcs := runtime.NumCPU()
	fmt.Println("Number of processors:", runtime.NumCPU())

	startTime := time.Unix(0, time.Now().UnixNano())
	timesStr := r.URL.Query().Get("times")

	if timesStr != "" {
		times, err := strconv.Atoi(timesStr)
		if err == nil {

			results := make(chan float64, times)
			computationsPerProc := times / numProcs

			var wg sync.WaitGroup
			for i := 0; i < numProcs; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					Alu(computationsPerProc, results)
				}()
			}

			wg.Wait()
			close(results)

			var total float64
			for result := range results {
				total += result
			}

			elapsed := time.Since(startTime)
			elapsedSec := fmt.Sprintf("%.8f", elapsed.Seconds())

			response := map[string]interface{}{
				"result":   total,
				"times":    times,
				"execTime": elapsedSec,
				"numProcs": numProcs,
			}

			responseJSON, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(responseJSON)
		}

	} else {
		message := "Error with times value passed"
		http.Error(w, message, http.StatusBadRequest)
	}
}

func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	// http.HandleFunc("/api/hello-handler", helloHandler)
	http.HandleFunc("/api/in-function-parallelism", inFunctionHandler)

	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}