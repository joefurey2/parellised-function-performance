package main

import (
    "encoding/json"
    "net/http"
	"log"
	"os"
    "fmt"
)

type MultiplyRequest struct {
    A     [][]int `json:"a"`
    B     [][]int `json:"b"`
    Start int     `json:"start"`
    End   int     `json:"end"`
}

type MultiplyResponse struct {
    Result [][]int `json:"result"`
}

func multiplyHandler(w http.ResponseWriter, r *http.Request) {
    // Decode the request
    var req MultiplyRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Check the dimensions of the matrices
    if len(req.A[0]) != len(req.B) {
        http.Error(w, "Number of columns in matrix A does not match number of rows in matrix B", http.StatusBadRequest)
        return
    }

    // Check the start and end values
    if req.Start < 0 || req.End > len(req.A) {
        http.Error(w, fmt.Sprintf("Start or end value is out of bounds. Start: %d, End: %d, Length of A: %d", req.Start, req.End, len(req.A)), http.StatusBadRequest)
        return
    }

    // Perform the multiplication
    result := make([][]int, req.End-req.Start)
    for i := range result {
        result[i] = make([]int, len(req.B[0]))
    }
    for i := req.Start; i < req.End; i++ {
        for j := 0; j < len(req.B[0]); j++ {
            for k := 0; k < len(req.B); k++ {
                result[i-req.Start][j] += req.A[i][k] * req.B[k][j]
            }
        }
    }

    // Encode and send the response
    resp := MultiplyResponse{Result: result}
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
    listenAddr := ":8080"
    if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
        listenAddr = ":" + val
    }

	http.HandleFunc("/api/multi-function-parallelism", multiplyHandler)

    log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
    log.Fatal(http.ListenAndServe(listenAddr, nil))
}