package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func Createuser(contract *client.Contract, Iduser1 string, name string, tanque string) uint64 {

	_, commit, err := contract.SubmitAsync("Createuser", client.WithArguments(Iduser1, name, tanque))
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

/*

func (s *SmartContract) Userexist(ctx contractapi.TransactionContextInterface, id string) (bool, error) {

	user, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return user != nil, nil

}
*/

func checkuserexist(contract *client.Contract, id string) bool {
	fmt.Println("\n--> Evaluate Transaction: checkuserexist")

	evaluateResult, err := contract.EvaluateTransaction("Userexist", "user"+id)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	if evaluateResult != nil {
		result := formatJSON(evaluateResult)
		fmt.Printf("*** Result:%s\n", result)
		boolValue, err := strconv.ParseBool(result)
		if err != nil {
			fmt.Printf("*** Erro a converter\n")
		}
		return boolValue
	} else {
		fmt.Printf("*** Não encotrou o usuario %s\n", id)
		return false

	}
}

func Getuser(contract *client.Contract, id string) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("Userget", "user"+id)
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

func GetAlluser(contract *client.Contract, id string) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAlluser", id)
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
