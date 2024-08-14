package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func logJson(data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bytes))
}

func TestEncode(t *testing.T) {
	logJson("1")
}

type Customer struct {
	Firstname string
	Lastname  string
	Hobbies   []string
}

func TestObject(t *testing.T) {
	customer := Customer{
		Firstname: "John",
		Lastname:  "Doe",
		Hobbies:   []string{"Reading", "Music"},
	}
	logJson(customer)
}

func decJson(data []byte) {
	var customer = &Customer{}
	err := json.Unmarshal(data, customer)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(customer.Firstname)
	fmt.Println(customer.Lastname)
	fmt.Println(customer.Hobbies)
}

func TestDecode(t *testing.T) {
	json := `{"Firstname":"John","Lastname":"Doe","Hobbies":["Reading","Music"]}`
	decJson([]byte(json))
}

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func TestObjectWithTag(t *testing.T) {
	product := Product{
		Name:  "Book",
		Price: 9.99,
	}
	logJson(product)
}

func TestMap(t *testing.T) {
	data := `{"products": [{"name":"Book","price":9.99},{"name":"Pen","price":1.99}]}`
	var products map[string]interface{}
	err := json.Unmarshal([]byte(data), &products)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range products["products"].([]interface{}) {
		product := v.(map[string]interface{})
		fmt.Println(product["name"])
		fmt.Println(product["price"])
	}
}

func TestDecoder(t *testing.T) {
	data, _ := os.Open("test.json")
	dec := json.NewDecoder(data)
	customer := &Customer{}
	err := dec.Decode(customer)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(customer.Firstname)
	fmt.Println(customer.Lastname)
	fmt.Println(customer.Hobbies)
}
