package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

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
