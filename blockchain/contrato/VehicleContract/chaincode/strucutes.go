package chaincode

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type User struct {
	Id         string  `json:"id"` //vin + cpf or CNPJ
	Name       string  `json:"name"`
	Criptmoeda float64 `json:"criptmoeda"`
	Score      float64 `json:"score"`
}

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

type Tuple struct {
	T    string  `json:"T"`
	Pos  string  `json:"Pos"` // long + lat
	Comb float64 `json:"Comb"`
}

type Path struct {
	DataVehicle [][]Tuple `json:"DataVehicle"`
	EventID     string    `json:"EventID"`
	Distance    float64   `json:"dist"`
	Fuel        float64   `json:"fuel"`
	Totaltime   float64   `json:"time"`
	DataR       string    `json:"dataR"`
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}
