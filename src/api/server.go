package api

import (
	"fmt"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/token"
	"github.com/Chien179/NMCBookstoreBE/src/util"
	"github.com/Chien179/NMCBookstoreBE/src/worker"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our bookstore.
type Server struct {
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	router          *gin.Engine
	elastic         *elasticsearch.Client
	taskDistributor worker.TaskDistributor
	uploader        util.MediaUpload
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store, elastic *elasticsearch.Client, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetrictKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		elastic:         elastic,
		taskDistributor: taskDistributor,
		uploader:        util.NewMediaUpload(&config),
	}

	server.setupRouter()
	return server, nil
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
