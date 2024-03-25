package chaincode

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//cria evento antigo

// Cria um evento no blockchain

// Cria um evento no blockchain
func (s *SmartContract) Createevent(ctx contractapi.TransactionContextInterface, Fsupi string, Dff string, Iduser1 string, Iduser2 string) error {
	layout := "2006-01-02 15:04:05"
	aux, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("/n Erro ao criar evento. %v", err)
	}

	/*
		exist, err := s.Eventexist(ctx, Iduser1, Iduser2)
		if err != nil {
			return fmt.Errorf("\n Erro ao verificar a existencia do evento. %v", err)
		}
	*/

	/*
		if exist {
			//return fmt.Errorf("\n Erro ao criar evento 2. %v", err)
			err = s.Closeevent(ctx, Iduser1, Iduser2, 0.0)
			if err != nil {
				return fmt.Errorf("\n Não pode ser fechado o evento. %v", err)
			}
		}
	*/

	var fsumpi, _ = strconv.ParseFloat(Fsupi, 64)
	var dffr, _ = strconv.ParseFloat(Dff, 64)

	event := Event{
		DocType:   "eventnew",
		Id:        "Event" + Iduser1 + Iduser2,
		Datai:     aux.AsTime().Format(layout),
		Dataiataf: "",
		Fsupi:     fsumpi,
		Fsupf:     0,
		Dff:       dffr,
		Fsupfd:    0.0,
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

func (s *SmartContract) Closeevent(ctx contractapi.TransactionContextInterface, id string, sum float64) error {

	eventjson, err := s.GetEventOpenSingle(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao recuperar evento: %v", err)
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

	trt := strings.Replace(eventjson.Datai, " ", "-", -1)
	trt = strings.Replace(trt, ":", "-", -1)

	temp, err := s.GetPathhOpen(ctx, id, trt)
	if err != nil {
		return err
	}

	ntotal := 0.0
	ntimeless := 0.0
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

	if math.IsNaN(ntimeless/ntotal) || math.IsNaN(eventjson.Dff/fuelsum) {
		return fmt.Errorf("/n erro Credibility value conf %f, %f, %f,%f, %d", ntimeless, ntotal, eventjson.Dff, fuelsum, len(temp))
	}

	//Atualiza confiança do usuario
	conf, err := s.Updatuser(ctx, "user"+id, ntimeless/ntotal, eventjson.Dff/fuelsum, valuevalids)
	if err != nil {
		return fmt.Errorf("erro ao atualizar user: %v", err)
	}

	eventjson.DocType = "eventpast"
	eventjson.Id = eventjson.Id + eventjson.Dataiataf + "eventpast"

	/*
		if eventjson.Fsupfd+sum <= eventjson.Dff {
			eventjson.Fsupfd = eventjson.Fsupf + sum
		} else {
			eventjson.Fsupfd = eventjson.Dff
		}*/

	eventjson.Fsupf = eventjson.Fsupi - fuelsum
	eventjson.Compl = eventjson.Fsupi / eventjson.Fsupf
	eventjson.Freq = ntimeless / ntotal
	eventjson.Confi = conf
	eventJSON, err := json.Marshal(eventjson)
	if err != nil {
		return fmt.Errorf("erro ao salvar evento antigo: %v", err)
	}

	return ctx.GetStub().PutState(eventjson.Id, eventJSON)

}

func (s *SmartContract) updatevent(ctx contractapi.TransactionContextInterface, id string, id2 string, minus1 string) error {
	minus, _ := strconv.ParseFloat(minus1, 64)
	var eventjson, err = s.GetEventOpenSingle(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao recuperar evento: %v", err)
	}
	eventjson.Fsupfd = eventjson.Fsupfd + minus

	if eventjson.Dff < eventjson.Fsupfd {

		return s.Closeevent(ctx, id, minus)
	}

	eventJS, err := json.Marshal(eventjson)
	if err != nil {
		return fmt.Errorf("erro ao compactar novo evento: %v", err)
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

func (s *SmartContract) GetEventOpenSingle(ctx contractapi.TransactionContextInterface, id string) (Event, error) {
	events := Event{}

	//queryString := fmt.Sprintf(`{"selector":{"docType":"event","vstatus":"false","iduser1":"%s","iduser2":"%s" }}`, id, id2)
	queryString := fmt.Sprintf(`{"selector":{"docType":"eventnew","iduser1":"%s"}}`, id)
	//arrayteststring := [3]string["true", ]
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return events, err
	}
	defer resultsIterator.Close()

	var eventobjs []*Event

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return events, err
		}
		var eventobj Event
		err = json.Unmarshal(queryResult.Value, &eventobj)
		if err != nil {
			return events, err
		}
		eventobjs = append(eventobjs, &eventobj)
	}

	events = *eventobjs[0]

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
func (s *SmartContract) GetallEventPast(ctx contractapi.TransactionContextInterface, id string) ([]*Event, error) {

	//queryString := fmt.Sprintf(`{"selector":{"docType":"event","vstatus":"false","iduser1":"%s","iduser2":"%s" }}`, id, id2)
	queryString := fmt.Sprintf(`{"selector":{"docType":"eventpast","iduser1":"%s" }}`, id)
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

func (s *SmartContract) GetallEventNew(ctx contractapi.TransactionContextInterface, id string) ([]*Event, error) {

	//queryString := fmt.Sprintf(`{"selector":{"docType":"event","vstatus":"false","iduser1":"%s","iduser2":"%s" }}`, id, id2)
	queryString := fmt.Sprintf(`{"selector":{"docType":"eventnew","iduser1":"%s" }}`, id)
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
