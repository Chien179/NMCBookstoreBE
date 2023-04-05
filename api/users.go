package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/token"
	"github.com/Chien179/NMCBookstoreBE/util"
	"github.com/Chien179/NMCBookstoreBE/val"
	"github.com/Chien179/NMCBookstoreBE/worker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
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
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Image             string    `json:"image"`
	PhoneNumber       string    `json:"phone_number"`
	Role              string    `json:"role"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		Image:             user.Image,
		PhoneNumber:       user.PhoneNumber,
		Role:              user.Role,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
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
		ctx.JSON(http.StatusBadRequest, ValidateCreateUserRequest(&req))
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

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:    req.Username,
			Password:    hashedPassword,
			FullName:    req.Fullname,
			Email:       req.Email,
			Image:       req.Image,
			PhoneNumber: req.PhoneNumber,
			Role:        "user",
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	user, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	rsp := newUserResponse(user.User)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}
type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
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

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
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
	Fullname    string `json:"full_name"`
	Email       string `json:"email" binding:"email"`
	Image       string `json:"image"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

// @Summary      Update user
// @Description  Use this API to update user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Request body updateUserRequest  false  "Update user"
// @Success      200  {object}  UserResponse
// @failure		 400
// @failure		 404
// @failure		 500
// @Router       /users/update [put]
func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ValidateUpdateUserRequest(&req))
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
		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
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

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
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

func ValidateCreateUserRequest(req *createUserRequest) (errs []errorCustom) {
	if err := val.ValidateUsername(req.Username); err != nil {
		errs = append(errs, errorCustom{
			"username",
			errorResponse(err),
		})
	}

	if err := val.ValidatePassword(req.Password); err != nil {
		errs = append(errs, errorCustom{
			"password",
			errorResponse(err),
		})
	}

	if err := val.ValidateFullName(req.Fullname); err != nil {
		errs = append(errs, errorCustom{
			"full_name",
			errorResponse(err),
		})
	}

	if err := val.ValidateEmail(req.Email); err != nil {
		errs = append(errs, errorCustom{
			"email",
			errorResponse(err),
		})
	}

	if err := val.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		errs = append(errs, errorCustom{
			"phone_number",
			errorResponse(err),
		})
	}

	return errs
}

func ValidateUpdateUserRequest(req *updateUserRequest) (errs []errorCustom) {
	if req.Password != "" {
		if err := val.ValidatePassword(req.Password); err != nil {
			errs = append(errs, errorCustom{
				"password",
				errorResponse(err),
			})
		}
	}

	if req.Fullname != "" {
		if err := val.ValidateFullName(req.Fullname); err != nil {
			errs = append(errs, errorCustom{
				"full_name",
				errorResponse(err),
			})
		}
	}

	if req.Email != "" {
		if err := val.ValidateEmail(req.Email); err != nil {
			errs = append(errs, errorCustom{
				"email",
				errorResponse(err),
			})
		}
	}

	if req.PhoneNumber != "" {
		if err := val.ValidatePhoneNumber(req.PhoneNumber); err != nil {
			errs = append(errs, errorCustom{
				"phone_number",
				errorResponse(err),
			})
		}
	}

	return errs
}
