package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"os"
	s "strings"
	"time"
)

var i = 0
var pTotalAmount *int = &i
var x = 0
var pTotalFiles *int = &x
var buffer bytes.Buffer

func checkFile(e error) {
	if e != nil {
		panic(e)
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func checkContainsScanner(tags []string) {
	defer timeTrack(time.Now(), "Scanner")

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input := scanner.Text()
		parseLine(input, tags)
	}	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func checkFiles(path string) {
	defer timeTrack(time.Now(), "Operation")
	mimeTxt := mime.TypeByExtension(".txt")
	tags := []string{"ERROR", "Error", "error", "Warning"}
	
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if(mime.TypeByExtension(file.Name()[len(file.Name())-4:]) == mimeTxt) {
			//fmt.Printf("Type: %s\n", file.Name())
			checkContainsReadString(fmt.Sprintf("%s\\%s",path, file.Name()), tags)
		}				
	}	
}

func checkContainsReadString(path string, tags []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	changeFlag := *pTotalAmount
	
	r := bufio.NewReader(file)
	for {
		recordRaw, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		parseLine(recordRaw, tags)
	}			
	if(changeFlag < *pTotalAmount) {
		buffer.WriteString(fmt.Sprintf("%s\n", file.Name()))
		*pTotalFiles++
	}
	file.Close()
}

func parseLine(input string, tags []string) {
	amount := 0
	for index, element := range tags {
		if s.Contains(input, element) {
			fmt.Printf("tag: %s | %s \n", tags[index], input)
			amount++
		}
	}
	*pTotalAmount = *pTotalAmount + amount
}


func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}	

	
	
	checkFiles(os.Args[1])
	
	t := time.Now()
	fmt.Printf("%d-%02d-%02d %d:%d:%d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	fmt.Printf("total amount of errors found: %d\n", *pTotalAmount)
	fmt.Printf("in %d files:\n", *pTotalFiles)
	fmt.Printf(buffer.String())
	

	// TXT file parse concept
}
