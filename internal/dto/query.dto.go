package dto

type QueryGlobal struct {
	Search string `form:"search"`
	SortBy string `form:"sort"`
	Order  string `form:"order"`
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=10"`
}
