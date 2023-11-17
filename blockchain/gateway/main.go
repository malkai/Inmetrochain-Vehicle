package main

import (
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
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
	//GetAllevents(contract, "event")

	//Createevent(contract, 93.00, 5.0, "1", "3")
	//GetAllevents(contract, "event")
	//GetAllevents(contract, "event1")

	//Createuser(contract, "1", "Malkai", "80")
	//GetAlluser(contract, "user")

	var t []Tuple
	layout := "2006-01-02 15:04:05"
	time2 := time.Now()

	var t1 = Tuple{T: time2.Format(layout), Pos: "1/2", Comb: 93.00}
	t = append(t, t1)
	time2 = time2.Add(10 * time.Second)
	//fmt.Println(time2.Format(layout))
	t1 = Tuple{T: time2.Format(layout), Pos: "2/3", Comb: 92.00}
	t = append(t, t1)
	time2 = time2.Add(10 * time.Second)
	//fmt.Println(time2.Format(layout))
	t1 = Tuple{T: time2.Format(layout), Pos: "4/5", Comb: 91.00}
	t = append(t, t1)
	time2 = time2.Add(10 * time.Second)
	//fmt.Println(time2.Format(layout))
	t1 = Tuple{T: time2.Format(layout), Pos: "5/6", Comb: 90.00}
	t = append(t, t1)

	CreatePath(contract, t, "1")

	//GetAllPath(contract, "Path")
	//GetAlluser(contract, "user")
	GetAllevents(contract, "event")

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
	//Createevent(contract, time.Now(), 50.00, 10.0, "1", "2")
	//GetAllevents(contract, "event1")

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

// Evaluate a transaction to query ledger state.

// Evaluate a transaction to query ledger state.
