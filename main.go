package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
	/*Command line parsing*/
	urlTarget := flag.String("url", "", "parser target url(without http:// or https://)")
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
	fmt.Println("[log]sending get request...")
	url_arr := strings.SplitN(url, "/", 2)
	userAgent := "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:65.0) Gecko/20100101 Firefox/65.0"
	getRequest := "GET /" + url_arr[len(url_arr)-1] + " HTTP/1.1\nHost: "
	getRequest += url_arr[0] + "\nUser-Agent: " + userAgent + "\n\n"

	conn, err := net.Dial("tcp", url_arr[0]+":80")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(getRequest))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var buffer bytes.Buffer
	io.Copy(&buffer, conn)

	return buffer.String()
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
	regular := "(?ims)<"
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
