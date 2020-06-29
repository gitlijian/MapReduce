package main

import (
	"hash/fnv"
	"io/ioutil"
	"log"
	"encoding/json"
	"os"
	"unicode"
	"strings"
	"strconv"
	"fmt"
)

// KeyValue is a type used to hold the key/value pairs passed to the map and
// reduce functions.
type KeyValue struct {
	Key   string
	Value string
}

// reduceName constructs the name of the intermediate file which map task
// <mapTask> produces for reduce task <reduceTask>.
func reduceName(mapID int, reduceID int) string {
	return "maptmp-" + strconv.Itoa(mapID)  + "-"  + strconv.Itoa(reduceID) + ".txt"
}

// The mapping function is called once for each piece of the input.
// In this framework, the key is the name of the file that is being processed,
// and the value is the file's contents. The return value should be a slice of
// key/value pairs, each represented by a mapreduce.KeyValue.
func mapF(document string, value string) (res []KeyValue) {
	// TODO: you have to write this function
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	words := strings.FieldsFunc(value, f)

	for _, key := range words {
		res = append(res, KeyValue{key, "1"})
	}
	return res

}

func doMap(mapID int, nReduce int) {

	//setp 1 read file
	inFile := strconv.Itoa(mapID) + ".txt"
	contents, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Fatal("do map error for inFile ",err)
	}
	//setp 2 call user user-map method ,to get kv
	kvResult := mapF(inFile, string(contents))

	/**
	 *   setp 3 use key of kv generator nReduce file ,partition
	 *      a. create tmpFiles
	 *      b. create encoder for tmpFile to write contents
	 *      c. partition by key, then write tmpFile
	 */

	var tmpFiles  [] *os.File = make([] *os.File, nReduce)
	var encoders    [] *json.Encoder = make([] *json.Encoder, nReduce)  // 可否不用json编码

	for i := 0; i < nReduce; i++ {
		tmpFileName := reduceName(mapID, i+1)
		tmpFiles[i],err = os.Create(tmpFileName)
		if err!=nil {
			log.Fatal(err)
		}

		defer tmpFiles[i].Close()
		encoders[i] = json.NewEncoder(tmpFiles[i])
		if err!=nil {
			log.Fatal(err)
		}
	}

	for _ , kv := range kvResult {
		hashKey := int(ihash(kv.Key)) % nReduce /* 将相同的key分配给同一个reducer处理 */
		err := encoders[hashKey].Encode(&kv)
		if err!=nil {
			log.Fatal("do map encoders ",err)
		}
	}

}

func ihash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: go run mapNode.go mapID reduce节点的个数\n")
	} else {
		mapID, _ := strconv.Atoi(os.Args[1])
		reduceN, _ := strconv.Atoi(os.Args[2])
		doMap(mapID, reduceN)  // 参数1是mapID 参数而是reduce节点的个数 
	                     // mapID决定打开文件的名称，因为分配给该map的分片文件命令是根据mapID来的
		             // reduce节点的个数决定了map运算会产生多少个分区文件，从而根据分区号分配给相应的reduce节点
	     }
}
