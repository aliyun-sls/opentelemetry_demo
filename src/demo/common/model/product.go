package model

type ProductsStatus uint

const (
	Shelve   ProductsStatus = 1
	Unshelve ProductsStatus = 2
)

// Product 商品
type Product struct {
	Model
	ProductsCate int `json:"products_cate" gorm:"index:ProductsCateIndex;comment:'商品类别'"`
	ProductBasicType
	ProductsDesc    string           `json:"products_desc" gorm:"size:1000;not null;default:'';comment:'商品简介'"`
	BrandId         int              `json:"brand_id" form:"brand_id" gorm:"not null;default:0;comment:'品牌ID'"`
	SellerId        int              `json:"seller_id" form:"seller_id" gorm:"not null;default:0;comment:'商家ID'"`
	ProductsStatus  ProductsStatus   `json:"products_status" gorm:"size:10;index;not null;default:0;comment:'商品状态'"`
	Inventory       *Inventory       `json:"inventory,omitempty" gorm:"foreignKey:ProductsId;"`
	ProductCategory *ProductCategory `json:"product_category,omitempty" gorm:"foreignKey:ProductsCate;"`
}

type ProductsIdType struct {
	ProductsId uint `json:"products_id" form:"products_id" gorm:"not null;index:ProductsIdIdx;comment:'商品ID'"`
}

type ProductBasicType struct {
	ProductsName string  `json:"products_name" gorm:"size:200;not null;default:'';index:ProductsNameIndex;comment:'商品名称'"`
	ProductsUnit int     `json:"products_unit" gorm:"size:10;not null;default:0;comment:'商品单位'"`
	UnitPrice    float64 `json:"unit_price" gorm:"decimal(10,2);not null;default:0;comment:'商品单价'"`
	ProductsPic  PicPath `json:"products_pic" gorm:"size:1000;not null;default:'';comment:'商品图'"`
}

// GetProductsDetailRequest 浏览请求
type GetProductsDetailRequest struct {
	ProductsIdType
	UserIdType
}

// GetProductsDetailResponse 浏览响应
type GetProductsDetailResponse struct {
	Product
	Like      bool  `json:"like"`
	CollectId *uint `json:"collect_id,omitempty"`
}
