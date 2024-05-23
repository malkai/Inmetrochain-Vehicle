package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Konstantin8105/pow"
	"github.com/sgreben/piecewiselinear"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/plotter"
)

const R = 6371 //raio da Terra em km

func estatistica(folder string, name string) {

	file2, err := os.Create("estatistica/" + name + "dados.csv")
	if err != nil {
		panic(err)
	}
	defer file2.Close()

	//var fuel []float64
	lst := [30]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29"}
	media_Array := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	Desvio_Array := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	numbertime_Array := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	numbertime := 0.0

	//lst := [1]string{"10"}

	for tt, ji := range lst {
		auxvalur := 0.0
		valores := []float64{}
		sum := 0.0

		media_Array[tt] = 0

		//fmt.Println(ji)
		f, err := os.Open(folder + ji + name)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		csvReader.FieldsPerRecord = -1

		data, err := csvReader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		for _, line := range data {

			for uu, field := range line {
				if field != line[len(line)-1] {
					s, err := strconv.ParseFloat(field, 64)
					if err != nil {
						log.Fatal(err)
					}
					s2, err := strconv.ParseFloat(line[uu+1], 64)
					if err != nil {
						log.Fatal(err)
					}
					if math.IsNaN(s2) {
						fmt.Println(line[uu+1], uu, field, line, ji)

					}
					if s-s2 > 0 {
						sum = sum + s - s2
					} else {
						sum = sum + (s-s2)*(-1)
					}

					auxvalur = auxvalur + 1.0

					if sum > 5.0 {
						//fmt.Println(auxvalur)
						valores = append(valores, auxvalur)
						//fmt.Println(auxvalur)
						//fmt.Println(DesvioPadrão(valores, media_Array[tt]))
						sum = 0.0
						auxvalur = 0

						numbertime = numbertime + 1.0
					}

					//if sum => 5.0 {
					//
					//}
				}

			}

		}

		if len(valores) > 0 {
			media_Array[tt] = media_Array[tt] + (math.Round(mediavector(valores)*1000) / 1000)
			//fmt.Println(media_Array)
			Desvio_Array[tt] = Desvio_Array[tt] + (math.Round(DesvioPadrão(valores, media_Array[tt])*1000) / 1000)
			//fmt.Println(Desvio_Array)

		}

		//media_Array[tt] = media_Array[tt]
		//fmt.Println(media_Array)
		//Desvio_Array[tt] = media_Array[tt]
		numbertime_Array[tt] = numbertime
		numbertime = 0.0

	}
	helpstring := []string{}
	auxmedia := ""
	writer := csv.NewWriter(file2)
	for _, a := range media_Array {

		auxmedia = strings.Replace(fmt.Sprintf("%f", math.Round(a)*10000/10000), ".", ",", -1)

		helpstring = append(helpstring, auxmedia)
	}
	writer.Write(helpstring)
	writer.Flush()
	helpstring = nil

	for _, a := range Desvio_Array {

		auxmedia = strings.Replace(fmt.Sprintf("%f", math.Round(a)*10000/10000), ".", ",", -1)

		helpstring = append(helpstring, auxmedia)
	}
	writer.Write(helpstring)
	writer.Flush()
	helpstring = nil

	for _, a := range numbertime_Array {

		auxmedia = strings.Replace(fmt.Sprintf("%f", math.Round(a)*10000/10000), ".", ",", -1)

		helpstring = append(helpstring, auxmedia)
	}
	writer.Write(helpstring)
	writer.Flush()
	helpstring = nil

	mediatudo := (math.Round(mediavector(media_Array)*1000) / 1000)
	desvitudo := (math.Round(DesvioPadrão(media_Array, mediatudo)*1000) / 1000)

	auxmedia = strings.Replace(fmt.Sprintf("%f", mediatudo), ".", ",", -1)
	auxdesv := strings.Replace(fmt.Sprintf("%f", desvitudo), ".", ",", -1)
	helpstring = append(helpstring, auxmedia)
	helpstring = append(helpstring, auxdesv)
	writer.Write(helpstring)
	writer.Flush()
	helpstring = nil
	fmt.Println("Media Conjunto", mediatudo, ",", "Desvio Padrão", desvitudo)

}

func mediavector(a []float64) float64 {
	media := 0.0
	for _, leitura := range a {
		media = media + leitura
	}

	media = media / float64(len(a))

	return media

}
func Timeliness(valuevalids []float64, k float64) float64 {

	k_aux := 1 / k
	fmt.Println(k, "k")

	f := mediavector(valuevalids)

	f_aux := 1 / f
	fmt.Println(f_aux, "f_aux")
	fmt.Println(k_aux, "k_aux")

	if f_aux >= k_aux {
		return 1.0
	}
	f_k := (f_aux / k_aux) / math.Log((f_aux / k_aux))

	fmt.Println(f_k, "f_k")
	res_2 := math.Exp(1)
	fmt.Println(res_2, "f_k")
	timeless := -math.Pow(res_2, f_k) + 1
	fmt.Println(timeless, "timeless")

	if math.IsNaN(timeless) {
		return 0.0
	}

	return timeless

}

func totaltime(s1, s2 string) float64 {
	layout := "2006-01-02 15:04:05"

	datet1, err := time.Parse(layout, s1)
	if err != nil {

		return 0.0
	}
	datet2, err := time.Parse(layout, s2)
	if err != nil {

		return 0.0
	}

	result := datet2.Sub(datet1)

	return float64(result.Seconds())
}

func DesvioPadrão(a []float64, media float64) float64 {
	var aux float64 = 0

	for _, leitura := range a {
		aux = aux + (math.Pow(leitura-media, 2))

	}
	aux = math.Sqrt((aux) / float64(len(a)))
	//fmt.Println(aux)
	return aux
}

func interpolation(datalist []float64) []float64 {

	medialen := mediavector(datalist)
	desvpa := DesvioPadrão(datalist, medialen)
	amplitude := 0.95 * desvpa

	m := 0.95

	gausseries := []float64{}

	regressiveseries := []float64{}

	for i := 0; i <= len(datalist)-1; i++ {

		gausseries = append(gausseries, 0.1+rand.Float64()*((amplitude)-0.1))

	}

	for i := 0; i <= len(datalist)-1; i++ {
		if i == 0 {
			regressiveseries = append(regressiveseries, 0)
		} else {
			regressiveseries = append(regressiveseries, (m)*regressiveseries[i-1]+(1-m)*gausseries[i])
		}

	}

	//fmt.Println(regressiveseries)

	for i := 0; i <= len(datalist)-1; i++ {

		if i != 0 && i != len(datalist)-1 {
			if datalist[i]+regressiveseries[i] < 100 && datalist[i]+regressiveseries[i] > 0 {
				//fmt.Println(datalist[i], regressiveseries[i])
				r := 0 + rand.Float64()*(1-0)
				if r < 0.5 {
					r = -1.0
				} else {
					r = 1.0
				}
				datalist[i] = math.Round((datalist[i]+regressiveseries[i]*r)*1000) / 1000
			} else if datalist[i]+regressiveseries[i] > 100 {
				datalist[i] = 100
			} else if datalist[i]+regressiveseries[i] < 0 {
				datalist[i] = 0
			}
		}

	}

	return datalist
}

func KalmanFilter(capacidade float64, medições []float64) ([]float64, error) {

	var media = mediavector(medições)

	//fmt.Println(media)

	Gerro := errovector(medições, math.Round(media*10000)/10000) // desvio global

	//fmt.Println(Gerro)

	auxe := []float64{}
	auxe = append(auxe, medições[0])
	Lerro := errovector(auxe, math.Round(media*10000)/10000) // desvio local

	fmt.Println(Lerro)

	//fmt.Println(Lerro, auxe, math.Ceil(media*100)/100, math.Ceil(math.Pow(Lerro, 2)*100)/100, math.Ceil(math.Pow(Gerro, 2)*10000)/10000)

	var estima float64
	leiturasPercentuaispos := []float64{}
	k := (math.Ceil(math.Pow(Lerro, 2)*100000) / 100000) / (math.Ceil(math.Pow(Lerro, 2)*100000)/100000 + math.Ceil(math.Pow(Gerro, 2)*100000)/100000)

	for _, leitura := range medições {
		estima = media + k*((math.Ceil(leitura*100000)/100000)-(math.Ceil(media*100000)/100000))
		if estima > 100 {
			estima = 100
		} else if estima < 0 {
			estima = 0
		}

		Lerro = (1.0 - k) * (math.Ceil(Lerro*100000) / 100000)
		media = estima
		k = (math.Ceil(Lerro*100000) / 100000) / (math.Ceil(Lerro*100000)/100000 + math.Ceil(math.Pow(Gerro, 2)*100000)/100000)

		leiturasPercentuaispos = append(leiturasPercentuaispos, estima)

	}

	//resultadotanque := ((medições[0] - estima) * capacidade) / 100
	return leiturasPercentuaispos, nil
}

func regress(datalist []float64) []float64 {

	p := 0.7
	o := 1.0

	v := distuv.Normal{Mu: 0, Sigma: o * math.Sqrt(1-(p*p))}

	e := []float64{0}

	rtt := []float64{}

	rtt2 := []float64{}

	copy(rtt, datalist)

	for i := 1; i <= len(datalist)-1; i++ {

		e = append(e, p*e[i-1]+v.Rand())

	}

	for i := 0; i <= len(datalist)-1; i++ {

		if datalist[i]+e[i] > 100 {

			rtt2 = append(rtt2, 100.0)
		} else if datalist[i]+e[i] < 0 {

			rtt2 = append(rtt2, 0.0)
		} else {
			//rtt[i] = rtt[i] + rtt[i]*e[i]
			rtt2 = append(rtt2, datalist[i]+e[i])
		}

	}
	/*
		tt := plot.New()

		tt.Title.Text = "Plotutil example"
		tt.X.Label.Text = "X"
		tt.Y.Label.Text = "Y"

		err := plotutil.AddLines(tt,
			"Comruído", convertToPtl(rtt2),
			"Semruído", convertToPtl(datalist),
		)
		if err != nil {
			panic(err)
		}

		// Save the plot to a PNG file.
		t := time.Now()
		ER := t.Format("20060102150405")
		if err := tt.Save(4*vg.Inch, 4*vg.Inch, "plot/"+ER+"points.png"); err != nil {
			panic(err)
		}
	*/

	return rtt2
}

func convertToPtl(a []float64) plotter.XYs {

	pts := make(plotter.XYs, len(a))
	for i := range a {
		pts[i].X = float64(i)

		pts[i].Y = a[i]
	}

	return pts

}

func ruido(simudata []Tuple) []Tuple {

	noise := []float64{}
	interpolationdata := []float64{}

	for _, data := range simudata {

		noise = append(noise, data.Comb)

	}

	fmt.Println("Normal	", noise[0]-noise[len(noise)-1])

	//fmt.Println(noise)
	if len(noise) > 0 {
		interpolationdata = regress(noise)
	}
	//fmt.Println(interpolationdata)
	fuel_vector := []float64{}
	time3 := []float64{0}
	time := 0.0

	for i, ts := range simudata {

		fuel_vector = append(fuel_vector, ts.Comb)
		if i < len(simudata)-1 {
			time1 := totaltime(simudata[i].T, simudata[i+1].T)

			time = time + time1

			//rtt, _ := strconv.ParseFloat(tuples[i].T, 64)
			time3 = append(time3, time)
		}

	}

	f := piecewiselinear.Function{Y: interpolationdata} // range: "hat" function
	f.X = time3
	rtt := [][2]float64{}

	for rty := range time3 {
		rtt1 := [2]float64{time3[rty], fuel_vector[rty]}
		rtt = append(rtt, rtt1)
	}
	a, b, _, err := Linear(rtt)
	if err != nil {
		fmt.Println("Erro em alguma coisa")
	}

	tui := []float64{}
	for i := range time3 {

		tui = append(tui, time3[i]*a+b)
	}

	fmt.Println("Regressão", tui[0]-tui[len(tui)-1])

	for i := range simudata {

		simudata[i].Comb = interpolationdata[i]

	}

	return simudata
}

// https://physics.stackexchange.com/questions/704367/how-to-quantify-the-uncertainty-of-the-time-series-average
func errovector(a []float64, media float64) float64 {
	var aux float64 = 0
	errovec := []float64{}
	for _, leitura := range a {
		errovec = append(errovec, math.Pow(leitura-media, 2))
	}
	for _, leitura := range errovec {
		aux = +leitura
	}
	aux = aux / float64(len(a)-2)
	return aux
}

func CoordenadasCartesianas(latitude, longitude float64) []float64 {
	theta := latitude * (math.Pi / 180.0) //teta
	phi := longitude * (math.Pi / 180.0)  //fi

	x := R * math.Cos(theta) * math.Cos(phi)
	y := R * math.Cos(theta) * math.Sin(phi)
	z := R * math.Sin(theta)
	var a []float64
	a = append(a, x)
	a = append(a, y)
	a = append(a, z)
	return a
}

func Distanceeucle(latlongA, latlongB string) (float64, error) {

	distancia := 0.0
	res1 := strings.Split(latlongA, "/")

	latitudeA, err := strconv.ParseFloat(res1[0], 64)
	if err != nil {
		return 0.0, fmt.Errorf("\n Erro checar LatA. %v", err)
	}
	longitudeA, _ := strconv.ParseFloat(res1[1], 64)
	if err != nil {
		return 0.0, fmt.Errorf("\n Erro checar LonA. %v", err)
	}

	//return 0.0, fmt.Errorf("\n %f %f ", latitudeA, longitudeA)

	res1 = strings.Split(latlongB, "/")
	latitudeB, err := strconv.ParseFloat(res1[0], 64)
	if err != nil {
		return 0.0, fmt.Errorf("\n Erro checar LatB. %v", err)
	}
	longitudeB, err := strconv.ParseFloat(res1[1], 64)
	if err != nil {
		return 0.0, fmt.Errorf("\n Erro checar LonB. %v", err)
	}
	//distancia = latitudeA + longitudeA + latitudeB + longitudeB
	//return 0.0, fmt.Errorf("\n %f %f ", latitudeB, longitudeB)

	a := CoordenadasCartesianas(latitudeA, longitudeA)
	b := CoordenadasCartesianas(latitudeB, longitudeB)
	//return 0.0, fmt.Errorf("\n %f %f ", latitudeB, longitudeB)

	x1, y1, z1 := a[0], a[1], a[2]
	x2, y2, z2 := b[0], b[1], b[2]
	distancia = math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2) + math.Pow(z2-z1, 2)) //calculando a distancia euclidiana entre os pontos A e B
	//return 0.0, fmt.Errorf("\n %f %f %f", distancia, a, b)
	return distancia, nil

}

func convertToString(a []float64) []string {
	b := []string{}
	for i := range a {
		s := fmt.Sprintf("%f", a[i])

		b = append(b, s)
	}

	return b

}

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
