package main

import (
	"bufio"
	"bytes"	
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	s "strings"
	"sync"
	"time"
	
	"github.com/InGaN/Go-ErrorParser/slack"
)

var i = 0
var pTotalTagged *int = &i
var x = 0
var pTotalFilesTagged *int = &x
var a = 0
var pTotalFiles *int = &a
var b int64 = 0 //this all seems redundant
var pTotalFileSize *int64 = &b
var buffer bytes.Buffer
var pTags *[]string

var pointerLock sync.Mutex
var mimeTxt = mime.TypeByExtension(".txt")	

type channelStruct struct {
	TotalFiles chan int
	TotalFileSize chan int64
	TotalTagged chan int	
}

var (	
	flagHelp1	 = flag.Bool("h", false, "help")
	flagHelp2	 = flag.Bool("help", false, "help")
	flagSearchR	 = flag.Bool("r", false, "Recursive search")
	flagEchoMsg	 = flag.Bool("e", false, "Echo messages containing tags")
	flagFileDir	 = flag.String("f", ".", "Directory containing files to parse")
	flagTags	 = flag.String("t", "tags.txt", "file containing tags")
	flagSlack	 = flag.Bool("s", false, "Send error messages to Slack")
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func visit(path string, f os.FileInfo, err error) error {
	//fmt.Printf("Visited: %s\n", path)
	arr := (s.Split(path,"."))
	mimeTxt := mime.TypeByExtension(".txt")
	if(len(arr)>1) {
		if(mime.TypeByExtension("."+arr[len(arr)-1]) == mimeTxt) {
			//checkContainsScanner(path, *pTags)
		}
	}  
  return nil
} 

func checkFiles(path string, channels channelStruct) {	
		
	if(*flagSearchR) {
		err := filepath.Walk(path, visit)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		files, _ := ioutil.ReadDir(path)
		for _, file := range files { // this should be async							
			go checkContainsScanner(fmt.Sprintf("%s\\%s",path, file.Name()), *pTags, channels)
			//<-(channels.TotalFiles)
			//<-(channels.TotalFileSize)
			//<-(channels.TotalTagged)
		}	
	}	
}

func checkContainsScanner(path string, tags []string, channels channelStruct) {	
	arr := (s.Split(path,"."))
	if(len(arr)>1) {		
		if(mime.TypeByExtension("."+arr[len(arr)-1]) == mimeTxt) {		
			
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
		
			//pointerLock.Lock()
			*pTotalFiles++		
			channels.TotalFiles <- *pTotalFiles
			
			stat, err := file.Stat()
			*pTotalFileSize = *pTotalFileSize + stat.Size()
			channels.TotalFileSize <- *pTotalFileSize
				
			defer file.Close()
			changeFlag := *pTotalTagged
			
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				input := scanner.Text()
				parseLine(input, tags)
			}	
			if(changeFlag < *pTotalTagged) {
				buffer.WriteString(fmt.Sprintf("%s (%d)\n", file.Name(), (*pTotalTagged-changeFlag)))
				*pTotalFilesTagged++		
			}
			channels.TotalTagged <- *pTotalFilesTagged
			//pointerLock.Unlock()
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}
		
	}
}

func checkContainsReadString(path string, tags []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	
	*pTotalFiles++
	stat, err := file.Stat()
	*pTotalFileSize = *pTotalFileSize + stat.Size()
	
	changeFlag := *pTotalTagged
	
	r := bufio.NewReader(file)
	for {
		recordRaw, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		parseLine(recordRaw, tags)
	}			
	if(changeFlag < *pTotalTagged) {
		buffer.WriteString(fmt.Sprintf("%s (%d)\n", file.Name(), (*pTotalTagged-changeFlag)))
		*pTotalFilesTagged++
	}
	file.Close()
}

func parseLine(input string, tags []string) {
	amount := 0
	for index, element := range tags {
		if s.Contains(input, element) {
			if(*flagEchoMsg) {
				fmt.Printf("tag: %s | %s \n", tags[index], input)
			}
			if(*flagSlack) {
				slack.SendJSONtoSlack(fmt.Sprintf("tag: %s | %s \n", tags[index], input))
			}
			amount++
		}
	}
	*pTotalTagged = *pTotalTagged + amount
}

func parseTagFile(path string) []string{
	extension := (s.Split(path,"."))
	if(mime.TypeByExtension("."+extension[len(extension)-1]) == mime.TypeByExtension(".txt")) {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		arr := s.Split(string(file), ",")		
		return arr
	}	
	return nil
}

func showHelp() {
	fmt.Println("This Go program parses plain text files containing certain tags.")
	fmt.Println("Paramaters that can be used:")
	fmt.Println("-e\t\t\techo messages containing tags")
	fmt.Println("-f=<path>\t\tpath to folder to search in")
	fmt.Println("-h, -help\t\tshow help")
	fmt.Println("-r\t\t\trecursive search")
	fmt.Println("-t=<path>\t\tpath to csv file containing tags")
}

func getByteSize(value int64) string {
	sizes := []string{"Byte", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	
	if(value < 0) {
		return "-"
	}
	i := 0
	for (value/1024 >= 1) {
		value = value / 1024
		i++
	}
	return fmt.Sprintf("%d%s", value, sizes[i])
}

func main() {
	flag.Parse()
	if(*flagHelp1 || *flagHelp2) {
		showHelp()
	} else {
		defer timeTrack(time.Now(), "Operation")		
		if *flagFileDir == "" {
			fmt.Fprintln(os.Stderr, "require a folder")
			flag.Usage()
			os.Exit(1)
		}	
		fmt.Println("\nStarting parse...")		
		
		tags := parseTagFile(*flagTags)
		pTags = &tags
		
		channels := channelStruct{}
		channels.TotalFiles = make(chan int)
		channels.TotalFileSize = make(chan int64)
		channels.TotalTagged = make(chan int)

		checkFiles(*flagFileDir, channels)			
		
		t := time.Now()
		fmt.Println("\n=== FILES ===")
		fmt.Printf(buffer.String())
		fmt.Println("\n=== RESULTS ===")
		fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
		fmt.Println("Directory: ", *flagFileDir)
		fmt.Println("Tags file: ", *flagTags)
		fmt.Printf("Tags: %v\n", parseTagFile(*flagTags))
		fmt.Println("Recursive: ", *flagSearchR)	
		fmt.Println("Send to Slack Channel: ", *flagSlack)
		fmt.Printf("Total files scanned: %d\n", channels.TotalFiles)
		fmt.Printf("total file size: %s\n", getByteSize(*pTotalFileSize))
		fmt.Printf("total amount of tags found: %d\n", *pTotalTagged)
		fmt.Printf("in %d files\n", *pTotalFilesTagged)					
	}
}
