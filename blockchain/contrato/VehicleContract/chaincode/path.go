package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/sgreben/piecewiselinear"
)

/*
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:7051
export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt")
*/

// peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n vehicle -c '{"function":"CreatePath","Args":[[91 123 34 84 34 58 34 49 48 34 44 34 80 111 115 34 58 34 49 47 50 34 44 34 67 111 109 98 34 58 57 51 125 44 123 34 84 34 58 34 49 49 34 44 34 80 111 115 34 58 34 49 47 50 34 44 34 67 111 109 98 34 58 57 50 125 44 123 34 84 34 58 34 49 50 34 44 34 80 111 115 34 58 34 49 47 50 34 44 34 67 111 109 98 34 58 57 49 125 93][34 49 34]]}'
// Cria um evento no blockchain

func (s *SmartContract) CreatPath(ctx contractapi.TransactionContextInterface, data string, id string, id2 string, k string) error {

	kfloat, _ := strconv.ParseFloat(strings.TrimSpace(k), 64)

	tuples := []Tuple{}
	stringByte := []byte(data)
	stringByte, err := Decompress(stringByte)
	if err != nil {
		return fmt.Errorf("\n Erro Path converter. %v", err)
	}
	tuples, err = DecodeToTuple(stringByte)
	if err != nil {
		return fmt.Errorf("\n Erro Path Decodar. %v", err)
	}

	h, err := s.Eventexist(ctx, id, id2)
	if err != nil {
		return fmt.Errorf("\n Erro checar evento. %v", err)
	}
	if h {
		return fmt.Errorf("\n Caminho não esta conectado a nenhum evento %v", err)
	}

	dist := 0.0
	fuel := 0.0
	var fuel_vector []float64
	time := 0.0
	var time2 []string
	time3 := []float64{}
	i := 0

	for i = range tuples {
		if i < len(tuples)-1 {

			dist1, err := Distanceeucle(tuples[i].Pos, tuples[i+1].Pos)
			if err != nil {
				return fmt.Errorf("\n Erro checar Distancia. %v %s %s ", err, tuples[i].Pos, tuples[i+1].Pos)
			}

			dist = dist + dist1

			time1, err := totaltime(tuples[i].T, tuples[i+1].T)
			if err != nil {
				return fmt.Errorf("\n Error checa Tempo. %v", err)
			}
			time = time + time1

			//rtt, _ := strconv.ParseFloat(tuples[i].T, 64)
			time3 = append(time3, time)

			time2 = append(time2, tuples[i].T)

			fuel_vector = append(fuel_vector, tuples[i].Comb)
		}

	}

	timeles, err := Timeliness(time2, kfloat)
	if err != nil {
		return fmt.Errorf("\n Error na metrica timeless. %v", err)
	}

	user, err := s.Userget(ctx, "user"+id)
	if err != nil {
		return fmt.Errorf("\n Error ao recuperar user %v %s", err, "user"+id)
	}

	/*
		fuel, err = KalmanFilter(user.Tanque, fuel_vector)
		if err != nil {
			return fmt.Errorf("\n Error ao aplicar filtro de kalman %v", err)
		}*/

	f := piecewiselinear.Function{Y: fuel_vector} // range: "hat" function
	f.X = time3
	rtt := [][2]float64{}

	for rty := range time3 {
		rtt1 := [2]float64{time3[rty], fuel_vector[rty]}
		rtt = append(rtt, rtt1)
	}
	_, _, R2, err := Linear(rtt)
	if err != nil {
		panic(err)
	}
	//fmt.Fprintf(os.Stdout, "y   = %.4f*x+%.4f\n", a, b)
	//fmt.Fprintf(os.Stdout, "R^2 = %.4f\n", R2)
	fuel = R2 * user.Tanque / 100

	aux, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("\n Erro checar data. %v", err)
	}

	layout := "2006-01-02 15:04:05"

	eventjson, err := s.GetEventOpenSingle(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao recuperar evento: %v", err)

	}

	trt := strings.Replace(eventjson.Datai, " ", "-", -1)
	trt = strings.Replace(trt, ":", "-", -1)

	path := Path{
		DocType:     "path",
		PathID:      "path" + id + aux.AsTime().Format(layout),
		DataEvent:   trt,
		DataVehicle: data,
		Distance:    dist,
		Fuel:        fuel,
		Totaltime:   time,
		Timeless:    timeles,
		DataR:       aux.AsTime().Format(layout),
		Iduser:      id,
		K:           kfloat,
	}

	//return fmt.Errorf("\n Sucesso  %v %f", err, fuel)

	patJSON, err := json.Marshal(path)
	if err != nil {
		return fmt.Errorf("\n Erro ao comparctar data. %v %+v", err, path.Iduser)
	}
	st := strconv.FormatFloat(fuel, 'E', -1, 64)
	err = s.updatevent(ctx, id, id2, st)
	if err != nil {
		return fmt.Errorf("\n Erro ao atualizar evento em path. %v %f", err, fuel)
	}

	return ctx.GetStub().PutState(path.PathID, patJSON)

}

func (s *SmartContract) ExistPath(ctx contractapi.TransactionContextInterface, id string) (bool, error) {

	event, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return event != nil, nil

}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetallPath(ctx contractapi.TransactionContextInterface, id string) ([]*Path, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange(id, "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var paths []*Path
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Path erro1: %v", err)
		}

		var path Path
		err = json.Unmarshal(queryResponse.Value, &path)
		if err != nil {
			return nil, fmt.Errorf("Path erro2: %v", err)
		}

		if path.PathID != " " {
			paths = append(paths, &path)
		}

	}
	return paths, nil
}

func (s *SmartContract) GetPathhOpen(ctx contractapi.TransactionContextInterface, id string, time string) ([]*Path, error) {

	//queryString := fmt.Sprintf(`{"selector":{"docType":"event","vstatus":"false","iduser1":"%s","iduser2":"%s" }}`, id, id2)
	queryString := fmt.Sprintf(`{"selector":{"docType":"path","dataEvent":"%s","iduser":"%s"  }}`, time, id)
	//arrayteststring := [3]string["true", ]

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var pathbjs []*Path

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var pathobj Path
		err = json.Unmarshal(queryResult.Value, &pathobj)
		if err != nil {
			return nil, err
		}
		pathbjs = append(pathbjs, &pathobj)
	}
	return pathbjs, nil

}

func (s *SmartContract) GetPathhAll(ctx contractapi.TransactionContextInterface, id string, time string) ([]*Path, error) {

	//queryString := fmt.Sprintf(`{"selector":{"docType":"event","vstatus":"false","iduser1":"%s","iduser2":"%s" }}`, id, id2)
	queryString := fmt.Sprintf(`{"selector":{"docType":"path","iduser":"%s"  }}`, id)
	//arrayteststring := [3]string["true", ]

	resultsIterator, _, err := ctx.GetStub().GetQueryResultWithPagination(queryString, 1000, "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var pathbjs []*Path

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var pathobj Path
		err = json.Unmarshal(queryResult.Value, &pathobj)
		if err != nil {
			return nil, err
		}
		pathbjs = append(pathbjs, &pathobj)
	}
	return pathbjs, nil

}
