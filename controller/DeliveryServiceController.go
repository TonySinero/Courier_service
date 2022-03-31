package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"stlab.itechart-group.com/go/food_delivery/courier_service/dao"
	"strconv"
)

// CreateDeliveryService godoc
// @Summary CreateDeliveryService
// @Security ApiKeyAuth
// @Description create a Delivery Service
// @Tags DeliveryService
// @Accept  json
// @Produce  json
// @Param input body dao.DeliveryService true "Delivery Service"
// @Success 200 {object} dao.DeliveryService
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {string} string
// @Router /deliveryservice [post]
func (h *Handler) CreateDeliveryService(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler CreateDeliveryService:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	var service dao.DeliveryService
	if err := ctx.ShouldBindJSON(&service); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	if service.Email == "" || service.Name == "" {
		log.Println(errors.New("empty fields"))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "empty fields"})
		return
	}
	idService, err := h.services.CreateDeliveryService(service)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error: %s", err)})
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"id": idService})
}

// GetDeliveryServiceById godoc
// @Summary GetDeliveryServiceById
// @Security ApiKeyAuth
// @Description get delivery service by id
// @Tags DeliveryService
// @Produce  json
// @Param id path int true "id"
// @Success 200 {object} dao.DeliveryService
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /deliveryservice/{id} [get]
func (h *Handler) GetDeliveryServiceById(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler GetDeliveryServiceById:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "expect an integer greater than 0"})
		return
	}

	service, err := h.services.GetDeliveryServiceById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error: %s", err)})
		return
	}
	ctx.JSON(http.StatusOK, service)
}

type listDeliveryServices struct {
	Data []dao.DeliveryService `json:"data"`
}

// GetAllDeliveryServices godoc
// @Summary GetAllDeliveryServices
// @Security ApiKeyAuth
// @Description get list of all delivery service
// @Tags DeliveryService
// @Produce  json
// @Success 200 {object} listDeliveryServices
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /deliveryservice [get]
func (h *Handler) GetAllDeliveryServices(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler GetAllDeliveryServices:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	services, err := h.services.GetAllDeliveryServices()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error: %s", err)})
		return
	}
	ctx.JSON(http.StatusOK, listDeliveryServices{Data: services})
}

// UpdateDeliveryService godoc
// @Summary UpdateDeliveryService
// @Security ApiKeyAuth
// @Description update delivery service information
// @Tags DeliveryService
// @Accept  json
// @Produce  json
// @Param id path int true "order_id"
// @Param input body dao.DeliveryService true "delivery service"
// @Success 204
// @Failure 400 {string} string
// @Router /deliveryservice/{id} [put]
func (h *Handler) UpdateDeliveryService(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler UpdateDeliveryService:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Error: %s", err)})
		return
	}
	var service dao.DeliveryService
	if err := ctx.ShouldBindJSON(&service); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	service.Id = id
	if err := h.services.UpdateDeliveryService(service); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error: %s", err)})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// SaveLogoController godoc
// @Summary SaveLogoController
// @Security ApiKeyAuth
// @Description set logo to DO Spaces, and it's way to DB
// @Tags DeliveryService
// @Accept  image/jpeg
// @Produce  json
// @Param id query int true "id delivery service"
// @Param logo  formData  file  true  "logo image"
// @Success 204
// @Failure 400 {string} string
// @Router /deliveryservice/logo [post]
func (h *Handler) SaveLogoController(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler SaveLogoController:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	id, er := strconv.Atoi(ctx.Query("id"))
	if er != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "expect an integer greater than 0"})
		return
	}
	if ctx.Request.Body == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "empty"})
	}
	defer ctx.Request.Body.Close()
	cover, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
	}
	if err := h.services.SaveLogoFile(cover, id); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error: %s", err)})
		return
	}
	ctx.Status(http.StatusNoContent)
}
