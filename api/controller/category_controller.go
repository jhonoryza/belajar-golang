package controller

import (
	"api/helper"
	"api/repository"
	"api/request"
	"api/service"
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type CategoryController interface {
	Store(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Destroy(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Show(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Index(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

type CategoryControllerImpl struct {
	catService service.CategoryService
}

func (c *CategoryControllerImpl) Store(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var catRequest = &request.CreateCategoryRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(catRequest)
	helper.PanicIfError(err)

	category := c.catService.Create(r.Context(), catRequest)
	resp := helper.WebResponse{
		Code:    http.StatusCreated,
		Message: "OK",
		Data:    category,
	}
	resp.ToJson(w)
}

func (c *CategoryControllerImpl) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	categoryId := p.ByName("categoryId")
	categoryIdInt, _ := strconv.Atoi(categoryId)

	var catRequest = &request.UpdateCategoryRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(catRequest)
	helper.PanicIfError(err)
	catRequest.Id = categoryIdInt

	category := c.catService.Update(r.Context(), catRequest)
	resp := helper.WebResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    category,
	}
	resp.ToJson(w)
}

func (c *CategoryControllerImpl) Destroy(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	categoryId := p.ByName("categoryId")
	categoryIdInt, _ := strconv.Atoi(categoryId)

	c.catService.Delete(r.Context(), &categoryIdInt)
	resp := helper.WebResponse{
		Code:    http.StatusNoContent,
		Message: "OK",
		Data:    nil,
	}
	resp.ToJson(w)
}

func (c *CategoryControllerImpl) Show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	categoryId := p.ByName("categoryId")
	categoryIdInt, _ := strconv.Atoi(categoryId)

	category := c.catService.FindById(r.Context(), &categoryIdInt)
	resp := helper.WebResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    category,
	}
	resp.ToJson(w)
}

func (c *CategoryControllerImpl) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	categories := c.catService.FindAll(r.Context())
	resp := helper.WebResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    categories,
	}
	resp.ToJson(w)
}

func NewCategoryController(db *sql.DB) CategoryController {
	catRepo := repository.NewCategoryRepository()
	validate := validator.New()
	catService := service.NewCategoryService(catRepo, db, validate)

	return &CategoryControllerImpl{
		catService: catService,
	}
}
