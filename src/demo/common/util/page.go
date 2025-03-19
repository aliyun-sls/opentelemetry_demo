package util

type Page struct {
	Limit  int  `json:"limit" form:"limit"`
	Offset int  `json:"offset" form:"offset"`
	Desc   bool `json:"desc" form:"desc"`
}
