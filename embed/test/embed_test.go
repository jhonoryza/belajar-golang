package test

import (
	"embed"
	"fmt"
	"os"
	"testing"
)

// harus diluar function
//
//go:embed file.txt
var file string

func TestEmbed(t *testing.T) {
	fmt.Println(file)
}

func TestTanpaEmbed(t *testing.T) {
	f, err := os.ReadFile("file.txt")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(f))
}

//go:embed lara.png
var logo []byte

func TestByte(t *testing.T) {
	err := os.WriteFile("new-lara.png", logo, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

//go:embed file.txt
var files embed.FS

func TestMultipleFile(t *testing.T) {
	a, _ := files.ReadFile("file.txt")
	fmt.Println(string(a))
}
