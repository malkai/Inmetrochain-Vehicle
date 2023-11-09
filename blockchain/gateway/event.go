package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func Createevent(contract *client.Contract, Datai time.Time, Fsupi float64, Dff float64, Iduser1 string, Iduser2 string) uint64 {
	s1 := strconv.FormatFloat(Fsupi, 'f', 2, 64)
	s3 := strconv.FormatFloat(Dff, 'f', 2, 64)

	_, commit, err := contract.SubmitAsync("Createevent", client.WithArguments(s1, s3, Iduser1, Iduser2))
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

func updatevent(contract *client.Contract, id string, minus float64) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")
	s1 := strconv.FormatFloat(minus, 'f', 2, 64)
	evaluateResult, err := contract.EvaluateTransaction("updatevent", id, s1)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

// Evaluate a transaction to query ledger state.
func GetAllevents(contract *client.Contract, id string) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAllevents", id)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}
