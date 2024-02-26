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
)

type vehicletruct struct {
	T    string  `json:"T"`
	Pos  string  `json:"Pos"` // long + lat
	Comb float64 `json:"Comb"`
}

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

	//fmt.Println("Desvio padrao", desvpa, amplitude, datalist[0])

	//fmt.Println("Calculo", medialen, desvpa, amplitude, len(datalist))

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

	fmt.Println(media)

	Gerro := errovector(medições, math.Round(media*10000)/10000) // desvio global

	fmt.Println(Gerro)

	auxe := []float64{}
	auxe = append(auxe, medições[0])
	Lerro := errovector(auxe, math.Round(media*10000)/10000) // desvio local

	fmt.Println(Lerro)

	//fmt.Println(Lerro, auxe, math.Ceil(media*100)/100, math.Ceil(math.Pow(Lerro, 2)*100)/100, math.Ceil(math.Pow(Gerro, 2)*10000)/10000)

	var estima float64
	leiturasPercentuaispos := []float64{}

	for _, leitura := range medições {
		//k := (math.Ceil(math.Pow(Lerro, 2)*100000) / 1000000) / (math.Ceil(math.Pow(Lerro, 2)*1000000)/1000000 + math.Ceil(math.Pow(Gerro, 2)*1000000)/1000000)
		k := (math.Round(Lerro*1000) / 1000) / ((math.Round(Lerro*1000) / 1000) + math.Round(Gerro*1000)/1000)

		estima = media + k*((math.Round(leitura*1000)/1000)-(math.Round(media*1000)/1000))
		if estima > 100 {
			estima = 100
		} else if estima < 0 {
			estima = 0
		}
		//fmt.Println(leitura, estima)
		Lerro = (1.0 - k) * (math.Round(Lerro*1000) / 1000)

		media = estima

		//fmt.Println(1-k, k)
		//fmt.Println("Kalman gain", k, "estima", estima, "erro local", Lerro, "erro glbal", Gerro)
		leiturasPercentuaispos = append(leiturasPercentuaispos, math.Round(estima*1000)/1000)

	}

	//fmt.Println(leiturasPercentuaispos[0], leiturasPercentuaispos[len(leiturasPercentuaispos)-1], len(leiturasPercentuaispos))

	//resultadotanque := ((medições[0] - estima) * capacidade) / 100
	return leiturasPercentuaispos, nil
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
	phi := latitude * (math.Pi / 180.0)   //fi

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
