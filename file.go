package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Copy file from source to destination
func copyFile(src, dest string) {
	bytesRead, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(dest, bytesRead, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Checks if file or directory exists
// Return true if exists, else false
func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getFilesFromPattern(dir, pattern string) []string {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), pattern) {
			files = append(files, e.Name())
		}
	}

	return files
}

func WriteFile(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		log.Println("File " + path + " already exists, delete it.")
		if r := os.Remove(path); r != nil {
			return fmt.Errorf("Unable to remove file.")
		}
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("Unable to write file")
	}

	return nil
}
