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

// Get obtém as informações de um funcionário pelo ID.
// @Summary Obtém as informações de um funcionário pelo ID
// @Description Retorna as informações de um funcionário com base no ID fornecido
// @Tags employees
// @Accept  json
// @Produce  json
// @Param id path int true "ID do funcionário a ser obtido"
// @Success 200 {object} domain.Employee
// @Failure 400 {string} string "invalid card id"
// @Failure 404 {string} string "invalid id"
// @Router /employees/{id} [get]
func (e *Employee) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, "invalid card id")
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

// GetAll obtém todas as informações dos funcionários.
// @Summary Obtém todas as informações dos funcionários
// @Description Retorna uma lista com todas as informações dos funcionários cadastrados
// @Tags employees
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Employee
// @Failure 400 {string} string "employee not found"
// @Router /employees [get]
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

// Create cria um novo funcionário.
// @Summary Cria um novo funcionário
// @Description Cria um novo funcionário com base nos dados fornecidos
// @Tags employees
// @Accept  json
// @Produce  json
// @Param employee body domain.Employee true "Novo funcionário a ser criado"
// @Success 201 {object} domain.Employee
// @Failure 400 {string} string "employee not created"
// @Failure 422 {string} string "employee card ID need to be only"
// @Router /employees [post]
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

// Update atualiza as informações de um funcionário.
// @Summary Atualiza as informações de um funcionário
// @Description Atualiza as informações de um funcionário com base no ID fornecido e nos dados enviados
// @Tags employees
// @Accept  json
// @Produce  json
// @Param id path int true "ID do funcionário a ser atualizado"
// @Param employee body domain.Employee true "Dados do funcionário a serem atualizados"
// @Success 200 {object} domain.Employee
// @Failure 404 {string} string "action could not be processed correctly due to invalid data provided"
// @Failure 409 {string} string "invalid id"
// @Router /employees/{id} [patch]
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

// Delete remove um funcionário.
// @Summary Remove um funcionário
// @Description Remove um funcionário com base no ID fornecido
// @Tags employees
// @Accept  json
// @Produce  json
// @Param id path int true "ID do funcionário a ser removido"
// @Success 204 "No Content"
// @Failure 404 {string} string "invalid id"
// @Failure 405 {string} string "employee not deleted"
// @Router /employees/{id} [delete]
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
