package bookstore_test

import (
	"bookstore"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuyBook(t *testing.T) {
	t.Run("Buys book and decreases copies by 1 on success", func(t *testing.T) {
		book := bookstore.Book{1, "Title", "Author", 5}
		boughtBook, err := bookstore.Buy(book)
		if err != nil {
			t.Errorf("Expected no error. Got %v", err)
		}
		if boughtBook.Copies != 4 {
			t.Errorf("Expected book copies to decrease by 1. Got %v", boughtBook.Copies)
		}
	})
	t.Run("Returns error when book has 0 copies", func(t *testing.T) {
		book := bookstore.Book{1, "Title", "Author", 0}
		_, err := bookstore.Buy(book)
		if err == nil {
			t.Error("Expected error but didn't get one")
		}
	})
}

func TestGetAllBook(t *testing.T) {
	wants := map[int]bookstore.Book{
		1: {1, "Title 1", "Author 1", 5},
		2: {2, "Title 2", "Author 2", 10},
	}

	gotBooks := bookstore.GetAllBooks()

	if !cmp.Equal(wants, gotBooks) {
		t.Errorf(cmp.Diff(wants, gotBooks))
	}
}

func TestGetBook(t *testing.T) {
	t.Run("Returns book if found in catalogue", func(t *testing.T) {
		catalogue := map[int]bookstore.Book{
			1: {1, "Book Title", "Author", 2},
			2: {2, "Book Title", "Author", 2},
		}
		gotBook, err := bookstore.GetBook(catalogue, 1)
		if err != nil {
			t.Errorf("Expected no error. Got %v", err)
		}
		want := catalogue[1]

		if !cmp.Equal(want, gotBook) {
			t.Errorf("want %v got %v", want, gotBook)
		}
	})

	t.Run("Returns error if book not found in catalogue", func(t *testing.T) {
		catalogue := map[int]bookstore.Book{
			1: {1, "Book Title", "Author", 2},
			2: {2, "Book Title", "Author", 2},
		}
		_, err := bookstore.GetBook(catalogue, 3)
		if err == nil {
			t.Error("Expected error but didn't get one")
		}
	})
}
