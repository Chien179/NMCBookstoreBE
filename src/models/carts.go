package models

type AmountData struct {
	Amount int32 `json:"amount" binding:"required,min=1"`
}

type AddToCartRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
	AmountData
}

type UpdateAmountCartRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
	AmountData
}

type DeleteMultiBookInCartRequest struct {
	IDs []int64 `form:"ids" binding:"required"`
}

type ListBooksInCartRespone struct {
	CartID   int64   `json:"cart_id"`
	BookID   int64   `json:"book_id"`
	BookName string  `json:"book_name"`
	Image    string  `json:"image"`
	Price    float64 `json:"price"`
	Amount   int32   `json:"amount"`
	Author   string  `json:"author"`
}
