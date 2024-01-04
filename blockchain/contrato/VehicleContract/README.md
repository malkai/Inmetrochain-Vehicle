./network.sh up createChannel -ca

./network.sh deployCC -ccn vehicle -ccp contrato/VehicleContract -ccl go

./network.sh deployCC -ccn vehicle -ccp contrato/VehicleContract -ccl go -ccv 1.1 -ccs 2

go run monetiza.go 

peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n vehicle -c '{"function":"Closeevent","Args":["event1"]}'

./network.sh down



