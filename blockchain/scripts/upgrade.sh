#!/bin/bash

source scripts/utils.sh

 #./network.sh upgrade -c mychannel -ccn vehicle -ccv contracto/VehicleContract  -ccp 3 -ccs 1.3

  # peer lifecycle chaincode package vehicle.out -p contracto/BasicCRUD --label vehicle_1.2
#./network.sh deployCC -ccn vehicle -ccp contracto/BasicCRUD -ccl go
CHANNEL_NAME=${1:-"mychannel"}
CC_NAME=${2}
CC_SRC_PATH=${3}
CC_VERSION=${4:-"1.0"}
CC_SEQUENCE=${5:-"1"}
CC_INIT_FCN=${6:-"NA"}
CC_END_POLICY=${7:-"NA"}
CC_COLL_CONFIG=${8:-"NA"}
DELAY=${9:-"3"}
MAX_RETRY=${10:-"5"}
VERBOSE=${11:-"false"}

println "executing with the following"
println "- CHANNEL_NAME: ${C_GREEN}${CHANNEL_NAME}${C_RESET}"
println "- CC_NAME: ${C_GREEN}${CC_NAME}${C_RESET}"
println "- CC_SRC_PATH: ${C_GREEN}${CC_SRC_PATH}${C_RESET}"
println "- CC_VERSION: ${C_GREEN}${CC_VERSION}${C_RESET}"
println "- CC_SEQUENCE: ${C_GREEN}${CC_SEQUENCE}${C_RESET}"
println "- CC_END_POLICY: ${C_GREEN}${CC_END_POLICY}${C_RESET}"
println "- CC_COLL_CONFIG: ${C_GREEN}${CC_COLL_CONFIG}${C_RESET}"
println "- CC_INIT_FCN: ${C_GREEN}${CC_INIT_FCN}${C_RESET}"
println "- DELAY: ${C_GREEN}${DELAY}${C_RESET}"
println "- MAX_RETRY: ${C_GREEN}${MAX_RETRY}${C_RESET}"
println "- VERBOSE: ${C_GREEN}${VERBOSE}${C_RESET}"

#User has not provided a name
if [ -z "$CC_NAME" ] || [ "$CC_NAME" = "NA" ]; then
  fatalln "Não foi providenciado o nome do seu contrato inteligente"

elif [ -z "$CC_VERSION" ] || [ "$CC_VERSION" = "NA" ]; then
  fatalln "Não foi providenciado a versão do seu contrato"


# User has not provided a path
elif [ -z "$CC_SRC_PATH" ] || [ "$CC_SRC_PATH" = "NA" ]; then
  fatalln "Não foi providenciado o caminho para o contrato inteligente"
fi

packageChaincodeOrg1() {

  export PATH=${PWD}/../bin:${PWD}:$PATH
  export FABRIC_CFG_PATH=$PWD/../config/
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="Org1MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
  export CORE_PEER_ADDRESS=localhost:7051
 

  infoln "Package ${CC_NAME} on Org1"
  peer lifecycle chaincode package ${CC_NAME}${CC_VERSION}.tar.gz --path contracto/BasicCRUD --label vehicle_${CC_VERSION}
  if [[ $? -ne 0 ]]; then
    errorln "Erro ao executar o lifecycle chaincode package ${CC_NAME} "
    exit 1
  fi

  export PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid ${CC_NAME}${CC_VERSION}.tar.gz)
 

  infoln "Install ${CC_NAME}${CC_VERSION}.tar.gz on Org1"
  peer lifecycle chaincode install  ${CC_NAME}${CC_VERSION}.tar.gz 
  if [[ $? -ne 0 ]]; then
    errorln "Erro ao executar o lifecycle chaincode install ${CC_NAME} "
    exit 1
  fi 
  #peer lifecycle chaincode package ${CC_NAME}.tar.gz -p ${CC_SRC_PATH} --label vehicle_${CC_VERSION}
 
  export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" )
  infoln "approveformyorg ${CC_NAME} on Org1  "
  peer lifecycle chaincode approveformyorg  "${TARGET_TLS_OPTIONS[@]}"  --channelID ${CHANNEL_NAME} --name  ${CC_NAME} --version ${CC_VERSION}  --package-id   ${PACKAGE_ID} --sequence ${CC_SEQUENCE}
  if [[ $? -ne 0 ]]; then
    errorln "Erro ao executar o lifecycle chaincode approveformyorg ${CC_NAME} "
    exit 1
  fi
  
  infoln "checkcommitreadiness ${CC_NAME} on Org1  "
  peer lifecycle chaincode checkcommitreadiness "${TARGET_TLS_OPTIONS[@]}" --channelID ${CHANNEL_NAME} --name  ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE}
  if [[ $? -ne 0 ]]; then
    errorln "Erro ao executar o lifecycle chaincode checkcommitreadiness ${CC_NAME} "
    exit 1
  fi
  

}

packageChaincodeOrg2() {
 PATH=${PWD}/../bin:${PWD}:$PATH
 FABRIC_CFG_PATH=$PWD/../config/
 export CORE_PEER_LOCALMSPID="Org2MSP"
 export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
 export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
 export CORE_PEER_ADDRESS=localhost:9051
 export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" )
 peer lifecycle chaincode package ${CC_NAME}.tar.gz -p ${CC_SRC_PATH} --label vehicle_${CC_VERSION}
PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid ${CC_NAME}.tar.gz)

 
 peer lifecycle chaincode package ${CC_NAME}.tar.gz -p ${CC_SRC_PATH} --label vehicle_${CC_VERSION}
 infoln "Package ${CC_NAME} on Org2"
 if [[ $? -ne 0 ]]; then
    errorln "Erro ao executar o lifecycle chaincode package ${CC_NAME} "
    exit 1
  fi
 peer lifecycle chaincode install ${CC_NAME}.tar.gz
 infoln "Install ${CC_NAME}.tar.gz on Org2"
 if [[ $? -ne 0 ]]; then
    errorln "Erro ao executar o lifecycle chaincode install ${CC_NAME} "
    exit 1
  fi
 peer lifecycle chaincode approveformyorg  "${TARGET_TLS_OPTIONS[@]}"  --channelID ${CHANNEL_NAME} --name  ${CC_NAME} --version ${CC_VERSION}  --package-id   ${PACKAGE_ID} --sequence ${CC_SEQUENCE}
 infoln "approveformyorg ${CC_NAME} on Org2  "
  if [[ $? -ne 0 ]]; then
    errorln "Erro ao executar o lifecycle chaincode approveformyorg ${CC_NAME} "
    exit 1
  fi
 peer lifecycle chaincode checkcommitreadiness "${TARGET_TLS_OPTIONS[@]}" --channelID ${CHANNEL_NAME} --name  ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE}
  infoln "checkcommitreadiness ${CC_NAME} on Org2  "
 if [[ $? -ne 0 ]]; then
    errorln "Erro ao executar o lifecycle chaincode checkcommitreadiness ${CC_NAME} "
    exit 1
  fi


}

Assinatura(){
 TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt")
 peer lifecycle chaincode commit "${TARGET_TLS_OPTIONS[@]}" --channelID ${CHANNEL_NAME} --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE}

 infoln "commit ${CC_NAME} on Org1 and Org2"
 if [[ $? -ne 0 ]]; then
    errorln "Erro ao executar o lifecycle chaincode commit ${CC_NAME} "
    exit 1
  fi
  TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt")
  peer lifecycle chaincode querycommitted --channelID  ${CHANNEL_NAME} --name  ${CC_NAME} --cafile "${TARGET_TLS_OPTIONS[@]}"
  TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"  --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt")
  peer lifecycle chaincode querycommitted --channelID  ${CHANNEL_NAME} --name  ${CC_NAME} --cafile "${TARGET_TLS_OPTIONS[@]}"

}

#packageChaincodeOrg1
#packageChaincodeOrg2
#Assinatura

testetudo(){

  export PATH=${PWD}/../bin:${PWD}:$PATH
  export FABRIC_CFG_PATH=$PWD/../config/
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="Org1MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
  export CORE_PEER_ADDRESS=localhost:7051
 

  peer lifecycle chaincode package  vehicle2.tar.gz --path contracto/BasicCRUD --label vehicle_2


  export PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid  vehicle2.tar.gz)
 

  peer lifecycle chaincode install vehicle2.tar.gz  

  export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" )
  peer lifecycle chaincode approveformyorg  "${TARGET_TLS_OPTIONS[@]}"  --channelID mychannel --name  vehicle --version 1.2  --package-id   ${PACKAGE_ID} --sequence 2
  peer lifecycle chaincode checkcommitreadiness "${TARGET_TLS_OPTIONS[@]}" --channelID mychannel --name  vehicle --version 1.2 --sequence 2


  export CORE_PEER_LOCALMSPID="Org2MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
  export CORE_PEER_ADDRESS=localhost:9051
  peer lifecycle chaincode install vehicle2.tar.gz  
  peer lifecycle chaincode package ${CC_NAME}.tar.gz -p ${CC_SRC_PATH} --label vehicle_${CC_VERSION}
  export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" )
  peer lifecycle chaincode approveformyorg  "${TARGET_TLS_OPTIONS[@]}"  --channelID mychannel --name  vehicle --version 1.2  --package-id   ${PACKAGE_ID} --sequence 2
  peer lifecycle chaincode checkcommitreadiness "${TARGET_TLS_OPTIONS[@]}" --channelID mychannel --name  vehicle --version 1.2 --sequence 2

  export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt")
  peer lifecycle chaincode commit "${TARGET_TLS_OPTIONS[@]}" --channelID mychannel --name vehicle --version 1.2  --sequence 2 

  export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt")
  peer lifecycle chaincode querycommitted --channelID  mychannel --name  vehicle --cafile "${TARGET_TLS_OPTIONS[@]}"
  export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"  --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt")
  peer lifecycle chaincode querycommitted --channelID  mychannel --name vehicle --cafile "${TARGET_TLS_OPTIONS[@]}"
  peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n vehicle -c '{"function":"InitLedger","Args":[]}'
 
}
