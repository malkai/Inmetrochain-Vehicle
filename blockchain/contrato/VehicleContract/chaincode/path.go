package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) Creatteste(ctx contractapi.TransactionContextInterface) error {
	fmt.Printf("oi")
	return nil

}

// Cria um evento no blockchain
func (s *SmartContract) CreatPath(ctx contractapi.TransactionContextInterface, tuples [][]float64, id string) error {

	exist, err := s.Eventexist(ctx, "Event"+id)
	if err != nil {
		return fmt.Errorf("\n Erro checar evento. %v", err)
	}
	if exist {
		dist := 0.0
		fuel := 0.0
		time := 0.0
       	//tt, err := ctx.GetStub().GetTxTimestamp()

			for i, aux := range tuples {
				if(len(tuples)>i){
				dist = +Distanceeucle(aux[2], aux[i], aux[i+1], aux[i+1])
				time = +totaltime(aux[i], a[i+1].T)
				}
				s.updatevent(ctx, "Event"+id, aux.Comb)
				fmt.Println(i, aux)
			}

			/*

				path := Path{
					DataVehicle: a,
					EventID:     "Path" + id + tt.AsTime().UTC().GoString(),
					Distance:    dist,
					Fuel:        fuel,
					Totaltime:   time,
				}

				s.updatevent(ctx, "Event"+id, path.Fuel)
				assetJSON, err := json.Marshal(path)
				if err != nil {
					return err
				}
				return ctx.GetStub().PutState("Path"+id, assetJSON)
			*/

		}

	}

	return fmt.Errorf("\n Caminho não esta conectado a nenhum evento %v", err)

}

func (s *SmartContract) ExistPath(ctx contractapi.TransactionContextInterface, id string) (bool, error) {

	event, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return event != nil, nil

}

func containsDuplicate(nums []int) bool {
	allKeys := make(map[int]bool)
	for _, number := range nums {
		if _, value := allKeys[number]; !value {
			return false
		}
	}
	return true

}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetallPath(ctx contractapi.TransactionContextInterface, id string) ([]*Path, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange(id, "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Path
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Path
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
