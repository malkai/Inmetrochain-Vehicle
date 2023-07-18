package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
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

	chaincodeName := "vehicle"
	channelName := "mychannel"

	network := gw.GetNetwork(channelName)

	contract := network.GetContract(chaincodeName)

	//initLedger(contract)

	// Context used for event listening
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	// Listen for events emitted by subsequent transactions
	//startChaincodeEventListening(ctx, network)
	/*
		Datai := time.Now()
		Dataiataf := time.Now()
		Fsupi := 50.00
		Fsupf := 20.00
		Dff := 10
		Vstatus := false
		Iduser1 := "1"
		Iduser2 := "2"
	*/
	Createevent(contract, time.Now(), 50.00, 10.0, "1", "2")
	GetAllevents(contract, "event1")

	//updatevent("event1", 1.0)
	//updatevent("event1", 1.0)
	//updatevent("event1", 1.0)
	//updatevent("event1", 1.0)
	//updatevent("event1", 1.0)
	//updatevent("event1", 1.0)
	//updatevent("event1", 1.0)
	//updatevent("event1", 1.0)

	//initLedger(contract)
	//getAllAssets(contract)

	//replayChaincodeEvents(ctx, network, firstBlockNumber)

}

func Createevent(contract *client.Contract, Datai time.Time, Fsupi float64, Dff float64, Iduser1 string, Iduser2 string) uint64 {
	s1 := strconv.FormatFloat(Fsupi, 'f', 2, 64)
	s3 := strconv.FormatFloat(Dff, 'f', 2, 64)

	_, commit, err := contract.SubmitAsync("Createevent", client.WithArguments(Datai.String(), s1, s3, Iduser1, Iduser2))
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
func GetAllevents(contract *client.Contract, id string) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAllevents", id)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

// Evaluate a transaction to query ledger state.
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

func Closeevent(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger \n")

	_, err := contract.SubmitTransaction("Closeevent")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate("../organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt")
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
	certificatePEM, err := os.ReadFile("../organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/cert.pem")

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
	files, err := os.ReadDir("../organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/")
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}

	privateKeyPEM, err := os.ReadFile(path.Join("../organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/", files[0].Name()))
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
