package chaincode

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

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

	var fsumpi, _ = strconv.ParseFloat(Fsupi, 64)
	var dffr, _ = strconv.ParseFloat(Dff, 64)

	event := Event{
		DocType: "eventnew",
		Id:      "Event" + Iduser1 + Iduser2,
		Datai:   aux.AsTime().Format(layout),
		Dataf:   "",
		Fsupi:   fsumpi,
		Fsupf:   0,
		Dff:     dffr,
		Fsupfd:  0.0,
		Vstatus: true,
		Iduser1: Iduser1,
		Iduser2: Iduser2,
		Compl:   0.0,
		Freq:    0.0,
		Confi:   0.0,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(event.Id, eventJSON)
}

func (s *SmartContract) GeteventOpensingle(ctx contractapi.TransactionContextInterface, id string, id2 string) ([]Event, error) {

	aux, err := ctx.GetStub().GetState("Event" + id + id2)
	if err != nil {
		return nil, fmt.Errorf("erro ao retornar evento em aberto: %v", err)
	}
	if aux == nil {
		return nil, fmt.Errorf("erro ao recuperar o evento em aberto: %v", err)
	}
	ue := []Event{}
	ueaux := Event{}
	err = json.Unmarshal(aux, &ueaux)
	if err != nil {
		return nil, fmt.Errorf("erro ao retornar evento em aberto: %v", err)
	}

	ue = append(ue, ueaux)
	return ue, nil

}

func (s *SmartContract) GeteventOpensingleext(ctx contractapi.TransactionContextInterface, id string, id2 string) int {

	aux, err := ctx.GetStub().GetState("Event" + id + id2)
	if err != nil {
		return 0
	}
	if aux == nil {
		return 0
	}
	ue := []Event{}
	ueaux := Event{}
	err = json.Unmarshal(aux, &ueaux)
	if err != nil {
		return 0
	}

	ue = append(ue, ueaux)
	return int(ue[0].Dff)

}

func (s *SmartContract) Eventexist(ctx contractapi.TransactionContextInterface, id string, id2 string) (bool, error) {

	aux, err := ctx.GetStub().GetState(id + id2)
	if err != nil {
		return false, fmt.Errorf("Erro ao retornar evento em aberto: %v", err)
	}

	return aux != nil, nil

}

func (s *SmartContract) Closeeventext(ctx contractapi.TransactionContextInterface, id string, sum float64) error {

	eventjson1, err := s.GeteventOpensingle(ctx, id, "9999")
	if err != nil {
		return fmt.Errorf("erro ao recuperar evento: %v", err)
	}
	eventjson := eventjson1[0]

	layout := "2006-01-02 15:04:05"

	aux, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("erro ao pegar data: %v", err)
	}

	eventjson.Dataf = aux.AsTime().Format(layout)
	eventjson.Vstatus = false

	trt := strings.Replace(eventjson.Datai, " ", "-", -1)
	trt = strings.Replace(trt, ":", "-", -1)

	temp, err := s.GetPathhIndex(ctx, id, trt)
	if err != nil {
		return err
	}

	ntotal := 0.0
	ntimeless := 0.0

	if len(temp) > 0 {

		for _, xx := range temp {

			ntimeless = ntimeless + xx.Timeless
			ntotal = ntotal + 1

		}

		ntimeless = ntimeless / ntotal

		if math.IsNaN(ntimeless/ntotal) || math.IsNaN(eventjson.Fsupfd/eventjson.Dff) {
			return fmt.Errorf("/n erro na completude externa  %f, %f, %f, %f", eventjson.Dff, eventjson.Fsupfd, ntimeless, ntotal)
		}

		//Atualiza confiança do usuario
		conf, err := s.Updatuser(ctx, "user"+id, ntimeless, eventjson.Fsupfd/eventjson.Dff)
		if err != nil {
			return fmt.Errorf("erro ao atualizar user externo: %v", err)
		}

		user, err := s.Userget(ctx, "user"+id)
		if err != nil {
			return fmt.Errorf("\n Error ao recuperar user %v %s", err, "user"+id)
		}

		eventjson.DocType = "eventpast"
		eventjson.Id = eventjson.Id + eventjson.Dataf + "eventpast"

		eventjson.Fsupf = eventjson.Fsupi - ((eventjson.Fsupfd) * 100 / user.Tank)
		if eventjson.Fsupfd/eventjson.Dff > 1 {
			eventjson.Compl = 1
		} else if eventjson.Fsupfd/eventjson.Dff < 0 {
			eventjson.Compl = 0
		} else if math.IsNaN(eventjson.Fsupfd / eventjson.Dff) {
			eventjson.Compl = 0

		} else {
			eventjson.Compl = eventjson.Fsupfd / eventjson.Dff
		}
		if ntimeless > 1 {
			eventjson.Freq = 1
		} else if ntimeless < 0 {
			eventjson.Freq = 0
		} else {
			eventjson.Freq = ntimeless
		}
		if conf > 1 {
			eventjson.Confi = 1
		} else if conf < 0 {
			eventjson.Confi = 0
		} else {
			eventjson.Confi = conf
		}

		eventJSON, err := json.Marshal(eventjson)
		if err != nil {
			return fmt.Errorf("erro ao salvar evento antigo externo1 : %v", err)
		}

		err = ctx.GetStub().PutState(eventjson.Id, eventJSON)
		if err != nil {
			return fmt.Errorf("erro ao salvar evento antigo externo2: %v", err)
		}
	} else {

		conf, err := s.Updatuser(ctx, "user"+id, 0, 0)
		if err != nil {
			return fmt.Errorf("erro ao atualizar user externo: %v", err)
		}

		eventjson.DocType = "eventpast"
		eventjson.Id = eventjson.Id + eventjson.Dataf + "eventpast"

		eventjson.Fsupf = 0

		eventjson.Compl = 0

		eventjson.Freq = 0

		eventjson.Confi = conf

		eventJSON, err := json.Marshal(eventjson)
		if err != nil {
			return fmt.Errorf("erro ao salvar evento antigo externo3 : %v", err)
		}

		err = ctx.GetStub().PutState(eventjson.Id, eventJSON)
		if err != nil {
			return fmt.Errorf("erro ao salvar evento antigo externo4: %v", err)
		}
	}

	event := Event{
		DocType: "eventnew",
		Id:      "Event" + eventjson.Iduser1 + eventjson.Iduser2,
		Datai:   "",
		Dataf:   "",
		Fsupi:   0,
		Fsupf:   0,
		Dff:     0,
		Fsupfd:  0.0,
		Vstatus: false,
		Iduser1: eventjson.Iduser1,
		Iduser2: eventjson.Iduser2,
		Compl:   0.0,
		Freq:    0.0,
		Confi:   0.0,
	}

	eventJSONnew, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(event.Id, eventJSONnew)

}

func (s *SmartContract) Closeevent(ctx contractapi.TransactionContextInterface, id string, sum float64, pathant Path) error {

	eventjson1, err := s.GeteventOpensingle(ctx, id, "9999")
	if err != nil {
		return fmt.Errorf("erro ao recuperar evento: %v", err)
	}
	eventjson := eventjson1[0]

	layout := "2006-01-02 15:04:05"

	var fuelsum float64 = 0.0
	aux, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("erro ao pegar data: %v", err)
	}

	eventjson.Dataf = aux.AsTime().Format(layout)
	eventjson.Vstatus = false
	eventjson.Fsupfd = eventjson.Fsupfd + sum

	trt := strings.Replace(eventjson.Datai, " ", "-", -1)
	trt = strings.Replace(trt, ":", "-", -1)

	temp, err := s.GetPathhIndex(ctx, id, trt)
	if err != nil {
		return err
	}

	ntotal := 0.0
	ntimeless := 0.0

	for _, xx := range temp {
		fuelsum = fuelsum + xx.Fuel
		ntimeless = ntimeless + xx.Timeless
		ntotal = ntotal + 1

	}
	fuelsum = fuelsum + pathant.Fuel
	ntimeless = ntimeless + pathant.Timeless
	ntotal = ntotal + 1

	if ntotal > 0 {

		ntimeless = ntimeless / ntotal

		if math.IsNaN(ntimeless/ntotal) || math.IsNaN(fuelsum/eventjson.Dff) {
			return fmt.Errorf("/n erro nas contas  interno %f, %f, %f,%f, %d", eventjson.Fsupfd, eventjson.Dff, eventjson.Dff, fuelsum, len(temp))
		}

		//Atualiza confiança do usuario
		conf, err := s.Updatuser(ctx, "user"+id, ntimeless, eventjson.Fsupfd/eventjson.Dff)
		if err != nil {
			return fmt.Errorf("erro ao atualizar user interno: %v", err)

		}

		user, err := s.Userget(ctx, "user"+id)
		if err != nil {
			return fmt.Errorf("\n Error ao recuperar user %v %s", err, "user"+id)
		}

		eventjson.DocType = "eventpast"
		eventjson.Id = eventjson.Id + eventjson.Dataf + "eventpast"

		eventjson.Fsupf = eventjson.Fsupi - ((eventjson.Fsupfd) * 100 / user.Tank)
		if eventjson.Fsupfd/eventjson.Dff > 1 {
			eventjson.Compl = 1
		} else if eventjson.Fsupfd/eventjson.Dff < 0 {
			eventjson.Compl = 0
		} else {
			eventjson.Compl = eventjson.Fsupfd / eventjson.Dff
		}
		if ntimeless > 1 {
			eventjson.Freq = 1
		} else if ntimeless < 0 {
			eventjson.Freq = 0
		} else {
			eventjson.Freq = ntimeless
		}
		if conf > 1 {
			eventjson.Confi = 1
		} else if conf < 0 {
			eventjson.Confi = 0
		} else {
			eventjson.Confi = conf
		}

		eventJSON, err := json.Marshal(eventjson)
		if err != nil {
			return fmt.Errorf("erro ao salvar evento antigo interno1: %v", err)
		}

		err = ctx.GetStub().PutState(eventjson.Id, eventJSON)
		if err != nil {
			return fmt.Errorf("erro ao salvar evento antigo interno2: %v", err)
		}
	} else {

		conf, err := s.Updatuser(ctx, "user"+id, 0, 0)
		if err != nil {
			return fmt.Errorf("erro ao atualizar user externo: %v", err)
		}

		eventjson.DocType = "eventpast"
		eventjson.Id = eventjson.Id + eventjson.Dataf + "eventpast"

		eventjson.Fsupf = 0

		eventjson.Compl = 0

		eventjson.Freq = 0

		eventjson.Confi = conf

		eventJSON, err := json.Marshal(eventjson)
		if err != nil {
			return fmt.Errorf("erro ao salvar evento antigo externo1 : %v", err)
		}

		err = ctx.GetStub().PutState(eventjson.Id, eventJSON)
		if err != nil {
			return fmt.Errorf("erro ao salvar evento antigo externo2: %v", err)
		}
	}

	event := Event{
		DocType: "eventnew",
		Id:      "Event" + eventjson.Iduser1 + eventjson.Iduser2,
		Datai:   "",
		Dataf:   "",
		Fsupi:   0,
		Fsupf:   0,
		Dff:     0,
		Fsupfd:  0.0,
		Vstatus: false,
		Iduser1: eventjson.Iduser1,
		Iduser2: eventjson.Iduser2,
		Compl:   0.0,
		Freq:    0.0,
		Confi:   0.0,
	}

	eventJSONnew, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(event.Id, eventJSONnew)

}

func (s *SmartContract) updatevent(ctx contractapi.TransactionContextInterface, id, id2, minus1 string, path Path) error {
	minus, _ := strconv.ParseFloat(minus1, 64)
	var eventjson1, err = s.GeteventOpensingle(ctx, id, id2)
	if err != nil {
		return fmt.Errorf("erro ao recuperar evento: %v", err)
	}
	eventjson := eventjson1[0]
	if minus > 0 {
		eventjson.Fsupfd = eventjson.Fsupfd + minus
	}

	if eventjson.Dff < eventjson.Fsupfd {

		err = s.Closeevent(ctx, id, minus, path)
		if err != nil {
			return fmt.Errorf("erro em salvar o banco: %v", err)
		}

		return nil

	}

	eventJS, err := json.Marshal(eventjson)
	if err != nil {
		return fmt.Errorf("erro ao compactar novo evento: %v", err)
	}

	err = ctx.GetStub().PutState(eventjson.Id, eventJS)
	if err != nil {
		return fmt.Errorf("erro em salvar o banco: %v", err)
	}

	return nil

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

func (s *SmartContract) GetIfEventOpen(ctx contractapi.TransactionContextInterface, id string, id2 string) (float64, error) {
	var events Event
	event, err := ctx.GetStub().GetState(id + id2)
	if err != nil {
		return 0.0, fmt.Errorf("failed to read from world state: %v", err)
	}

	if event != nil {
		err = json.Unmarshal(event, &events)
		if err != nil {
			return 0.0, fmt.Errorf("falha na leitura do evento : %v %s", err, event)
		}
		return events.Dff, nil
	}
	return 0.0, nil

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
