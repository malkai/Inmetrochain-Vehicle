package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func Createevent(contract *client.Contract, Fsupi float64, Dff float64, Iduser1 string, Iduser2 string) uint64 {
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

	fmt.Println("Create Event successfully", Iduser1)

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
func GetStatusEvent(contract *client.Contract, id, id2 string) bool {
	fmt.Println("\n--> Evaluate Transaction: GetStatusEvent")

	evaluateResult, err := contract.EvaluateTransaction("GetIfEventOpen", "Event"+id, id2)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	if evaluateResult != nil {
		result := formatJSON(evaluateResult)
		fmt.Printf("*** Result:%s", result)
		boolValue, _ := strconv.ParseBool(result)
		return boolValue
	} else {
		return false

	}

}

// Evaluate a transaction to query ledger state.
func GetopenEvent(contract *client.Contract, id, id2 string) bool {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetEventOpen", id, id2)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	if evaluateResult != nil {
		result := formatJSON(evaluateResult)
		fmt.Printf("*** Result:%s\n", result)
		boolValue, _ := strconv.ParseBool(result)
		return boolValue
	} else {
		return false
	}

}

// Evaluate a transaction to query ledger state.
func GetAlleventsPast(contract *client.Contract, id string) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetallEventPast", id)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	if evaluateResult != nil {
		result := formatJSON(evaluateResult)
		fmt.Printf("*** Result:%s\n", result)
	} else {
		fmt.Printf("*** Erro não encontrou nada\n")
	}
}

// Evaluate a transaction to query ledger state.
func CloseEvent(contract *client.Contract, id string, minus string) {
	fmt.Println("\n--> Evaluate Transaction: CloseEvent, function returns all the current assets on the ledger")
	_, commit, err := contract.SubmitAsync("Closeevent", client.WithArguments(id, minus))
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

	fmt.Println("Close Event successfully", id)

}

// Evaluate a transaction to query ledger state.
func GetAlleventsNews(contract *client.Contract, id string) {
	fmt.Println("\n--> Evaluate Transaction: GetAlleventsNews, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetallEventNew", id)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	if evaluateResult != nil {
		result := formatJSON(evaluateResult)
		fmt.Printf("*** Result:%s\n", result)
	} else {
		fmt.Printf("*** Erro não encontrou nada\n")
	}
}

// Evaluate a transaction to query ledger state.
func GetAllevents(contract *client.Contract, id string) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAllevents", id)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	if evaluateResult != nil {
		result := formatJSON(evaluateResult)
		fmt.Printf("*** Result:%s\n", result)
	} else {
		fmt.Printf("*** Erro não encontrou nada\n")
	}
}
