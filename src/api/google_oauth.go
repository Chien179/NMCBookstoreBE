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
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func (server *Server) GoogleOAuth(ctx *gin.Context) {
	var code models.GoogleOauthRequest
	if err := ctx.ShouldBindQuery(&code); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if code.Code == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Authorization code not provided!")))
		return
	}

	tokenRes, err := util.GetGoogleOauthToken(code.Code)

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

func (server *Server) GoogleOAuthURL(ctx *gin.Context) {
	// Developer Console (https://console.developers.google.com).
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	conf := &oauth2.Config{
		ClientID:     config.GoogleOauthClientID,
		ClientSecret: config.GoogleOauthClientSecret,
		RedirectURL:  config.GoogleOAuthRedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/bigquery",
			"https://www.googleapis.com/auth/blogger",
		},
		Endpoint: google.Endpoint,
	}
	// Redirect user to Google's consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state")

	ctx.JSON(http.StatusOK, map[string]string{"url": url})
}
