package function

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func matrixMultiplication(rows, cols int) [][]int {
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

	// Perform matrix multiplication
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			for k := 0; k < cols; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}

	return result
}

func Handle(w http.ResponseWriter, r *http.Request) {
	startTime := time.Unix(0, time.Now().UnixNano())
	rowStr := r.URL.Query().Get("rows")
	colStr := r.URL.Query().Get("cols")

	if rowStr != "" && colStr != "" {
		rows, err := strconv.Atoi(rowStr)
		cols, err := strconv.Atoi(colStr)

		fmt.Printf("rows, cols = %d, %d", rows, cols)

		if err == nil {

			matrixMultiplication(rows, cols)

			elapsed := time.Since(startTime)
			elapsedSec := fmt.Sprintf("%.8f", elapsed.Seconds())
			response := map[string]interface{}{
				// "result":   temp,
				"matrix size": rows * cols,
				"execTime":    elapsedSec,
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
