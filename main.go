package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

var (
	url         string
	dir         string
	e           = 1
	showVersion = false
	version     = "dev"
)

func readFlags() {
	flag.StringVar(&url, "url", "https://clients3.google.com/cast/chromecast/home", "the chromecast homepage")
	flag.StringVar(&dir, "dir", "", "the directory to download the images to")
	flag.BoolVar(&showVersion, "version", false, "show the version")
	flag.Parse()
}

func main() {

	readFlags()

	if showVersion {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

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

func checkForFile(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	print("\033[32m")
	fmt.Printf("%s exists, skipping\n", filename)
	print("\033[0m")
	return true
}

func download(url string, filename string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("Downloading: %s\n", url)

		// Create the file
		out, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}

		e = 0
	} else {
		fmt.Printf("Error downloading %s, %d", filename, resp.StatusCode)
	}

	return nil
}

func unescape(b []byte) []byte {
	b = bytes.Replace(b, []byte("\\/"), []byte("/"), -1)
	return b
}

func safeEncoding(b []byte) []byte {
	b = bytes.Replace(b, []byte("\\x5b"), []byte("["), -1)
	b = bytes.Replace(b, []byte("\\x5d"), []byte("]"), -1)
	b = bytes.Replace(b, []byte("\\u003d"), []byte("="), -1)
	b = bytes.Replace(b, []byte("\\u2215"), []byte("/"), -1)
	b = bytes.Replace(b, []byte("\\x22"), []byte("\""), -1)
	b = bytes.Replace(b, []byte("\\n"), []byte("\n"), -1)
	b = bytes.Replace(b, []byte("\\"), []byte(""), -1)
	return b
}
