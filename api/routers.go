package api

import (
	docs "github.com/Chien179/NMCBookstoreBE/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(corsMiddleware())

	docs.SwaggerInfo.BasePath = "/"

	server.publicRouter(router)
	server.userAuth(router)
	server.adminAuth(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(cors.Default())
	server.router = router
}

func (server *Server) publicRouter(router *gin.Engine) {
	router.POST("/signup", server.createUser)
	router.POST("/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	router.GET("/verify_email", server.verifyEmail)
	router.POST("/forgot_password", server.forgotPassword)
	router.PUT("/reset_password", server.resetPassword)

	router.GET("/searchs", server.fullSearch)

	bookRoutes := router.Group("/books")
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
	addressRoutes.DELETE("/:id", server.deleteAddress)
	addressRoutes.GET("/cities/:id", server.getCity)
	addressRoutes.GET("/cities", server.listCities)
	addressRoutes.GET("/districts/:city_id", server.listDistricts)

	reviewRoutes := usersRoutes.Group("/reviews").Use(authMiddleware(server.tokenMaker))
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

	bookRoutes := adminRoutes.Group("/books").Use(authMiddleware(server.tokenMaker), isAdmin())
	bookRoutes.GET("/", server.listAllBook)
	bookRoutes.POST("/", server.createBook)
	bookRoutes.PUT("/:id", server.updateBook)
	bookRoutes.DELETE("/:id", server.deleteBook)

	genreRoutes := adminRoutes.Group("/genres").Use(authMiddleware(server.tokenMaker), isAdmin())
	genreRoutes.POST("/", server.createGenre)
	genreRoutes.PUT("/:id", server.updateGenre)
	genreRoutes.DELETE("/:id", server.deleteGenre)

	subgenreRoutes := adminRoutes.Group("/subgenres").Use(authMiddleware(server.tokenMaker), isAdmin())
	subgenreRoutes.POST("/", server.createSubgenre)
	subgenreRoutes.PUT("/:id", server.updateSubgenre)
	subgenreRoutes.DELETE("/:id", server.deleteSubgenre)

	revenueRoutes := adminRoutes.Group("/revenues").Use(authMiddleware(server.tokenMaker), isAdmin())
	revenueRoutes.GET("/days", server.revenueDays)
	revenueRoutes.GET("/months", server.revenueMonths)
	revenueRoutes.GET("/quarters", server.revenueQuarters)
	revenueRoutes.GET("/years", server.revenueYears)

	orderRoutes := adminRoutes.Group("/orders").Use(authMiddleware(server.tokenMaker))
	orderRoutes.GET("/", server.listOrder)
}
