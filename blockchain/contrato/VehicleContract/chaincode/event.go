package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/*

type Event struct {
	Id        string  `json:"id"`     //Iduser1+Iduser2+ttimestamp
	Datai     string  `json:"datai"`  //data inicial
	Dataiataf string  `json:"dataf"`  //data final do acordo
	Fsupi     float64 `json:"fsupi"`  //combustivel inicial
	Fsupf     float64 `json:"fsupf"`  //combustivel final
	Fsupfd    float64 `json:"fsupfd"` //i+constant k times
	Dff       float64 `json:"dff"`
	Vstatus   bool    `json:"vstatus"` //sentinela de controle
	Iduser1   string  `json:"iduser1"` //identificação do usuario 1
	Iduser2   string  `json:"iduser2"` //identificação do usuario 2
}
*/

// Cria um evento no blockchain
func (s *SmartContract) Createevent(ctx contractapi.TransactionContextInterface, Fsupi float64, Dff float64, Iduser1 string, Iduser2 string) error {

	exist, err := s.Eventexist(ctx, "event"+Iduser1)
	if err != nil {
		return fmt.Errorf("\n Erro ao criar evento. %v", err)
	}
	if exist {
		s.Closeevent(ctx, "event"+Iduser1)
	}
	aux, err := ctx.GetStub().GetTxTimestamp()

	event := Event{
		Id:        "event" + Iduser1,
		Datai:     aux.AsTime().UTC().GoString(),
		Dataiataf: "",
		Fsupi:     Fsupi,
		Fsupf:     0,
		Dff:       Dff,
		Vstatus:   false,
		Iduser1:   Iduser1,
		Iduser2:   Iduser2,
	}
	assetJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("event"+Iduser1, assetJSON)
}

func (s *SmartContract) Eventexist(ctx contractapi.TransactionContextInterface, id string) (bool, error) {

	event, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return event != nil, nil

}

func (s *SmartContract) Closeevent(ctx contractapi.TransactionContextInterface, id string) error {

	event, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	var asset Event
	err = json.Unmarshal(event, &asset)
	if err != nil {
		return err
	}

	aux, err := ctx.GetStub().GetTxTimestamp()
	asset.Dataiataf = aux.AsTime().UTC().GoString()
	asset.Vstatus = true

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(id+asset.Dataiataf, assetJSON)
	return nil

}

func (s *SmartContract) updatevent(ctx contractapi.TransactionContextInterface, id string, minus float64) error {

	event, err := ctx.GetStub().GetState(id)

	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	var asset Event
	err = json.Unmarshal(event, &asset)
	if err != nil {
		return err
	}

	asset.Fsupf = +minus

	if asset.Fsupf >= asset.Dff {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}
		ctx.GetStub().PutState(id, assetJSON)
		s.Closeevent(ctx, id)
	}
	if asset.Fsupf < asset.Dff {

		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}
		ctx.GetStub().PutState(id, assetJSON)
	}
	return nil

}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllevents(ctx contractapi.TransactionContextInterface, id string) ([]*Event, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange(id, "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var events []*Event
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var event Event
		err = json.Unmarshal(queryResponse.Value, &event)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}
