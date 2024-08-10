package belajar_database

import (
	"belajar_database/entity"
	"belajar_database/repository"
	"context"
	"fmt"
	"log"
	"testing"
)

func TestInsertComment(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	repo := repository.New(db)
	comment := entity.Comment{
		Email:   "jgQpQ@example.com",
		Comment: "test comment",
	}
	result, err := repo.Insert(ctx, comment)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)
}

func TestFindByIdComment(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	repo := repository.New(db)
	result, err := repo.FindById(ctx, 87)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)
}

func TestFindByIdNotFoundComment(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	repo := repository.New(db)
	_, err := repo.FindById(ctx, 9999)
	if err == nil {
		log.Println("expected error")
	}
}

func TestFindAllComment(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	repo := repository.New(db)
	result, err := repo.FindAll(ctx)
	if err != nil {
		log.Println(err)
	}
	for _, comment := range result {
		fmt.Println(comment)
	}
}
