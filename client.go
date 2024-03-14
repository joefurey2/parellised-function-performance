package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"os"
	"strconv"
	"encoding/json"
)

type Response struct {
    Result   float64 `json:"result"`
    Times    int     `json:"times"`
    ExecTime string `json:"execTime"`
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Please provide the number of times and threads as a command-line argument.")
        return
    }

    numTimes, err := strconv.Atoi(os.Args[1])
    if err != nil {
        fmt.Println("Invalid number of times:", os.Args[1])
        return
    }

	numProcs, err := strconv.Atoi(os.Args[2])
    if err != nil {
        fmt.Println("Invalid number of processes:", os.Args[2])
        return
    }

	timesPerProc := numTimes / numProcs
	fmt.Printf("times per proc = %d \n", timesPerProc)

	startTime := time.Now()


	// Simultaneously call the endpoint with a fraction of the numbers specified
	var wg sync.WaitGroup
	results := make(chan float64, numTimes)

	for i := 0; i < numProcs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(fmt.Sprintf("http://20.26.236.208:8080/function/sequential-function?times=%d", timesPerProc))
			// https://in-function-parallelism.azurewebsites.net/api/in-function-parallelism?
			// https://seqential-parallelism.azurewebsites.net/api/sequential-processing?
			// http://20.26.236.208:8080/function/sequential-function
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close()

			// Read the response and convert it to an integer
			// Assuming the response body contains a single integer value
			var respBody Response
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			result := respBody.Result

			results <- result
		}()
	}

	wg.Wait()
	close(results)

	// Sum the results returned by the function
	totalScore := 0.0
	for result := range results {
		totalScore += result
	}

	elapsed := time.Since(startTime)
	elapsedSec := fmt.Sprintf("%.8f", elapsed.Seconds())

	fmt.Printf("Time elapsed: %s\n", elapsedSec)
	fmt.Printf("Total score: %f\n", totalScore)
	fmt.Printf("Initial number of times: %d\n", numTimes)
}
