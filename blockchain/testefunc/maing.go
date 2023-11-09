package main

import (
	"encoding/json"
	"fmt"
)

type Tuple struct {
	T    string  `json:"T"`
	Pos  string  `json:"Pos"` // long + lat
	Comb float64 `json:"Comb"`
}

func teste(tuples []Tuple, id string) ([]byte, []byte) {

	t, err := json.Marshal(tuples)
	if err != nil {
		panic(err)
	}

	i, err := json.Marshal(id)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(t)
	fmt.Println()
	fmt.Println(i)

	return t, i

}

func teste2(args ...[]byte) ([]Tuple, string) {

	var t []Tuple
	var id string

	_ = json.Unmarshal(args[0], &t)
	_ = json.Unmarshal(args[1], &id)
	fmt.Println()
	fmt.Println(t)
	fmt.Println()
	fmt.Println(id)

	return t, id

}

func main() {

	var t []Tuple

	var t1 = Tuple{T: "10", Pos: "1/2", Comb: 93.00}
	t = append(t, t1)
	t1 = Tuple{T: "11", Pos: "1/2", Comb: 92.00}
	t = append(t, t1)
	t1 = Tuple{T: "12", Pos: "1/2", Comb: 91.00}
	t = append(t, t1)

	var a, b = teste(t, "1")
	teste2(a, b)

}
