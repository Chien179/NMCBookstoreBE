package api

import (
	"net/http"
	"time"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/helper"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/Chien179/NMCBookstoreBE/src/token"
	"github.com/Chien179/NMCBookstoreBE/src/val"
	"github.com/Chien179/NMCBookstoreBE/src/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

func (server *Server) verifyEmail(ctx *gin.Context) {
	var req models.VerifyEmailRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, validateVerifyEmailRequest(&req))
		return
	}

	txResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.EmailID,
		SecretCode: req.SecretCode,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, txResult.User.IsEmailVerified)
}

func (server *Server) sendEmailVerify(ctx *gin.Context) {
	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	taskPayload := &worker.PayloadSendVerifyEmail{
		Username: authPayLoad.Username,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	err := server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, "email will send after 10s")
}

func validateVerifyEmailRequest(req *models.VerifyEmailRequest) (errs []helper.ErrorCustom) {
	if err := val.ValidateEmailId(req.EmailID); err != nil {
		errs = append(errs, helper.ErrorCustom{
			Field:   "email_id",
			Message: errorResponse(err),
		})
	}

	if err := val.ValidateSecretCode(req.SecretCode); err != nil {
		errs = append(errs, helper.ErrorCustom{
			Field:   "secret_code",
			Message: errorResponse(err),
		})
	}

	return errs
}
