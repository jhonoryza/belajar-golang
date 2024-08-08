package bookstore

import (
	"errors"
)

type Book struct {
	ID     int
	Title  string
	Author string
	Copies uint
}

func Buy(b Book) (Book, error) {
	if b.Copies == 0 {
		return Book{}, errors.New(`"no copies left to buy"`)
	}
	b.Copies--
	return b, nil
}

func GetAllBooks() []Book {
	var books = []Book{
		{ID: 1, Title: "Title 1", Author: "Author 1", Copies: 5},
		{ID: 2, Title: "Title 2", Author: "Author 2", Copies: 10},
	}
	return books
}

func GetBook(catalogue map[int]Book, id int) (Book, error) {
	book, ok := catalogue[id]
	if !ok {
		return Book{}, errors.New("book not found")
	}
	return book, nil
}
