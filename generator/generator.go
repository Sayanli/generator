package generator

import (
	"math"
	"math/rand"
)

var N = 10000

type Generator struct {
	weights []float64
	cdf     []float64
	rtp     float64
}

func NewGenerator(RTP float64) *Generator {
	gamma := findGammaForTarget(RTP, N)
	weights := weightsExpGamma(gamma, N)

	cdf := make([]float64, N)
	accum := 0.0
	for i, w := range weights {
		accum += w
		cdf[i] = accum
	}
	return &Generator{
		weights: weights,
		cdf:     cdf,
		rtp:     RTP,
	}
}

func (g *Generator) GenerateNumber() float64 {
	r := rand.Float64()

	// бинарный поиск по CDF
	lo, hi := 0, len(g.cdf)-1
	for lo < hi {
		mid := (lo + hi) / 2
		if r > g.cdf[mid] {
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
