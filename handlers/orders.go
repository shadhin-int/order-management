package handlers

import (
	"log"
	"net/http"
	"order-management/config"
	"order-management/models"
	"order-management/utils"
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
