package chaincode

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Insere um usuario na blockchain
func (s *SmartContract) Createuser(ctx contractapi.TransactionContextInterface, id, name, tanq, tipo string) error {
	exists, err := s.Userexist(ctx, "user"+id)
	if err != nil {
		return fmt.Errorf("erro checar se o usuario existe. %v", err)
	}
	if exists {
		return fmt.Errorf("o usuario %s já existe", id)
	}

	tanque, err := strconv.ParseFloat(tanq, 64)
	if err != nil {
		return fmt.Errorf(" Erro ao converter. %v", err)
	}

	user := User{
		DocType: "user",
		Id:      "user" + id,
		Name:    name,
		Coin:    0.0,
		Score:   0.5,
		Tank:    tanque,
		Typee:   tipo,
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("erro ao criar usuario. %v", err)
	}

	return ctx.GetStub().PutState(user.Id, userJSON)
}

func (s *SmartContract) Userget(ctx contractapi.TransactionContextInterface, id string) (User, error) {
	var users User
	user, err := ctx.GetStub().GetState(id)
	if err != nil {
		return users, fmt.Errorf("erro em acessar a informação na blockchain: %v", err)
	}
	if user != nil {
		err = json.Unmarshal(user, &users)
		if err != nil {
			return users, fmt.Errorf("Falha na leitura do usuario : %v %s", err, users.Id)
		}

		return users, nil
	}
	return users, nil

}

func (s *SmartContract) Userexist(ctx contractapi.TransactionContextInterface, id string) (bool, error) {

	user, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return user != nil, nil

}

//s.Updatuser(ctx, "user"+id, 1*score, ntimeless/ntotal, eventjson.Dff, fuelsum, valuevalids )

func (s *SmartContract) Updatuser(ctx contractapi.TransactionContextInterface, id string, timeless, completness float64) (float64, error) {

	user, err := s.Userget(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("erro ao acessar a blockchain para atualizar: %v", err)
	}

	/*
		var users User
		err = json.Unmarshal(user, &users)
		if err != nil {
			return 0, fmt.Errorf("Erro ao acessar o usuario para atualizar: %v", err)
		}
	*/

	m := 0.9
	score := (user.Score * m) + (((timeless + completness) / 2) * (1 - m))

	if math.IsNaN(score) {
		return 0.0, fmt.Errorf("\n Erro ao atualizar o score %f, %f, %f", timeless, completness, score)
	}

	if score > 1 {
		score = 1
	} else if score < 0 {
		score = 0

	}

	user.Coin = user.Coin + 1*score

	user.Score = score

	userJSON, err := json.Marshal(user)
	if err != nil {
		return 0, fmt.Errorf("/n Erro ao compactar json para atualiazar o usuario. %v %s", err, userJSON)
	}
	ctx.GetStub().PutState(id, userJSON)
	return score, nil

}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAlluser(ctx contractapi.TransactionContextInterface, id string) ([]*User, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("u", "v")
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
		if user.Id != "" {
			users = append(users, &user)
		}

	}

	return users, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) StatisticsUser(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.

	/*
		a, err := s.Userget(ctx, id)
		if err != nil {
			return "", err
		}

		event, err := s.GetEventOpen(ctx, id, "9999")
		if err != nil {
			return "", err
		}

		//media

		//moda

		//mediana

		//Percentis

		path, err := s.GetPathhOpen(ctx, "Path"+id)
		if err != nil {
			return "", err
		}

	*/

	/*
		resultsIterator, err := ctx.GetStub().GetStateByRange("u", "v")
		if err != nil {
			return "", err
		}
		defer resultsIterator.Close()

		var users []*User
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return "", err
			}

			var user User
			err = json.Unmarshal(queryResponse.Value, &user)
			if err != nil {
				return "", err
			}
			if user.Id != "" {
				users = append(users, &user)
			}

		}
	*/

	return "", nil
}
