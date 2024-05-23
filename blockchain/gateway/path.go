package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func stripCtlFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r != 127 {
			return r
		}
		return -1
	}, str)
}

func CreatePath(contract *client.Contract, tuples []Tuple, id string, id2 string, k float64) uint64 {

	s := strconv.FormatFloat(k, 'f', -1, 64)

	//fmt.Println(uii)

	/*
		dst := &bytes.Buffer{}
		if err := json.Compact(dst, []byte(tu)); err != nil {
			panic(err)
		}
	*/

	aa := Compress(tuples)

	//fmt.Println(buf.String())

	_, commit, err := contract.SubmitAsync("CreatPath", client.WithArguments(string(aa), id, id2, s))

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

	fmt.Println("\n Create Path committed successfully", id, len(tuples))

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

// Evaluate a transaction to query ledger state.
func GetPathhOpen(contract *client.Contract, id string, date string) []byte {
	fmt.Println("\n--> Evaluate Transaction: Get All Path, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetPathhIndex", id, date)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	if evaluateResult != nil {
		//result := formatJSON(evaluateResult)
		//fmt.Printf("*** Result:%s\n", result)
	} else {
		fmt.Printf("*** Erro não encontrou nada\n")
	}

	return evaluateResult

}

func GetPathhAll(contract *client.Contract, id string) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetPathhAll", id, id)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}
