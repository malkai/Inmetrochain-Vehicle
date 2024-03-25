package chaincode

import (
	"fmt"

	"github.com/Konstantin8105/pow"
	"gonum.org/v1/gonum/mat"
)

// Linear regression model:
//
//	y   = a*x+b
//	R^2 - the relative predictive power of a linear model
func Linear(data [][2]float64) (
	a, b float64,
	R2 float64,
	err error,
) {
	if len(data) < 2 {
		err = fmt.Errorf("not enought data for regression")
		return
	}
	var x2, x1 float64
	n := float64(len(data))
	for i := range data {
		x1 += data[i][0]
		x2 += pow.E2(data[i][0])
	}
	A := mat.NewDense(2, 2, []float64{
		x2, x1,
		x1, n,
	})
	var b1, b2 float64
	for i := range data {
		b1 += data[i][0] * data[i][1]
		b2 += data[i][1]
	}
	right := mat.NewDense(2, 1, []float64{b1, b2})
	var res mat.Dense
	if err = res.Solve(A, right); err != nil {
		return
	}
	a = res.At(0, 0)
	b = res.At(1, 0)

	// the relative predictive power of a quadratic model
	var xMean, yMean float64
	for i := range data {
		xMean += data[i][0]
		yMean += data[i][1]
	}
	xMean = xMean / float64(len(data))
	yMean = yMean / float64(len(data))

	var SPxy, SSx float64
	for i := range data {
		xi, yi := data[i][0], data[i][1]
		SPxy += (xi - xMean) * (yi - yMean)
		SSx += pow.E2(xi - xMean)
	}
	bb1 := SPxy / SSx
	bb0 := yMean - bb1*xMean

	var SSE, SST float64
	for i := range data {
		xi, yi := data[i][0], data[i][1]
		SSE += pow.E2((bb1*xi + bb0) - yMean)
		SST += pow.E2(yi - yMean)
	}
	R2 = SSE / SST
	return
}

// Quadratic regression model:
//
//	y   = a*x^2+b*x+c
//	R^2 - the relative predictive power of a quadratic model
func Quadratic(data [][2]float64) (
	a, b, c float64,
	R2 float64,
	err error,
) {
	if len(data) < 3 {
		err = fmt.Errorf("not enought data for regression")
		return
	}
	var x4, x3, x2, x1 float64
	n := float64(len(data))
	for i := range data {
		x1 += data[i][0]
		x2 += pow.E2(data[i][0])
		x3 += pow.E3(data[i][0])
		x4 += pow.E4(data[i][0])
	}
	A := mat.NewDense(3, 3, []float64{
		x4, x3, x2,
		x3, x2, x1,
		x2, x1, n,
	})
	var b1, b2, b3 float64
	for i := range data {
		b1 += pow.E2(data[i][0]) * data[i][1]
		b2 += data[i][0] * data[i][1]
		b3 += data[i][1]
	}
	right := mat.NewDense(3, 1, []float64{b1, b2, b3})
	var res mat.Dense
	if err = res.Solve(A, right); err != nil {
		return
	}
	a = res.At(0, 0)
	b = res.At(1, 0)
	c = res.At(2, 0)

	// the relative predictive power of a quadratic model
	var SSE, SST, yMean float64
	for i := range data {
		yMean += data[i][1]
	}
	yMean = yMean / float64(len(data))
	for i := range data {
		xi, yi := data[i][0], data[i][1]
		SSE += pow.E2(yi - (a*xi*xi + b*xi + c))
		SST += pow.E2(yi - yMean)
	}
	R2 = 1 - SSE/SST
	return
}
