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
	"sync"
)


func multiply(a, b [][]int, result [][]int, start, end int) {
	for i := start; i < end; i++ {
		for j := 0; j < len(b[0]); j++ {
			for k := 0; k < len(b); k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
}

func matrixMultiplication(rows, cols, numProcesses int) [][]int {
	a := make([][]int, rows)
	for i := range a {
		a[i] = make([]int, cols)
		for j := range a[i] {
			a[i][j] = rand.Intn(91) + 10
		}
	}

	b := make([][]int, rows)
	for i := range b {
		b[i] = make([]int, cols)
		for j := range b[i] {
			b[i][j] = rand.Intn(91) + 10
		}
	}

	// Initialize result matrix
	result := make([][]int, rows)
	for i := range result {
		result[i] = make([]int, cols)
	}

	var wg sync.WaitGroup
	wg.Add(numProcesses)

	// Perform matrix multiplication
	for i := 0; i < numProcesses; i++ {
		go func(i int) {
			defer wg.Done()
			start := (rows / numProcesses) * i
			end := start + (rows / numProcesses)
			if i == numProcesses-1 {
				end = rows
			}
			multiply(a, b, result, start, end)
		}(i)
	}

	wg.Wait()

	return result
}


func inFunctionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("entering inFunctionHandler")
	numProcs := runtime.NumCPU()
	fmt.Println("Number of processors:", numProcs)

	startTime := time.Unix(0, time.Now().UnixNano())
	rowStr := r.URL.Query().Get("rows")
	colStr := r.URL.Query().Get("cols")

	if rowStr != "" && colStr != "" {
		rows, err1 := strconv.Atoi(rowStr)
		if err1 != nil {
			http.Error(w, "Error with rows value passed", http.StatusBadRequest)
			return
		}

		cols, err2 := strconv.Atoi(colStr)
		if err2 != nil {
			http.Error(w, "Error with cols value passed", http.StatusBadRequest)
			return
		}
		fmt.Printf("entering multiplication")
		matrixMultiplication(rows, cols, numProcs)

		elapsed := time.Since(startTime)
		elapsedSec := fmt.Sprintf("%.8f", elapsed.Seconds())

		response := map[string]interface{}{
			"matrix size": rows * cols,
			"execTime":    elapsedSec,
			"numProcs":    numProcs,
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)

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

	http.HandleFunc("/api/in-function-parallelism", inFunctionHandler)

    log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
    log.Fatal(http.ListenAndServe(listenAddr, nil))
}