package service

import (
	"api/entity"
	"api/exception"
	"api/helper"
	"api/repository"
	"api/request"
	"api/response"
	"context"
	"database/sql"
	"github.com/go-playground/validator/v10"
)

type CategoryService interface {
	Create(ctx context.Context, req *request.CreateCategoryRequest) *response.CategoryResponse
	Update(ctx context.Context, req *request.UpdateCategoryRequest) *response.CategoryResponse
	Delete(ctx context.Context, categoryId *int)
	FindById(ctx context.Context, categoryId *int) *response.CategoryResponse
	FindAll(ctx context.Context) *[]response.CategoryResponse
}

type CategoryServiceImpl struct {
	catRepo  repository.CategoryRepository
	DB       *sql.DB
	Validate *validator.Validate
}

func (c *CategoryServiceImpl) Create(ctx context.Context, req *request.CreateCategoryRequest) *response.CategoryResponse {
	// validation
	err := c.Validate.Struct(req)
	helper.PanicIfError(err)

	// transaction
	tx, err := c.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// map to model
	model := &entity.Category{
		Name: req.Name,
	}

	// create category
	category := c.catRepo.Create(ctx, tx, model)
	return &response.CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
	}
}

func (c *CategoryServiceImpl) Update(ctx context.Context, req *request.UpdateCategoryRequest) *response.CategoryResponse {
	// validation
	err := c.Validate.Struct(req)
	helper.PanicIfError(err)

	// transaction
	tx, err := c.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// model mapper
	model := &entity.Category{
		Id:   req.Id,
		Name: req.Name,
	}

	// find category
	category, err := c.catRepo.FindById(ctx, tx, &model.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	// update category
	category = c.catRepo.Update(ctx, tx, model)
	return &response.CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
	}
}

func (c *CategoryServiceImpl) Delete(ctx context.Context, categoryId *int) {
	// transaction
	tx, err := c.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// find category
	category, err := c.catRepo.FindById(ctx, tx, categoryId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	// delete category
	c.catRepo.Delete(ctx, tx, category)
}

func (c *CategoryServiceImpl) FindById(ctx context.Context, categoryId *int) *response.CategoryResponse {
	// transaction
	tx, err := c.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// find category
	category, err := c.catRepo.FindById(ctx, tx, categoryId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return &response.CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
	}
}

func (c *CategoryServiceImpl) FindAll(ctx context.Context) *[]response.CategoryResponse {
	// transaction
	tx, err := c.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// get all category
	categories := c.catRepo.FindAll(ctx, tx)
	var catResponses []response.CategoryResponse
	for _, category := range *categories {
		catResponses = append(catResponses, response.CategoryResponse{
			Id:   category.Id,
			Name: category.Name,
		})
	}
	return &catResponses
}

func NewCategoryService(catRepo repository.CategoryRepository, db *sql.DB, validate *validator.Validate) CategoryService {
	return &CategoryServiceImpl{
		catRepo:  catRepo,
		DB:       db,
		Validate: validate,
	}
}
