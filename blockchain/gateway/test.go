package main

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// Create gRPC client connection, which should be shared by all gateway connections to this endpoint.
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// Create client identity and signing implementation based on X.509 certificate and private key.
	id := NewIdentity()
	sign := NewSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),

		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	chaincodeName := "vehicle3"
	channelName := "mychannel"

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	//initLedger(contract)

	// Context used for event listening
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	// Listen for events emitted by subsequent transactions
	//startChaincodeEventListening(ctx, network)

	Datai := time.Now()
	Dataiataf := time.Now()
	Fsupi := 50.00
	Fsupf := 20.00
	Dff := 10
	Vstatus := false
	Iduser1 := "1"
	Iduser2 := "2"

	Createevent(contract, "1", Datai, Dataiataf, Fsupi, Fsupf, Dff, Vstatus, Iduser1, Iduser2)
	getAllAssets(contract)

	//replayChaincodeEvents(ctx, network, firstBlockNumber)

}

func startChaincodeEventListening(ctx context.Context, network *client.Network) {
	fmt.Println("\n*** Start chaincode event listening")
	chaincodeName := "vehicle3"

	events, err := network.ChaincodeEvents(ctx, chaincodeName)
	if err != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
	}

	go func() {
		for event := range events {
			asset := formatJSON(event.Payload)
			fmt.Printf("\n<-- Chaincode event received: %s - %s\n", event.EventName, asset)
		}
	}()
}

func replayChaincodeEvents(ctx context.Context, network *client.Network, startBlock uint64) {
	fmt.Println("\n*** Start chaincode event replay")
	chaincodeName := "vehicle3"

	events, err := network.ChaincodeEvents(ctx, chaincodeName, client.WithStartBlock(startBlock))
	if err != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
	}

	for {
		select {
		case <-time.After(50 * time.Second):
			panic(errors.New("timeout waiting for event replay"))

		case event := <-events:
			asset := formatJSON(event.Payload)
			fmt.Printf("\n<-- Chaincode event replayed: %s - %s\n", event.EventName, asset)

			if event.EventName == "DeleteAsset" {
				// Reached the last submitted transaction so return to stop listening for events
				return
			}
		}
	}
}

func Createevent(contract *client.Contract, Id string, Datai time.Time, Dataiataf time.Time, Fsupi float64, Fsupf float64, Dff int, Vstatus bool, Iduser1 string, Iduser2 string) uint64 {
	s1 := strconv.FormatFloat(Fsupi, 'E', -1, 64)
	s2 := strconv.FormatFloat(Fsupf, 'E', -1, 64)
	s3 := strconv.Itoa(Dff)
	s4 := strconv.FormatBool(Vstatus)
	_, commit, err := contract.SubmitAsync("Createevent", client.WithArguments(Id, Datai.String(), Dataiataf.String(), s1, s2, s3, s4, Iduser1, Iduser2))
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
func getAllAssets(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAllevents")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func initLedger(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger \n")

	_, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate("../blockchain/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt")
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, "peer0.org1.example.com")

	connection, err := grpc.Dial("localhost:7051", grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// NewIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func NewIdentity() *identity.X509Identity {
	certificatePEM, err := os.ReadFile("../blockchain/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/cert.pem")

	//cryptoPath   = "../../test-network/organizations/peerOrganizations/org1.example.com"
	//keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"

	panicOnError(err)

	certificate, err := identity.CertificateFromPEM(certificatePEM)

	panicOnError(err)

	id, err := identity.NewX509Identity("Org1MSP", certificate)
	panicOnError(err)

	return id
}

// NewSign creates a function that generates a digital signature from a message digest using a private key.
func NewSign() identity.Sign {
	files, err := os.ReadDir("../blockchain/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/")
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}

	privateKeyPEM, err := os.ReadFile(path.Join("../blockchain/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/", files[0].Name()))
	panicOnError(err)

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	panicOnError(err)

	sign, err := identity.NewPrivateKeySign(privateKey)
	panicOnError(err)

	return sign
}

func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
