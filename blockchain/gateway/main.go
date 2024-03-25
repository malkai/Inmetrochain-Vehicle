package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type Results struct {
	Data []datavehicle
}
type datavehicle struct {
	Id  int `json:"id"`
	Pos struct {
		Long float64 `json:"long"`
		Lat  float64 `json:"lat"`
	} `json:"pos"`
	Distancia     float64 `json:"Dist\u00e2ncia"`
	Combustivel   float64 `json:"combustivel"`
	TamanhoTanque int     `json:"TamanhoTanque"`
	Time          string  `json:"time"`
	Contrato      int     `json:"Contrato"`
	Tipo          string  `json:"Tipo"`

	//  FooBar  string `json:"foo.bar"`

}

type datahelp struct {
	id   string
	Data []datavehicle
}

type datafile struct {
	id   string
	file []string
}

func main() {

	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// Create client identity and signing implementation based on X.509 certificate and private key.
	id := NewIdentity()
	sign := NewSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),

		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	chaincodeName := "vehicle"
	channelName := "mychannel"

	network := gw.GetNetwork(channelName)

	contract := network.GetContract(chaincodeName)
	//GetPathhAll(contract, "1")
	//GetAlleventsNews(contract, "1")
	//GetAlleventsPast(contract, "1")
	//send_data(contract, "0")
	//	GetAlluser(contract, "user100")

	processing_Data(contract, 0.33)
	//GetAlluser(contract, "user1")
	//GetAlleventsOpens(contract, "1", "9999")
	GetAlleventsNews(contract, "1")
	GetAlleventsPast(contract, "1")
	//GetAllevents(contract, "event")

	//2024-03-22 01:21:27

	//"dataEvent": "2024-03-22-10-48-50",

	//GetPathhAll(contract, "1")

	//	name1 := "datawithoutnoise.csv"
	//	name2 := "datanoise.csv"
	//	name3 := "datatimstump.csv"
	//	name4 := "filter.csv"

	//GetAllevents(contract, "event")
	//openhugejson(contract)

	//fmt.Println("oi")

	//GetAlluser(contract, "user100")

	//estatistica("Kalmanfilter/", "filter.csv")
}

// estou inserindo um ruido branco guaisiano com o sigma
// linear interpolation

func send_data(contract *client.Contract, name string, file []string, alpha float64) {

	datapack := []Tuple{}

	sentfile := file[0]
	result := 0

	comb := 0.0
	timeacerto := ""
	timeaprox := ""

	meta := 0
	postoid := "9999"

	postocontrato := " "

	trr := datavehicle{}

	for _, filesss := range file {

		file, err := os.Open("datavehicle/" + name + "/" + filesss)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			tt := scanner.Text()

			err := json.Unmarshal([]byte(tt), &trr)
			if err != nil {
				fmt.Println(err.Error())

			}

			if err != nil {
				panic(err)
			}

			comb = trr.Combustivel * 100 / float64(trr.TamanhoTanque)

			meta = rand.Intn(4000-1000) + 1000

			cont := strconv.Itoa(trr.Contrato)

			if len(datapack) > 0 {
				timeacerto = datapack[len(datapack)-1].T
				timeaprox = trr.Time
				layout := "2006-01-02 15:04:05"
				datet1, _ := time.Parse(layout, timeacerto)
				datet2, _ := time.Parse(layout, timeaprox)
				result = int(datet2.Sub(datet1).Seconds())
			}
			if trr.Contrato != 0 {

				if postocontrato == " " {

					convertid := strconv.Itoa(trr.Id)
					insertuser(contract, convertid, strconv.Itoa(trr.TamanhoTanque))
					Createevent(contract, comb, float64(trr.Contrato), name, postoid)
					postocontrato = cont

				}

				if postocontrato != cont && len(datapack) > 500 {

					//send data
					fmt.Println("caso2 Contrato diferente")
					if GetStatusEvent(contract, name, postoid) {

						//(contract, datapack, name, postoid)
						datapack = ruido(datapack)
						aux := EncodeToBytes(datapack)
						aux = Compress(aux)

						CreatePath(contract, aux, name, postoid, alpha)
						CloseEvent(contract, name, "0.0")
					}

					//resetvetortuplas
					datapack = nil

					//create new event
					Createevent(contract, comb, float64(trr.Contrato), name, postoid)

					postocontrato = cont

				} else if len(datapack) >= meta && sentfile == filesss && result >= 0 {

					//send data

					if GetStatusEvent(contract, name, postoid) {

						datapack = ruido(datapack)
						aux := EncodeToBytes(datapack)
						aux = Compress(aux)

						//fmt.Println(datapack[0].Comb, datapack[0].T, datapack[0].Pos)

						CreatePath(contract, aux, name, postoid, alpha)
					}

					//resetvetortuplas
					datapack = nil
				} else if len(datapack) >= 500 && sentfile != filesss && result >= 0 {
					sentfile = filesss

					if GetStatusEvent(contract, name, postoid) {

						datapack = ruido(datapack)
						aux := EncodeToBytes(datapack)
						aux = Compress(aux)

						CreatePath(contract, aux, name, postoid, alpha)
					}

					datapack = nil
				} else if result < 0 && len(datapack) >= 500 {
					sentfile = filesss

					if GetStatusEvent(contract, name, postoid) {

						datapack = ruido(datapack)
						aux := EncodeToBytes(datapack)
						aux = Compress(aux)
						CreatePath(contract, aux, name, postoid, alpha)
					}

					datapack = nil
				} else if result < 0 && len(datapack) <= 500 {

					fmt.Println("caso6 ")
					datapack = nil
				}

				r := 0 + rand.Float64()*(1-0)
				if r <= alpha {

					lat := fmt.Sprintf("%f", trr.Pos.Lat)
					long := fmt.Sprintf("%f", trr.Pos.Long)
					datapack = append(datapack, Tuple{T: trr.Time, Pos: lat + "/" + long, Comb: math.Round(comb*1000) / 1000})
				}
			}

			/*
				if postocontrato != datasplit[9] && len(datapack) > 0 {

					datakalman := []float64{}

					for _, n := range datapack {
						datakalman = append(datakalman, n.Comb)

					}

					datakalman, _ = KalmanFilter(50, datakalman)

					datakalmanstrin := []string{}

					for _, n := range datakalman {
						s := fmt.Sprintf("%v", n)
						datakalmanstrin = append(datakalmanstrin, s)

					}

					datapack = nil
			*/

		}
		datapack = nil
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}

}

func insertuser(contract *client.Contract, name, tanque string) {

	if !checkuserexist(contract, name) {
		Createuser(contract, name, name, tanque)
	}
}

func processing_Data(contract *client.Contract, k float64) {

	files, err := os.ReadDir("datavehicle")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var dir []string
	var datafile []string
	var m = make(map[string][]string)
	for _, file := range files {

		dir = append(dir, file.Name())
	}

	/*

		for i := range dir {
			go insertuser(contract, dir[i])
			time.Sleep(100 * time.Millisecond)
		}
	*/

	for _, dire := range dir {
		files, err = os.ReadDir("datavehicle/" + dire)

		for _, file := range files {

			datafile = append(datafile, file.Name())
		}

		m[dire] = datafile

		if len(m) == 1 {
			for i, ui := range m {
				send_data(contract, i, ui, k)
			}
			/*

				go
				time.Sleep(10 * time.Second)


				for k := range m {
					delete(m, k)
				}
			*/
			//60
			//time.Sleep(10 * time.Second)

		}

		datafile = nil

	}

	//lst := [30]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29"}

	//for i, file := range m {

	/*
		if contains(lst, i) {
			fmt.Println("Inicio", i)
			csv_write(contract, i, file, 1)
			fmt.Println("Termino", i)

		}
	*/

	//time.Sleep(1 * time.Second)

	//}

}

func contains(s [30]string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
