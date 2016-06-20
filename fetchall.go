// fetchall
package main

import (
	//	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)

	for _, url := range os.Args[1:] {
		go fetch(url, ch) //start a goroutine

	}
	for range os.Args[1:] {
		fmt.Println(<-ch) //recieve from ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())

}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	filename := "test.txt"
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	// create file, or append the new data to the existing file.
	//var f *os.File
	_, err2 := os.Stat(filename)
	if err2 != nil {
		f, err1 := os.Create(filename)
		if err1 != nil {
			ch <- fmt.Sprint(err1)
		}
		f.WriteString(time.Now().String() + " file created")
		f.Close()
	}

	f, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		ch <- fmt.Sprint(err)
		fmt.Printf("error at openning the file")
	}

	//get the webpage content to f
	nbytes, err := io.Copy(f, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	f.WriteString("opened")
	//f.Close()

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)

}
