package tagparser

import (
	"bufio"
	"bytes"	
	"encoding/json"	
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	s "strings"
	"sync"
	"time"
)

var pTags *[]string
var mimeTxt = mime.TypeByExtension(".txt")	
var mimeJSON = mime.TypeByExtension(".json")

type counter struct {
	TotalTagged int
	TotalFileSize int64
	FileName string
} 
type Response struct {
	Date string `json:"date"`
	Directory string `json:"directory"`
	TagsFile string `json:"tagsFile"`
	Tags []string `json:"tags"`	
	Scanned int `json:"scanned"`
	Size string `json:"size"`
	AmountTags int `json:"amountTags"`
	AmountFiles int `json:"amountFiles"`
}

var channels chan counter
var waitGroups *sync.WaitGroup = &sync.WaitGroup{}	

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func visit(path string, f os.FileInfo, err error) error {
	arr := (s.Split(path,"."))
	if(len(arr)>1) {
		if(mime.TypeByExtension("."+arr[len(arr)-1]) == mimeTxt || mime.TypeByExtension("."+arr[len(arr)-1]) == mimeJSON) {			
			go checkContainsScanner(path, *pTags, channels, waitGroups)
		}
	}  
  return nil
} 

func searchMethod(path string, c chan counter, wg *sync.WaitGroup, r bool) {			
	if(r) {
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
		go checkContainsScanner(fmt.Sprintf("%s\\%s",path, file.Name()), *pTags, c, wg)	
	}	
}

func checkContainsScanner(path string, tags []string, c chan counter, wg *sync.WaitGroup) {		
	wg.Add(1)
	ctr := counter{}
	Tagged := 0
	TotalTags := 0
	var FileSize int64 = 0
	var FileName string
	arr := (s.Split(path,"."))
	if(len(arr)>1) {		
		if(mime.TypeByExtension("."+arr[len(arr)-1]) == mimeTxt || mime.TypeByExtension("."+arr[len(arr)-1]) == mimeJSON) {					
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
				
			stat, err := file.Stat()
			FileSize = stat.Size()
			
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				input := scanner.Text()
				Tagged = parseLine(input, tags)
				TotalTags += Tagged
			}	
			if(TotalTags > 0) {
				FileName = path
			}
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}		
	}
	ctr.TotalTagged = TotalTags
	ctr.TotalFileSize = FileSize
	ctr.FileName = FileName
	c <- ctr
	wg.Done()
}

func parseLine(input string, tags []string) int {
	amount := 0
	for _, element := range tags { //_ = index
		if s.Contains(input, element) {			
			//if(e) {
			//	fmt.Printf("tag: %s | %s \n", tags[index], input)
			//}	
			amount++
		}
	}
	return amount
}

func parseTagFile(path string) []string{
	arr := (s.Split(path,"."))
	if(mime.TypeByExtension("."+arr[len(arr)-1]) == mimeTxt || mime.TypeByExtension("."+arr[len(arr)-1]) == mimeJSON) {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		val := s.Split(string(file), ",")		
		return val
	}	
	return nil
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

// Parameters:
// r - recursive search
// e - echo each single line with tags
// pr - print summary
// dir - directory to search .txt files
// tags - .txt file containing comma separated tags
// tagArray - []string array directly instead of text file
func Parse(r bool, pr bool, dir string, tags string, tagArray []string) ([]byte, error)  {
	if (pr) { 
		defer timeTrack(time.Now(), "Operation")
		fmt.Println("\nStarting parse...")
	}
	if (dir == "") {		
		dir = "."
	}				
	
	if(tagArray == nil) {
		if(tags == "") {
			return nil, errors.New("No tags specified")
		} else {
			t := parseTagFile(tags)
			pTags = &t
		}		
	} else {
		pTags = &tagArray
	}
	
	channels = make(chan counter)
	
	searchMethod(dir, channels, waitGroups, r)			
	
	go func(channels chan counter, waitGroups *sync.WaitGroup) {
		waitGroups.Wait()
		close(channels)
	}(channels, waitGroups)
	
	TotalTagged := 0
	TotalFilesTagged := 0
	TotalFiles := 0
	var TotalFileSize int64 = 0
	var buffer bytes.Buffer
	for i := range channels {
		TotalTagged += i.TotalTagged
		if(i.TotalTagged > 0) {
			TotalFilesTagged++
			buffer.WriteString(fmt.Sprintf("%s (%d)\n", i.FileName, i.TotalTagged))
		}
		TotalFileSize += i.TotalFileSize
		if(i.TotalFileSize > 0) {
			TotalFiles++
		}
	}	
	
	t := time.Now()
	
	response := &Response {
		Date: fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()),
		Directory: dir,
		TagsFile: tags,
		Tags: *pTags,
		Scanned: TotalFiles,
		Size: getByteSize(TotalFileSize),
		AmountTags: TotalTagged,
		AmountFiles: TotalFilesTagged,			
	} 
	rsp, err := json.Marshal(response)	
	if(err != nil) { return nil, errors.New("JSON error") }	
	
	if (pr) { 
		fmt.Println("\n=== FILES ===")
		fmt.Printf(buffer.String())
		fmt.Println("\n=== RESULTS ===")
		fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
		fmt.Println("Directory: ", dir)
		fmt.Println("Tags file: ", tags)
		fmt.Printf("Tags: %v\n", parseTagFile(tags))
		fmt.Println("Recursive: ", r)			
		fmt.Printf("Total files scanned: %d\n", TotalFiles)
		fmt.Printf("total file size: %s\n", getByteSize(TotalFileSize))
		fmt.Printf("total amount of tags found: %d\n", TotalTagged)
		fmt.Printf("in %d files\n", TotalFilesTagged)				
	}
	return rsp, err
}
