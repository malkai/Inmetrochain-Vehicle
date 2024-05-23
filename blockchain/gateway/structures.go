package main

type User struct {
	DocType string  `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Id      string  `json:"id"`      //vin + cpf or CNPJ
	Name    string  `json:"name"`
	Coin    float64 `json:"coin"`
	Tank    float64 `json:"Tank"`  //porcentagem
	Typee   string  `json:"Typee"` //porcentagem
	Score   float64 `json:"conf"`
}

type Tuple struct {
	T    string  `json:"T"`
	Pos  string  `json:"Pos"` // long + lat
	Comb float64 `json:"Comb"`
}

type vehicletruct struct {
	T    string  `json:"T"`
	Pos  string  `json:"Pos"` // long + lat
	Comb float64 `json:"Comb"`
}

type Results struct {
	Data []datavehicle
}
type datavehicle struct {
	Id  int `json:"id"`
	Pos struct {
		Long float64 `json:"long"`
		Lat  float64 `json:"lat"`
	} `json:"pos"`
	Distancia         float64 `json:"Dist\u00e2ncia"`
	Combustivel       float64 `json:"combustivel"`
	TamanhoTanque     int     `json:"TamanhoTanque"`
	Time              string  `json:"time"`
	Contrato          int     `json:"Contrato"`
	Tipo              string  `json:"Tipo"`
	Novoabastecimento string  `json:"Novoabastecimento"`

	//  FooBar  string `json:"foo.bar"`

}

type datahelp struct {
	id   string
	Data []datavehicle
}

type datafile struct {
	id   string
	file []string
}

type Event struct {
	DocType string  `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Id      string  `json:"id"`      //Iduser1+Iduser2+ttimestamp
	Vstatus bool    `json:"vstatus"` //sentinela de controle
	Iduser1 string  `json:"iduser1"` //identificação do usuario 1
	Iduser2 string  `json:"iduser2"` //identificação do usuario 2
	Datai   string  `json:"datai"`   //data inicial
	Dataf   string  `json:"dataf"`   //data final do acordo
	Fsupi   float64 `json:"fsupi"`   //combustivel inicial
	Fsupf   float64 `json:"fsupf"`   //combustivel final
	Fsupfd  float64 `json:"fsupfd"`  //i+constant k times
	Dff     float64 `json:"dff"`     //quantidade de combustivel acordado
	Compl   float64 `json:"compl"`   //identificação do usuario 2
	Freq    float64 `json:"freq"`    //data inicial
	Confi   float64 `json:"Confi"`   //data final do acordo
}
type Path struct {
	DocType     string  `json:"docType"` //docType is used to distinguish the various types of objects in state database
	PathID      string  `json:"EventID"`
	DataVehicle string  `json:"DataVehicle"` //`json:"DataVehicle,omitempty" metadata:"DataVehicle,optional"`
	Distance    float64 `json:"dist"`
	Fuel        float64 `json:"fuel"`
	Totaltime   float64 `json:"time"`
	Timeless    float64 `json:"Timeless"`
	DataR       string  `json:"dataR"`
	DataEvent   string  `json:"dataEvent"`
	K           float64 `json:"k"`
	Iduser      string  `json:"iduser"` //identificação do usuario
	Ntuples     int     `json:"ntuples"`
}
