package main

import (
	"embed"
	"fmt"
)

//go:embed files/*
var path embed.FS

func main() {
	dirFiles, _ := path.ReadDir("files")

	for _, file := range dirFiles {
		fmt.Println(file.Name())
	}
}
