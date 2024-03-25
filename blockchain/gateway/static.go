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

func csv_write(contract *client.Contract, name string, file []string, alpha float64) {

	datapack := []Tuple{}
	postocontrato := " "

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
