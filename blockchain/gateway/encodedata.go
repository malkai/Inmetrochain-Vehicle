package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
)

type Path struct {
	DocType     string  `json:"docType"` //docType is used to distinguish the various types of objects in state database
	PathID      string  `json:"EventID"`
	DataVehicle string  `json:"DataVehicle"` //`json:"DataVehicle,omitempty" metadata:"DataVehicle,optional"`
	Distance    float64 `json:"dist"`
	Fuel        float64 `json:"fuel"`
	Totaltime   float64 `json:"time"`
	Timeless    float64 `json:"Timeless"`
	DataR       string  `json:"dataR"`
	DataEvent   string  `json:"dataEvent"`
	K           float64 `json:"k"`
	Iduser      string  `json:"iduser"` //identificação do usuario
}

func EncodeToBytes(p interface{}) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("uncompressed size (bytes): ", len(buf.Bytes()))
	return buf.Bytes()
}

func Compress(s []byte) []byte {

	zipbuf := bytes.Buffer{}
	zipped := gzip.NewWriter(&zipbuf)
	zipped.Write(s)
	zipped.Close()
	fmt.Println("compressed size (bytes): ", len(zipbuf.Bytes()))
	return zipbuf.Bytes()
}

func Decompress(s []byte) []byte {

	rdr, _ := gzip.NewReader(bytes.NewReader(s))
	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		log.Fatal(err)
	}
	rdr.Close()
	fmt.Println("uncompressed size (bytes): ", len(data))
	return data
}

func DecodeToTuple(s []byte) []Tuple {

	p := []Tuple{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
