#!/bin/bash

source scripts/utils.sh

CHANNEL_NAME=${1:-"mychannel"}
CC_NAME=${2}
CC_SRC_PATH=${3}
CC_SRC_LANGUAGE=${4}
CC_VERSION=${5:-"1.0"}
CC_SEQUENCE=${6:-"1"}
CC_INIT_FCN=${7:-"NA"}
CC_END_POLICY=${8:-"NA"}
CC_COLL_CONFIG=${9:-"NA"}
DELAY=${10:-"3"}
MAX_RETRY=${11:-"5"}
VERBOSE=${12:-"false"}

println "executing with the following"
println "- CHANNEL_NAME: ${C_GREEN}${CHANNEL_NAME}${C_RESET}"
println "- CC_NAME: ${C_GREEN}${CC_NAME}${C_RESET}"
println "- CC_SRC_PATH: ${C_GREEN}${CC_SRC_PATH}${C_RESET}"
println "- CC_SRC_LANGUAGE: ${C_GREEN}${CC_SRC_LANGUAGE}${C_RESET}"
println "- CC_VERSION: ${C_GREEN}${CC_VERSION}${C_RESET}"
println "- CC_SEQUENCE: ${C_GREEN}${CC_SEQUENCE}${C_RESET}"
println "- CC_END_POLICY: ${C_GREEN}${CC_END_POLICY}${C_RESET}"
println "- CC_COLL_CONFIG: ${C_GREEN}${CC_COLL_CONFIG}${C_RESET}"
println "- CC_INIT_FCN: ${C_GREEN}${CC_INIT_FCN}${C_RESET}"
println "- DELAY: ${C_GREEN}${DELAY}${C_RESET}"
println "- MAX_RETRY: ${C_GREEN}${MAX_RETRY}${C_RESET}"
println "- VERBOSE: ${C_GREEN}${VERBOSE}${C_RESET}"

packageChaincode() {

  PATH=${PWD}/../bin:${PWD}:$PATH
  FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org1.example.com/
  FABRIC_CFG_PATH=$PWD/../config/
  CORE_PEER_TLS_ENABLED=true
  CORE_PEER_LOCALMSPID="Org1MSP"
  CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp
  CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
  CORE_PEER_ADDRESS=localhost:7051
  TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" )  
  peer chaincode upgrade "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n vehicle -v 2.0 -p ../contract/ -c '{"Args":["d", "e", "f"]}'
  successln "Chaincode is packaged"
}

packageChaincode

export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:7051

#--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" )

peer chaincode upgrade "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n vehicle -v 1.2   -c '{"Args":["init"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"

peer chaincode install "${TARGET_TLS_OPTIONS[@]}"  -n vehicle -p ../contract/ -v 1.1
peer chaincode upgrade  -C mychannel -p  remote-repo.com/username/repository -n vehicle 

 peer chaincode list  "${TARGET_TLS_OPTIONS[@]}"  --installed

  peer chaincode instantiate "${TARGET_TLS_OPTIONS[@]}" -C mychannel  -n vehicle -v 1.1  -c '{"Args":["init"]}'
exit 0

peer chaincode instantiate "${TARGET_TLS_OPTIONS[@]}" -C mychannel  -n vehicle -v 1.1  -c '{"Args":["GetAllAssets"]}'

peer chaincode package  "${TARGET_TLS_OPTIONS[@]}" vehicle.out  -n vehicle -p ../contract/ -v 1.1 -s -S

peer chaincode upgrade "${TARGET_TLS_OPTIONS[@]}"  -C mychannel -p  ../contract/  -n vehicle  -v 1.1 -c '{"Args":[""]}' 

peer chaincode instantiate "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n vehicle -v 1.1 -c '{"Args":[""]}' 

peer chaincode install "${TARGET_TLS_OPTIONS[@]}" -n vehicle -v 2.0 -p ../contract/

peer chaincode upgrade "${TARGET_TLS_OPTIONS[@]}" -n vehicle -v 2.0 -C mychannel -c '{"Args":[]}'

peer lifecycle chaincode commit "${TARGET_TLS_OPTIONS[@]}" -n vehicle -v 2.0 -C mychannel 

