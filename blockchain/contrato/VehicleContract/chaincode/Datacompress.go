package chaincode

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
)

func EncodeToBytes(p interface{}) []byte {

	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println("error:", err)
	}

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(b)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func Compress(s []byte) []byte {

	zipbuf := bytes.Buffer{}
	zipped := gzip.NewWriter(&zipbuf)
	zipped.Write(s)
	zipped.Close()

	return zipbuf.Bytes()
}

func Decompress(s string) ([]Tuple, error) {

	aa := []byte(s)
	zipbuf := []byte{}

	err := json.Unmarshal(aa, &zipbuf)
	if err != nil {
		return nil, fmt.Errorf("\n Erro ao descompactar o texto JSON. %v", err)
	}
	rdr, err := gzip.NewReader(bytes.NewReader(zipbuf))
	if err != nil {
		return nil, fmt.Errorf("\n Erro ao descompactar o byte para gzip. %v", err)
	}

	err = json.NewDecoder(rdr).Decode(&zipbuf)
	if err != nil {
		return nil, fmt.Errorf("\n Erro ao recuperar as tuplas. %v", err)
	}

	rdr.Close()

	/*

		dd := []byte{}
		dec := gob.NewDecoder(bytes.NewReader(zipbuf))
		err = dec.Decode(&dd)
		if err != nil {
			return nil, fmt.Errorf("\n Erro ao utilizar gob. %v", err)
		}
	*/

	p := []Tuple{}
	err = json.Unmarshal(zipbuf, &p)
	if err != nil {
		return nil, fmt.Errorf("\n Erro ao atribuir info. %v", err)

	}

	return p, nil
}
