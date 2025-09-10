package main

import (
	"flag"
	"fmt"
	"generator/generator"
	"image/color"
	"log"
	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func simulateAndPlot(pTarget float64, n int, trials int) error {
	generator := generator.NewGenerator(pTarget)

	wins := 0
	points := make(plotter.XYs, trials)

	for t := 1; t <= trials; t++ {
		y := generator.GenerateNumber()
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
	p.X.Label.Text = "Количество попыток"
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

	if err := p.Save(8*vg.Inch, 4*vg.Inch, fmt.Sprintf("convergence_%f.png", pTarget)); err != nil {
		return err
	}
	return nil
}

func main() {
	RTP := flag.Float64("rtp", 0.5, "RTP value")
	flag.Parse()
	if *RTP > 1.0 || *RTP < 0.0 {
		log.Fatal("rtp must be in (0, 1.0]")
	}

	n := 10000
	trials := 10000

	err := simulateAndPlot(*RTP, n, trials)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("График сохранён в convergence_float.png")
}
