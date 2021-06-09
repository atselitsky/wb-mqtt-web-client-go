package rules

import (
	"fmt"
	"log"
	"os"
)

func GetRulesFiles() {
	file, err := os.Open("/home/artem/hello-sql")
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	list, _ := file.Readdirnames(0)
	for _, name := range list {
		fmt.Println(name)
	}
}
