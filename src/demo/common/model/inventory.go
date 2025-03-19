package model

// Inventory 库存
type Inventory struct {
	InventoryId uint `json:"inventory_id" gorm:"primaryKey"`
	ProductsIdType
	InventoryName    string `json:"inventory_name" gorm:"size:1000;not null;default:'';comment:'商品名'"`
	InventoryNum     int    `json:"inventory_num" form:"inventory_num" gorm:"not null;default:0;comment:'库存数量'"`
	InventoryUnit    int    `json:"inventory_unit" gorm:"not null;default:0;comment:'库存单位'"`
	InventoryAddress string `json:"inventory_address" gorm:"size:1000;not null;default:'';comment:'地址'"`
}
