package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type APIResponse struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type PaginatedResponse struct {
	Data        interface{} `json:"data"`
	Total       int64       `json:"total"`
	CurrentPage int         `json:"current_page"`
	PerPage     int         `json:"per_page"`
	TotalInPage int         `json:"total_in_page"`
	LastPage    int         `json:"last_page"`
}

const (
	OrderStatusPending   = "Pending"
	OrderStatusDelivered = "Delivered"
	OrderStatusCancelled = "Cancelled"
)

var validOrderStatuses = []string{
	OrderStatusPending,
	OrderStatusDelivered,
	OrderStatusCancelled,
}

type Order struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	StoreID            int       `gorm:"column:store_id;not null;" json:"store_id" binding:"required"`
	MerchantOrderID    string    `gorm:"column:merchant_order_id;size:100" json:"merchant_order_id,omitempty"`
	RecipientName      string    `gorm:"column:recipient_name;size:255;not null" json:"recipient_name" binding:"required"`
	RecipientPhone     string    `gorm:"column:recipient_phone;size:20;not null" json:"recipient_phone" binding:"required"`
	RecipientAddress   string    `gorm:"column:recipient_address;size:500;not null" json:"recipient_address" binding:"required"`
	RecipientCity      int       `gorm:"column:recipient_city;not null;default:1" json:"recipient_city"`
	RecipientZone      int       `gorm:"column:recipient_zone;not null;default:1" json:"recipient_zone"`
	RecipientArea      int       `gorm:"column:recipient_area;not null;default:1" json:"recipient_area"`
	DeliveryType       int       `gorm:"column:delivery_type;not null;default:48" json:"delivery_type"`
	ItemType           int       `gorm:"column:item_type;not null;default:2" json:"item_type"`
	SpecialInstruction string    `gorm:"column:special_instruction;size:255" json:"special_instruction,omitempty"`
	ItemQuantity       int       `gorm:"column:item_quantity;not null;default:1" json:"item_quantity"`
	ItemWeight         float64   `gorm:"column:item_weight;not null;default:0.5" json:"item_weight"`
	AmountToCollect    float64   `gorm:"column:amount_to_collect;not null" json:"amount_to_collect" binding:"required"`
	ItemDescription    string    `gorm:"column:item_description;size:255" json:"item_description,omitempty"`
	OrderStatus        string    `gorm:"column:order_status;size:50;not null" json:"order_status"`
	OrderConsignmentID string    `gorm:"column:order_consignment_id;size:50" json:"order_consignment_id"`
	OrderCreatedAt     time.Time `gorm:"column:order_created_at;not null" json:"order_created_at"`
	OrderTypeID        int       `gorm:"column:order_type_id;not null" json:"order_type_id"`
	OrderType          string    `gorm:"column:order_type;size:50;not null" json:"order_type"`
	TotalFee           float64   `gorm:"column:total_fee;not null" json:"total_fee"`
	CODFee             float64   `gorm:"column:cod_fee;not null" json:"cod_fee"`
	PromoDiscount      float64   `gorm:"column:promo_discount;not null" json:"promo_discount"`
	Discount           float64   `gorm:"column:discount;not null" json:"discount"`
	DeliveryFee        float64   `gorm:"column:delivery_fee;not null" json:"delivery_fee"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Order) TableName() string {
	return "orders"
}

func ValidateOrderStatus(status string) error {
	for _, validStatus := range validOrderStatuses {
		if status == validStatus {
			return nil
		}
	}
	return errors.New("invalid order status")
}

func (o *Order) BeforeSave(tx *gorm.DB) (err error) {
	if err := ValidateOrderStatus(o.OrderStatus); err != nil {
		return err
	}
	return nil
}
