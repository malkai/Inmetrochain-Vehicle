package chaincode

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type User struct {
	Id         string `json:"id"` //vin + cpf or CNPJ
	Name       string `json:"name"`
	Criptmoeda int    `json:"criptmoeda"`
}

type Event struct {
	Id        string  `json:"id"` //Iduser1+Iduser2+ttimestamp
	Datai     string  `json:"datai"`
	Dataiataf string  `json:"dataf"`
	Fsupi     float64 `json:"fsupi"`
	Fsupf     float64 `json:"fsupf"`
	Dff       float64 `json:"dff"`
	Vstatus   bool    `json:"vstatus"`
	Iduser1   string  `json:"iduser1"`
	Iduser2   string  `json:"iduser2"`
}

type Tuple struct {
	T    string  `json:"T"`
	Pos  string  `json:"Pos"` // long + lat
	Comb float64 `json:"Comb"`
}

type Path struct {
	DataVehicle []Tuple `json:"DataVehicle"`
	EventID     string  `json:"EventID"`
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}
