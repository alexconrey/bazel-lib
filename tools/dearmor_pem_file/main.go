package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
)

type Key struct {
	ArmoredContents string
	Contents        []byte
	File            *os.File
}

func (k *Key) WriteFile(fp *os.File) error {
	k.File = fp

	if len(k.Contents) == 0 {
		err := errors.New("no key contents loaded")
		return err
	}

	_, err := k.File.Write([]byte(k.Contents))
	return err

}

func loadKeysFromFile(filePath string) ([]Key, error) {
	var keys []Key
	data, err := os.ReadFile(filePath)
	if err != nil {
		return keys, err
	}

	dataStr := string(data)
	armoredContents := dataStr

	// Trim:
	// * the known PEM start/end blocks from the input
	trimLines := []string{
		"^-----BEGIN .*-----\n(.*)\n",
		"-----END .*-----",
	}

	for i := range trimLines {
		line := trimLines[i]
		re := regexp.MustCompile(line)
		dataStr = re.ReplaceAllString(dataStr, "")
	}

	// If a checksum was provided in the file contents, parse it out
	checksumLine := `\n=(?P<checksum>.+)\n`
	re := regexp.MustCompile(checksumLine)
	if re.SubexpIndex("checksum") != -1 {
		matches := re.FindStringSubmatch(dataStr)
		if len(matches) > 0 {
			dataStr = re.ReplaceAllString(dataStr, "")
		}
	}

	var dearmoredData []byte

	dearmoredData, err = base64.RawStdEncoding.DecodeString(dataStr)
	if err != nil {
		// If the length of the dearmoredData is 0, there was no valid base64 data found
		// otherwise, if length is greater than 0 - there was data found and we will use it.
		if len(dearmoredData) == 0 {
			return keys, err
		}
	}

	// RFC says that the byte sequence must end a newline, enforce it if missing
	lastByte := dearmoredData[len(dearmoredData)-1]
	if lastByte != 14 {
		dearmoredData = append(dearmoredData, 14)
	}

	key := Key{
		ArmoredContents: armoredContents,
		Contents:        dearmoredData,
	}

	keys = append(keys, key)

	return keys, nil
}

func main() {
	inputFile := flag.String("input-file", "", "armored input file")
	outputFile := flag.String("output-file", "key.gpg", "File name of key to save as on disk")

	flag.Parse()

	if *inputFile == "" {
		err := errors.New("file is required")
		fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
		os.Exit(1)
	}

	keys, err := loadKeysFromFile(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
		os.Exit(1)
	}

	keyFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
		os.Exit(1)
	}
	defer keyFile.Close()

	for i := range keys {
		err := keys[i].WriteFile(keyFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
			os.Exit(1)
		}
	}
}
