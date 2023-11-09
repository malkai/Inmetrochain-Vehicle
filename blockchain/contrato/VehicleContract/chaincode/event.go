package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Cria um evento no blockchain
func (s *SmartContract) Createevent(ctx contractapi.TransactionContextInterface, Fsupi string, Dff string, Iduser1 string, Iduser2 string) error {

	exist, err := s.Eventexist(ctx, "event"+Iduser1)
	if err != nil {
		return fmt.Errorf("\n Erro ao criar evento. %v", err)
	}
	if exist {
		s.Closeevent(ctx, "event"+Iduser1)
	}
	aux, err := ctx.GetStub().GetTxTimestamp()

	var fsump, _ = strconv.ParseFloat(Fsupi, 64)
	var dffr, _ = strconv.ParseFloat(Dff, 64)

	event := Event{
		Id:        "event" + Iduser1,
		Datai:     aux.AsTime().UTC().GoString(),
		Dataiataf: "",
		Fsupi:     fsump,
		Fsupf:     0,
		Dff:       dffr,
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
		return fmt.Errorf("Erro ao recuperar evento: %v", err)
	}

	var asset Event
	err = json.Unmarshal(event, &asset)
	if err != nil {
		return err
	}
	var valuevalids []string
	var fuelsum float64 = 0.0
	aux, err := ctx.GetStub().GetTxTimestamp()
	asset.Dataiataf = aux.AsTime().UTC().GoString()
	asset.Vstatus = true
	temp, err := s.GetallPath(ctx, "Path"+id)
	var ntotal float64 = 0.0
	var ntimeless float64 = 0.0

	for _, xx := range temp {
		layout := "2006-01-02 15:04:05.000000"
		datet1, _ := time.Parse(layout, xx.DataR)

		datet2, _ := time.Parse(layout, asset.Datai)

		datet3, _ := time.Parse(layout, asset.Dataiataf)

		if datet1.After(datet2) && datet1.Before(datet3) {
			ntotal = ntotal + 1
			fuelsum = fuelsum + xx.Fuel
			ntimeless = ntimeless + xx.Timeless

		}

	}
	var users, _ = s.Userget(ctx, "user"+id)
	score := Credibility(users.Score, ntimeless/ntotal, asset.Dff, fuelsum, valuevalids)
	s.Updatuser(ctx, "user"+id, 1*score, score)
	asset.Fsupf = asset.Fsupi - fuelsum
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("event"+id+asset.Dataiataf, assetJSON)

}

func (s *SmartContract) updatevent(ctx contractapi.TransactionContextInterface, id string, minus float64) error {

	event, err := ctx.GetStub().GetState(id)

	if err != nil {
		return fmt.Errorf("erro ao atualizar evento: %v", err)
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
		return s.Closeevent(ctx, id)
	}
	if asset.Fsupf < asset.Dff {

		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}
		return ctx.GetStub().PutState(id, assetJSON)
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
			return nil, fmt.Errorf("event erro1: %v", err)
		}

		var event Event
		err = json.Unmarshal(queryResponse.Value, &event)
		if err != nil {
			return nil, fmt.Errorf("event erro2: %v", err)
		}
		if event.Id != "" {
			events = append(events, &event)
		}

	}

	return events, nil
}
