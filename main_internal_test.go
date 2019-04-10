package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"testing"
)

func TestMain(t *testing.T) {

	dir := "test"
	url := "http://localhost:8000/backgrounds.html"

	if !checkForFile(dir) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`(https:\\\/\\\/lh)([0-9])(.googleusercontent.com\\\/)(proxy\\\/)(([A-Za-z0-9_\\-])*)(\\u003d)(([A-Za-z0-9_\\-])*)(\\x22,)`)

	s := re.FindAllStringSubmatch(string(body), -1)
	fmt.Printf("%d images found\n", len(s))
	dir = string(append([]byte(dir), "/"...))

	for _, match := range s {

		url := []byte(match[0])
		url = bytes.TrimSuffix(url, []byte("\\x22,"))
		url = unescape(url)
		url = safeEncoding(url)

		filename := []byte(match[5])
		filename = unescape(filename)
		filename = safeEncoding(filename)

		filename = append([]byte(dir), []byte(filename)...)
		filename = append([]byte(filename), []byte(".jpg")...)

		if !(checkForFile(string(filename))) {
			err := download(string(url), string(filename))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	os.Exit(e)

}
