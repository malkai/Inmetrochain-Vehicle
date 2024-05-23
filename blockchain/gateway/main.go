package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

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
	//processing_Data(contract, 3)
	gera_dados(contract)

}

// estou inserindo um ruido branco guaisiano com o sigma
// linear interpolation

func gera_dados(contract *client.Contract) {

	tipo := ""

	writer := &csv.Writer{}

	for i := 1; i <= 151; i++ {
		info := []string{}
		help := strconv.Itoa(i)
		if checkuserexist(contract, strconv.Itoa(i)) {
			a := Getuser(contract, help)
			user := User{}
			err := json.Unmarshal([]byte(a), &user)
			if err != nil {
				fmt.Println(err.Error())

			}
			if user.Id != "" {
				if tipo != user.Typee {
					tipo = user.Typee
					file, err := os.Create("tabelauser/" + tipo + ".csv")
					if err != nil {
						panic(err)
					} //s
					defer file.Close()
					writer = csv.NewWriter(file)
					info = append(info, "Id", "Tipo", "tamanho do tanque", "Km",
						"L", "Km/l", "litros esperado", "tempo", "Número de Eventos",
						"Número de Trajetos", "média de Trajetos por evento", "desvio de Trajetos por evento",
						"média de Km por trajeto", "desvio de Km por trajeto",
						"média de L por trajeto", "desvio de L por trajeto",
						"Total de Pontos da Completude", "média Completude", "desvio Completude",
						"Total de Pontos Frequencia", "média Frequencia  por evento", "desvio Frequencia por evento",
						"média Frequencia   por trajeto", "desvio Frequencia  por trajeto",
						"Total confiança", "média confiança por evento", "desvio confiança por evento", "moeda")
					writer.Write(info)
					writer.Flush()

					info = nil

					tu := GetAlleventsPast(contract, help)
					time.Sleep(500 * time.Millisecond)
					events := []Event{}
					err = json.Unmarshal([]byte(tu), &events)
					if err != nil {
						fmt.Println(err.Error())

					}

					nev := len(events) //numero de eventos

					ntkm := 0.0
					ntl := 0.0
					ntt := 0.0

					nt := 0             //numero total de trajetos
					ate := []float64{}  //array armazena trajetos por evento
					atkm := []float64{} //array armazena trajetos por km
					atl := []float64{}  //array armazena trajetos por

					gatkm := []float64{} //global array armazena trajetos por km
					gatl := []float64{}  //global array armazena trajetos por litros

					datkm := []float64{} //global array armazena trajetos por km
					datl := []float64{}  //global array armazena trajetos por litros

					mte := 0.0  //media de trajetos
					dte := 0.0  //desvio de trajetos
					mtkm := 0.0 //media de trajetos
					dtkm := 0.0 //desvio de trajetos
					mtl := 0.0  //media de trajetos
					dtl := 0.0  //desvio de trajetos

					tpc := 0.0
					gtpc := []float64{}
					tcpm := 0.0
					dcpm := 0.0

					tf := 0.0
					gfe := []float64{}
					gft := []float64{}
					gftg := []float64{}
					gftd := []float64{}
					mfe := 0.0
					dfe := 0.0
					mft := 0.0
					dft := 0.0

					aconf := []float64{}

					lesperados := 0.0

					for _, er := range events {

						lesperados = lesperados + er.Dff

						trt := strings.Replace(er.Datai, " ", "-", -1)
						trt = strings.Replace(trt, ":", "-", -1)

						tu1 := GetPathhOpen(contract, help, trt)
						paths := []Path{}
						err = json.Unmarshal([]byte(tu1), &paths)
						if err != nil {
							fmt.Println(err.Error())

						}

						for _, ps := range paths {

							/*

								info2 := []string{}

								//fmt.Println(ps.Fuel)

								tuples := Decompress(ps.DataVehicle)

								fuel_vector := []float64{}
								time3 := []float64{0}
								time12 := []float64{0}
								time := 0.0
								//yuis := [][]float64{}

								//fmt.Println(ps.Timeless, "frequencia")

								for i, ts := range tuples {

									fuel_vector = append(fuel_vector, ts.Comb)
									if i < len(tuples)-1 {
										time1 := totaltime(tuples[i].T, tuples[i+1].T)

										time12 = append(time12, time1)

										time = time + time1

										//rtt, _ := strconv.ParseFloat(tuples[i].T, 64)
										time3 = append(time3, time)
									}

								}

								//timeles := Timeliness(time12, 3)

								//					fmt.Println(timeles)

								f := piecewiselinear.Function{Y: fuel_vector} // range: "hat" function
								f.X = time3
								rtt := [][2]float64{}

								for rty := range time3 {
									rtt1 := [2]float64{time3[rty], fuel_vector[rty]}
									rtt = append(rtt, rtt1)
								}
								a, b, R2, err := Linear(rtt)
								if err != nil {
									fmt.Println("Erro em alguma coisa")
								}

								fmt.Println(a, b, R2, tuples[0].Comb, tuples[len(tuples)-1].Comb)

								tui := []float64{}
								for i := range time3 {

									tui = append(tui, time3[i]*a+b)
								}

								///	fmt.Println(tui[0], tui[len(tuples)-1], tui[0]-tui[len(tuples)-1])

								allhelp := []string{}
								allvalues := [][]string{}

								for i := range time3 {
									allhelp = append(allhelp, fmt.Sprintf("%f", time3[i]), fmt.Sprintf("%f", tuples[i].Comb), fmt.Sprintf("%f", tui[i]))
									allvalues = append(allvalues, allhelp)
									allhelp = nil
								}

								file, err := os.Create("dadosblockchain/" + ps.PathID + ".csv")
								if err != nil {
									panic(err)
								}
								defer file.Close()
								writer = csv.NewWriter(file)
								info2 = append(info, "Tempototal", "KM", "Fuel")
								writer.Write(info2)
								writer.Flush()
								info2 = append(info, fmt.Sprintf("%f", ps.Totaltime), fmt.Sprintf("%f", ps.Distance), fmt.Sprintf("%f", ps.Fuel))
								writer.Write(info2)
								writer.Flush()
								info2 = nil
								info2 = append(info, "tempo", "comb", "Regression Linear")
								writer.Write(info2)
								writer.Flush()
								writer.WriteAll(allvalues)
							*/

							atkm = append(atkm, ps.Distance)
							atl = append(atl, ps.Fuel)

							ntkm = ntkm + ps.Distance

							ntl = ntl + ps.Fuel

							ntt = ntt + ps.Totaltime

							gft = append(gft, ps.Timeless)

						}

						nt = nt + len(paths)
						ate = append(ate, float64(len(paths)))

						mtkm = mediavector(atkm)
						dtkm = DesvioPadrão(atkm, mtkm)

						gatkm = append(gatkm, mtkm)
						datkm = append(datkm, dtkm)

						mtl = mediavector(atl)
						dtl = DesvioPadrão(atl, mtl)
						gatl = append(gatl, mtl)
						datl = append(datl, dtl)

						atkm = []float64{}
						atl = []float64{}

						tpc = tpc + er.Compl
						gtpc = append(gtpc, er.Compl)

						tf = tf + er.Freq
						gfe = append(gfe, tf)

						mft = mediavector(gft)
						dft = DesvioPadrão(gft, mft)
						gftg = append(gftg, mft)
						gftd = append(gftd, dft)
						gft = []float64{}

						aconf = append(aconf, er.Confi)

					}

					mfe = mediavector(gfe)
					dfe = DesvioPadrão(gfe, mfe)

					mte = mediavector(ate)
					dte = DesvioPadrão(ate, mte)

					mtkm = mediavector(gatkm)
					dtkm = mediavector(datkm)

					mtl = mediavector(gatl)
					dtl = mediavector(datl)

					tcpm = mediavector(gtpc)
					dcpm = DesvioPadrão(gtpc, tcpm)

					mftuples := mediavector(gftg)
					dftuples := mediavector(gftd)

					mconf := mediavector(aconf)
					dconf := DesvioPadrão(aconf, mconf)

					info = append(
						info, user.Name, user.Typee, fmt.Sprintf("%f", user.Tank), fmt.Sprintf("%f", ntkm),
						fmt.Sprintf("%f", ntl), fmt.Sprintf("%f", ntkm/ntl), fmt.Sprintf("%f", lesperados), fmt.Sprintf("%f", ntt), fmt.Sprintf("%d", nev),
						fmt.Sprintf("%d", nt), fmt.Sprintf("%f", mte), fmt.Sprintf("%f", dte),
						fmt.Sprintf("%f", mtkm), fmt.Sprintf("%f", dtkm),
						fmt.Sprintf("%f", mtl), fmt.Sprintf("%f", dtl),
						fmt.Sprintf("%f", tpc), fmt.Sprintf("%f", tcpm), fmt.Sprintf("%f", dcpm),
						fmt.Sprintf("%f", tf), fmt.Sprintf("%f", mfe), fmt.Sprintf("%f", dfe),
						fmt.Sprintf("%f", mftuples), fmt.Sprintf("%f", dftuples),
						fmt.Sprintf("%f", user.Score), fmt.Sprintf("%f", mconf), fmt.Sprintf("%f", dconf),
						fmt.Sprintf("%f", user.Coin))

					writer.Write(info)
					writer.Flush()

				} else if tipo == user.Typee {
					events := []Event{}

					tu := GetAlleventsPast(contract, help)
					err = json.Unmarshal([]byte(tu), &events)
					if err != nil {
						fmt.Println(err.Error())

					}

					nev := len(events) //numero de eventos

					ntkm := 0.0
					ntl := 0.0
					ntt := 0.0

					nt := 0             //numero total de trajetos
					ate := []float64{}  //array armazena trajetos por evento
					atkm := []float64{} //array armazena trajetos por km
					atl := []float64{}  //array armazena trajetos por

					gatkm := []float64{} //global array armazena trajetos por km
					gatl := []float64{}  //global array armazena trajetos por litros

					datkm := []float64{} //global array armazena trajetos por km
					datl := []float64{}  //global array armazena trajetos por litros

					mte := 0.0  //media de trajetos
					dte := 0.0  //desvio de trajetos
					mtkm := 0.0 //media de trajetos
					dtkm := 0.0 //desvio de trajetos
					mtl := 0.0  //media de trajetos
					dtl := 0.0  //desvio de trajetos

					tpc := 0.0
					gtpc := []float64{}
					tcpm := 0.0
					dcpm := 0.0

					tf := 0.0
					gfe := []float64{}
					gft := []float64{}
					gftg := []float64{}
					gftd := []float64{}
					mfe := 0.0
					dfe := 0.0
					mft := 0.0
					dft := 0.0

					aconf := []float64{}

					lesperados := 0.0

					for _, er := range events {
						lesperados = lesperados + er.Dff

						trt := strings.Replace(er.Datai, " ", "-", -1)
						trt = strings.Replace(trt, ":", "-", -1)

						tu1 := GetPathhOpen(contract, help, trt)
						time.Sleep(500 * time.Millisecond)
						paths := []Path{}
						err = json.Unmarshal([]byte(tu1), &paths)
						if err != nil {
							fmt.Println(err.Error())

						}

						for _, ps := range paths {

							//nt = nt + ps.Ntuples

							//ate = append(ate, float64(ps.Ntuples))
							atkm = append(atkm, ps.Distance)
							atl = append(atl, ps.Fuel)

							ntkm = ntkm + ps.Distance

							ntl = ntl + ps.Fuel

							ntt = ntt + ps.Totaltime

							gft = append(gft, ps.Timeless)

						}

						nt = nt + len(paths)
						ate = append(ate, float64(len(paths)))

						mtkm = mediavector(atkm)
						dtkm = DesvioPadrão(atkm, mtkm)

						gatkm = append(gatkm, mtkm)
						datkm = append(datkm, dtkm)

						mtl = mediavector(atl)
						dtl = DesvioPadrão(atl, mtl)
						gatl = append(gatl, mtl)
						datl = append(datl, dtl)

						atkm = []float64{}
						atl = []float64{}

						tpc = tpc + er.Compl
						gtpc = append(gtpc, er.Compl)

						tf = tf + er.Freq
						gfe = append(gfe, tf)

						mft = mediavector(gft)
						dft = DesvioPadrão(gft, mft)
						gftg = append(gftg, mft)
						gftd = append(gftd, dft)
						gft = []float64{}

						aconf = append(aconf, er.Confi)

						/*


							tf := 0.0
							gfe := []float64{}
							gft := []float64{}
							gftg := []float64{}
							mfe := 0.0
							dfe := 0.0
							mft := 0.0
							dft := 0.0*
						*/

					}
					mfe = mediavector(gfe)
					dfe = DesvioPadrão(gfe, mfe)

					mte = mediavector(ate)
					dte = DesvioPadrão(ate, mte)
					//fmt.Println(ate, mte, dte)
					//ate = []float64{}

					mtkm = mediavector(gatkm)
					dtkm = mediavector(datkm)

					mtl = mediavector(gatl)
					dtl = mediavector(datl)

					tcpm = mediavector(gtpc)
					dcpm = DesvioPadrão(gtpc, tcpm)

					mftuples := mediavector(gftg)
					dftuples := mediavector(gftd)

					mconf := mediavector(aconf)
					dconf := DesvioPadrão(aconf, mconf)

					info = append(
						info, user.Name, user.Typee, fmt.Sprintf("%f", user.Tank), fmt.Sprintf("%f", ntkm),
						fmt.Sprintf("%f", ntl), fmt.Sprintf("%f", ntkm/ntl), fmt.Sprintf("%f", lesperados), fmt.Sprintf("%f", ntt), fmt.Sprintf("%d", nev),
						fmt.Sprintf("%d", nt), fmt.Sprintf("%f", mte), fmt.Sprintf("%f", dte),
						fmt.Sprintf("%f", mtkm), fmt.Sprintf("%f", dtkm),
						fmt.Sprintf("%f", mtl), fmt.Sprintf("%f", dtl),
						fmt.Sprintf("%f", tpc), fmt.Sprintf("%f", tcpm), fmt.Sprintf("%f", dcpm),
						fmt.Sprintf("%f", tf), fmt.Sprintf("%f", mfe), fmt.Sprintf("%f", dfe),
						fmt.Sprintf("%f", mftuples), fmt.Sprintf("%f", dftuples),
						fmt.Sprintf("%f", user.Score), fmt.Sprintf("%f", mconf), fmt.Sprintf("%f", dconf),
						fmt.Sprintf("%f", user.Coin))
					writer.Write(info)
					writer.Flush()

				}
				//info = nil
			}
		}
	}
}

func send_data(contract *client.Contract, name string, file []string, k float64) {

	var mn = [][]datavehicle{}

	trr := datavehicle{}

	//fmt.Println(file)

	//	for ty := range file {
	//		fmt.Println(file[ty])
	//	}
	//antfile := ""
	for _, filesss := range file {

		//datapack := []Tuple{}
		//result := 0
		trrA := []datavehicle{}

		file, err := os.Open("datavehicle/" + name + "/" + filesss)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)

		antcontrato := -1

		//antfuel := -1.0
		for scanner.Scan() {
			tt := scanner.Text()
			err := json.Unmarshal([]byte(tt), &trr)
			if err != nil {
				fmt.Println(err.Error())

			}

			//	fmt.Println(trr.Novoabastecimento)
			if trr.Contrato != 0 {
				trr.Novoabastecimento = "old"

				if antcontrato == -1 {
					antcontrato = trr.Contrato
					if len(trrA) > 0 {
						mn = append(mn, trrA)

						trrA = nil
					}

					//antfuel = trr.Combustivel
				}
				if len(trrA) > 0 {
					if math.Round(trrA[len(trrA)-1].Combustivel*100)/100+(math.Round(trrA[len(trrA)-1].Combustivel*100)/100*0.40) < math.Round(trr.Combustivel*100)/100 {

						trr.Novoabastecimento = "new"
						//fmt.Println(trr.Novoabastecimento)

						//fmt.Println((trrA[len(trrA)-1].Combustivel), trr.Combustivel, trr.Contrato, "Maior2")
						//antfuel = trr.Combustivel
						antcontrato = trr.Contrato

						//fmt.Println(antcontrato, trr.Contrato, "Depois")

						mn = append(mn, trrA)

						trrA = nil

					}
				}

				if int(trr.Contrato) != int(antcontrato) && len(trrA) > 0 {

					//fmt.Println(antcontrato, trr.Contrato, "Antes")

					antcontrato = trr.Contrato

					//fmt.Println(antcontrato, trr.Contrato, "Depois")

					mn = append(mn, trrA)

					trrA = nil
				}
				trr.Combustivel = math.Round(trr.Combustivel*10000) / 10000

				if trr.Combustivel > 100 {
					trr.Combustivel = 100
				}
				trrA = append(trrA, trr)
			}

		}

		/*
			if len(trrA) > 0 {
				fmt.Println(antcontrato, trr.Contrato, "Antes")

				antcontrato = trr.Contrato

				fmt.Println(antcontrato, trr.Contrato, "Depois")
				mn[op] = trrA
				op = op + 1

				trrA = nil
			}
		*/

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
		if len(trrA) > 0 {
			mn = append(mn, trrA)

			trrA = nil
		}

	}

	//fmt.Println(op, len(mn))
	//	b := 0.05
	a := rand.Float64()*((k+2)-(k-2)) + (k - 2)
	alpha := float64(1 - 1/math.Round(a*1)/1)

	c := rand.Float64()*((0.95)-(0.85)) + (0.85)
	//fmt.Println(u)

	fmt.Println(math.Round(a*1)/1, alpha)
	//meta := 0
	postoid := "9999"
	//sentcomb := -1.0
	antcontrato := -1
	//recupera := -1
	for i, ors := range mn {

		fmt.Println(len(mn[i]), len(ors), ors[0].Contrato, i)

		//fmt.Println(len(ors), ors[0].Combustivel, ors[len(ors)-1].Combustivel)

		sentinela := 0

		meta := 0

		if ors[0].Contrato != 0 && (len(ors)) > 100 {

			tuplestam := len(ors)

			//fmt.Println(tuplestam, "teste", len(ors), "tam")

			if tuplestam > 1000 {
				if tuplestam < 12000 {
					meta = rand.Intn(2000-1000) + 1000
					if meta > tuplestam {
						meta = tuplestam
						tuplestam = 0
					}
				} else {
					meta = rand.Intn(4000-3000) + 3000
				}

			} else {
				meta = tuplestam
				tuplestam = tuplestam - meta
			}
			//fmt.Println(tuplestam, "tuplas", meta, "meta", sentinela, "sentinela")

			tuples := []Tuple{}

			if antcontrato == -1 {

				fmt.Println(ors[0].Novoabastecimento)

				insertuser(contract, strconv.Itoa(ors[0].Id), strconv.Itoa(ors[0].TamanhoTanque), ors[0].Tipo)
				antcontrato = GetStatusEvent(contract, name, postoid)

				if antcontrato != 0 {

					CloseEvent(contract, name, "0.0")
				}

				Createevent(contract, math.Round(ors[0].Combustivel*100/float64(ors[0].TamanhoTanque)*1000)/1000, float64(ors[0].Contrato), name, postoid)

			} else if ors[0].Novoabastecimento == "new" {

				fmt.Println(ors[0].Novoabastecimento)

				insertuser(contract, strconv.Itoa(ors[0].Id), strconv.Itoa(ors[0].TamanhoTanque), ors[0].Tipo)

				recupera := GetStatusEvent(contract, name, postoid)

				if recupera != 0 {

					CloseEvent(contract, name, "0.0")
				}

				Createevent(contract, math.Round(ors[0].Combustivel*100/float64(ors[0].TamanhoTanque)*1000)/1000, float64(ors[0].Contrato), name, postoid)

			}

			for j := range ors {

				if len(tuples) <= meta {

					//		r := rand.Float64()

					//if r >= alpha {
					//	fmt.Println(ors[j].Combustivel, ors[j].TamanhoTanque)
					comb := ors[j].Combustivel * 100 / float64(ors[j].TamanhoTanque)
					//fmt.Println(comb)
					lat := fmt.Sprintf("%f", ors[j].Pos.Lat)
					long := fmt.Sprintf("%f", ors[j].Pos.Long)
					tuples = append(tuples, Tuple{T: ors[j].Time, Pos: lat + "/" + long, Comb: math.Round(comb*1000) / 1000})
					//	} else {
					//sentinela = sentinela + 1
					//	}
				}

				if (len(tuples) == meta || j == len(ors)-1) && len(tuples) > 100 {

					fmt.Println(tuplestam, "tuplas", meta, "meta", sentinela, "sentinela", ors[j].Contrato, "len", len(tuples))

					if GetStatusEvent(contract, name, postoid) == ors[0].Contrato {
						tuples = ruido(tuples)

						p := rand.Float64()

						if p > c {
							fmt.Println("Não envia")
						} else {
							CreatePath(contract, tuples, name, postoid, k)
						}
						time.Sleep(800 * time.Millisecond)

					}

					tuples = nil

					if tuplestam > meta+sentinela {
						tuplestam = tuplestam - (meta + sentinela)

						if tuplestam < 12000 {
							meta = rand.Intn(2000-500) + 1000
							if meta > tuplestam {
								meta = tuplestam
								tuplestam = 0
							}
						} else if tuplestam > 12000 {
							meta = rand.Intn(4000-3000) + 3000
						}
					} else if tuplestam <= 2000 && tuplestam > 0 {
						meta = tuplestam
						tuplestam = 0

					} else if tuplestam <= 0 {
						meta = 0

					}

					sentinela = 0

				}

			}

		}

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

	for _, dire := range dir {
		files, err = os.ReadDir("datavehicle/" + dire)

		for _, file := range files {

			datafile = append(datafile, file.Name())
		}
		// || dire == "20" || dire == "41" || dire == "61" || dire == "81" || dire == "111" || dire == "120" dire == "1" || dire == "20" || dire == "41" || dire == "61" ||  || dire == "111" || dire == "120" || dire == "82" || dire == "83" || dire == "84" || dire == "85" || dire == "86" || dire == "87" || dire == "88" || dire == "89" || dire == "90" || dire == "91"

		m[dire] = datafile

		var wg sync.WaitGroup
		if len(m) == 15 { //15
			for i, ui := range m {
				wg.Add(1)
				go func(i string, ui []string) {
					defer wg.Done()
					send_data(contract, i, ui, k)
				}(i, ui)
				time.Sleep(1 * time.Second)

			}
			wg.Wait()
			for k := range m {
				delete(m, k)
			}

		}

		//}

		datafile = nil

	}

}

func insertuser(contract *client.Contract, name, tanque, tipo string) {

	if !checkuserexist(contract, name) {
		Createuser(contract, name, name, tanque, tipo)
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
