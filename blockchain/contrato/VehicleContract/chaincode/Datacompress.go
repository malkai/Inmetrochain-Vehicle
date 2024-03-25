package chaincode

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"log"
)

func EncodeToBytes(p interface{}) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
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

func Decompress(s []byte) ([]byte, error) {

	rdr, err := gzip.NewReader(bytes.NewReader(s))
	if err != nil {
		return nil, fmt.Errorf("\n Erro ao Ler. %v %s", err, s)
	}
	data, err := io.ReadAll(rdr)
	if err != nil {
		return nil, fmt.Errorf("\n Erro ao descomprimir. %v", err)
	}
	rdr.Close()
	return data, nil
}

func DecodeToTuple(s []byte) ([]Tuple, error) {

	p := []Tuple{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		return nil, fmt.Errorf("\n Erro ao decodar. %v", err)
	}
	return p, nil
}
