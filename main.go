package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"generator/generator"
	"log"
	"net/http"
)

var RTP float64
var WEIGHTS []float64
var CDF []float64
var N = 10000

type Response struct {
	Result float64 `json:"result"`
}

func main() {
	RTP := flag.Float64("rtp", 0.5, "RTP value")
	flag.Parse()
	if *RTP > 1.0 || *RTP < 0.0 {
		log.Fatal("rtp must be in (0, 1.0]")
	}
	generator := generator.NewGenerator(*RTP)

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		multiplier := generator.GenerateNumber()

		resp := Response{Result: multiplier}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("Server started at http://localhost:64333/get")
	log.Fatal(http.ListenAndServe(":64333", nil))
}
