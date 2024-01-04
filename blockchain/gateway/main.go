package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type Results struct {
	Data []datavehicle
}
type datavehicle struct {
	Uservehicle struct {
		Id   string `json:"id"`
		Path string `json:"path"`

		Tuple struct {
			Pos struct {
				Long float64 `json:"long"`
				Lat  float64 `json:"lat"`
			} `json:"pos"`
			Time        string `json:"time"`
			Combustivel string `json:"combustivel"`
		} `json:"Tuple"`

		//  FooBar  string `json:"foo.bar"`
	} `json:"uservehicle"`
}

type datahelp struct {
	id   string
	Data []datavehicle
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
	//GetAllevents(contract, "event")
	openhugejson(contract)
	//GetAlluser(contract, "user100")

}

func inserigoroutine(contract *client.Contract, aaa []Tuple, id string) {

	aaa = selectionSort(aaa, len(aaa)-1)
	//idposto := "9999"
	//fmt.Println(tt)
	layout := "2006-01-02 15:04:05"
	acc := 0.0
	if aaa != nil {
		//fmt.Println(tt)
		//res1 := strings.Split(id, "_")
		//t := strconv.Itoa(res1[0])
		//dadosbrutos(contract, aaa, res1[0], idposto)
		fmt.Println(len(aaa))
		for help := range aaa {
			if aaa[help] != aaa[len(aaa)-1] {
				date, _ := time.Parse(layout, aaa[help].T)
				date2, _ := time.Parse(layout, aaa[help+1].T)
				acc = acc + float64(date2.Sub(date).Seconds())
			}
		}
		fmt.Println(acc)
		fmt.Println(acc / 3600)

	}

	//idposto := "9999"

	//fmt.Println("tuple", aaa)
	//fmt.Println("id", id)

	//	fmt.Println(len(aaa))

	/*
		if len(aaa) > 1 {
			//fmt.Println(m)
			for k := range aaa {
				//fmt.Println(k)
				tt = append(tt, aaa[k])
			}

			tt = selectionSort(tt, len(tt)-1)
			//fmt.Println(tt)
			if tt != nil {
				//fmt.Println(tt)
				t := strconv.Itoa(id)
				dadosbrutos(contract, tt, t, idposto)
			}
		}
	*/
	/*
		helppath = h.Uservehicle.Path
		res1 = strings.Split(h.Uservehicle.Vehicle_data.Tuple[2].Combustivel, "%")

		s, _ := strconv.ParseFloat(res1[0], 64)

		if s < 0 {
			s = s * -1
		}

		if GetStatusEvent(contract, h.Uservehicle.Id, idposto) == true {
			Createevent(contract, s, 10.0, h.Uservehicle.Id, idposto)
		}

		if j > 1 {
			//fmt.Println(m)
			for k := range m {
				//fmt.Println(k)
				tt = append(tt, m[k])
			}

			tt = selectionSort(tt, len(tt)-1)
			//fmt.Println(tt)
			if tt != nil {
				//fmt.Println(tt)
				dadosbrutos(contract, tt, h.Uservehicle.Id, idposto)
			}
		}
		j = 0
		tt = nil
		for k := range m {
			delete(m, k)
		}
		m[j] = Tuple{T: date.Format(layout), Pos: lat + "/" + long, Comb: s}
		j++
	*/

	/*
		s, _ := strconv.ParseFloat(res1[0], 64)

					if checkuserexist(contract, h.Uservehicle.Id) == true {
						Createuser(contract, h.Uservehicle.Id, h.Uservehicle.Id, "40")
					}

					if s < 0 {
						s = s * -1
					}
					if GetStatusEvent(contract, h.Uservehicle.Id, idposto) == false {
						Createevent(contract, s, 10.0, h.Uservehicle.Id, idposto)
					}

					if j > 1 {

						for k := range m {
							//fmt.Println(k)
							tt = append(tt, m[k])
						}

						tt = selectionSort(tt, len(tt)-1)

						if tt != nil {
							//fmt.Println(tt)
							dadosbrutos(contract, tt, h.Uservehicle.Id, idposto)
						}
					}
	*/
}

func openhugejson(contract *client.Contract) {

	content, err := os.Open("vehicles_data2.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	defer content.Close()
	i := 0

	dec := json.NewDecoder(content)

	// reacontentd open bracket
	t, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)
	idposto := "9999"
	j := 0
	o := 0
	soma_part_trajeto := 0
	helpid := ""
	helppath := ""
	layout := "2006-01-02 15:04:05"
	//var tt []Tuple
	//var sizedata uintptr
	trajetos_num := 0
	m := make(map[int]Tuple)
	m2 := make(map[string][]Tuple)
	m3 := make(map[int]int)
	var oo []Tuple
	// while the array contains values
	for dec.More() {
		var h datavehicle
		// decode an array value (Message)
		err := dec.Decode(&h)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(h.Uservehicle.Id)

		res1 := strings.Split(h.Uservehicle.Id, "_")
		h.Uservehicle.Id = res1[0]
		h.Uservehicle.Path = res1[1]

		//inseri tuplas
		if i == 0 {
			helpid = h.Uservehicle.Id
			helppath = h.Uservehicle.Path
			//print(checkuserexist(contract, h.Uservehicle.Id))

			if checkuserexist(contract, h.Uservehicle.Id) == false {
				Createuser(contract, h.Uservehicle.Id, h.Uservehicle.Id, "40")
			}

			// checar se há um evento aberto se não tiver abrir um

			s, _ := strconv.ParseFloat(h.Uservehicle.Tuple.Combustivel, 64)
			if GetStatusEvent(contract, h.Uservehicle.Id, idposto) == false {
				Createevent(contract, s, 10.0, h.Uservehicle.Id, idposto)
			}

			//sizedata = unsafe.Sizeof(m) + unsafe.Sizeof(make([]Tuple, int(j)))
		}

		i++

		//res1 = strings.Split(h.Uservehicle.Vehicle_data.Tuple[2].Combustivel, "%")
		//s, err := strconv.ParseFloat(res1[0], 64)
		if err != nil {
			fmt.Printf("Erro na leitura d combustivel: %v", err)
		}
		lat := fmt.Sprintf("%f", h.Uservehicle.Tuple.Pos.Lat)

		long := fmt.Sprintf("%f", h.Uservehicle.Tuple.Pos.Long)

		date, error := time.Parse(layout, h.Uservehicle.Tuple.Time)
		if error != nil {
			fmt.Println(error)
		}
		if o >= 10 || dec.More() == false {
			//fmt.Println(o, m2)

			//fmt.Println(m3)

			for lo, url := range m2 {
				fmt.Println(lo, len(url))
				go inserigoroutine(contract, url, lo)

			}
			time.Sleep(2 * time.Second)

			o = 0

			for k := range m2 {
				delete(m2, k)
			}
			for k := range m3 {
				delete(m3, k)
			}
			// && j <= 5000
		} else if h.Uservehicle.Id == helpid && h.Uservehicle.Path == helppath {

			res1 = strings.Split(h.Uservehicle.Tuple.Combustivel, "%")

			s, _ := strconv.ParseFloat(res1[0], 64)
			if s < 0 {
				s = s * -1
			}
			m[j] = Tuple{T: date.Format(layout), Pos: lat + "/" + long, Comb: s}
			j++

		} else if h.Uservehicle.Id == helpid && (h.Uservehicle.Path != helppath) || j > 5000 {
			trajetos_num++

			if h.Uservehicle.Path == helppath {

				soma_part_trajeto++
			} else {
				trajetos_num++
				soma_part_trajeto = 0
			}

			fmt.Println("tam", j, "Caminho diferente path1 ", h.Uservehicle.Path, "pathwatch ", helppath)
			//fmt.Println(len(m))
			helppath = h.Uservehicle.Path
			idd, _ := strconv.Atoi(h.Uservehicle.Id)

			for _, kk := range m {
				//fmt.Println(k)
				oo = append(oo, kk)
			}

			t := strconv.Itoa(soma_part_trajeto)
			m2["_"+h.Uservehicle.Id+"_"+h.Uservehicle.Path+"_"+t] = oo

			m3[o] = idd

			for k := range m {
				delete(m, k)
			}
			s, _ := strconv.ParseFloat(res1[0], 64)

			if s < 0 {
				s = s * -1
			}

			if GetStatusEvent(contract, h.Uservehicle.Id, idposto) == true {
				Createevent(contract, s, 10.0, h.Uservehicle.Id, idposto)
			}

			oo = nil
			j = 0
			m[j] = Tuple{T: date.Format(layout), Pos: lat + "/" + long, Comb: s}
			j++
			o++

		} else if h.Uservehicle.Id != helpid {

			//checar se id existe
			//fmt.Println("ID diferente id1", h.Uservehicle.Id, "id2", helpid)
			//t := strconv.Itoa(soma_part_trajeto)
			res1 = strings.Split(h.Uservehicle.Tuple.Combustivel, "%")
			s, _ := strconv.ParseFloat(res1[0], 64)

			if s < 0 {
				s = s * -1
			}

			if checkuserexist(contract, h.Uservehicle.Id) == false {
				Createuser(contract, h.Uservehicle.Id, h.Uservehicle.Id, "40")
			}

			if GetStatusEvent(contract, h.Uservehicle.Id, idposto) == false {
				Createevent(contract, s, 10.0, h.Uservehicle.Id, idposto)
			}

			helpid = h.Uservehicle.Id
			helppath = h.Uservehicle.Path
			//tt = nil
			j = 0
			idd, _ := strconv.Atoi(h.Uservehicle.Id)
			for _, kk := range m {
				//fmt.Println(k)
				oo = append(oo, kk)
			}
			t := strconv.Itoa(soma_part_trajeto)
			m2["_"+h.Uservehicle.Id+"_"+h.Uservehicle.Path+"_"+t] = oo

			m3[o] = idd
			oo = nil
			soma_part_trajeto = 0
			for k := range m {
				delete(m, k)
			}

			m[j] = Tuple{T: date.Format(layout), Pos: lat + "/" + long, Comb: s}
			j++
			o++

		}

	}
	fmt.Println("Numero de trajetos", trajetos_num)
	t, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)
}

/*
func openhugejson2(contract *client.Contract) {

		content, err := os.Open("vehicles_data2.json")
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}
		defer content.Close()
		i := 0

		dec := json.NewDecoder(content)

		// reacontentd open bracket
		t, err := dec.Token()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%T: %v\n", t, t)
		idposto := "9999"
		j := 0
		helpid := ""
		helppath := ""
		layout := "2006-01-02 15:04:05"
		var tt []Tuple
		//var sizedata uintptr
		m := make(map[int]Tuple)
		// while the array contains values
		for dec.More() {
			var h datavehicle
			// decode an array value (Message)
			err := dec.Decode(&h)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println(h.Uservehicle.Id)

			res1 := strings.Split(h.Uservehicle.Id, "_")
			h.Uservehicle.Id = res1[0]
			h.Uservehicle.Path = res1[1]

			//inseri tuplas
			if i == 0 {
				helpid = h.Uservehicle.Id
				helppath = h.Uservehicle.Path
				//print(checkuserexist(contract, h.Uservehicle.Id))
				if checkuserexist(contract, h.Uservehicle.Id) == false {
					Createuser(contract, h.Uservehicle.Id, h.Uservehicle.Id, "40")
				}

				// checar se há um evento aberto se não tiver abrir um

				//s, _ := strconv.ParseFloat(h.Uservehicle.Vehicle_data.Tuple[2].Combustivel, 64)
				if GetStatusEvent(contract, h.Uservehicle.Id, idposto) == false {
					Createevent(contract, s, 10.0, h.Uservehicle.Id, idposto)
				}

				//sizedata = unsafe.Sizeof(m) + unsafe.Sizeof(make([]Tuple, int(j)))
			}

			i++

			//res1 = strings.Split(h.Uservehicle.Vehicle_data.Tuple[2].Combustivel, "%")
			//s, err := strconv.ParseFloat(res1[0], 64)
			if err != nil {
				fmt.Printf("Erro na leitura d combustivel: %v", err)
			}
			//lat := fmt.Sprintf("%f", h.Uservehicle.Vehicle_data.Tuple[0].Pos.Lat)

			//long := fmt.Sprintf("%f", h.Uservehicle.Vehicle_data.Tuple[0].Pos.Long)

			//date, error := time.Parse(layout, h.Uservehicle.Vehicle_data.Tuple[1].Time)
			if error != nil {
				fmt.Println(error)
			}

			if h.Uservehicle.Id == helpid && h.Uservehicle.Path == helppath && j < 5000 {
				res1 = strings.Split(h.Uservehicle.Vehicle_data.Tuple[2].Combustivel, "%")

				s, _ := strconv.ParseFloat(res1[0], 64)
				if s < 0 {
					s = s * -1
				}
				m[j] = Tuple{T: date.Format(layout), Pos: lat + "/" + long, Comb: s}
				j++

			} else if h.Uservehicle.Id == helpid && (h.Uservehicle.Path != helppath || j > 5000) {
				fmt.Print(j > 5000, h.Uservehicle.Path != helppath, (h.Uservehicle.Path != helppath || j > 5000))
				fmt.Println("id", h.Uservehicle.Id, "Caminho  path1 ", h.Uservehicle.Path, "pathwatch ", helppath)
				helppath = h.Uservehicle.Path
				res1 = strings.Split(h.Uservehicle.Vehicle_data.Tuple[2].Combustivel, "%")

				s, _ := strconv.ParseFloat(res1[0], 64)

				if s < 0 {
					s = s * -1
				}

				if GetStatusEvent(contract, h.Uservehicle.Id, idposto) == true {
					Createevent(contract, s, 10.0, h.Uservehicle.Id, idposto)
				}

				if j > 1 {
					//fmt.Println(m)
					for k := range m {
						//fmt.Println(k)
						tt = append(tt, m[k])
					}

					tt = selectionSort(tt, len(tt)-1)
					//fmt.Println(tt)
					if tt != nil {
						//fmt.Println(tt)
						dadosbrutos(contract, tt, h.Uservehicle.Id, idposto)
					}
				}
				j = 0
				tt = nil
				for k := range m {
					delete(m, k)
				}
				m[j] = Tuple{T: date.Format(layout), Pos: lat + "/" + long, Comb: s}
				j++

			} else if h.Uservehicle.Id != helpid {

				//checar se id existe
				fmt.Println("ID diferente id1", h.Uservehicle.Id, "id2", helpid)
				res1 = strings.Split(h.Uservehicle.Vehicle_data.Tuple[2].Combustivel, "%")

				s, _ := strconv.ParseFloat(res1[0], 64)

				if checkuserexist(contract, h.Uservehicle.Id) == true {
					Createuser(contract, h.Uservehicle.Id, h.Uservehicle.Id, "40")
				}

				if s < 0 {
					s = s * -1
				}
				if GetStatusEvent(contract, h.Uservehicle.Id, idposto) == false {
					Createevent(contract, s, 10.0, h.Uservehicle.Id, idposto)
				}

				if j > 1 {

					for k := range m {
						//fmt.Println(k)
						tt = append(tt, m[k])
					}

					tt = selectionSort(tt, len(tt)-1)

					if tt != nil {
						//fmt.Println(tt)
						dadosbrutos(contract, tt, h.Uservehicle.Id, idposto)
					}
				}
				helpid = h.Uservehicle.Id
				helppath = h.Uservehicle.Path
				tt = nil
				j = 0
				for k := range m {
					delete(m, k)
				}
				m[j] = Tuple{T: date.Format(layout), Pos: lat + "/" + long, Comb: s}
				j++

			}

		}

		t, err = dec.Token()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%T: %v\n", t, t)
	}
*/
func dadosbrutos(contract *client.Contract, tuples []Tuple, id string, id2 string) {
	/*
		Análise 0: programação dos contratos inteligentes para inserção,
		cálculo de créditos, e consulta de carteiras;
	*/
	CreatePath(contract, tuples, id, id2)
}

func analise1(contract *client.Contract, tuples []Tuple, id string, id2 string) {

	/*
		Análise 1: influência de cada métrica no total de créditos,
		inserindo perda seletiva de dados para a Frequência (fixar em 10%);
		antes da inserção (go e excel )
	*/
	CreatePath(contract, tuples, id, id2)

}

func analise2(contract *client.Contract, tuples []Tuple, id string, id2 string) {

	/*
		Análise 2: avaliar apenas a Frequência
		(variando a perda seletiva entre 5% a 35%, e podendo também variar k),
		verificando como ela se comporta nestes cenários; (go e excel)
	*/
	CreatePath(contract, tuples, id, id2)

}
func selectionSort(array []Tuple, size int) []Tuple {

	layout := "2006-01-02 15:04:05"
	for ind := 0; ind <= size; ind++ {
		min_index := ind

		for j := ind + 1; j <= size; j++ {
			// select the minimum element in every iteration
			datet1, err := time.Parse(layout, array[ind].T)
			if err != nil {
				fmt.Println("Erro")
			}
			datet2, err := time.Parse(layout, array[j].T)
			if err != nil {
				fmt.Println("Erro")
			}
			if datet2.Compare(datet1) == -1 {

				min_index = j
				v := array[ind]
				array[ind] = array[min_index]
				array[min_index] = v

			}

		}

	}

	return array
}

/*
	content, err := os.Open("vehicles_data2.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	byteValueJSON, _ := io.ReadAll(content)

	ts := []datavehicle{}
	err = json.Unmarshal(byteValueJSON, &ts)
	if err != nil {
		log.Fatal(err)
	}
*/

//{{9_3  {[{{-22.802550583789813 -43.202926866082976}  } {{0 0} 2023-12-18 14:49:57.775409 } {{0 0}  98.966920%}]}}}

//GetAllPath(contract, "Path")
//GetAlluser(contract, "user")
//GetAllevents(contract, "Event")

// Create gRPC client connection, which should be shared by all gateway connections to this endpoint.

//GetAllevents(contract, "event")

//Createevent(contract, 93.00, 5.0, "1", "3")

//GetAllevents(contract, "event1")

//Createuser(contract, "1", "Malkai1", "80")
//GetAlluser(contract, "user")
//GetAllevents(contract, "Event")

/*
	var t []Tuple
	layout := "2006-01-02 15:04:05"
	time2 := time.Now()

	var t1 = Tuple{T: time2.Format(layout), Pos: "1/2", Comb: 93.00}
	t = append(t, t1)
	time2 = time2.Add(10 * time.Second)
	//fmt.Println(time2.Format(layout))
	t1 = Tuple{T: time2.Format(layout), Pos: "2/3", Comb: 92.00}
	t = append(t, t1)
	time2 = time2.Add(10 * time.Second)
	//fmt.Println(time2.Format(layout))
	t1 = Tuple{T: time2.Format(layout), Pos: "4/5", Comb: 91.00}
	t = append(t, t1)
	time2 = time2.Add(10 * time.Second)
	//fmt.Println(time2.Format(layout))
	t1 = Tuple{T: time2.Format(layout), Pos: "5/6", Comb: 90.00}
	t = append(t, t1)

	CreatePath(contract, t, "1")
*/
//GetAllPath(contract, "Path")
//GetAlluser(contract, "user")
//GetAllevents(contract, "Event")

//initLedger(contract)

// Context used for event listening
//ctx, cancel := context.WithCancel(context.Background())
//defer cancel()

// Listen for events emitted by subsequent transactions
//startChaincodeEventListening(ctx, network)
/*
	Datai := time.Now()
	Dataiataf := time.Now()
	Fsupi := 50.00
	Fsupf := 20.00
	Dff := 10
	Vstatus := false
	Iduser1 := "1"
	Iduser2 := "2"
*/
//Createevent(contract, time.Now(), 50.00, 10.0, "1", "2")
//GetAllevents(contract, "event1")

//updatevent("event1", 1.0)
//updatevent("event1", 1.0)
//updatevent("event1", 1.0)
//updatevent("event1", 1.0)
//updatevent("event1", 1.0)
//updatevent("event1", 1.0)
//updatevent("event1", 1.0)
//updatevent("event1", 1.0)

//initLedger(contract)
//getAllAssets(contract)

//replayChaincodeEvents(ctx, network, firstBlockNumber)

//GetAlluser(contract, "user")
//Createuser(contract, "0", "0", "80")
//checkuserexist(contract, "0")
//Getuser(contract, "0")
//fmt.Println(checkuserexist(contract, "0"))

//fmt.Println()

//fmt.Println()
//fmt.Println(array)
//fmt.Println("Sucesso")

///	byteValueJSON, _ := io.ReadAll(content)

//ts := []datavehicle{}
/*
	err = json.Unmarshal(byteValueJSON, &ts)
	if err != nil {
		log.Fatal(err)
	}
*/

func memUsage(m1, m2 *runtime.MemStats) {
	fmt.Println("Alloc:", m2.Alloc-m1.Alloc,
		"TotalAlloc:", m2.TotalAlloc-m1.TotalAlloc,
		"HeapAlloc:", m2.HeapAlloc-m1.HeapAlloc)
}
