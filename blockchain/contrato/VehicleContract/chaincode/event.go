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
	if exist == false {
		//return fmt.Errorf("\n Erro ao criar evento 2. %v", err)
		err = s.Closeevent(ctx, Iduser1, Iduser2, 0.0)
	}

	var fsump, _ = strconv.ParseFloat(Fsupi, 64)
	var dffr, _ = strconv.ParseFloat(Dff, 64)

	event := Event{
		DocType:   "event",
		Id:        "Event" + Iduser1 + Iduser2,
		Datai:     aux.AsTime().Format(layout),
		Dataiataf: "",
		Fsupi:     fsump,
		Fsupf:     0,
		Dff:       dffr,
		Vstatus:   true,
		Iduser1:   Iduser1,
		Iduser2:   Iduser2,
		Compl:     0.0,
		Freq:      0.0,
		Confi:     0.0,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(event.Id, eventJSON)
}

func (s *SmartContract) Eventexist(ctx contractapi.TransactionContextInterface, id string, id2 string) (bool, error) {

	aux, err := ctx.GetStub().GetState(id + id2)
	if err != nil {
		return false, fmt.Errorf("Erro ao retornar evento em aberto: %v", err)
	}

	return aux != nil, nil

}

func (s *SmartContract) Closeevent(ctx contractapi.TransactionContextInterface, id string, id2 string, sum float64) error {

	eventjson, err := s.GetEventOpenSingle(ctx, id, id2)
	if err != nil {
		return fmt.Errorf("Erro ao recuperar evento: %v", err)
	}

	layout := "2006-01-02 15:04:05"

	var valuevalids []string
	var fuelsum float64 = 0.0
	aux, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("erro ao pegar data: %v", err)
	}

	eventjson.Dataiataf = aux.AsTime().Format(layout)
	eventjson.Vstatus = false
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
	conf, err := s.Updatuser(ctx, "user"+id, ntimeless/ntotal, eventjson.Dff/fuelsum, valuevalids)
	if err != nil {
		return fmt.Errorf("erro ao atualizar user: %v", err)
	}
	eventjson.Id = eventjson.Id + eventjson.Dataiataf
	eventjson.Fsupi = eventjson.Fsupi + sum
	eventjson.Fsupfd = eventjson.Fsupi - fuelsum
	eventjson.Compl = eventjson.Dff / fuelsum
	eventjson.Freq = ntimeless / ntotal
	eventjson.Confi = conf
	eventJSON, err := json.Marshal(eventjson)
	if err != nil {
		return fmt.Errorf("erro ao atualizar evento: %v", err)
	}

	err = ctx.GetStub().PutState(eventjson.Id+eventjson.Dataiataf, eventJSON)
	if err != nil {
		return fmt.Errorf("erro ao inserir evento: %v", err)
	}

	events := Event{
		DocType:   "event",
		Id:        "Event" + eventjson.Iduser1 + eventjson.Iduser2,
		Datai:     "",
		Dataiataf: "",
		Fsupi:     eventjson.Fsupi,
		Fsupf:     0,
		Dff:       0.0,
		Vstatus:   false,
		Iduser1:   eventjson.Iduser1,
		Iduser2:   eventjson.Iduser2,
		Compl:     0.0,
		Freq:      0.0,
		Confi:     0.0,
	}

	eventJSON, err = json.Marshal(events)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(events.Id, eventJSON)
}

func (s *SmartContract) updatevent(ctx contractapi.TransactionContextInterface, id string, id2 string, minus float64) error {

	var eventjson, err = s.GetEventOpenSingle(ctx, "event"+id, id2)
	if err != nil {
		return fmt.Errorf("erro ao atualizar evento: %v", err)
	}
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

	return ctx.GetStub().PutState(eventjson.Id, eventJS)

}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllevents(ctx contractapi.TransactionContextInterface, id string) ([]*Event, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("E", "f")
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

func (s *SmartContract) GetEventOpenSingle(ctx contractapi.TransactionContextInterface, id string, id2 string) (Event, error) {
	var events Event
	event, err := ctx.GetStub().GetState(id + id2)
	if err != nil {
		return events, fmt.Errorf("Erro em acessar a informação na blockchain: %v", err)
	}
	err = json.Unmarshal(event, &events)
	if err != nil {
		return events, fmt.Errorf("Falha na leitura do evento singular : %v %s", err, events.Id)
	}

	return events, nil
}

func (s *SmartContract) GetIfEventOpen(ctx contractapi.TransactionContextInterface, id string, id2 string) (bool, error) {
	var events Event
	event, err := ctx.GetStub().GetState(id + id2)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	if event != nil {
		err = json.Unmarshal(event, &events)
		if err != nil {
			return false, fmt.Errorf("Falha na leitura do evento : %v %s", err, event)
		}
		return events.Vstatus, nil
	}
	return false, nil

}

func (s *SmartContract) GetEventOpen(ctx contractapi.TransactionContextInterface, id string, id2 string) ([]*Event, error) {

	//queryString := fmt.Sprintf(`{"selector":{"docType":"event","vstatus":"false","iduser1":"%s","iduser2":"%s" }}`, id, id2)
	queryString := fmt.Sprintf(`{"selector":{"docType":"event","vstatus":false,"iduser1":"%s","iduser2":"%s" }}`, id, id2)
	//arrayteststring := [3]string["true", ]
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var eventobjs []*Event

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var eventobj Event
		err = json.Unmarshal(queryResult.Value, &eventobj)
		if err != nil {
			return nil, err
		}
		eventobjs = append(eventobjs, &eventobj)
	}
	return eventobjs, nil

}
