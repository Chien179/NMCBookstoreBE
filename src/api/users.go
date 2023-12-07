package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/helper"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/Chien179/NMCBookstoreBE/src/token"
	"github.com/Chien179/NMCBookstoreBE/src/util"
	"github.com/Chien179/NMCBookstoreBE/src/val"
	"github.com/Chien179/NMCBookstoreBE/src/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
)

func newUserResponse(user db.User) models.UserResponse {
	return models.UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		Image:             user.Image,
		Age:               user.Age,
		Sex:               user.Sex,
		PhoneNumber:       user.PhoneNumber,
		Role:              user.Role,
		IsEmailVerified:   user.IsEmailVerified,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req models.CreateUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	imgUrl := "https://res.cloudinary.com/doqhasjec/image/upload/v1681990980/samples/NMC%20Bookstore/Default_ct9xzk.png"

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:    req.Username,
			Password:    hashedPassword,
			FullName:    "",
			Email:       req.Email,
			Image:       imgUrl,
			Age:         0,
			Sex:         "",
			PhoneNumber: "",
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user.User)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req models.LoginUserRequest
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

	rsp := models.LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}

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

func (server *Server) updateUser(ctx *gin.Context) {
	var req models.UpdateUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ValidateUpdateUserRequest(&req))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.UpdateUserParams{
		Username: authPayLoad.Username,
		FullName: sql.NullString{
			String: req.Fullname,
			Valid:  req.Fullname != "",
		},
		Email: sql.NullString{
			String: req.Email,
			Valid:  req.Email != "",
		},
		Age: sql.NullInt32{
			Int32: req.Age,
			Valid: req.Age > 0 && req.Age < 200,
		},
		Sex: sql.NullString{
			String: req.Sex,
			Valid:  req.Sex != "",
		},
		PhoneNumber: sql.NullString{
			String: req.PhoneNumber,
			Valid:  req.PhoneNumber != "",
		},
	}

	if req.Image != nil {
		imgUrl, err := server.uploadFile(ctx, req.Image, "NMCBookstore/Image/Users", authPayLoad.Username)
		if err != nil {
			return
		}
		arg.Image = sql.NullString{
			String: imgUrl,
			Valid:  true,
		}
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

func (server *Server) listUser(ctx *gin.Context) {
	users, err := server.store.ListUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func ValidateCreateUserRequest(req *models.CreateUserRequest) (errs []helper.ErrorCustom) {
	if err := val.ValidateUsername(req.Username); err != nil {
		errs = append(errs, helper.ErrorCustom{
			Field:   "username",
			Message: errorResponse(err),
		})
	}

	if err := val.ValidatePassword(req.Password); err != nil {
		errs = append(errs, helper.ErrorCustom{
			Field:   "password",
			Message: errorResponse(err),
		})
	}

	if err := val.ValidatePassword(req.Password); err != nil {
		errs = append(errs, helper.ErrorCustom{
			Field:   "password",
			Message: errorResponse(err),
		})
	}

	if err := val.ValidateEmail(req.Email); err != nil {
		errs = append(errs, helper.ErrorCustom{
			Field:   "email",
			Message: errorResponse(err),
		})
	}

	return errs
}

func ValidateUpdateUserRequest(req *models.UpdateUserRequest) (errs []helper.ErrorCustom) {
	if req.Password != "" {
		if err := val.ValidatePassword(req.Password); err != nil {
			errs = append(errs, helper.ErrorCustom{
				Field:   "password",
				Message: errorResponse(err),
			})
		}
	}

	if req.Fullname != "" {
		if err := val.ValidateFullName(req.Fullname); err != nil {
			errs = append(errs, helper.ErrorCustom{
				Field:   "full_name",
				Message: errorResponse(err),
			})
		}
	}

	if req.Email != "" {
		if err := val.ValidateEmail(req.Email); err != nil {
			errs = append(errs, helper.ErrorCustom{
				Field:   "email",
				Message: errorResponse(err),
			})
		}
	}

	if req.PhoneNumber != "" {
		if err := val.ValidatePhoneNumber(req.PhoneNumber); err != nil {
			errs = append(errs, helper.ErrorCustom{
				Field:   "phone_number",
				Message: errorResponse(err),
			})
		}
	}

	return errs
}
