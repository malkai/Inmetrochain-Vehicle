package chaincode_test

import (
	"encoding/json"
	"testing"
	"vehiclecontract/chaincode"
	"vehiclecontract/chaincode/mocks"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/stretchr/testify/require"
)

//go:generate counterfeiter -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

func TestPath(t *testing.T) {

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	//testa a inserção de um usuario na rede
	assetTransfer := chaincode.SmartContract{}
	err := assetTransfer.CreatPath(transactionContext)
	require.NoError(t, err)

	/*
		chaincodeStub.GetStateReturns([]byte{}, nil)
		err = assetTransfer.Createuser(transactionContext, "1", "Malkai")
		require.EqualError(t, err, "O usuario 1 já existe")
	*/

	/*
		//define um asset para retornar da busca
		user := &chaincode.User{Id: "1", Name: "Malkai"}
		bytes, err := json.Marshal(user)
		require.NoError(t, err)

		iterator := &mocks.StateQueryIterator{}
		iterator.HasNextReturnsOnCall(0, true)
		iterator.HasNextReturnsOnCall(1, false)
		iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)
		chaincodeStub.GetStateByRangeReturns(iterator, nil)
		users, err := assetTransfer.GetAlluser(transactionContext, "")
		require.NoError(t, err)
		require.Equal(t, []*chaincode.User{user}, users)
	*/

}

func TestUser(t *testing.T) {

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	//testa a inserção de um usuario na rede
	assetTransfer := chaincode.SmartContract{}
	err := assetTransfer.Createuser(transactionContext, "1", "Malkai")
	require.NoError(t, err)

	//testa o case quando o usuario existe
	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = assetTransfer.Createuser(transactionContext, "1", "Malkai")
	require.EqualError(t, err, "O usuario 1 já existe")

	//define um user para retornar da busca
	user := &chaincode.User{Id: "1", Name: "Malkai"}
	bytes, err := json.Marshal(user)
	require.NoError(t, err)

	//define um array de users que será retornado
	iterator := &mocks.StateQueryIterator{}
	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	chaincodeStub.GetStateByRangeReturns(iterator, nil)
	users, err := assetTransfer.GetAlluser(transactionContext, "")
	require.NoError(t, err)
	require.Equal(t, []*chaincode.User{user}, users)

}

/*
func TestGetAllAssets(t *testing.T) {
	asset := &chaincode.Asset{ID: "asset1"}
	bytes, err := json.Marshal(asset)
	require.NoError(t, err)

	iterator := &mocks.StateQueryIterator{}
	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	chaincodeStub.GetStateByRangeReturns(iterator, nil)
	assetTransfer := &chaincode.SmartContract{}
	assets, err := assetTransfer.GetAllAssets(transactionContext)
	for aux, i := range assets {
		fmt.Print(aux, i)

	}
	require.NoError(t, err)
	require.Equal(t, []*chaincode.Asset{asset}, assets)

	iterator.HasNextReturns(true)
	iterator.NextReturns(nil, fmt.Errorf("failed retrieving next item"))
	assets, err = assetTransfer.GetAllAssets(transactionContext)
	require.EqualError(t, err, "failed retrieving next item")
	require.Nil(t, assets)

	chaincodeStub.GetStateByRangeReturns(nil, fmt.Errorf("failed retrieving all assets"))
	assets, err = assetTransfer.GetAllAssets(transactionContext)
	require.EqualError(t, err, "failed retrieving all assets")
	require.Nil(t, assets)
}
*/

/*
func TestEvent(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}
	err := assetTransfer.Createevent(transactionContext, 51, )
	require.NoError(t, err)

	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err = assetTransfer.InitLedger(transactionContext)
	require.EqualError(t, err, "failed to put to world state. failed inserting key")
}
*/
