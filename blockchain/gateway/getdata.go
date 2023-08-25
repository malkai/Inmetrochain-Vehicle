package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func isAnagram(s string, t string) bool {
	arr := []rune(s)
	arr2 := []rune(t)
	allKeys := make(map[string]int)

	if len(s) == len(t) {

		for i := 0; i < len(s); i++ {

			//		fmt.Printf("\n%s %s", arr[i], arr2[i])

			_, ok := allKeys[string(arr[i])]
			if ok {

				allKeys[string(arr[i])] += 1

			} else {

				allKeys[string(arr[i])] = 1
			}

			_, ok = allKeys[string(arr2[i])]
			if ok {

				allKeys[string(arr2[i])] += 1

			} else {

				allKeys[string(arr2[i])] = 1
			}

		}

		for _, element := range allKeys {

			if element%2 != 0 {

				return false

			}

		}

	} else {
		return false
	}
	return true

}

func CompletenessTimelinessCredibility(data []string) (string, string, string) {

	return "", "", ""
}

func completness(data []string) string {
	return ""
}

func Timeliness(data []string) string {
	return ""
}

func Credibility(data []string) string {
	return ""
}

func main() {

	var a bool
	a = isAnagram("rat", "car")
	fmt.Printf("\n%t", a)

}

func matrix() [][]string {
	jsonFile, err := os.Open("malkai.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	fmt.Println("Successfully Opened users.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		fmt.Println(err)
	}
	//out := make([]string, len(result))

	b := make([]string, len(result))

	var matriz [][]string
	var ant = ""
	temp := make([]string, 0)

	for i, v := range result {

		help := fmt.Sprint(v)
		//fmt.Printf("\n")
		//fmt.Println(help)
		r, _ := regexp.Compile("((?:id)[^\\s|^%|^\\]]+).*((?:lat)[^\\s|^%|^\\]]+).*((?:long)[^\\s|^%|^\\]]+).*((?:time)[^%|^\\]]+).*((?:Combustivel)[^\\s|^%|^\\]]+)")

		matches := r.FindAllStringSubmatch(help, -1)
		s1 := strings.Split(matches[0][1], ":")
		if ant == "" {
			ant = s1[1]
		}
		s2 := strings.Split(matches[0][2], ":")
		s3 := strings.Split(matches[0][3], ":")
		s4 := strings.Split(matches[0][4], "time:")
		s5 := strings.Split(matches[0][5], "Combustivel:")

		//fmt.Println(s1[1], s2[1], s3[1], s4[1])

		temp = append(temp, s1[1]+"/"+s2[1]+"+"+s3[1]+"/"+s4[1]+"/"+s5[1])

		if ant != s1[1] {
			ant = s1[1]
			//fmt.Printf("\n")
			//print(temp[0])
			matriz = append(matriz, temp)
			temp = nil
		}

		//for i, v := range matches {
		//	var yt = matches[i][4]
		//var aux interface{ yt }
		//result = append(result, aux)

		//}
		//fmt.Printf(help)

		b[i] = fmt.Sprint(v)
	}

	return matriz

}

//(?:id|lat|long|time|Combustivel?)[^\s|^%|^\]]+
