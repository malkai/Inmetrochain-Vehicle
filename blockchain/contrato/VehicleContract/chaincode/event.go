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
	layout := "2006-01-02 15:04:05"
	aux, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("\n Erro ao criar evento. %v", err)
	}
	exist, err := s.Eventexist(ctx, Iduser1, Iduser2)
	if err != nil {
		return fmt.Errorf("\n Erro ao verificar a existencia do evento. %v", err)
	}
	if exist {
		//return fmt.Errorf("\n Erro ao criar evento 2. %v", err)
		return s.Closeevent(ctx, Iduser1, Iduser2, 0.0)
	}

	var fsump, _ = strconv.ParseFloat(Fsupi, 64)
	var dffr, _ = strconv.ParseFloat(Dff, 64)

	event := Event{
		Id:        Iduser1 + Iduser2 + aux.AsTime().Format(layout),
		Datai:     aux.AsTime().Format(layout),
		Dataiataf: "",
		Fsupi:     fsump,
		Fsupf:     0,
		Dff:       dffr,
		Vstatus:   false,
		Iduser1:   Iduser1,
		Iduser2:   Iduser2,
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("event", eventJSON)
}

func (s *SmartContract) Eventexist(ctx contractapi.TransactionContextInterface, id string, id2 string) (bool, error) {

	aux, err := GetEventOpen(ctx, id, id2)
	if err != nil {
		return false, fmt.Errorf("Erro ao retornar evento em aberto: %v", err)
	}

	return aux != nil, nil

}

func (s *SmartContract) Closeevent(ctx contractapi.TransactionContextInterface, id string, id2 string, sum float64) error {

	event, err := GetEventOpen(ctx, id, id2)
	if err != nil {
		return fmt.Errorf("Erro ao recuperar evento: %v", err)
	}
	layout := "2006-01-02 15:04:05"

	var eventjson = event[0]

	var valuevalids []string
	var fuelsum float64 = 0.0
	aux, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("erro ao pegar data: %v", err)
	}

	eventjson.Dataiataf = aux.AsTime().Format(layout)
	eventjson.Vstatus = true
	temp, err := s.GetallPath(ctx, "Path"+id)
	if err != nil {
		return err
	}

	/*

	   {
	     "id": "event1",
	     "datai": "2023-11-14 18:16:27",
	     "dataf": "",
	     "fsupi": 93,
	     "fsupf": 3.3777777777777827,
	     "fsupfd": 0,
	     "dff": 5,
	     "vstatus": false,
	     "iduser1": "1",
	     "iduser2": "3"
	   },
	*/

	var ntotal float64 = 0.0
	var ntimeless float64 = 0.0
	if len(temp) > 0 {
		for _, xx := range temp {
			if xx.DataR != "" {
				datet1, err := time.Parse(layout, xx.DataR)
				if err != nil {
					return fmt.Errorf("erro1: %v", err)
				}

				datet2, err := time.Parse(layout, eventjson.Datai)
				if err != nil {
					return fmt.Errorf("erro2: %v", err)
				}

				datet3, err := time.Parse(layout, eventjson.Dataiataf)
				if err != nil {
					return fmt.Errorf("erro3: %v", err)
				}

				if datet1.After(datet2) && datet1.Before(datet3) {
					ntotal = ntotal + 1
					fuelsum = fuelsum + xx.Fuel
					ntimeless = ntimeless + xx.Timeless

				}
			}

		}
	}
	//return fmt.Errorf("\n Erro ao criar evento 2. %v %s", err, temp[0].DataR)

	s.Updatuser(ctx, "user"+id, ntimeless/ntotal, eventjson.Dff, fuelsum, valuevalids)
	eventjson.Fsupi = eventjson.Fsupi + sum
	eventjson.Fsupfd = eventjson.Fsupi - fuelsum
	eventJSON, err := json.Marshal(eventjson)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("event", eventJSON)

}

func (s *SmartContract) updatevent(ctx contractapi.TransactionContextInterface, id string, id2 string, minus float64) error {

	event, err := GetEventOpen(ctx, id, id2)

	if err != nil {
		return fmt.Errorf("erro ao atualizar evento: %v", err)
	}

	var eventjson = event[0]

	eventjson.Fsupf = eventjson.Fsupf + minus
	//return fmt.Errorf("\n Sucesso  %v", err)

	if eventjson.Dff < eventjson.Fsupf {
		//return fmt.Errorf("erro verificar %v", err)
		//return fmt.Errorf("\n Sucesso  %v", err)
		return s.Closeevent(ctx, id, id2, minus)
	}

	eventJS, err := json.Marshal(eventjson)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, eventJS)

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

func GetEventOpen(ctx contractapi.TransactionContextInterface, id string, id2 string) ([]*Event, error) {

	queryString := fmt.Sprintf(`{"selector":{"docType":"event","Vstatus":"true","Iduser1":"%s","Iduser2":"%s" }}`, id, id2)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Event
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset Event
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil

}
