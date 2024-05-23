<h2>Introduçção</h2>
 Este repositorio propoe uma plataforma de monetização que utiliza a tecnologia blockchain para criar um ecossistema seguro e transparente para a coleta, armazenamento e comercialização de dados de telemetria veicular. A telemetria veicular envolve a coleta de dados em tempo real sobre o desempenho e o comportamento dos veículos. Esses dados podem incluir informações sobre velocidade, consumo de combustível, localização, diagnósticos de motor, entre outros Tradicionalmente, esses dados são subutilizados, mas a plataforma permite que os proprietários de veículos monetizem essas informações valiosas. Além disso, desenvolvemos um conjunto de métodos baseados em técnicas de crowdsourcing.


Os dados coletados consistem principalmente em informações essenciais para estimar a distância percorrida e o consumo de combustível de um veículo em um trajeto específico, com o objetivo de avaliar sua eficiência energética. Essas informações são transmitidas para a rede blockchain, onde são processadas por meio de contratos inteligentes. 
 
Os contratos inteligentes, automatizados e seguros, são responsáveis por transformar os dados do veículo em ativos monetizáveis. Eles garantem que todas as transações e acordos sejam realizados de forma automática, transparente e sem a necessidade de intermediários, proporcionando aos proprietários de veículos uma nova fonte de renda e promovendo uma utilização mais eficiente dos dados de telemetria.


<h2>A Hyperledger Fabric</h2>

A plataforma Hyperledger Fabric, é uma blockchain permissionada que se destaca
por sua estrutura de nós diferenciados em categorias distintas, que incluem clientes, nós e
ordenadores. Os clientes desempenham um papel fundamental, agindo em nome dos usuários
finais que solicitam transações e estabelecendo comunicação tanto com os nós quanto com os
ordenadores. Os peers, também conhecidos como nós, têm a responsabilidade crucial de manter o
livro-razão. Eles recebem mensagens de atualização cuidadosamente ordenadas dos ordenadores
para validar e confirmar as novas transações registradas no livro-razão.

O modelo de consenso padrão utilizado pela plataforma Hyperledger é conhecido como
Raft. O algoritmo Raft implementa o consenso elegendo inicialmente um líder, atribuindo-
lhe total responsabilidade total responsabilidade pelo gerenciamento dos dados replicados. A rede emprega identidades digitais com atributos adicionais que desempenham um papel
fundamental na determinação de permissões. Esses atributos são vinculados a uma identidade
por meio de um identificador especial e podem abranger uma ampla variedade de características
da identidade de um participante, como a organização à qual pertence, a unidade organizacional,
a função desempenhada ou até mesmo a identidade específica do participante. Quando se
trata desses identificadores, são esses atributos que influenciam diretamente as permissões
correspondentes. 


 <h2>Guia de Instalação do Hyperledger Fabric</h2>


Caso a plataforma não esteja instalada na maquina, existem alguns requisitos minimos necessarios para utilizar está rede

```
sudo apt-get install git curl docker-compose -y

# Make sure the Docker daemon is running.
sudo systemctl start docker


# Add your user to the Dockcder group.
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker
sudo shutdown -r now


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


Utilizamos o comando git clone para clonar o repositorio que contém o projeto. O projeto deve ser clonado em areas que não tenham problemas de permissão root. Exemplo Area de trabalho, documentos ou qualquer local que não seja a raiz do seu user. 

```
git clone https://github.com/malkai/Inmetrochain-Vehicle 
cd Inmetrochain-Vehicle
```


Dentro do seu projeto faço o download dos Docker images, and binaries. 

```
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh

./install-fabric.sh docker binary
or
./install-fabric.sh d b
```

Para mais detalhes a documentação será um otimo guia com tutoriais de como a rede funciona -  <link>https://hyperledger-fabric.readthedocs.io/en/latest/</link>



 <h2>Guia de execução do Hyperledger Fabric</h2>

Para iniciar a rede é necessario adentrar na pasta do projeto que contém o projeto já configurado. 
```
cd blockchain
```

Este script traz uma rede Hyperledger Fabric para testar contratos inteligentes  e aplicativos. A rede de teste consiste em duas organizações com uma peer cada, e um serviço de pedido Raft de nó único. Os usuários também podem usar este script para criar um canal implanta um chaincode no canal. 

```
./network.sh up createChannel -ca
```

caso deseja criar um canal com nome personalidade utilize o comando. 

```
./network.sh up createChannel -c meucanal
```

Antes de criar uma rede, cada organ
ização precisa gerar a criptografia material que vai definir aquela organização na rede. Porque Hyperledger Fabric é um blockchain autorizado, cada nó e usuário na rede precisa use certificados e chaves para assinar e verificar suas ações. Além disso, cada usuário precisa pertencer a uma organização reconhecida como membro da rede. Você pode usar a ferramenta Cryptogen ou Fabric CAs para gerar a criptografia da organizaçãomateriais.
 
Por padrão, a rede de amostra usa cryptogen. Cryptogen é uma ferramenta que é destinado ao desenvolvimento e teste que pode criar rapidamente os certificados e chaves que pode ser consumido por uma rede Fabric. A ferramenta cryptogen consome uma série de arquivos de configuração para cada organização no "organizations/cryptogen" diretório. Cryptogen usa os arquivos para gerar o material criptográfico para cada org no diretório "organizações".

Você também pode usar Fabric CAs para gerar o material criptográfico. CAs assinam os certificados e as chaves que eles geram para criar uma raiz de confiança válida para cada organização.O script usa o Docker Compose para trazer três CAs, uma para cada organização de mesmo nível e a organização do pedido. O arquivo de configuração para criar o Fabric CA servidores estão no diretório "organizations/fabric-ca". No mesmo diretório, o script "registerEnroll.sh" usa o cliente Fabric CA para criar as identidades, certificados e pastas MSP necessários para criar a rede de teste no Diretório "organizations/ordererOrganizations".

Após montar nossa rede fazemos o deploy do nosso contrato inteligente. Do ponto de vista de um desenvolvedor de aplicativos, um contrato inteligente, juntamente com o livro-razão, formam o coração de um sistema blockchain Hyperledger Fabric. Enquanto um livro-razão contém fatos sobre o estado atual e histórico de um conjunto de objetos de negócios, um contrato inteligente define a lógica executável que gera novos fatos que são adicionados ao livro-razão. Um chaincode é normalmente usado por administradores para agrupar contratos inteligentes relacionados para implantação, mas também pode ser usado para programação de sistema de baixo nível do Fabric. 

```
./network.sh deployCC -ccn vehicle -ccp contrato/BasicCRUD -ccl go
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
cd gateway
go run test.go 
```

Upgrade do contrato inteligente

```
./network.sh deployCC -ccn nomedocontrato -ccp path -ccl go  -ccv versão -ccs sequencia
./network.sh deployCC -ccn vehicle -ccp contrato/VehicleContract -ccl go -ccv 1.1 -ccs 2

```

Para fazer a rede parar utilize os comandos abaixo:

```
cd .. 
cd blockchain
./network.sh down
```

 <h2>Execução do experimento</h2>


Para validar as métricas de monetização, foi gerada uma simulação de
36 horas. Os arquivos estão disponiveis no [Link](https://drive.google.com/file/d/17bCUjP_sxY6WPvoaEXvpwyyK5zacWWSy/view?usp=sharing), 
que contém todas as informações necessariais para a simulação.  

Se você já completou o guia de instalção dos componentes e se familarizou, 
iniciaremos com o guia de como executar a rede. Inicialmente executa-se os seguintes comandos


```
cd .. 
cd blockchain
./network.sh up createChannel -ca && ./network.sh deployCC -ccn vehicle -ccp contrato/VehicleContract -ccl go
```


Um dos contratempos para o condutor é a inconsistência na captura de dados veiculares,
uma vez que depende de diversos fatores além do seu controle direto. Teoricamente, pode não
haver nenhum problema, porém, devido a limitações tecnológicas, o condutor pode não ser
capaz de obter informações de forma constante do OBD. Isso ocorre devido à variação na
disponibilidade dos dispositivos para a captura, na diversidade de aplicativos disponíveis e na
forma como foi feita a programação para a aquisição e armazenamento de dados, além das
tecnologias presentes no veículo.

Para simular o problema anterior, foi proposto que antes do envio das tuplas veiculares,
fosse utiliza a variável baseada na métrica da frequência k. Por exemplo, caso k, seja igual a
três, isso significa que o condutor envia somente um terço da informação, perdendo dois terços
das informações. Esta abordagem auxilia na compreensão do comportamento das métricas para
determinados tipos de veículos.

Além disso, notamos que as informações do Sistema de Monitoramento de Consumo
de Combustível (SUMO) são geradas e enviadas sem considerar o comportamento analisado a
partir dos dados reais. Para abordar essa questão, desenvolvemos um algoritmo para introduzir
ruído branco nas informações do tanque de combustível . Antes de serem enviados, os dados
passam por um processo no qual ruído branco é inserido nas tuplas de dados.

A inserção de dados na blockchain ocorre por meio do Fabric Gateway, um serviço
incorporado aos nós do Hyperledger Fabric v2.4. Essa ferramenta oferece uma API simplificada
e minimalista, facilitando o envio de transações para uma rede Fabric. Com esse serviço, é
possível enviar dados veiculares de forma simultânea, otimizando o tempo de inserção das
informações. No entanto, com o aumento do número de tuplas no trajeto e, consequentemente,
do tempo de resposta do endorser para processar a transação, constatou-se que a rede blockchain
não conseguia lidar com a inserção simultânea dos 150 veículos simulados. Mesmo com a
compressão das informações, a rede ainda precisava processar esses dados para gerar uma
pontuação para o condutor, o que limitava a capacidade de acomodação tanto no envio quanto no
número de tuplas. Por isso, foi estabelecido um tamanho máximo de tuplas, variando entre 1000
e 4000, e a análise de 15 veículos por vez