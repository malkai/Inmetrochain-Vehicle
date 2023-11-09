package chaincode

import (
	"math"
	"strconv"
	"strings"
	"time"
)

const R = 6371 //raio da Terra em km

func CoordenadasCartesianas(latitude, longitude float64) (x, y, z float64) {
	theta := latitude * (math.Pi / 180.0) //teta
	phi := latitude * (math.Pi / 180.0)   //fi

	x = R * math.Cos(theta) * math.Cos(phi)
	y = R * math.Cos(theta) * math.Sin(phi)
	z = R * math.Sin(theta)

	return x, y, z
}

func Distanceeucle(latlongA, latlongB string) float64 {

	res1 := strings.Split(latlongA, "/")
	latitudeA, _ := strconv.ParseFloat(res1[0], 64)
	longitudeA, _ := strconv.ParseFloat(res1[1], 64)
	res1 = strings.Split(latlongB, "/")
	latitudeB, _ := strconv.ParseFloat(res1[0], 64)
	longitudeB, _ := strconv.ParseFloat(res1[1], 64)

	x1, y1, z1 := CoordenadasCartesianas(latitudeA, longitudeA)
	x2, y2, z2 := CoordenadasCartesianas(latitudeB, longitudeB)
	distancia := math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2) + math.Pow(z2-z1, 2)) //calculando a distancia euclidiana entre os pontos A e B
	return distancia
}

func KalmanFilter(capacidade float64, medições []float64) float64 {

	var media = mediavector(medições)

	Gerro := errovector(medições, media) // erro global

	auxe := []float64{}
	auxe = append(auxe, medições[0])
	Lerro := errovector(auxe, media) // erro local

	var estima float64
	leiturasPercentuaispos := []float64{}

	for _, leitura := range medições {
		var k = math.Pow(Lerro, 2)/math.Pow(Lerro, 2) + math.Pow(Gerro, 2)
		estima = media + k*(leitura-media)

		Lerro = (1 - k) * Lerro
		leiturasPercentuaispos = append(leiturasPercentuaispos, estima)

	}

	resultadotanque := (medições[1] - estima) * capacidade

	return resultadotanque

}

func mediavector(a []float64) float64 {
	media := 0.0
	for _, leitura := range a {
		media = media + leitura
	}

	media = media / float64(len(a))
	return media

}

func errovector(a []float64, m float64) float64 {
	var aux float64
	errovec := []float64{}
	for _, leitura := range a {
		errovec = append(errovec, leitura-m)
	}
	for _, leitura := range errovec {
		aux = +leitura
	}
	aux = aux / float64(len(a))
	return aux
}

func totaltime(s1, s2 string) float64 {
	layout := "2006-01-02 15:04:05.000000"

	datet1, error := time.Parse(layout, s1)
	if error != nil {

		return 0.0
	}
	datet2, error := time.Parse(layout, s2)
	if error != nil {

		return 0.0
	}

	result := datet2.Sub(datet1)
	convert, _ := strconv.ParseFloat(result.String(), 64)
	return convert
}

func Timeliness(valuevalids []string) float64 {
	var k = 5.0
	var vectortime = 0.0
	var vectortotal = 0.0
	for i, test := range valuevalids {
		if i != len(valuevalids) {
			layout := "2006-01-02 15:04:05.000000"
			datet1, _ := time.Parse(layout, test)
			datet2, _ := time.Parse(layout, valuevalids[i+1])
			result := datet2.Sub(datet1)
			convert, _ := strconv.ParseFloat(result.String(), 64)
			if vectortime < k {
				vectortime = vectortime + convert
			}
			vectortotal = vectortotal + convert

		}

	}
	var prop = vectortime / vectortotal
	var f_k = (prop / k) / math.Log((prop / k))
	res_2 := math.Exp(1)
	var timeless = math.Pow(res_2, f_k) + 1

	return timeless
}

func Credibility(scoren1 float64, timelesstotal float64, fuelcheck float64, fuelsum float64, valuevalids []string) float64 {
	helpc := fuelcheck / fuelsum //completness
	m := 0.9
	var conf = scoren1*m + (timelesstotal+helpc)/2*(1-m)
	return conf
}
