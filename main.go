package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	"log"

	"github.com/lrstanley/go-ytdlp"
)

func download(chunk []string, wg *sync.WaitGroup) {
	defer wg.Done()
	//fmt.Println(place)
	//log.Debug("I am", "place", place, "I have", chunk)

	for _, link := range chunk {
		//fmt.Println(place, ":", link)
		dl := ytdlp.New().
			FormatSort("res,ext:mp4:m4a").
			RecodeVideo("mp4").
			NoPlaylist().
			Progress().
			ProgressTemplate("download-title:%(info.id)s-%(progress.eta)s").
			NoOverwrites().
			Continue().
			Output("/output/%(extractor)s - %(title)s.%(ext)s")

		_, err := dl.Run(context.TODO(), link)
		if err != nil {
			panic(err)
		}
	}
}

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func main() {
	ytdlp.MustInstall(context.TODO(), nil)

	urlFile, err := os.Open("URLS.txt")
	if err != nil {
		log.Fatalf("error opening file: %s", err)
	}
	defer urlFile.Close()

	scanner := bufio.NewScanner(urlFile)
	var URLS []string

	for scanner.Scan() {
		URLS = append(URLS, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error while reading file: %s", err)
	}

	threads := len(URLS)

	if len(URLS) > 5 {
		threads = 5
	}

	chunkSize := (len(URLS) + threads - 1) / threads
	chunks := chunkSlice(URLS, chunkSize)
	fmt.Println(chunks)
	fmt.Println("chunk size is")
	fmt.Println(chunkSize + 1)

	var wg sync.WaitGroup
	for i := 0; i < chunkSize+1; i++ {
		wg.Add(1)
		go download(chunks[i], &wg)
	}

	wg.Wait()
}
