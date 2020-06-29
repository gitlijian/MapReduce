package main

import (
	"fmt"
	"io/ioutil"
	"unicode"
	"strings"
	"os"
	"strconv"
	"log"
	"encoding/json"
)

func DataSlice(document string, value string, mapN int) {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	words := strings.FieldsFunc(value, f)
       //everySlice :=  len(words) / mapN   // 每一个分片的大小

	var tmpFiles  [] *os.File = make([] *os.File, mapN)
	var encoders    [] *json.Encoder = make([] *json.Encoder, mapN)  // 可否不用json编码

	for i := 0; i < mapN; i++ {
		tmpFileName := strconv.Itoa(i+1) + ".txt"
		var err error
		tmpFiles[i],err= os.Create(tmpFileName)
		if err!=nil {
			log.Fatal(err)
		}

		defer tmpFiles[i].Close()
		encoders[i] = json.NewEncoder(tmpFiles[i])
		if err!=nil {
			log.Fatal(err)
		}
	}

	i := 0
	for _ , word := range words {
		index := i % mapN
		err := encoders[index].Encode(&word)
		if err!=nil {
			fmt.Println("do client encoders ",err)
		}
		i++
	}
	
 }

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: go run clientNode.go 待分片的文件名 map节点个数\n")
	} else {
		
		inFile := os.Args[1]
		contents, err := ioutil.ReadFile(inFile)
		if err != nil {
		fmt.Println("do map error for inFile ",err)
		}

		mapN, _ := strconv.Atoi(os.Args[2])

		DataSlice(inFile, string(contents), mapN)
	}
}
