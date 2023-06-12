package handler

import (
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/employee"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"

	"net/http"

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
	return func(c *gin.Context) {}
}

func (e *Employee) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		employees, err := e.employeeService.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "warehouse not found")
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
	return func(c *gin.Context) {}
}

func (e *Employee) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
