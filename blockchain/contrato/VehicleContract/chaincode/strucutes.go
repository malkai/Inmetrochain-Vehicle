package chaincode

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//DocType        string `json:"docType"` //docType is used to distinguish the various types of objects in state database

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type User struct {
	DocType    string  `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Id         string  `json:"id"`      //vin + cpf or CNPJ
	Name       string  `json:"name"`
	Criptmoeda float64 `json:"criptmoeda"`
	Tanque     float64 `json:"Tanque"` //porcentagem
	Score      float64 `json:"score"`
}

type Event struct {
	DocType   string  `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Id        string  `json:"id"`      //Iduser1+Iduser2+ttimestamp
	Vstatus   bool    `json:"vstatus"` //sentinela de controle
	Iduser1   string  `json:"iduser1"` //identificação do usuario 1
	Iduser2   string  `json:"iduser2"` //identificação do usuario 2
	Datai     string  `json:"datai"`   //data inicial
	Dataiataf string  `json:"dataf"`   //data final do acordo
	Fsupi     float64 `json:"fsupi"`   //combustivel inicial
	Fsupf     float64 `json:"fsupf"`   //combustivel final
	Fsupfd    float64 `json:"fsupfd"`  //i+constant k times
	Dff       float64 `json:"dff"`
}

type Tuple struct {
	T    string  `json:"T"`
	Pos  string  `json:"Pos"` // long + lat
	Comb float64 `json:"Comb"`
}

type Path struct {
	DocType     string  `json:"docType"` //docType is used to distinguish the various types of objects in state database
	EventID     string  `json:"EventID"`
	DataVehicle []Tuple `json:"DataVehicle,omitempty" metadata:"DataVehicle,optional"`
	Distance    float64 `json:"dist"`
	Fuel        float64 `json:"fuel"`
	Totaltime   float64 `json:"time"`
	Timeless    float64 `json:"Timeless"`
	DataR       string  `json:"dataR"`
}
