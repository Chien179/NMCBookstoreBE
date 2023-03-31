package api

import (
	docs "github.com/Chien179/NMCBookstoreBE/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func (server *Server) setupRouter() {
	router := gin.Default()

	corsMiddleware(router)

	docs.SwaggerInfo.BasePath = "/"

	server.publicRouter(router)
	server.userAuth(router)
	server.adminAuth(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.router = router
}

func (server *Server) publicRouter(router *gin.Engine) {
	router.POST("/signup", server.createUser)
	router.POST("/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	router.GET("/books/:id", server.getBook)
	router.GET("/books", server.listBook)

	router.GET("/genres", server.listGenre)
	router.GET("/subgenres/:genre_id", server.listSubgenre)

	router.GET("/reviews/:book_id", server.listReview)
}

func (server *Server) userAuth(router *gin.Engine) {
	userRoutes := router.Group("/users").Use(authMiddleware(server.tokenMaker))
	userRoutes.GET("/me", server.getUser)
	userRoutes.PUT("/update", server.updateUser)
	userRoutes.DELETE("/delete", server.deleteUser)

	userRoutes.POST("/add_to_cart/:id", server.addToCart)
	userRoutes.DELETE("/delete_book_in_cart/:id", server.deleteBookInCart)
	userRoutes.GET("/list_book_in_cart", server.listBookInCart)

	userRoutes.POST("/add_to_wishlist/:id", server.addToWishlist)
	userRoutes.DELETE("/delete_book_in_wishlist/:id", server.deleteBookInWishlist)
	userRoutes.GET("/list_book_in_wishlist", server.listBookInWishlist)

	userRoutes.POST("/addresses", server.createAddress)
	userRoutes.GET("/addresses/:id", server.getAddress)
	userRoutes.GET("/addresses", server.listAddress)
	userRoutes.PUT("/addresses/update/:id", server.updateAddress)
	userRoutes.DELETE("/addresses/delete/:id", server.deleteAddress)

	userRoutes.POST("/reviews/:book_id", server.createReview)
	userRoutes.DELETE("/reviews/delete/:id", server.deleteReview)

	userRoutes.POST("/orders", server.createOrder)
	userRoutes.GET("/orders", server.listOrder)
	userRoutes.DELETE("/orders/delete/:id", server.deleteOrder)
}

func (server *Server) adminAuth(router *gin.Engine) {
	adminRoutes := router.Group("/admin").Use(authMiddleware(server.tokenMaker), isAdmin())
	adminRoutes.POST("/books", server.createBook)
	adminRoutes.PUT("/books/update/:id", server.updateBook)
	adminRoutes.DELETE("/books/delete/:id", server.deleteBook)

	adminRoutes.POST("/genres", server.createGenre)
	adminRoutes.PUT("/genres/update/:id", server.updateGenre)
	adminRoutes.DELETE("/genres/delete/:id", server.deleteGenre)

	adminRoutes.POST("/subgenres", server.createSubgenre)
	adminRoutes.PUT("/subgenres/update/:id", server.updateSubgenre)
	adminRoutes.DELETE("/subgenres/delete/:id", server.deleteSubgenre)
}
