package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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

func (s *SmartContract) CreatPath(ctx contractapi.TransactionContextInterface, data string, id string, id2 string, k string) error {

	kfloat, _ := strconv.ParseFloat(strings.TrimSpace(k), 64)

	tuples, err := Decompress(data)
	if err != nil {
		return fmt.Errorf("\n Erro ao ler tuplas. %v", err)

	}

	if len(tuples) > 0 {

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
		var time2 []float64
		time3 := []float64{}

		for i := range tuples {
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

				time2 = append(time2, time1)

				fuel_vector = append(fuel_vector, tuples[i].Comb)
			}

		}

		tamTuples := len(tuples)

		timeles, err := Timeliness(time2, kfloat)
		if err != nil {
			return fmt.Errorf("\n Error na metrica timeless. %v", err)
		}

		user, err := s.Userget(ctx, "user"+id)
		if err != nil {
			return fmt.Errorf("\n Error ao recuperar user %v %s", err, "user"+id)
		}

		rtt := [][2]float64{}

		for rty := range time3 {
			rtt1 := [2]float64{time3[rty], fuel_vector[rty]}
			rtt = append(rtt, rtt1)
		}
		a, b, _, err := Linear(rtt)
		if err != nil {
			return fmt.Errorf("\n Error na regressão %v", err)
		}

		tui := []float64{}
		for i := range time3 {

			tui = append(tui, time3[i]*a+b)
		}

		fuel = ((tui[0] - tui[len(tui)-1]) * user.Tank) / 100

		aux, err := ctx.GetStub().GetTxTimestamp()
		if err != nil {
			return fmt.Errorf("\n Erro checar data. %v", err)
		}

		layout := "2006-01-02 15:04:05"

		eventjson1, err := s.GeteventOpensingle(ctx, id, id2)
		if err != nil {
			return fmt.Errorf("erro ao recuperar evento: %v", err)

		}
		eventjson := eventjson1[0]

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
			Ntuples:     tamTuples,
		}

		//return fmt.Errorf("\n Sucesso  %v %f", err, fuel)

		patJSON, err := json.Marshal(path)
		if err != nil {
			return fmt.Errorf("\n Erro ao comparctar data. %v %+v", err, path.Iduser)
		}
		st := strconv.FormatFloat(fuel, 'E', -1, 64)

		err = ctx.GetStub().PutState(path.PathID, patJSON)
		if err != nil {
			return fmt.Errorf("erro ao salvar caminho: %v", err)
		}

		//  Create an index to enable color-based range queries, e.g. return all blue assets.
		//  An 'index' is a normal key-value entry in the ledger.
		//  The key is a composite key, with the elements that you want to range query on listed first.
		//  In our case, the composite key is based on indexName~color~name.
		//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
		colorNameIndexKey, err := ctx.GetStub().CreateCompositeKey("path", []string{path.DataEvent, path.Iduser, path.PathID})
		if err != nil {
			return err
		}
		//  Save index entry to world state. Only the key name is needed, no need to store a duplicate copy of the asset.
		//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
		value := []byte{0x00}
		err = ctx.GetStub().PutState(colorNameIndexKey, value)
		if err != nil {
			return fmt.Errorf("\n Erro ao salvar index. %v %f", err, fuel)
		}

		err = s.updatevent(ctx, id, id2, st, path)
		if err != nil {
			return fmt.Errorf("\n Erro ao atualizar evento em path. %v %f", err, fuel)
		}

		tuples = nil
	}
	return nil

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
	queryString := fmt.Sprintf(`{"selector":{"iduser":"%s","dataEvent":"%s","docType":"path"}}`, id, time)
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

func (t *SmartContract) ReadPath(ctx contractapi.TransactionContextInterface, pathID string) (*Path, error) {
	PathBytes, err := ctx.GetStub().GetState(pathID)
	if err != nil {
		return nil, fmt.Errorf("failed to get path %s: %v", pathID, err)
	}
	if PathBytes == nil {
		return nil, fmt.Errorf("path %s does not exist", pathID)
	}

	var path Path
	err = json.Unmarshal(PathBytes, &path)
	if err != nil {
		return nil, fmt.Errorf("erro Unmarshal %s", pathID)
	}

	return &path, nil
}

func (s *SmartContract) GetPathhIndex(ctx contractapi.TransactionContextInterface, id string, time string) ([]*Path, error) {
	// Execute a key range query on all keys starting with 'color'
	coloredAssetResultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("path", []string{time, id})
	if err != nil {
		return nil, fmt.Errorf("erro ao recuperar compostion Key")
	}
	defer coloredAssetResultsIterator.Close()
	var paths []*Path

	for coloredAssetResultsIterator.HasNext() {
		responseRange, err := coloredAssetResultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("erro ao interar compostion Key")
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, fmt.Errorf("erro ao assimilar compostion Key")
		}

		if len(compositeKeyParts) > 1 {

			returnedAssetID := compositeKeyParts[2]

			path, err := s.ReadPath(ctx, returnedAssetID)
			if err != nil {
				return nil, fmt.Errorf("erro ao recuperar path %v %s", err, compositeKeyParts)
			}
			paths = append(paths, path)

		}
	}

	return paths, nil
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
