package function

import (
	"encoding/json"
	"net/http"
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
