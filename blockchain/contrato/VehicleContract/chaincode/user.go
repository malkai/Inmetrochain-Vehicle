package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Insere um usuario na blockchain
func (s *SmartContract) Createuser(ctx contractapi.TransactionContextInterface, id string, name string) error {
	exists, err := s.Userexist(ctx, "user"+id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("O usuario %s já existe", id)
	}

	user := User{
		Id:         "user" + id,
		Name:       name,
		Criptmoeda: 0.0,
		Score:      0.5,
	}
	assetJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("user"+id, assetJSON)
}

func (s *SmartContract) Userexist(ctx contractapi.TransactionContextInterface, id string) (bool, error) {

	user, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return user != nil, nil

}

func (s *SmartContract) Updatuser(ctx contractapi.TransactionContextInterface, id string, coin float64) error {

	user, err := ctx.GetStub().GetState(id)

	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	var users User
	err = json.Unmarshal(user, &users)
	if err != nil {
		return err
	}

	users.Criptmoeda = +coin

	assetJSON, err := json.Marshal(users)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(id, assetJSON)

	return nil

}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAlluser(ctx contractapi.TransactionContextInterface, id string) ([]*User, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var users []*User
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var user User
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
