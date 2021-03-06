package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func main() {
	/*Command line parsing*/
	urlTarget := flag.String("url", "", "parser target url(with http:// or https://)")
	classFind := flag.String("class", "", "parsing by some class")
	idFind := flag.String("id", "", "search by some class")
	elemFind := flag.String("el", "", "search by element type")
	fileToSave := flag.String("save_file", "parsing.txt", "save file for parsing results")
	//proxyFlag := flag.Bool("proxy", false, "use proxy(in dev!)")
	flag.Parse()
	if *urlTarget == "" {
		printHelp()
		return
	}
	response := send_get(*urlTarget)
	data := parse(response, *classFind, *idFind, *elemFind)
	saveToFile(*fileToSave, data)
	fmt.Println(data)
	fmt.Println("************\nEnd of work!\n************")
}

//Sends get request to target server and returns response
func send_get(url string) string {
	userAgent := "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:65.0) Gecko/20100101 Firefox/65.0"
	fmt.Println("[log]sending get request...")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(data)
}

//prints help
func printHelp() {
	fmt.Println("Help:\n-url=*targetUrl*\n-class=*need class or nothing*\n-id=*need id or nothing*\n-el=*need element or nothing*\n-save_file=*result file path*")
	return
}

//Writes data to file
func saveToFile(fileName, data string) {
	fmt.Println("[log]saving...")
	f, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(data)

	if err != nil {
		fmt.Println(err)
		return
	}
}

//Parses html code using regexp
func parse(data, class, id, elemType string) string {
	fmt.Println("[log]Parsing...")
	if data == "" {
		return "Error: html is empty!"
	}
	regular := "<"
	if elemType != "" {
		regular += elemType
	} else {
		regular += ".+"
	}

	if id != "" {
		regular += " id=\"" + id + "\""

		if class != "" {
			regular += " "
		}
	}
	if class != "" {
		regular += "class=\"" + class + "\""
	}

	regular += ".*>.*</"
	if elemType != "" {
		regular += elemType + ">"
	} else {
		regular += ".+>"
	}
	fmt.Println("[log]Regular: " + regular)
	r, _ := regexp.Compile(regular)
	parsed_out := r.FindAllString(data, -1)

	var parsed_string string
	for _, val := range parsed_out {
		parsed_string += val + "\n"
	}
	return parsed_string
}
