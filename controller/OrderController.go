package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"stlab.itechart-group.com/go/food_delivery/courier_service/dao"
	"strconv"
)

type listOrders struct {
	Data []dao.DetailedOrder `json:"data"`
}

type text struct {
	Status string `json:"status"`
}

// GetOrders godoc
// @Summary GetOrders
// @Security ApiKeyAuth
// @Description get orders by courier ID
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param id path int true "Courier ID"
// @Success 200 {object} dao.Order
// @Failure 400 {string} string
// @Failure 500 {string} err
// @Router /orders/{id} [get]
func (h *Handler) GetOrders(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler GetOrders:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	var Orders []dao.Order
	idQuery := ctx.Param("id")
	id, err := strconv.Atoi(idQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	Orders, err = h.services.GetOrders(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"No such orders": err})
		return
	}
	ctx.JSON(http.StatusOK, Orders)
}

// GetOrder godoc
// @Summary GetOrder
// @Description check auth information
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} dao.Order
// @Failure 400 {string} string
// @Failure 500 {string} err
// @Router /order/{id} [get]
func (h *Handler) GetOrder(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler GetOrder:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	var Order dao.Order
	idQuery := ctx.Param("id")
	id, err := strconv.Atoi(idQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error with query parameter": err})
		return
	}
	Order, err = h.services.GetOrder(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"No such order": err})
		return
	}
	ctx.JSON(http.StatusOK, Order)
}

// ChangeOrderStatus godoc
// @Summary ChangeOrderStatus
// @Security ApiKeyAuth
// @Description put order status by order ID
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Param input body string true "status"
// @Success 200 {object} dao.Order
// @Failure 400 {string} string
// @Failure 500 {string} err
// @Router /order/status_change/{id} [put]
func (h *Handler) ChangeOrderStatus(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler ChangeOrderStatus:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	idQuery := ctx.Param("id")
	var txt text
	var status string

	if err := ctx.ShouldBindJSON(&txt); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	status = txt.Status

	id, err := strconv.Atoi(idQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error with query parameter": err})
		return
	}
	orderId, err := h.services.ChangeOrderStatus(status, uint16(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"No such order": err})
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"Order id": orderId})
}

// GetCourierCompletedOrders godoc
// @Summary GetCourierCompletedOrders
// @Security ApiKeyAuth
// @Description get list of completed orders by courier id
// @Tags Orders
// @Produce json
// @Param page query int true "page"
// @Param limit query int true "limit"
// @Param idcourier query int true "idcourier"
// @Success 200 {object} listOrders
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /orders/completed [get]
func (h *Handler) GetCourierCompletedOrders(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler GetCourierCompletedOrders:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	page, er := strconv.Atoi(ctx.Query("page"))
	if er != nil || page == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "page query param is wrong. Expected an integer greater than 0"})
		return
	}
	limit, er1 := strconv.Atoi(ctx.Query("limit"))
	if er1 != nil || limit == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "limit query param is wrong. Expected an integer greater than 0"})
		return
	}
	idCourier, er := strconv.Atoi(ctx.Query("idcourier"))
	if er != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "expect an integer greater than 0"})
		return
	}

	DetOrders, err := h.services.GetCourierCompletedOrders(limit, page, idCourier)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Error: %s", err)})
		return
	}
	ctx.JSON(http.StatusOK, listOrders{Data: DetOrders})

}

type listShortOrders struct {
	Data []dao.Order `json:"data"`
}
type listDetailedOrders struct {
	Data []dao.DetailedOrder `json:"data"`
}

// GetAllOrdersOfCourierService godoc
// @Summary GetAllOrdersOfCourierService
// @Security ApiKeyAuth
// @Description get list of all orders by courier service id
// @Tags Orders
// @Produce json
// @Param page query int true "page"
// @Param limit query int true "limit"
// @Param iddeliveryservice query int true "iddeliveryservice"
// @Success 200 {object} listDetailedOrders
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /orders [get]
func (h *Handler) GetAllOrdersOfCourierService(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler GetAllOrdersOfCourierService:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	page, er := strconv.Atoi(ctx.Query("page"))
	if er != nil || page == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "page query param is wrong. Expected an integer greater than 0"})
		return
	}
	limit, er1 := strconv.Atoi(ctx.Query("limit"))
	if er1 != nil || limit == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "limit query param is wrong. Expected an integer greater than 0"})
		return
	}
	idService, er := strconv.Atoi(ctx.Query("iddeliveryservice"))
	if er != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "expect an integer greater than 0"})
		return
	}

	Orders, err := h.services.GetAllOrdersOfCourierService(limit, page, idService)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Error: %s", err)})
		return
	}
	ctx.JSON(http.StatusOK, listDetailedOrders{Data: Orders})

}

// GetCourierCompletedOrdersByMonth godoc
// @Summary GetCourierCompletedOrdersByMonth
// @Security ApiKeyAuth
// @Description get list of completed orders by courier id sorted by month
// @Tags Orders
// @Produce json
// @Param page query int true "page"
// @Param limit query int true "limit"
// @Param idcourier query int true "idcourier"
// @Param month query int true "month"
// @Param year query int true "year"
// @Success 200 {object} listShortOrders
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /orders/bymonth [get]
func (h *Handler) GetCourierCompletedOrdersByMonth(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler GetCourierCompletedOrdersByMonth:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	page, er := strconv.Atoi(ctx.Query("page"))
	if er != nil || page == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "page query param is wrong. Expected an integer greater than 0"})
		return
	}
	limit, er1 := strconv.Atoi(ctx.Query("limit"))
	if er1 != nil || limit == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": " limit query param is wrong. Expected an integer greater than 0"})
		return
	}
	idCourier, er := strconv.Atoi(ctx.Query("idcourier"))
	if er != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "expect an integer greater than 0"})
		return
	}
	Month, er := strconv.Atoi(ctx.Query("month"))
	if er != nil || Month == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": " expect an integer from 1 to 12"})
		return
	}
	Year, er := strconv.Atoi(ctx.Query("year"))
	if er != nil || Month == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "expect an integer greater than 2021"})
		return
	}
	Orders, err := h.services.GetCourierCompletedOrdersByMonth(limit, page, idCourier, Month, Year)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Error: %s", err)})
		return
	}
	ctx.JSON(http.StatusOK, listShortOrders{Data: Orders})

}

// UpdateOrder godoc
// @Summary UpdateOrder
// @Security ApiKeyAuth
// @Tags Orders
// @Description assign order to courier
// @ID UpdateOrder
// @Accept  json
// @Produce json
// @Param id path int true "order_id"
// @Param input body dao.Order true "id courier"
// @Success 204
// @Failure 400 {string} string
// @Router /orders/{id} [put]
func (h *Handler) UpdateOrder(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler UpdateOrder:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Error: %s", err)})
		return
	}
	var order dao.Order
	if err := ctx.ShouldBindJSON(&order); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	order.Id = id
	if err := h.services.AssigningOrderToCourier(order); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Error: %s", err)})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// GetDetailedOrderById godoc
// @Summary GetDetailedOrderById
// @Security ApiKeyAuth
// @Description get detailed order by id
// @Tags order
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} dao.AllInfoAboutOrder
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /order/detailed/{id} [get]
func (h *Handler) GetDetailedOrderById(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler GetDetailedOrderById:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "expect an integer greater than 0"})
		return
	}

	DetOrder, err := h.services.GetDetailedOrderById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error: %s", err)})
		return
	}
	ctx.JSON(http.StatusOK, DetOrder)
}

// GetCompletedOrdersOfCourierService godoc
// @Summary GetCompletedOrdersOfCourierService
// @Security ApiKeyAuth
// @Description get list of completed orders by courier service id
// @Tags order
// @Produce json
// @Param limit query int true "limit"
// @Param page query int true "page"
// @Param iddeliveryservice query int true "iddeliveryservice"
// @Param sort query string false "sort"
// @Success 200 {object} listShortOrders
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /orders/service/completed [get]
func (h *Handler) GetCompletedOrdersOfCourierService(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler GetCompletedOrdersOfCourierService:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	page, er := strconv.Atoi(ctx.Query("page"))
	if er != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "page query param is wrong. Expected an integer greater than 0"})
		return
	}
	limit, er1 := strconv.Atoi(ctx.Query("limit"))
	if er1 != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "limit query param is wrong. Expected an integer greater than 0"})
		return
	}
	idService, er := strconv.Atoi(ctx.Query("iddeliveryservice"))
	if er != nil || idService <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "expect an integer greater than 0"})
		return
	}
	Sort := ctx.Query("sort")
	if Sort == "date" {
		Orders, err := h.services.GetCompletedOrdersOfCourierServiceByDate(limit, page, idService)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error: %s", err)})
			return
		}
		ctx.JSON(http.StatusOK, listShortOrders{Data: Orders})
	}
	if Sort == "courier" {
		Orders, err := h.services.GetCompletedOrdersOfCourierServiceByCourierId(limit, page, idService)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error: %s", err)})
			return
		}
		ctx.JSON(http.StatusOK, listShortOrders{Data: Orders})
	} else {
		Orders, err := h.services.GetCompletedOrdersOfCourierService(limit, page, idService)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error: %s", err)})
			return
		}
		ctx.JSON(http.StatusOK, listShortOrders{Data: Orders})
	}
}

// GetOrdersOfCourierServiceForManager godoc
// @Summary GetOrdersOfCourierServiceForManager
// @Security ApiKeyAuth
// @Description get list of all orders by courier service id with custom status
// @Tags order
// @Produce json
// @Param page query int true "page"
// @Param limit query int true "limit"
// @Param iddeliveryservice query int true "iddeliveryservice"
// @Success 200 {object} listDetailedOrders
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /orders/manager [get]
func (h *Handler) GetOrdersOfCourierServiceForManager(ctx *gin.Context) {
	necessaryRole := []string{"Superadmin", "Courier manager"}
	if err := h.services.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		log.Println("Handler GetOrdersOfCourierServiceForManager:not enough rights")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not enough rights"})
		return
	}
	page, er := strconv.Atoi(ctx.Query("page"))
	if er != nil || page == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "page query param is wrong. Expected an integer greater than 0"})
		return
	}
	limit, er1 := strconv.Atoi(ctx.Query("limit"))
	if er1 != nil || limit == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "limit query param is wrong. Expected an integer greater than 0"})
		return
	}
	idService, er := strconv.Atoi(ctx.Query("iddeliveryservice"))
	if er != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "expect an integer greater than 0"})
		return
	}

	Orders, err := h.services.GetOrdersOfCourierServiceForManager(limit, page, idService)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error: %s", err)})
		return
	}
	ctx.JSON(http.StatusOK, listDetailedOrders{Data: Orders})

}
