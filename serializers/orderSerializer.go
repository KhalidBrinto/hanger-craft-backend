package serializers

import (
	"time"

	"github.com/lib/pq"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type User struct {
	ID          uint    `gorm:"primarykey"`
	Name        string  `gorm:"size:100;not null"`
	Email       string  `gorm:"size:100;unique;not null"`
	PhoneNumber *string `gorm:"size:15"`
}

type Product struct {
	gorm.Model
	Name        string         `gorm:"size:150;not null"`
	Description string         `gorm:"type:text"`
	SKU         string         `gorm:"size:150;not null;unique;index"`
	Barcode     *string        `gorm:"size:150"`
	Price       float64        `gorm:"not null"`
	Currency    string         `gorm:"size:3; not null"`
	Images      pq.StringArray `gorm:"type:varchar[]"`
}

type OrderItem struct {
	ID              uint          `gorm:"primaryKey"`
	OrderID         uint          `gorm:"not null" json:"-"`
	Order           OrderResponse `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"-"`
	ProductID       uint          `gorm:"not null" json:"-"`
	Product         Product       `gorm:"foreignKey:ProductID"`
	Quantity        int           `gorm:"not null"`
	PriceAtPurchase float64       `gorm:"not null"`
}

type Payment struct {
	ID             uint    `gorm:"primarykey"`
	PaymentMethod  string  `gorm:"size:50;not null;check:payment_method IN ('credit_card', 'paypal', 'bank_transfer', 'cash_on_delivery')"`
	PaymentStatus  string  `gorm:"size:50;not null;check:payment_status IN ('pending', 'completed', 'failed')"`
	Amount         float64 `gorm:"not null"`
	TransanctionID *string `gorm:"size:11;not null"`
	PaymentDate    *time.Time
	OrderID        uint          `gorm:"not null" json:"-"`
	Order          OrderResponse `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"-"`
}

type ShippingAddress struct {
	gorm.Model
	OrderID      uint   `gorm:"not null" json:"-"`
	AddressLine1 string `gorm:"size:255;not null"`
	AddressLine2 string `gorm:"size:255"`
	City         string `gorm:"size:100;not null"`
	State        string `gorm:"size:100"`
	PostalCode   string `gorm:"size:20;not null"`
	Country      string `gorm:"size:100;not null"`

	Order OrderResponse `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"-"`
}

type OrderResponse struct {
	gorm.Model
	OrderIdentifier      string           `gorm:"type:varchar(8); not null;unique;index"`
	UserID               uint             `gorm:"not null" json:"-"`
	User                 User             `gorm:"foreignKey:UserID" json:"Buyer"`
	OrderStatus          string           `gorm:"size:50;not null;check:order_status IN ('pending', 'shipped', 'delivered', 'cancelled')"`
	TotalPrice           float64          `gorm:"not null"`
	OrderItems           []OrderItem      `gorm:"foreignKey:OrderID"`
	OrderShippingAddress *ShippingAddress `gorm:"foreignKey:OrderID"`
	PaymentDetails       *Payment         `gorm:"foreignKey:OrderID"`
}

type ReviewResponse struct {
	gorm.Model
	UserID    uint    `gorm:"not null" json:"-"`
	User      User    `gorm:"foreignKey:UserID"`
	ProductID uint    `gorm:"not null" json:"-"`
	Product   Product `gorm:"foreignKey:ProductID" json:"-"`
	Rating    int     `gorm:"check:rating >= 1 AND rating <= 5"`
	Comment   string  `gorm:"type:text"`
}

type InventoryResponse struct {
	gorm.Model
	ProductID         uint    `gorm:"not null" json:"-"`
	Product           Product `gorm:"foreignKey:ProductID"`
	StockLevel        int     `gorm:"not null"`
	InOpen            int     `gorm:"not null"`
	AvailableQuantity int
	ChangeType        string `gorm:"size:50;not null;check:change_type IN ('restock', 'purchase')"`
	ChangeDate        time.Time
}

type SubCategory struct {
	Name         null.String `binding:"required"`
	CategoryType null.String
	ParentID     *uint
}
type CategoryCreateSerializer struct {
	Name         null.String `binding:"required"`
	CategoryType null.String `binding:"required"`
	ParentID     *uint
	SubCatergory *SubCategory `json:"SubCategory"`
}
