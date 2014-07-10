/*
Curl is Simple http download and readline lib by Golang. Vesion 0.0.1

Website https://github.com/kenshin/curl

Copyright (c) 2014 Kenshin Wang <kenshin@ksria.com>
*/
package curl

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Read line use callback Process
// Line by line to obtain content and line num
type processFunc func(content string, line int) bool

// Get url method
//
//  url e.g. http://nodejs.org/dist/v0.10.0/node.exe
//
// Return code
//   0: success
//  -1: status code != 200
//
// Return res, err
//
// For example:
//  code, res, _ := curl.Get("http://nodejs.org/dist/")
//  if code != 0 {
//      return
//  }
//  defer res.Body.Close()
func Get(url string) (code int, res *http.Response, err error) {

	// get res
	res, err = http.Get(url)

	// err
	if err != nil {
		panic(err)
	}

	// check state code
	if res.StatusCode != 200 {
		fmt.Printf("URL [%v] an [%v] error occurred, please check.\n", url, res.StatusCode)
		return -1, res, err
	}

	return 0, res, err

}

// Read line from io.ReadCloser
//
// For example:
//  versionFunc := func(content string, line int) bool {
//    // TO DO
//    return false
//  }
//
//  if err := curl.ReadLine(res.Body, versionFunc); err != nil && err != io.EOF {
//    //TO DO
//  }
func ReadLine(body io.ReadCloser, process processFunc) error {

	var content string
	var err error
	var line int = 1

	// set buff
	buff := bufio.NewReader(body)

	for {
		content, err = buff.ReadString('\n')

		if line > 1 && (err != nil || err == io.EOF) {
			break
		}

		if ok := process(content, line); ok {
			break
		}

		line++
	}

	return err
}

// Download method
//
// Parameter
//  url : download url e.g. http://nodejs.org/dist/v0.10.0/node.exe
//  name: download file name e.g. node.exe
//  dst : download path
//
// Return code
//   0: success
//  -2: create file error
//  -3: download node.exe error
//  -4: content length = -1
//
// For example:
//  curl.New("http://nodejs.org/dist/", "0.10.28", "v0.10.28")
//
//  Console show
//  Start download [0.10.28] from http://nodejs.org/dist/.
//  1% 5% 10% 20% 30% 40% 50% 60% 70% 80% 90% 100%
//  End download.
func New(url, name, dst string) int {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("CURL Error: Download %v from %v an error has occurred. \nError: %v", name, url, err)
			panic(msg)
		}
	}()

	// get url
	code, res, err := Get(url)
	if code != 0 {
		return code
	}

	// close
	defer res.Body.Close()

	// create file
	file, createErr := os.Create(dst)
	if createErr != nil {
		fmt.Println("Create file error, Error: " + createErr.Error())
		return -2
	}
	defer file.Close()

	if res.ContentLength == -1 {
		fmt.Printf("Download %v fail from %v.\n", name, url)
		return -4
	}

	fmt.Printf("Start download [%v] from %v.\n%v", name, url, "1% ")

	// loop buff to file
	buf := make([]byte, res.ContentLength)
	var m float32
	isShow, oldCurrent := false, 0
	for {
		n, err := res.Body.Read(buf)

		// write complete
		if n == 0 && err.Error() == "EOF" {
			fmt.Println("100% \nEnd download.")
			break
		}

		//error
		if err != nil && err.Error() != "EOF" {
			panic(err)
		}

		m = m + float32(n)
		current := int(m / float32(res.ContentLength) * 100)

		switch {
		case current > 0 && current < 6:
			current = 5
		case current > 5 && current < 11:
			current = 10
		case current > 10 && current < 21:
			current = 20
		case current > 20 && current < 31:
			current = 30
		case current > 30 && current < 41:
			current = 40
		case current > 40 && current < 51:
			current = 50
		case current > 60 && current < 71:
			current = 60
		case current > 70 && current < 81:
			current = 70
		case current > 80 && current < 91:
			current = 80
		case current > 90 && current < 101:
			current = 90
		}

		if current > oldCurrent {
			switch current {
			case 5, 10, 20, 30, 40, 50, 60, 70, 80, 90:
				isShow = true
			}

			if isShow {
				fmt.Printf("%d%v", current, "% ")
			}

			isShow = false
		}

		oldCurrent = current

		file.WriteString(string(buf[:n]))
	}

	// valid download exe
	fi, err := file.Stat()
	if err == nil {
		if fi.Size() != res.ContentLength {
			fmt.Printf("Error: Downlaod [%v] size error, please check your network.\n", name)
			return -3
		}
	}

	return 0
}
