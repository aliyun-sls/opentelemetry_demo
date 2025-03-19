package model

// Pay 支付
type Pay struct {
	Model
	UserIdType
	OrderIdType
	CardIdType
	TotalAmount float64 `json:"total_amount" gorm:"decimal(10,2);not null;comment:'支付金额'"`
	Card        *Card   `json:"card" gorm:"foreignKey:CardId;references:PayType"`
}
