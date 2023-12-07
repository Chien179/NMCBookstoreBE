package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(corsMiddleware())

	server.publicRouter(router)
	server.userAuth(router)
	server.adminAuth(router)

	router.Use(cors.Default())
	server.router = router
}

func (server *Server) publicRouter(router *gin.Engine) {
	router.POST("/signup", server.createUser)
	router.POST("/login", server.loginUser)
	router.GET("/login/oauth/google", server.GoogleOAuth)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	router.GET("/verify_email", server.verifyEmail)
	router.POST("/forgot_password", server.forgotPassword)
	router.PUT("/reset_password", server.resetPassword)

	router.GET("/searchs", server.fullSearch)

	router.GET("/just_for_you", server.justForYou)

	bookRoutes := router.Group("/books")
	bookRoutes.GET("/recommend", server.recommend)
	bookRoutes.GET("/:id", server.getBook)
	bookRoutes.GET("/", server.listBook)
	bookRoutes.GET("/the_best", server.listTheBestBook)
	bookRoutes.GET("/newest", server.listNewestBook)

	genreRoutes := router.Group("/genres")
	genreRoutes.GET("/:id", server.getGenre)
	genreRoutes.GET("/", server.listGenre)

	subgenreRoutes := router.Group("/subgenres")
	subgenreRoutes.GET("/one/:id", server.getSubgenre)
	subgenreRoutes.GET("/:genre_id", server.listSubgenre)
	subgenreRoutes.GET("/", server.listAllSubgenre)
	subgenreRoutes.GET("/noticeable", server.listSubgenresNoticeable)

	reviewRoutes := router.Group("/reviews")
	reviewRoutes.GET("/:book_id", server.listReview)
}

func (server *Server) userAuth(router *gin.Engine) {
	usersRoutes := router.Group("/users")

	userRoutes := usersRoutes.Use(authMiddleware(server.tokenMaker))
	userRoutes.GET("/", server.getUser)
	userRoutes.PUT("/", server.updateUser)
	userRoutes.DELETE("/", server.deleteUser)
	usersRoutes.POST("/send_verify_email", server.sendEmailVerify)

	cartRoutes := usersRoutes.Group("/carts").Use(authMiddleware(server.tokenMaker))
	cartRoutes.POST("/:id", server.addToCart)
	cartRoutes.DELETE("/", server.deleteBookInCart)
	cartRoutes.PUT("/:id", server.upatdeAmountCart)
	cartRoutes.GET("/", server.listBookInCart)

	wishlistRoutes := usersRoutes.Group("/wishlists").Use(authMiddleware(server.tokenMaker))
	wishlistRoutes.POST("/:id", server.addToWishlist)
	wishlistRoutes.DELETE("/", server.deleteBookInWishlist)
	wishlistRoutes.GET("/", server.listBookInWishlist)

	addressRoutes := usersRoutes.Group("/addresses").Use(authMiddleware(server.tokenMaker))
	addressRoutes.POST("/", server.createAddress)
	addressRoutes.GET("/:id", server.getAddress)
	addressRoutes.GET("/", server.listAddress)
	addressRoutes.PUT("/:id", server.updateAddress)
	addressRoutes.DELETE("/", server.deleteAddress)
	addressRoutes.GET("/cities/:id", server.getCity)
	addressRoutes.GET("/cities", server.listCities)
	addressRoutes.GET("/districts/:city_id", server.listDistricts)

	reviewRoutes := usersRoutes.Group("/reviews").Use(authMiddleware(server.tokenMaker))
	reviewRoutes.GET("/like", server.getLikeReview)
	reviewRoutes.GET("/action/like", server.likeReview)
	reviewRoutes.GET("/action/dislike", server.dislikeReview)
	reviewRoutes.GET("/dislike", server.getDislikeReview)
	reviewRoutes.POST("/:book_id", server.createReview)
	reviewRoutes.DELETE("/:id", server.deleteReview)

	orderRoutes := usersRoutes.Group("/orders").Use(authMiddleware(server.tokenMaker))
	orderRoutes.POST("/", server.createOrder)
	orderRoutes.GET("/", server.listOrder)
	orderRoutes.GET("/paid", server.listOrderPaid)
	orderRoutes.GET("/cancelled", server.listOrderCancelled)
	orderRoutes.PUT("/:id", server.cancelOrder)
	orderRoutes.DELETE("/:id", server.deleteOrder)
}

func (server *Server) adminAuth(router *gin.Engine) {
	adminRoutes := router.Group("/admin")

	userRoutes := adminRoutes.Group("/users").Use(authMiddleware(server.tokenMaker), isAdmin())
	userRoutes.GET("/", server.listUser)

	bookRoutes := adminRoutes.Group("/books").Use(authMiddleware(server.tokenMaker), isAdmin())
	bookRoutes.GET("/", server.listBook)
	bookRoutes.POST("/", server.createBook)
	bookRoutes.PUT("/:id", server.updateBook)
	bookRoutes.DELETE("/:id", server.deleteBook)
	bookRoutes.DELETE("soft/:id", server.softDeleteBook)

	genreRoutes := adminRoutes.Group("/genres").Use(authMiddleware(server.tokenMaker), isAdmin())
	genreRoutes.POST("/", server.createGenre)
	genreRoutes.PUT("/:id", server.updateGenre)
	genreRoutes.DELETE("/:id", server.deleteGenre)
	genreRoutes.DELETE("soft/:id", server.softDeleteGenre)

	subgenreRoutes := adminRoutes.Group("/subgenres").Use(authMiddleware(server.tokenMaker), isAdmin())
	subgenreRoutes.POST("/", server.createSubgenre)
	subgenreRoutes.PUT("/:id", server.updateSubgenre)
	subgenreRoutes.DELETE("/:id", server.deleteSubgenre)
	subgenreRoutes.DELETE("soft/:id", server.softDeleteSubgenre)

	revenueRoutes := adminRoutes.Group("/revenues").Use(authMiddleware(server.tokenMaker), isAdmin())
	revenueRoutes.GET("/days", server.revenueDays)
	revenueRoutes.GET("/months", server.revenueMonths)
	revenueRoutes.GET("/quarters", server.revenueQuarters)
	revenueRoutes.GET("/years", server.revenueYears)

	orderRoutes := adminRoutes.Group("/orders").Use(authMiddleware(server.tokenMaker))
	orderRoutes.GET("/", server.listOrder)
}
