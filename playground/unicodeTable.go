package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

//func main() {
//	loadUnicodeTable()
//}
func loadUnicodeTable() {
	f, err := os.Open("unicodeTable.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var name string
	var from, to uint32
	//var prevTo uint32
	var count int
	for scanner.Scan() {

		line := scanner.Text()
		parts := strings.Split(line, ";")
		_, e := fmt.Sscanf(parts[0], "%x..%x", &from, &to)
		if e != nil {
			fmt.Printf("error: %v\n", e)
			break

		}
		if count > 0 {
			fmt.Printf(",\t")
		} else {
			fmt.Printf("{\n\t")
		}
		count++
		name = strings.TrimSpace(parts[1])
		fmt.Printf("\t\"%s\": [%d", name, from/0x80)
		for block := 1 + from/0x80; block <= to/0x80; block++ {
			fmt.Printf(", %d", block)
		}
		fmt.Printf("]\n")

		//if from > prevTo+1 {
		//	fmt.Printf("%x\t%x\t%x-%x\t%x-%x\t!!!!!SKIPPED!!!!!\n", (from-prevTo)/0x80, (to - from), from/0x80, to/0x80, from, to)
		//}
		//fmt.Printf("%x  %x\n", from&0xFFFFFF00, to&0xFFFFFF00)
		//if from&0xFFFFFF00 != to&0xFFFFFF00 {
		//	fmt.Printf("%x - %x %s\n", from, to, name)
		//}

		//prevTo = to
	}
	fmt.Printf("}\n")
}
