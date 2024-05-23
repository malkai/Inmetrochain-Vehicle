package chaincode

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const R = 6371 //raio da Terra em km

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

func KalmanFilter(capacidade float64, medições []float64) (float64, error) {

	var media = mediavector(medições)

	Gerro := errovector(medições, media) // desvio global

	auxe := []float64{}
	auxe = append(auxe, medições[0])
	Lerro := errovector(auxe, media) // desvio local

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

	resultadotanque := ((medições[0] - estima) * capacidade) / 100

	if math.IsNaN(resultadotanque) {
		return 0.0, fmt.Errorf("\n erro NaN value conf")
	}
	return resultadotanque, nil
}

func mediavector(a []float64) float64 {
	media := 0.0
	for _, leitura := range a {
		media = media + leitura
	}

	media = media / float64(len(a)-1)
	return media

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
	if len(a) == 1 {
		aux = aux / float64(len(a))
	} else {

		aux = aux / float64(len(a)-1)

	}

	return aux
}

func totaltime(s1, s2 string) (float64, error) {
	layout := "2006-01-02 15:04:05"

	datet1, err := time.Parse(layout, s1)
	if err != nil {

		return 0.0, err
	}
	datet2, err := time.Parse(layout, s2)
	if err != nil {

		return 0.0, err
	}

	result := datet2.Sub(datet1)

	return float64(result.Seconds()), nil
}

func Timeliness(valuevalids []float64, k float64) (float64, error) {

	k_aux := 1 / k

	f := mediavector(valuevalids)

	f_aux := 1 / f

	if f_aux >= k_aux {
		return 1.0, nil
	}
	f_k := (f_aux / k_aux) / math.Log((f_aux / k_aux))

	res_2 := math.Exp(1)

	timeless := -math.Pow(res_2, f_k) + 1

	if math.IsNaN(timeless) {
		return 0.0, fmt.Errorf("\n erro NaN value timeless")
	}

	return timeless, nil

}
