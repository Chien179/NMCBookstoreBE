package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/Chien179/NMCBookstoreBE/src/util"
	"github.com/Chien179/NMCBookstoreBE/src/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

func (server *Server) GoogleOAuth(ctx *gin.Context) {
	code := ctx.Query("code")

	if code == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Authorization code not provided!")))
		return
	}

	tokenRes, err := util.GetGoogleOauthToken(code)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, errorResponse(err))
		return
	}

	google_user, err := util.GetGoogleUser(tokenRes.AccessToken, tokenRes.IdToken)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, errorResponse(err))
		return
	}

	email := strings.ToLower(google_user.Email)

	user, err := server.store.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			arg := db.CreateUserTxParams{
				CreateUserParams: db.CreateUserParams{
					Username:    google_user.Name,
					Password:    "google",
					FullName:    "",
					Email:       email,
					Image:       google_user.Picture,
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

			newUser, err := server.store.CreateUserTx(ctx, arg)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			user = newUser.User
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
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
