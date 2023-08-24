package models

import db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"

type AddToWishlistRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type DeleteMultiBookInWishlistRequest struct {
	IDs []int64 `form:"ids" binding:"required"`
}

type ListBooksInWishlistRespone struct {
	WishlistID int64   `json:"wishlist_id"`
	Book       db.Book `json:"book"`
}
