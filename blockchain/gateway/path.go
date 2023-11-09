package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type Tuple struct {
	T    string  `json:"T"`
	Pos  string  `json:"Pos"` // long + lat
	Comb float64 `json:"Comb"`
}

func CreatePath(contract *client.Contract, tuples []Tuple, id string) uint64 {

	//test := Path{DataVehicle: tuples}

	t, err := json.Marshal(tuples)
	if err != nil {
		panic(err)
	}
	_, commit, err := contract.SubmitAsync("CreatPath", client.WithArguments(string(t), id))

	/*
		var ue []byte
		if len(ex) == 0 {
			fmt.Println("Error: Empty JSON input")
		} else {
			err = json.Unmarshal(ex, &ue)
			if err != nil {
				panic(err)
			}
		}
	*/

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	status, err := commit.Status()

	if err != nil {
		panic(fmt.Errorf("failed to get transaction commit status: %w", err))
	}

	if !status.Successful {
		panic(fmt.Errorf("failed to commit transaction with status code %v", status.Code))
	}

	fmt.Println("\n*** CreateAsset committed successfully")

	return status.BlockNumber

}

// Evaluate a transaction to query ledger state.
func GetAllPath(contract *client.Contract, id string) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetallPath", id)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}
