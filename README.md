Caso a Api não esteja instalada estes são os requisitos minimos necessarios para utilizar está rede
```
sudo apt-get install git curl docker-compose -y

# Make sure the Docker daemon is running.
sudo systemctl start docker

# Add your user to the Docker group.
sudo usermod -a -G docker <username>

# Check version numbers  
docker --version
docker-compose --version

# Optional: If you want the Docker daemon to start when the system starts, use the following:
sudo systemctl enable docker


```

Go
Optional: Install the latest version of Go (only required if you will be writing Go chaincode or SDK applications).

JQ
Optional: Install the latest version of jq (only required for the tutorials related to channel configuration transactions).

Download Fabric samples, Docker images, and binaries

```
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh

./install-fabric.sh docker binary
or
./install-fabric.sh d b
```

Para mais detalhes a documentação será um otimo guia com tutoriais de como a rede funciona. 

<link>https://hyperledger-fabric.readthedocs.io/en/latest/</link>

Entre na pasta para levantar a rede
```
cd blokchain
```

Este comnado vai levantar a rede e criar um canal 
```
./network.sh up createChannel -ca
```

caso deseja criar um canal com nome personalidade utilize o comando. 
```
./network.sh up createChannel -ca
```
ARRUMAR

explicar sobre o que significa o Certificate Autoraty

Após montar nossa rede fazemos o deploy do nosso contrato inteligente. 
```
./network.sh deployCC -ccn vehicle -ccp ../contract/ -ccl go
```

Após realizar o deploy sem erros iremos na pasta que está o nosso gateway. 

EXPLICAR O QUE È GATEWAY 

Explicar quais inguagens utilizam o gateway


```
cd.. 
cd gateway
go run test.go 
```

