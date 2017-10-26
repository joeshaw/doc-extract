// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
)

// +extract
// FORMAT: 1A
//
// # Sums and Averages
//
// This API computes sums and averages (with standard deviation).

type jsonRequest struct {
	Values []int `json:"values"`
}

// +extract
// ## POST /sum
// + Request (application/json)
//
//         {
//             "values": [ 10, 20, 30 ]
//         }
//
// + Response 200 (application/json)
//
//         {
//             "sum": 60
//         }

func sum(w http.ResponseWriter, r *http.Request) {
	var j jsonRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&j); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	var sum int
	for _, v := range j.Values {
		sum += v
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"sum": "%d"}`, sum)
}

// +extract
// ## POST /average
// + Request (application/json)
//
//         {
//             "values": [ 10, 20, 30 ]
//         }
//
// + Response 200 (application/json)
//
//         {
//             "average": 20,
//             "stddev": 10
//         }

func average(w http.ResponseWriter, r *http.Request) {
	var j jsonRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&j); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	var sum, sqSum float64
	for _, v := range j.Values {
		sum += float64(v)
		sqSum += float64(v * v)
	}

	var stddev float64
	n := len(j.Values)
	avg := float64(sum) / float64(n)

	if n > 1 {
		stddev = math.Sqrt((float64(n)*sqSum - sum*sum) / float64(n*(n-1)))
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"average": "%f", "stddev": "%f"}`, avg, stddev)
}

func main() {
	http.Handle("/sum", http.HandlerFunc(sum))
	http.Handle("/average", http.HandlerFunc(average))
	log.Fatal(http.ListenAndServe(":9090", nil))
}
