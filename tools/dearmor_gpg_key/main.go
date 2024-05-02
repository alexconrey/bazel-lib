package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
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

	trimLines := []string{
		"-----BEGIN PGP PUBLIC KEY BLOCK-----",
		"-----END PGP PUBLIC KEY BLOCK-----",
	}

	dataStr = strings.TrimSpace(dataStr)
	dataStr = strings.Replace(dataStr, "\n", "", -1)

	for i := range trimLines {
		line := trimLines[i]
		dataStr = strings.TrimPrefix(dataStr, line)
		dataStr = strings.TrimSuffix(dataStr, line)
	}

	dataStr = dataStr[:len(dataStr)-5]

	dearmoredData, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		return keys, err
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
