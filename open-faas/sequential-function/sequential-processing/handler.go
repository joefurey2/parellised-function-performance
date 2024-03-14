package function

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func Alu(times int) float64 {
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
		} else if i%4 == 3 {
			temp = float64(a) / float64(b)
		}
	}
	return temp
}

func Handle(w http.ResponseWriter, r *http.Request) {
	startTime := time.Unix(0, time.Now().UnixNano())
	timesStr := r.URL.Query().Get("times")

	if timesStr != "" {
		times, err := strconv.Atoi(timesStr)
		fmt.Printf("times = %d", times)

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
