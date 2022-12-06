package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var fileMap map[string][]string

func main() {
	fileMap = make(map[string][]string)
	if len(os.Args) >= 2 {
		path := os.Args[1]

		if !strings.HasSuffix(path, "\\") {
			path = path + "\\"
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("ERROR!!! path is not exist")
		} else {
			listfiles := getAllFilesInDir(&path)
			for _, file := range listfiles {
				proccesFile(&file)
			}
			//вставить ожидания конца запуска всех горутин
		}
	} else {
		fmt.Println("ERROR!!! No such argument")
	}
	for key, _ := range fileMap {
		if len(fileMap[key]) > 1 {
			fmt.Println(fileMap[key])
		}
	}
	fmt.Println("main stopped")
}

func proccesFile(file *string) {
	hashSum, err := getSha256FileHashSum(file)
	if err != nil {
		fmt.Println(*file, ":", err)
	} else {
		tempSlice := fileMap[hashSum]
		tempSlice = append(tempSlice, *file)
		fileMap[hashSum] = tempSlice
		//fmt.Println(*file, ":", hashSum)
	}
}

func getSha256FileHashSum(path *string) (string, error) {
	f, err := os.Open(*path)
	if err != nil {
		return "error", err
	}
	defer f.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "error", err
	}
	result := fmt.Sprintf("%x", hash.Sum(nil))
	return result, nil
}

func getAllFilesInDir(path *string) []string {
	var result []string
	files, err := os.ReadDir(*path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fullname := *path + file.Name()
		if !file.IsDir() {
			result = append(result, fullname)
		} else {
			fullname := fullname + "\\"
			result = append(result, getAllFilesInDir(&fullname)...)
		}
	}
	return result
}
