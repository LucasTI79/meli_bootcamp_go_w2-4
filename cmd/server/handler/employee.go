package handler

import (
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/employee"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Employee struct {
	employeeService employee.Service
}

func NewEmployee(e employee.Service) *Employee {
	return &Employee{
		employeeService: e,
	}
}

func (e *Employee) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, "invalid id")
			return
		}

		employee, err := e.employeeService.Get(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, "invalid id")
			return
		}
		web.Success(c, http.StatusOK, employee)
	}
}

func (e *Employee) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		employees, err := e.employeeService.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "employee not found")
			return
		}
		web.Success(c, http.StatusOK, employees)
	}
}

func (e *Employee) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var employee domain.Employee
		if err := c.ShouldBindJSON(&employee); err != nil {
			web.Error(c, http.StatusBadRequest, "employee not created")
			return
		}
		if employee.CardNumberID == "" {
			web.Error(c, http.StatusBadRequest, "employee card ID need to be only")
			return
		}

		employee, err := e.employeeService.Create(c, employee)
		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, "employee not created")
			return
		}
		web.Success(c, http.StatusCreated, employee)
	}
}

func (e *Employee) Update() gin.HandlerFunc {
	return func(c *gin.Context) {

		var employee domain.Employee
		if err := c.ShouldBindJSON(&employee); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, "action could not be processed correctly due to invalid data provided")
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusNotFound, "invalid id")
			return
		}
		employee.ID = id
		employee, err = e.employeeService.Update(c, employee)
		if err != nil {
			web.Error(c, http.StatusConflict, "employee code id must be unique")
			return
		}
		web.Success(c, http.StatusOK, employee)
	}
}

func (e *Employee) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusNotFound, "invalid id")
			return
		}
		err = e.employeeService.Delete(c, id)
		if err != nil {
			web.Error(c, http.StatusMethodNotAllowed, "employee not deleted")
			return
		}
		web.Success(c, http.StatusNoContent, nil)
	}
}
