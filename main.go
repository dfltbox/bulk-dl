package main

import (
	"bufio"
	"context"
	"flag"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/lrstanley/go-ytdlp"
)

var urlsFilePath string
var outputPath string
var audioOnly bool
var fileStructure string

func download(chunk []string, wg *sync.WaitGroup) {
	defer wg.Done()
	//fmt.Println(place)
	//log.Debug("I am", "place", place, "I have", chunk)

	for _, link := range chunk {
		log.Info("Attempting to download", "url", link)

		dl := ytdlp.New().
			FormatSort("res,ext:mp4:m4a").
			RecodeVideo("mp4").
			NoPlaylist().
			Progress().
			ProgressTemplate("download-title:%(info.id)s-%(progress.eta)s").
			NoOverwrites().
			Continue().
			Output(outputPath + fileStructure)

		if audioOnly {
			dl = dl.ExtractAudio()
			dl = dl.AudioFormat("mp3")
			dl = dl.UnsetRecodeVideo()
		}

		_, err := dl.Run(context.TODO(), link)
		if err != nil {
			panic(err)
		}
		log.Info("Completed Downloading", "url", link)
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
	urlsFilePath := flag.String("u", "!", "The file containing the URLs")
	outputPath := flag.String("o", "output", "The directory the files will be downloaded to")
	audioOnly := flag.Bool("a", false, "'true' will only download the audio")
	fileStructure := flag.String("fs", "%(extractor)s - %(title)s.%(ext)s", "'true' will only download the audio")
	flag.Parse()

	if *urlsFilePath == "!" {
		log.Fatal("No URL file specified")
		os.Exit(3)
	}

	log.Info("Audio Only?", "audioOnly", *audioOnly)
	log.Info("Output Path", "outputPath", *outputPath)
	log.Info("File Structure", "fileStructure", *fileStructure)
	log.Info("URLs File Path", "filePath", *urlsFilePath)

	log.Info("Downloading yt-dlp if you don't already have it")
	ytdlp.Install(context.Background(), nil)

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

	if len(URLS) == 1 {
		chunkSize = 0
	}

	var wg sync.WaitGroup
	for i := 0; i < chunkSize+1; i++ {
		wg.Add(1)
		go download(chunks[i], &wg)
	}

	wg.Wait()
}
