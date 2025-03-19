package model

import (
	"sls-mall-go/common/util"
)

// Comment 评论
type Comment struct {
	Model
	CommentsContent string  `json:"comments_content" gorm:"size:1000;not null;default:'';comment:'内容'"`
	CommentsScore   int     `json:"comments_score" gorm:"not null;default:0;comment:'分数'"`
	CommentsPic     PicPath `json:"comments_pic" gorm:"size:1000;not null;default:'';comment:'评论图'"`
	ProductsIdType
	UserIdType
	OrderIdType
}

type ListCommentsRequest struct {
	Comment
	IsEffective bool `json:"is_effective" form:"is_effective"`
	util.Page
}

type ListCommentsResponse struct {
	Model
	Comment
	Product
}
