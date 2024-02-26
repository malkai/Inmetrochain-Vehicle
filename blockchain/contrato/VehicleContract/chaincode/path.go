package chaincode

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
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

func (s *SmartContract) CreatPath(ctx contractapi.TransactionContextInterface, data string, id string, id2 string) error {

	var tuples []Tuple

	stringByte := []byte(data)

	err := json.Unmarshal(stringByte, &tuples)
	if err != nil {
		return fmt.Errorf("\n Erro nas tuplas. %v", err)
	}

	h, err := s.Eventexist(ctx, id, id2)
	if err != nil {
		return fmt.Errorf("\n Erro checar evento. %v", err)
	}
	if h == true {
		return fmt.Errorf("\n Caminho não esta conectado a nenhum evento %v", err)
	}

	dist := 0.0
	fuel := 0.0
	var fuel_vector []float64
	time := 0.0
	var time2 []string
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

			time2 = append(time2, tuples[i].T)

			fuel_vector = append(fuel_vector, tuples[i].Comb)
		}

	}

	timeles, err := Timeliness(time2, 1)
	if err != nil {
		return fmt.Errorf("\n Error na metrica timeless. %v", err)
	}
	fmt.Println(timeles)

	user, err := s.Userget(ctx, "user"+id)
	if err != nil {
		return fmt.Errorf("\n Error ao recuperar user %v %s", err, "user"+id)
	}
	/*

		fuel, err = KalmanFilter(user.Tanque, fuel_vector)
		if err != nil {
			return fmt.Errorf("\n Error ao aplicar filtro de kalman %v", err)
		}
		fmt.Println(fuel)
	*/
	fuel = (fuel_vector[0] - fuel_vector[len(fuel_vector)-1]) * user.Tanque / 100

	aux, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("\n Erro checar data. %v", err)
	}
	fmt.Println(aux)
	//layout := "2006-01-02 15:04:05"
	if math.IsNaN(fuel) {
		return fmt.Errorf("\n Checando: %f", fuel)
	}

	layout := "2006-01-02 15:04:05"

	path := Path{
		DocType:     "path",
		PathID:      id + aux.AsTime().Format(layout),
		DataVehicle: tuples,
		Distance:    dist,
		Fuel:        fuel,
		Totaltime:   time,
		Timeless:    timeles,
		DataR:       aux.AsTime().Format(layout),
		Iduser:      id,
	}

	//return fmt.Errorf("\n Sucesso  %v %f", err, fuel)

	patJSON, err := json.Marshal(path)
	if err != nil {
		return fmt.Errorf("\n Erro ao comparctar data. %v %+v", err, path)
	}
	s.updatevent(ctx, id, id2, fuel)

	return ctx.GetStub().PutState("path"+path.PathID, patJSON)

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

func (s *SmartContract) GetPathhOpen(ctx contractapi.TransactionContextInterface, id string) ([]*Path, error) {

	//queryString := fmt.Sprintf(`{"selector":{"docType":"event","vstatus":"false","iduser1":"%s","iduser2":"%s" }}`, id, id2)
	queryString := fmt.Sprintf(`{"selector":{"docType":"path","Iduser":"%s"}}`, id)
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
