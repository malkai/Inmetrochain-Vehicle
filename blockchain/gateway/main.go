package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type Results struct {
	Data []datavehicle
}
type datavehicle struct {
	Uservehicle struct {
		Id  int `json:"id"`
		Pos struct {
			Long float64 `json:"long"`
			Lat  float64 `json:"lat"`
		} `json:"pos"`
		Time        string `json:"time"`
		Combustivel string `json:"combustivel"`

		//  FooBar  string `json:"foo.bar"`
	} `json:"uservehicle"`
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
	/*
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
		processing_Data(contract)



			name1 := "datawithoutnoise.csv"
			name2 := "datanoise.csv"
			name3 := "datatimstump.csv"
			name4 := "filter.csv"
	*/
	//GetAllevents(contract, "event")
	//openhugejson(contract)

	//fmt.Println("oi")

	//send_data(contract,"0")
	//GetAlluser(contract, "user100")

	estatistica("Kalmanfilter/", "filter.csv")
}

// estou inserindo um ruido branco guaisiano com o sigma
// linear interpolation

func ruido(simudata []Tuple) []Tuple {

	noise := []float64{}
	interpolationdata := []float64{}

	for _, data := range simudata {

		noise = append(noise, data.Comb)

	}

	//fmt.Println(noise)
	if len(noise) > 0 {
		interpolationdata = interpolation(noise)
	}
	//fmt.Println(interpolationdata)

	for i := range simudata {

		simudata[i].Comb = interpolationdata[i]

	}

	return simudata
}

func csv_write(contract *client.Contract, name string, file []string, alpha float64) {

	//postoid := "9999"
	datapack := []Tuple{}
	postocontrato := " "
	//timeacerto := ""
	//timeaprox := ""
	//sentfile := file[0]
	//result := 0
	//contratocomb := 0.0
	//meta := 0
	comb := 0.0
	combant := 0.0
	datacomb := []string{}
	datacomb2 := []float64{}
	datatime := []string{}
	file2, err := os.Create("datacomb/" + name + "datawithoutnoise.csv")
	if err != nil {
		panic(err)
	}
	defer file2.Close()
	file3, err := os.Create("datacombnoise/" + name + "datanoise.csv")
	if err != nil {
		panic(err)
	}
	file4, err := os.Create("datatimstump/" + name + "datatimstump.csv")
	if err != nil {
		panic(err)
	}

	file5, err := os.Create("Kalmanfilter/" + name + "filter.csv")
	if err != nil {
		panic(err)
	}
	defer file2.Close()
	defer file3.Close()
	defer file4.Close()
	defer file5.Close()
	for _, filesss := range file {

		file, err := os.Open("datavehicle/" + name + "/" + filesss)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		//analise1
		//meta = (rand.Intn(4000-1000) + 1000)

		scanner := bufio.NewScanner(file)
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {

			reg := regexp.MustCompile("[^\n0-9.-]+")
			datasplit := reg.Split(scanner.Text(), 11)
			a := reg.Split(datasplit[4], 2)
			combant = comb
			comb, err = strconv.ParseFloat(a[0], 64)
			if err != nil {
				panic(err)
			}

			data1 := reg.Split(datasplit[5], 2)
			data2 := reg.Split(datasplit[6], 2)
			data3 := reg.Split(datasplit[7], 2)
			data4 := reg.Split(datasplit[8], 2)

			//correctfloat := reg.Split(datasplit[9], 2)

			//contratocomb, err = strconv.ParseFloat(correctfloat[0], 64)
			if err != nil {
				panic(err)
			}
			//fmt.Println(contratocomb)

			comb = (comb * 100) / 50
			//fmt.Println(comb)

			lat := reg.Split(datasplit[2], 2)
			long := reg.Split(datasplit[3], 2)

			if postocontrato != datasplit[9] && len(datapack) > 0 {

				fmt.Println(combant, len(datapack))

				writer := csv.NewWriter(file2)
				for _, row := range datapack {
					s := fmt.Sprintf("%v", row.Comb)
					d := fmt.Sprintf("%v", row.T)
					datacomb = append(datacomb, s)
					datatime = append(datatime, d)
					datacomb2 = append(datacomb2, row.Comb)

				}

				writer.Write(datacomb)
				writer.Flush()
				datacomb = nil
				datapack = ruido(datapack)

				writer = csv.NewWriter(file3)
				for _, row := range datapack {
					s := fmt.Sprintf("%v", row.Comb)
					datacomb = append(datacomb, s)

				}
				writer.Write(datacomb)
				writer.Flush()

				writer = csv.NewWriter(file4)
				writer.Write(datatime)
				writer.Flush()

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
				writer = csv.NewWriter(file5)
				writer.Write(datakalmanstrin)
				writer.Flush()

				datacomb = nil
				datatime = nil
				datapack = nil
				datacomb2 = nil
			}
			postocontrato = datasplit[9]

			//r := 0 + rand.Float64()*(1-0)
			//if r < alpha {
			datapack = append(datapack, Tuple{T: data1[0] + " " + data2[0] + ":" + data3[0] + ":" + data4[0], Pos: lat[0] + "/" + long[0], Comb: math.Round(comb*1000) / 1000})
			//}

		}
		datapack = nil
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}

}

func send_data(contract *client.Contract, name string, file []string, alpha float64) {

	//postoid := "9999"
	datapack := []Tuple{}
	postocontrato := " "
	//timeacerto := ""
	//timeaprox := ""
	//sentfile := file[0]
	//result := 0
	//contratocomb := 0.0
	//meta := 0
	comb := 0.0
	combant := 0.0
	datacomb := []string{}
	datacomb2 := []float64{}
	datatime := []string{}
	file2, err := os.Create("datacomb/" + name + "datawithoutnoise.csv")
	if err != nil {
		panic(err)
	}
	defer file2.Close()
	file3, err := os.Create("datacombnoise/" + name + "datanoise.csv")
	if err != nil {
		panic(err)
	}
	file4, err := os.Create("datatimstump/" + name + "datatimstump.csv")
	if err != nil {
		panic(err)
	}

	file5, err := os.Create("Kalmanfilter/" + name + "filter.csv")
	if err != nil {
		panic(err)
	}
	defer file2.Close()
	defer file3.Close()
	defer file4.Close()
	defer file5.Close()
	for _, filesss := range file {

		file, err := os.Open("datavehicle/" + name + "/" + filesss)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		//analise1
		//meta = (rand.Intn(4000-1000) + 1000)

		scanner := bufio.NewScanner(file)
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {

			reg := regexp.MustCompile("[^\n0-9.-]+")
			datasplit := reg.Split(scanner.Text(), 11)
			a := reg.Split(datasplit[4], 2)
			combant = comb
			comb, err = strconv.ParseFloat(a[0], 64)
			if err != nil {
				panic(err)
			}

			data1 := reg.Split(datasplit[5], 2)
			data2 := reg.Split(datasplit[6], 2)
			data3 := reg.Split(datasplit[7], 2)
			data4 := reg.Split(datasplit[8], 2)

			//correctfloat := reg.Split(datasplit[9], 2)

			//contratocomb, err = strconv.ParseFloat(correctfloat[0], 64)
			if err != nil {
				panic(err)
			}
			//fmt.Println(contratocomb)

			comb = (comb * 100) / 50
			//fmt.Println(comb)

			lat := reg.Split(datasplit[2], 2)
			long := reg.Split(datasplit[3], 2)

			//contratocomb
			/*
				if len(datapack) > 0 {
					timeacerto = datapack[len(datapack)-1].T
					timeaprox = data1[0] + " " + data2[0] + ":" + data3[0] + ":" + data4[0]
					layout := "2006-01-02 15:04:05"
					datet1, _ := time.Parse(layout, timeacerto)

					datet2, _ := time.Parse(layout, timeaprox)

					result = int(datet2.Sub(datet1).Seconds())
				}
			*/
			/*

				//fmt.Println(data1[0] + " " + data2[0] + ":" + data3[0] + ":" + data4[0])

				if postocontrato == " " {

						a := reg.Split(datasplit[4], 2)

						comb, err := strconv.ParseFloat(a[0], 64)
						if err != nil {
							// ... handle error
							panic(err)
						}
						Createevent(contract, comb, contratocomb, name, postoid)

					//fmt.Println("Criando contrato " + name)
					postocontrato = datasplit[9]

				}
				if postocontrato != datasplit[9] {
					//fmt.Println("pode enviar contrato diferente " + name)
					//send data

						if GetStatusEvent(contract, name, postoid) {
							datapack = ruido(datapack)
							//(contract, datapack, name, postoid)
							CreatePath(contract, datapack, name, postoid)
						}

						meta = rand.Intn(4000-1000) + 1000

					//resetvetortuplas
					datapack = nil

						//create new event
						correctfloat := reg.Split(datasplit[9], 2)

						a := reg.Split(datasplit[4], 2)

						comb, err := strconv.ParseFloat(a[0], 64)
						if err != nil {
							panic(err)
						}

						contratocomb, err := strconv.ParseFloat(correctfloat[0], 64)
						if err != nil {
							panic(err)
						}

						Createevent(contract, comb, contratocomb, name, postoid)

					postocontrato = datasplit[9]

				} else if len(datapack) >= meta && sentfile == filesss && result >= 0 {

					//send data

						if GetStatusEvent(contract, name, postoid) {

							datapack = ruido(datapack)

							CreatePath(contract, datapack, name, postoid)
						}

					meta = rand.Intn(4000-1000) + 1000

					//resetvetortuplas
					datapack = nil
				}


					else if len(datapack) >= 500 && sentfile != filesss && result >= 0 {
						sentfile = filesss

						if GetStatusEvent(contract, name, postoid) {

							datapack = ruido(datapack)

							CreatePath(contract, datapack, name, postoid)
						}
						datapack = nil
					} else if result < 0 && len(datapack) >= 500 {
						sentfile = filesss

						if GetStatusEvent(contract, name, postoid) {

							datapack = ruido(datapack)

							CreatePath(contract, datapack, name, postoid)
						}
						meta = rand.Intn(4000-1000) + 1000
						datapack = nil
					} else if result < 0 && len(datapack) <= 500 {

						datapack = nil
					}
			*/
			if postocontrato != datasplit[9] && len(datapack) > 0 {

				fmt.Println(combant, len(datapack))

				writer := csv.NewWriter(file2)
				for _, row := range datapack {
					s := fmt.Sprintf("%v", row.Comb)
					d := fmt.Sprintf("%v", row.T)
					datacomb = append(datacomb, s)
					datatime = append(datatime, d)
					datacomb2 = append(datacomb2, row.Comb)

				}

				writer.Write(datacomb)
				writer.Flush()
				datacomb = nil
				datapack = ruido(datapack)

				writer = csv.NewWriter(file3)
				for _, row := range datapack {
					s := fmt.Sprintf("%v", row.Comb)
					datacomb = append(datacomb, s)

				}
				writer.Write(datacomb)
				writer.Flush()

				writer = csv.NewWriter(file4)
				writer.Write(datatime)
				writer.Flush()

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
				writer = csv.NewWriter(file5)
				writer.Write(datakalmanstrin)
				writer.Flush()

				datacomb = nil
				datatime = nil
				datapack = nil
				datacomb2 = nil
			}
			postocontrato = datasplit[9]

			//r := 0 + rand.Float64()*(1-0)
			//if r < alpha {
			datapack = append(datapack, Tuple{T: data1[0] + " " + data2[0] + ":" + data3[0] + ":" + data4[0], Pos: lat[0] + "/" + long[0], Comb: math.Round(comb*1000) / 1000})
			//}

		}
		datapack = nil
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}

}

func insertuser(contract *client.Contract, name string) {

	fmt.Println(name)
	if !checkuserexist(contract, name) {
		Createuser(contract, name, name, "50")
	}
}

func processing_Data(contract *client.Contract) {

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

		//if len(m) == 15 {
		//for k := range m {
		//	delete(m, k)
		//}
		//60
		//time.Sleep(5 * time.Second)

		//}

		datafile = nil

	}

	lst := [30]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29"}

	for i, file := range m {

		if contains(lst, i) {
			fmt.Println("Inicio", i)
			csv_write(contract, i, file, 1)
			fmt.Println("Termino", i)

		}
		//go send_data(contract, i, file, 0.33)
		//time.Sleep(1 * time.Second)

	}

}

func contains(s [30]string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
