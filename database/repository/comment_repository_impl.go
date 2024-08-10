package repository

import (
	"belajar_database/entity"
	"context"
	"database/sql"
	"errors"
	"log"
)

type commentRepositoryImpl struct {
	DB *sql.DB
}

func New(db *sql.DB) CommentRepository {
	return &commentRepositoryImpl{DB: db}
}

func (repository *commentRepositoryImpl) Insert(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	result, err := repository.DB.ExecContext(ctx, query, comment.Email, comment.Comment)
	if err != nil {
		return comment, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return comment, err
	}
	comment.Id = int(id)
	return comment, nil
}

func (repository *commentRepositoryImpl) FindById(ctx context.Context, id int) (entity.Comment, error) {
	query := "SELECT id, email, comment FROM comments WHERE id = ? LIMIT 1"

	result, err := repository.DB.QueryContext(ctx, query, id)
	comment := entity.Comment{}
	if err != nil {
		return comment, err
	}
	defer result.Close()
	if result.Next() {
		err := result.Scan(&comment.Id, &comment.Email, &comment.Comment)
		if err != nil {
			return comment, err
		}
		return comment, nil
	} else {
		return entity.Comment{}, errors.New("comment not found")
	}
}

func (repository *commentRepositoryImpl) FindAll(ctx context.Context) ([]entity.Comment, error) {
	query := "SELECT id, email, comment FROM comments"

	result, err := repository.DB.QueryContext(ctx, query)
	var comments []entity.Comment
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		comment := entity.Comment{}
		err2 := result.Scan(&comment.Id, &comment.Email, &comment.Comment)
		if err2 != nil {
			log.Println(err2)
		}
		comments = append(comments, comment)
	}
	return comments, err
}
