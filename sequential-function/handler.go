package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type jsonResponse struct {
    Result   float64 `json:"result"`
    Times    int     `json:"times"`
    ExecTime int     `json:"execTime"`
}

func Alu(times int) float64 {
    a := rand.Intn(91) + 10
    b := rand.Intn(91) + 10
    var temp float64
    for i := 0; i < times; i++ {
        if i % 4 == 0 {
            temp = float64(a + b)
        } else if i % 4 == 1 {
            temp = float64(a - b)
        } else if i % 4 == 2 { 
            temp = float64(a * b)
        } else if i % 4 == 1 { 
            temp = float64(a) / float64(b)
        }
    }
    return temp
}

func sequentialHandler(w http.ResponseWriter, r *http.Request) {
    startTime := time.Unix(0, time.Now().UnixNano())
	timesStr := r.URL.Query().Get("times")

	if timesStr != "" {
		times, err := strconv.Atoi(timesStr)
		if err == nil {

			temp := Alu(times)
            
            elapsed := time.Since(startTime)
            elapsedSec := fmt.Sprintf("%.8f", elapsed.Seconds())
            response := map[string]interface{}{
                "result":   temp,
                "times":    times,
                "execTime": elapsedSec,
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
    http.HandleFunc("/api/hello-handler", helloHandler)
	http.HandleFunc("/api/sequential-processing", sequentialHandler)

    log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
    log.Fatal(http.ListenAndServe(listenAddr, nil))
}


func helloHandler(w http.ResponseWriter, r *http.Request) {
    message := "This HTTP triggered function executed successfully. Pass a name in the query string for a personalized response.\n"
    name := r.URL.Query().Get("name")
    if name != "" {
        message = fmt.Sprintf("Hello, %s. This HTTP triggered function executed successfully.\n", name)
    }
    fmt.Fprint(w, message)
}