package main

import (
	"bufio"
	"bytes"	
	"flag"
	"fmt"
//	"io"
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

var buffer bytes.Buffer
var pTags *[]string
var mimeTxt = mime.TypeByExtension(".txt")	

type counter struct {
	TotalTagged int
	TotalFileSize int64
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

func searchMethod(path string, c chan counter, wg *sync.WaitGroup) {			
	if(*flagSearchR) {
		err := filepath.Walk(path, visit)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		files, _ := ioutil.ReadDir(path)
		checkFiles(path, files, c, wg)	
	}	
}

func checkFiles(path string, files []os.FileInfo, c chan counter, wg *sync.WaitGroup) {
	for _, file := range files {
		wg.Add(1)
		go checkContainsScanner(fmt.Sprintf("%s\\%s",path, file.Name()), *pTags, c, wg)	
	}	
}

func checkContainsScanner(path string, tags []string, c chan counter, wg *sync.WaitGroup) {		
	ctr := counter{}
	Tagged := 0
	TotalTags := 0
	var FileSize int64 = 0
	arr := (s.Split(path,"."))
	if(len(arr)>1) {		
		if(mime.TypeByExtension("."+arr[len(arr)-1]) == mimeTxt) {					
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			
			//Files++		
			stat, err := file.Stat()
			FileSize = stat.Size()
			
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				input := scanner.Text()
				Tagged = parseLine(input, tags)
				TotalTags += Tagged
			}	
			if(Tagged > 0) {
				buffer.WriteString(fmt.Sprintf("%s (%d)\n", file.Name(), Tagged))	
				
				//FilesTagged++		
			}
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}		
	}
	/*c.TotalFiles <- Files	
	c.TotalFilesTagged <- FilesTagged
	c.TotalTagged <- Tagged	*/	
	ctr.TotalTagged = TotalTags
	ctr.TotalFileSize = FileSize
	c <- ctr
	//c.TotalFileSize <- FileSize
	wg.Done()
}

func parseLine(input string, tags []string) int {
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
	return amount
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
		
		wg := &sync.WaitGroup{}	
		c := make(chan counter)
		
		searchMethod(*flagFileDir, c, wg)			
		
		go func(c chan counter, wg *sync.WaitGroup) {
			wg.Wait()
			close(c)
		}(c, wg)
		
		TotalTagged := 0
		TotalFilesTagged := 0
		TotalFiles := 0
		var TotalFileSize int64 = 0
		for i := range c {
			TotalTagged += i.TotalTagged
			if(i.TotalTagged > 0) {
				TotalFilesTagged++
			}
			TotalFileSize += i.TotalFileSize
			if(i.TotalFileSize > 0) {
				TotalFiles++
			}
		}	
		
		t := time.Now()
		fmt.Println("\n=== FILES ===")
		//fmt.Printf(buffer.String())
		fmt.Println("\n=== RESULTS ===")
		fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
		fmt.Println("Directory: ", *flagFileDir)
		fmt.Println("Tags file: ", *flagTags)
		fmt.Printf("Tags: %v\n", parseTagFile(*flagTags))
		fmt.Println("Recursive: ", *flagSearchR)	
		fmt.Println("Send to Slack Channel: ", *flagSlack)
		fmt.Printf("Total files scanned: %d\n", TotalFiles)
		fmt.Printf("total file size: %s\n", getByteSize(TotalFileSize))
		fmt.Printf("total amount of tags found: %d\n", TotalTagged)
		fmt.Printf("in %d files\n", TotalFilesTagged)				
	}
}
