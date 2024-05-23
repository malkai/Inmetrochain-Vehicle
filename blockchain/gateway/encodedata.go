package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
)

//conjunto de codigos para fazer a compressão e descompressão dos dados.

func Compress(p []Tuple) string {

	tu, err := json.Marshal(p)
	if err != nil {
		fmt.Println("error:", err)
	}
	var buf bytes.Buffer
	/*

		enc := gob.NewEncoder(&buf)
		err = enc.Encode(p)
		if err != nil {
			log.Fatal(err)
		}
	*/
	buf = bytes.Buffer{}
	gz := gzip.NewWriter(&buf)
	err = json.NewEncoder(gz).Encode(tu)
	if err != nil {
		fmt.Println("error:", err)
	}
	gz.Close()

	aa, err := json.Marshal(buf.Bytes())
	if err != nil {
		fmt.Println("error:", err)
	}

	return (string(aa))

}

func Decompress(s string) []Tuple {
	aa := []byte(s)
	zipbuf := []byte{}
	rty := []Tuple{}
	err := json.Unmarshal(aa, &zipbuf)
	if err != nil {
		fmt.Println("error:", err)
	}
	rdr, err := gzip.NewReader(bytes.NewReader(zipbuf))
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewDecoder(rdr).Decode(&zipbuf)
	if err != nil {
		fmt.Println("error:", err)
	}

	rdr.Close()
	err = json.Unmarshal(zipbuf, &rty)
	if err != nil {
		fmt.Println("error:", err)
	}

	return rty
}
