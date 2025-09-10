package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var RTP float64
var WEIGHTS []float64
var CDF []float64
var N = 10000

// Структура ответа
type Response struct {
	Result float64 `json:"result"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <rtp>")
	}

	var err error
	RTP, err = strconv.ParseFloat(os.Args[1], 64)
	if err != nil || RTP <= 0 || RTP > 1 {
		log.Fatal("rtp must be in (0, 1.0]")
	}

	gamma := findGammaForTarget(RTP, N)
	WEIGHTS = weightsExpGamma(gamma, N)

	CDF = make([]float64, N)
	accum := 0.0
	for i, w := range WEIGHTS {
		accum += w
		CDF[i] = accum
	}

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/get", getHandler)

	fmt.Println("Server started at http://localhost:64333/get")
	log.Fatal(http.ListenAndServe(":64333", nil))
}

func probWinGivenX(weights []float64, X int) float64 {
	var sum float64
	for i := X; i < len(weights); i++ {
		sum += weights[i]
	}
	return sum
}

func avgProbWin(gamma float64, n int) float64 {
	weights := weightsExpGamma(gamma, n)
	var total float64
	for X := 1; X <= n; X++ {
		total += probWinGivenX(weights, X)
	}
	return total / float64(n)
}

// Подбор gamma под нужный шанс
func findGammaForTarget(pTarget float64, n int) float64 {
	lo, hi := -200.0, 200.0
	tol := 1e-6
	for i := 0; i < 60; i++ {
		mid := (lo + hi) / 2
		pMid := avgProbWin(mid, n)
		if math.Abs(pMid-pTarget) < tol {
			return mid
		}
		if pMid < pTarget {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}

// Построение весов
func weightsExpGamma(gamma float64, n int) []float64 {
	weights := make([]float64, n)
	maxExp := gamma * (1 - 0.5)

	for i := 0; i < n; i++ {
		x := float64(i+1) / float64(n)
		exps := gamma*(x-0.5) - maxExp
		weights[i] = math.Exp(exps)
	}

	// нормализация
	var sum float64
	for _, w := range weights {
		sum += w
	}
	for i := range weights {
		weights[i] /= sum
	}
	return weights
}

func generateNumber() float64 {
	r := rand.Float64()

	// бинарный поиск по CDF
	lo, hi := 0, len(CDF)-1
	for lo < hi {
		mid := (lo + hi) / 2
		if r > CDF[mid] {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	res := float64(lo) + rand.Float64()
	if res > float64(N) {
		res = float64(N)
	}
	return res
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	multiplier := generateNumber()

	resp := Response{Result: multiplier}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
