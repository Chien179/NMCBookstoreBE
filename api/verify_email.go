package api

import (
	"net/http"
	"time"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/token"
	"github.com/Chien179/NMCBookstoreBE/val"
	"github.com/Chien179/NMCBookstoreBE/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

type verifyEmailRequest struct {
	EmailID    int64  `form:"email_id" binding:"required,min=1"`
	SecretCode string `form:"secret_code" binding:"required"`
}

func (server *Server) verifyEmail(ctx *gin.Context) {
	var req verifyEmailRequest
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

func validateVerifyEmailRequest(req *verifyEmailRequest) (errs []errorCustom) {
	if err := val.ValidateEmailId(req.EmailID); err != nil {
		errs = append(errs, errorCustom{"email_id", errorResponse(err)})
	}

	if err := val.ValidateSecretCode(req.SecretCode); err != nil {
		errs = append(errs, errorCustom{"secret_code", errorResponse(err)})
	}

	return errs
}
