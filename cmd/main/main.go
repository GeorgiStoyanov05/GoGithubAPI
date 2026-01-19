package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func read_file(fileName string) ([]string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	names := strings.Split(string(data), "\n")
	return names, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found (ok if env is set another way)")
	}

	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename as argument.")
		return
	}
	fileName := os.Args[1]
	usernames, err := read_file(fileName)
	if err != nil {
		panic(err)
	}

	for i := range len(usernames) {
		username := usernames[i]
		rep, err := GetUserReport(username)
		if err != nil {
			continue
		}
		data, _ := json.MarshalIndent(rep, "", "  ")
		fmt.Println(string(data))
		fmt.Println("-------------")
	}
}
