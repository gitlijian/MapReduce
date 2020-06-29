package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// KeyValue is a type used to hold the key/value pairs passed to the map and
// reduce functions.
type KeyValue struct {
	Key   string
	Value string
}

// mergeName constructs the name of the output file of reduce task <reduceTask>
func mergeName(reduceID int) string {
	return "reduceResult-"+ strconv.Itoa(reduceID) + ".txt" 
}

// merge combines the results of the many reduce jobs into a single output file
// XXX use merge sort
func merge(nReduce int) {
	kvs := make(map[string]string)
	for i := 0; i < nReduce; i++ {
		p := mergeName(i+1)
		fmt.Printf("Merge: read %s\n", p)
		file, err := os.Open(p)
		if err != nil {
			fmt.Println("Merge: ", err)
		}
		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			err = dec.Decode(&kv)
			if err != nil {
				break
			}
			kvs[kv.Key] = kv.Value
		}
		file.Close()
	}

	var keys []string
	for k := range kvs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	file, err := os.Create("mergeResult.txt")
	if err != nil {
		fmt.Println("Merge: create ", err)
	}
	w := bufio.NewWriter(file)
	for _, k := range keys {
		fmt.Fprintf(w, "%s %s\n", k, kvs[k])
	}
	w.Flush()
	file.Close()

}

func main() {
	if len(os.Args) < 1 {
		fmt.Printf("usage: reduceN\n")
	} else {
		reduceN, _ := strconv.Atoi(os.Args[1])
		merge(reduceN)   // 参数是reduce 的个数，这里有三个reduce
	}
}
