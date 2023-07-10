O Hyperledger Fabric é uma plataforma para soluções de contabilidade distribuída sustentadas por uma arquitetura modular que oferece altos graus de confidencialidade, resiliência, flexibilidade e escalabilidade. Ele foi projetado para oferecer suporte a implementações conectáveis ​​de diferentes componentes e acomodar a complexidade e os meandros existentes em todo o ecossistema econômico.

Utilizamos o comando git clone para clonar o repositorio que contém o projeto. O projeto deve ser clonado em areas que não tenham problemas de permissão root. Exemplo Area de trabalho, documentos ou qualquer local que não seja a raiz do seu user. 
```
apt install git
git clone https://github.com/malkai/Inmetrochain-Vehicle
cd Inmetrochain-Vehicle

```


Caso a plataforma não esteja instalada na maquina, existem alguns requisitos minimos necessarios para utilizar está rede
```
sudo apt-get install git curl docker-compose -y

# Make sure the Docker daemon is running.
sudo systemctl start docker


# Add your user to the Docker group.
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker
sudo shutdown -r


# Check version numbers  
docker --version
docker-compose --version

# Optional: If you want the Docker daemon to start when the system starts, use the following:
sudo systemctl enable docker


```
Go : Instale a versão mais recente do Go (necessário apenas se você estiver escrevendo aplicativos Go chaincode ou SDK).

```
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go
```

JQ : Instale a versão mais recente do jq (necessário apenas para os tutoriais relacionados às transações de configuração do canal).
```
sudo apt update
sudo apt install -y jq
jq --version
```

Dentro do seu projeto faço o download dos Docker images, and binaries. 

```
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh

./install-fabric.sh docker binary
or
./install-fabric.sh d b
```

Para mais detalhes a documentação será um otimo guia com tutoriais de como a rede funciona -  <link>https://hyperledger-fabric.readthedocs.io/en/latest/</link>

Para iniciar a rede é necessario adentrar na pasta do projeto que contém o projeto já configurado. 
```
cd blokchain
```

Este script traz uma rede Hyperledger Fabric para testar contratos inteligentes  e aplicativos. A rede de teste consiste em duas organizações com uma peer cada, e um serviço de pedido Raft de nó único. Os usuários também podem usar este script para criar um canal implanta um chaincode no canal. 
```
./network.sh up createChannel -ca
```

caso deseja criar um canal com nome personalidade utilize o comando. 

```
./network.sh createChannel -c meucanal
```

Antes de criar uma rede, cada organização precisa gerar a criptografia material que vai definir aquela organização na rede. Porque Hyperledger Fabric é um blockchain autorizado, cada nó e usuário na rede precisa use certificados e chaves para assinar e verificar suas ações. Além disso, cada usuário precisa pertencer a uma organização reconhecida como membro da rede. Você pode usar a ferramenta Cryptogen ou Fabric CAs para gerar a criptografia da organizaçãomateriais.
 
Por padrão, a rede de amostra usa cryptogen. Cryptogen é uma ferramenta que é destinado ao desenvolvimento e teste que pode criar rapidamente os certificados e chaves que pode ser consumido por uma rede Fabric. A ferramenta cryptogen consome uma série de arquivos de configuração para cada organização no "organizations/cryptogen" diretório. Cryptogen usa os arquivos para gerar o material criptográfico para cada org no diretório "organizações".

Você também pode usar Fabric CAs para gerar o material criptográfico. CAs assinam os certificados e as chaves que eles geram para criar uma raiz de confiança válida para cada organização.O script usa o Docker Compose para trazer três CAs, uma para cada organização de mesmo nível e a organização do pedido. O arquivo de configuração para criar o Fabric CA servidores estão no diretório "organizations/fabric-ca". No mesmo diretório, o script "registerEnroll.sh" usa o cliente Fabric CA para criar as identidades, certificados e pastas MSP necessários para criar a rede de teste no Diretório "organizations/ordererOrganizations".

Após montar nossa rede fazemos o deploy do nosso contrato inteligente. Do ponto de vista de um desenvolvedor de aplicativos, um contrato inteligente, juntamente com o livro-razão, formam o coração de um sistema blockchain Hyperledger Fabric. Enquanto um livro-razão contém fatos sobre o estado atual e histórico de um conjunto de objetos de negócios, um contrato inteligente define a lógica executável que gera novos fatos que são adicionados ao livro-razão. Um chaincode é normalmente usado por administradores para agrupar contratos inteligentes relacionados para implantação, mas também pode ser usado para programação de sistema de baixo nível do Fabric. 
```
./network.sh deployCC -ccn vehicle -ccp ../contracto/BasicCRUD -ccl go
```

Após realizar o deploy sem erros iremos na pasta que está o nosso gateway. O Fabric Gateway é um serviço, introduzido nos pares Hyperledger Fabric v2.4, que fornece uma API simplificada e mínima para enviar transações para uma rede Fabric. Os requisitos anteriormente colocados nos SDKs do cliente, como reunir endossos de transações de pares de várias organizações, são delegados ao serviço Fabric Gateway executado em um ponto para permitir o desenvolvimento simplificado de aplicativos e o envio de transações na v2.4.

A partir do Fabric v2.4, os aplicativos clientes devem usar uma das APIs do cliente Fabric Gateway (Go, Node ou Java), que são otimizadas para interagir com o Fabric Gateway. Essas APIs expõem o mesmo modelo de programação de alto nível que foi inicialmente introduzido no Fabric v1.4.

O Fabric Gateway (também conhecido como gateway) gerencia as seguintes etapas de transação:

Avalie uma proposta de transação: Isso invocará uma função de contrato inteligente (chaincode) em um único ponto e retornará o resultado ao cliente. Isso normalmente é usado para consultar o estado atual do razão sem fazer nenhuma atualização do razão. O gateway selecionará preferencialmente um par na mesma organização que o par do gateway e escolherá o par com a maior altura de bloco de razão. Se nenhum par estiver disponível na organização do gateway, ele escolherá um par de outra organização.

Aprovar uma proposta de transação: Isso reunirá respostas de endosso suficientes para satisfazer as políticas de assinatura combinadas (veja abaixo) e retornará um envelope de transação preparado ao cliente para assinatura.

Envie uma transação: Isso enviará um envelope de transação assinado ao serviço de pedidos para confirmar no livro-razão.

Aguarde os eventos de status de confirmação: Isso permite que o cliente espere que a transação seja confirmada no registro e obtenha o código de status de confirmação (validação/invalidação).

Receber eventos chaincode: Isso permitirá que os aplicativos clientes respondam a eventos emitidos por uma função de contrato inteligente quando uma transação é confirmada no registro.

As APIs do cliente Fabric Gateway combinam as ações Endorse/Submit/CommitStatus em uma única função de bloqueio SubmitTransaction para suportar o envio da transação com uma única linha de código. Alternativamente, as ações individuais podem ser invocadas para suportar padrões de aplicativos flexíveis.


```
cd .. 
cd gateway
go run test.go 
```

Para fazer a rede parar utilize os comandos abaixo:

```
cd .. 
cd blockchain
./network.sh down
```
