package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/token"
	"github.com/Chien179/NMCBookstoreBE/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password    string `json:"password" binding:"required,min=8"`
	Fullname    string `json:"full_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Image       string `json:"image" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type UserResponse struct {
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	Image       string    `json:"image"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

func newUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username:    user.Username,
		FullName:    user.FullName,
		Email:       user.Email,
		Image:       user.Image,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
	}
}

// @Summary      Create new user
// @Description  Use this API to create a new user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Request body createUserRequest  true  "Create user"
// @Success      200  {object}  UserResponse
// @failure	 	 400
// @failure		 403
// @failure		 500
// @Router       /signup [post]
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, err)
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateUserParams{
		Username:    req.Username,
		Password:    hashedPassword,
		FullName:    req.Fullname,
		Email:       req.Email,
		Image:       req.Image,
		PhoneNumber: req.PhoneNumber,
		Role:        "user",
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}
type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

// @Summary      Login user
// @Description  Use this API to login user and get access token & refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Request body loginUserRequest  true  "Login user"
// @Success      200  {object}  loginUserResponse
// @failure	 	 400
// @failure	 	 401
// @failure		 403
// @failure		 404
// @failure		 500
// @Router       /login [post]
func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}

// @Summary      Get user
// @Description  Use this API to get user with token access
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  UserResponse
// @failure	 	 400
// @failure	 	 401
// @failure		 403
// @failure		 404
// @failure		 500
// @Router       /users/me [get]
func (server *Server) getUser(ctx *gin.Context) {
	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUser(ctx, authPayLoad.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type updateUserRequest struct {
	Fullname    string `json:"full_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Image       string `json:"image" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

// @Summary      Update user
// @Description  Use this API to update user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Request body updateUserRequest  true  "Update user"
// @Success      200  {object}  UserResponse
// @failure		 400
// @failure		 404
// @failure		 500
// @Router       /users/update [put]
func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.UpdateUserParams{
		Username: authPayLoad.Username,
		FullName: sql.NullString{
			String: req.Fullname,
			Valid:  true,
		},
		Email: sql.NullString{
			String: req.Email,
			Valid:  true,
		},
		Image: sql.NullString{
			String: req.Image,
			Valid:  true,
		},
		PhoneNumber: sql.NullString{
			String: req.PhoneNumber,
			Valid:  true,
		},
	}

	if req.Password != "" {
		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		arg.Password = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

// @Summary      Delete user
// @Description  Use this API to delete user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200
// @failure		 404
// @failure		 500
// @Router       /users/delete [delete]
func (server *Server) deleteUser(ctx *gin.Context) {
	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err := server.store.DeleteUser(ctx, authPayLoad.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "User deleted successfully")
}
