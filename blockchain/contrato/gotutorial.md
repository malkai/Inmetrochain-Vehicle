Neste documento iremos iniciar alguns tutorias iremos propor alguns tutorias para entender um pouco da linguagem Go. Go é uma linguagem de programação procedural. Foi desenvolvido em 2007 por Robert Griesemer, Rob Pike e Ken Thompson no Google, mas lançado em 2009 como uma linguagem de programação de código aberto. Os programas são montados usando pacotes, para gerenciamento eficiente de dependências. 

O passo inicial é criar uma pasta que vai armazenar os codigos em go.
```
mkdir go-exemplo
cd go-exemplo
go mod init go-exemplo
```

A partir desses comandos estamos em nosso pasta em go. Criamos um arquivo.go que vai armazenar o nosso codigo inicial. 

```
package main
import ("fmt")


// Este é um comentario em GO
func main() {
  fmt.Println("Hello World!")
}
```

Na linha 1: Em Go, todo programa faz parte de um pacote. Definimos isso usando a palavra-chave package. Neste exemplo, o programa pertence ao pacote principal.

Na linha 2: import ("fmt") nos permite importar arquivos incluídos no pacote fmt.

Na linha 3: Uma linha em branco. Ir ignora o espaço em branco. Ter espaços em branco no código torna-o mais legível.

Na linha 4: func main() {} é uma função. Qualquer código dentro das chaves {} será executado.

Na linha 5: fmt.Println() é uma função disponibilizada pelo pacote fmt. É usado para enviar/imprimir texto. Em nosso exemplo, a saída será "Hello World!".

Para executar o arquivo que criamos, realizamos o seguinte comando:

```
go run arquivo.go 
```

Na linguagem existem tipos primitiso como:

int- armazena números inteiros (números inteiros), como 123 ou -123

float32- armazena números de ponto flutuante, com decimais, como 19,99 ou -19,99

string - armazena texto, como "Hello World". Os valores de string são colocados entre aspas duplas

bool- armazena valores com dois estados: true ou false