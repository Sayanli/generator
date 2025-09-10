package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func weightsExpGamma(gamma float64, n int) []float64 {
	weights := make([]float64, n)
	var maxExp float64

	maxExp = gamma * (1 - 0.5)
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

func findGammaForTarget(pTarget float64, n int) float64 {
	lo, hi := -200.0, 200.0
	tol := 1e-6
	for i := 0; i < 100; i++ {
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

func generateY(weights []float64, n int) float64 {
	cdf := make([]float64, len(weights))
	accum := 0.0
	for i, w := range weights {
		accum += w
		cdf[i] = accum
	}

	r := rand.Float64()

	lo, hi := 0, len(cdf)-1
	for lo < hi {
		mid := (lo + hi) / 2
		if r > cdf[mid] {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	return float64(lo+1) + rand.Float64()
}

func simulateAndPlot(pTarget float64, n int, trials int, outFile string) error {
	gamma := findGammaForTarget(pTarget, n)
	weights := weightsExpGamma(gamma, n)

	fmt.Printf("Подобранный γ=%.4f, теоретический шанс победы ≈ %.4f\n",
		gamma, avgProbWin(gamma, n))

	wins := 0
	points := make(plotter.XYs, trials)

	for t := 1; t <= trials; t++ {
		y := generateY(weights, n)
		x := 1 + rand.Float64()*float64(n-1)
		if y > x {
			wins++
		}
		ratio := float64(wins) / float64(t)
		points[t-1].X = float64(t)
		points[t-1].Y = ratio
	}

	p := plot.New()
	p.Title.Text = fmt.Sprintf("Сходимость к %.2f", pTarget)
	p.X.Label.Text = "Количество прокруток"
	p.Y.Label.Text = "Доля побед"

	line, err := plotter.NewLine(points)
	if err != nil {
		return err
	}
	p.Add(line)

	targetLine := plotter.XYs{
		{X: 0, Y: pTarget},
		{X: float64(trials), Y: pTarget},
	}
	lineTarget, _ := plotter.NewLine(targetLine)
	lineTarget.LineStyle.Color = color.Black
	p.Add(lineTarget)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, outFile); err != nil {
		return err
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	n := 10000
	pTarget := 0.5
	trials := 10000

	err := simulateAndPlot(pTarget, n, trials, "convergence_float.png")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("График сохранён в convergence_float.png")
}
