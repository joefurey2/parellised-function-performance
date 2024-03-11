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
	"runtime"
)


func GetTime() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}



func Alu(times int) float64 {
    a := rand.Intn(90) + 10
    b := rand.Intn(90) + 10
    var temp float64
    for i := 0; i < times; i++ {
        switch i % 4 {
        case 0:
            temp = float64(a + b)
        case 1:
            temp = float64(a - b)
        case 2:
            temp = float64(a * b)
        case 3:
            temp = float64(a) / float64(b)
        }
    }
    // fmt.Println(times)
    return temp
}

func inFunctionHandler(w http.ResponseWriter, r *http.Request) {
	numProcs := runtime.NumCPU()
	fmt.Println("Number of processors:", runtime.NumCPU())
	

	startTime := GetTime()
	timesStr := r.URL.Query().Get("times")

	if timesStr != "" {
		times, err := strconv.Atoi(timesStr)
		if err == nil {

			results := make(chan float64, times)
			computationsPerProc := times / numProcs

            for i := 0; i < numProcs; i++ {
                go func() {
                    results <- Alu(1)
                }()
            }

            var temp float64
            for i := 0; i < times; i++ {
                temp += <-results
            }

            response := map[string]interface{}{
                "result":   temp,
                "times":    times,
                "execTime": GetTime() - startTime,
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


func helloHandler(w http.ResponseWriter, r *http.Request) {
    message := "This HTTP triggered function executed successfully. Pass a name in the query string for a personalized response.\n"
    name := r.URL.Query().Get("name")
    if name != "" {
        message = fmt.Sprintf("Hello, %s. This HTTP triggered function executed successfully.\n", name)
    }
    fmt.Fprint(w, message)
}