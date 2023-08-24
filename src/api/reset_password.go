package api

import (
	"net/http"
	"time"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/Chien179/NMCBookstoreBE/src/util"
	"github.com/Chien179/NMCBookstoreBE/src/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
)

func (server *Server) forgotPassword(ctx *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	taskPayload := &worker.PayloadSendResetPassword{
		Email: req.Email,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	err := server.taskDistributor.DistributeTaskSendResetPassword(ctx, taskPayload, opts...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "send email after 10 seconds")
}

func (server *Server) resetPassword(ctx *gin.Context) {
	var req models.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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

	txResult, err := server.store.ResetPasswordTx(ctx, db.ResetPasswordTxParams{
		Id:        req.ID,
		ResetCode: req.ResetCode,
		Password:  hashedPassword,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(txResult.User)
	ctx.JSON(http.StatusOK, rsp)
}
