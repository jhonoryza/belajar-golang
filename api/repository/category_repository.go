package repository

import (
	"api/entity"
	"api/helper"
	"context"
	"database/sql"
	"errors"
)

type CategoryRepository interface {
	Create(ctx context.Context, tx *sql.Tx, category *entity.Category) *entity.Category
	Update(ctx context.Context, tx *sql.Tx, category *entity.Category) *entity.Category
	Delete(ctx context.Context, tx *sql.Tx, category *entity.Category)
	FindById(ctx context.Context, tx *sql.Tx, categoryId *int) (*entity.Category, error)
	FindAll(ctx context.Context, tx *sql.Tx) *[]entity.Category
}

type CategoryRepositoryImpl struct{}

func (c *CategoryRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, category *entity.Category) *entity.Category {
	query := "INSERT INTO categories (`name`) VALUES(?)"
	result, err := tx.ExecContext(ctx, query, category.Name)
	helper.PanicIfError(err)

	id, _ := result.LastInsertId()
	category.Id = int(id)
	return category
}

func (c *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category *entity.Category) *entity.Category {
	query := "UPDATE categories SET NAME = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, category.Name, category.Id)
	helper.PanicIfError(err)

	return category
}

func (c *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, category *entity.Category) {
	query := "DELETE FROM categories WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, category.Id)
	helper.PanicIfError(err)
}

func (c *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, categoryId *int) (*entity.Category, error) {
	query := "SELECT * FROM categories WHERE id = ? LIMIT 1"
	rows, err := tx.QueryContext(ctx, query, *categoryId)
	helper.PanicIfError(err)
	defer rows.Close()

	var category entity.Category
	if rows.Next() {
		err = rows.Scan(&category.Id, &category.Name)
		helper.PanicIfError(err)
		return &category, nil
	} else {
		return &category, errors.New("category not found")
	}
}

func (c *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) *[]entity.Category {
	query := "SELECT * FROM categories"
	rows, err := tx.QueryContext(ctx, query)
	helper.PanicIfError(err)
	defer rows.Close()

	var categories []entity.Category

	for rows.Next() {
		category := entity.Category{}
		err = rows.Scan(&category.Id, &category.Name)
		helper.PanicIfError(err)
		categories = append(categories, category)
	}
	return &categories
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
}
