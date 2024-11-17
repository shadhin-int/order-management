package handlers

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"net/http"
	"order-management/config"
	"order-management/models"
	"order-management/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var order models.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.APIResponse{
			Message: "Please fix the given errors",
			Type:    "error",
			Code:    422,
			Errors:  utils.GetValidationErrors(err),
		})
		return
	}

	if !utils.ValidatePhoneNumber(order.RecipientPhone) {
		c.JSON(http.StatusUnprocessableEntity, models.APIResponse{
			Message: "Please fix the given errors",
			Type:    "error",
			Code:    422,
			Errors: map[string][]string{
				"recipient_phone": {"Invalid phone number format"},
			},
		})
		return
	}

	consignmentID := utils.GenerateConsignmentID()

	deliveryFee := 60
	codFee := 12
	totalFee := deliveryFee + codFee

	orderObj := models.Order{
		StoreID:            order.StoreID,
		MerchantOrderID:    order.MerchantOrderID,
		RecipientName:      order.RecipientName,
		RecipientPhone:     order.RecipientPhone,
		RecipientAddress:   "banani, gulshan 2, dhaka, bangladesh",
		RecipientCity:      1,
		RecipientZone:      1,
		RecipientArea:      1,
		DeliveryType:       48,
		ItemType:           2,
		SpecialInstruction: order.SpecialInstruction,
		ItemQuantity:       order.ItemQuantity,
		ItemWeight:         order.ItemWeight,
		AmountToCollect:    order.AmountToCollect,
		ItemDescription:    order.ItemDescription,
		OrderStatus:        models.OrderStatusPending,
		OrderConsignmentID: consignmentID,
		OrderCreatedAt:     time.Now(),
		OrderTypeID:        1,
		OrderType:          "Regular",
		TotalFee:           float64(totalFee),
		CODFee:             float64(codFee),
		DeliveryFee:        float64(deliveryFee),
		PromoDiscount:      0,
		Discount:           0,
	}
	if config.DB == nil {
		log.Fatal("Database connection is not initialized")
	}
	orderCreated := config.DB.Create(&orderObj)
	if orderCreated.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Message: "Failed to create order",
			Type:    "error",
			Code:    500,
		})

		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Message: "Order created successfully",
		Type:    "success",
		Code:    200,
		Data: map[string]interface{}{
			"consignment_id":    consignmentID,
			"merchant_order_id": orderObj.MerchantOrderID,
			"order_status":      orderObj.OrderStatus,
			"delivery_fee":      deliveryFee,
		},
	})
}

func GetAllOrders(c *gin.Context) {
	orderStatus := c.DefaultQuery("order_status", models.OrderStatusPending)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var orders []models.Order
	var total int64

	query := config.DB.Model(&models.Order{}).Where("LOWER(order_status::text) = LOWER(?)", orderStatus)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("order_created_at DESC").Find(&orders).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Message: "Failed to fetch orders",
			Type:    "error",
			Code:    500,
		})
		return
	}

	lastPage := (int(total) + limit - 1) / limit
	totalInPage := len(orders)
	var orderData []map[string]interface{}

	for _, order := range orders {
		orderData = append(orderData, map[string]interface{}{
			"order_consignment_id": order.OrderConsignmentID,
			"order_created_at":     order.OrderCreatedAt.Format("2006-01-02 15:04:05"),
			"order_description":    order.ItemDescription,
			"merchant_order_id":    order.MerchantOrderID,
			"recipient_name":       order.RecipientName,
			"recipient_address":    order.RecipientAddress,
			"recipient_phone":      order.RecipientPhone,
			"order_amount":         order.AmountToCollect,
			"total_fee":            order.TotalFee,
			"instruction":          order.SpecialInstruction,
			"order_type_id":        order.OrderTypeID,
			"cod_fee":              order.CODFee,
			"promo_discount":       order.PromoDiscount,
			"discount":             order.Discount,
			"delivery_fee":         order.DeliveryFee,
			"order_status":         order.OrderStatus,
			"order_type":           order.OrderType,
			"item_type":            utils.GetItemType(order.ItemType),
		})
	}

	paginatedData := models.PaginatedResponse{
		Data:        orderData,
		Total:       total,
		CurrentPage: page,
		PerPage:     limit,
		TotalInPage: totalInPage,
		LastPage:    lastPage,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "Orders fetched successfully",
		Type:    "success",
		Code:    200,
		Data:    paginatedData,
	})
}

func CancelOrder(c *gin.Context) {
	consignmentID := c.Param("consignment_id")

	if consignmentID == "" {
		c.JSON(http.StatusUnprocessableEntity, models.APIResponse{
			Message: "Consignment ID is required",
			Type:    "error",
			Code:    422,
		})
		return
	}

	var order models.Order
	err := config.DB.Where("order_consignment_id = ?", consignmentID).First(&order).Error

	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Message: "Order not found",
			Type:    "error",
			Code:    404,
		})
		return
	}

	if order.OrderStatus == models.OrderStatusCancelled {
		c.JSON(http.StatusUnprocessableEntity, models.APIResponse{
			Message: "Order already cancelled",
			Type:    "error",
			Code:    422,
		})
		return
	}

	if order.OrderStatus != models.OrderStatusPending {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "Please contact cx to cancel the order",
			Type:    "error",
			Code:    400,
		})
		return
	}

	order.OrderStatus = models.OrderStatusCancelled
	err = config.DB.Save(&order).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Message: "Failed to cancel order",
			Type:    "error",
			Code:    500,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "Order cancelled successfully",
		Type:    "success",
		Code:    200,
	})
}
