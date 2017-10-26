// +build darwin windows
// +build amd64

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var lineNum = uint64(0)
	var fileNum = uint64(0)
	var line string
	var maxLines uint64
	var outFile *os.File
	var headerLine string
	var firstLine = bool(true)
	var lineSeparator = "\n" // we assume this is running under Windows OS
	flag.Uint64Var(&maxLines, "maxLines", uint64(49000), "How many lines before we split the file")
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("current dir = " + dir)
	log.Println(fmt.Sprintf("maxLines=%d", maxLines))
	f, err := os.OpenFile("splitFile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.Println("Logfile=" + f.Name())

	log.SetOutput(f)

	var inputDir = filepath.Join(dir, "input")
	var outputDir = filepath.Join(dir, "output")
	var archiveDir = filepath.Join(dir, "archive")

	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		os.Mkdir(inputDir, os.ModePerm)
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	if _, err := os.Stat(archiveDir); os.IsNotExist(err) {
		os.Mkdir(archiveDir, os.ModePerm)
	}

	files, err := ioutil.ReadDir(inputDir)

	for _, singleFileInfo := range files {
		if singleFileInfo.IsDir() {
			continue
		}
		s := singleFileInfo.Name()
		lineNum = 0
		fileNum = 0
		// open a file
		notice("Processing file " + s)
		if file, err := os.Open(filepath.Join(inputDir, s)); err == nil {

			// make sure it gets closed
			defer file.Close()

			// create a new scanner and read the file line by line
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				if lineNum%maxLines == 0 {
					name := strings.TrimSuffix(s, filepath.Ext(s))
					suffix := filepath.Ext(s)
					fileName := filepath.Join(outputDir, fmt.Sprintf("%s_%02d%s", name, fileNum, suffix))
					if outFile != nil {
						outFile.Sync()
						outFile.Close()
					}
					f, err := os.Create(fileName)

					if err != nil {
						log.Fatal(err)
						os.Exit(1)
					}
					outFile = f
					notice("Split to : " + fileName)
					if fileNum > 0 {
						outFile.WriteString(headerLine + lineSeparator)
					}
					fileNum++
				}

				line = scanner.Text()
				if firstLine {
					headerLine = line
					firstLine = false
				}

				// log.Println(line)
				outFile.WriteString(line + lineSeparator)
				lineNum++
			}

			// check for errors
			if err = scanner.Err(); err != nil {
				log.Fatal(err)
			}

			if outFile != nil {
				outFile.Sync()
				outFile.Close()
			}

			file.Close()
			os.Rename(filepath.Join(inputDir, s), filepath.Join(archiveDir, s))

		} else {
			log.Fatal(err)
		}
	}
}

func notice(message string) {
	log.Println(message)
	fmt.Println(message)
}
