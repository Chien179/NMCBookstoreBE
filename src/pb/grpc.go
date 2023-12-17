package pb

import (
	"context"

	"github.com/Chien179/NMCBookstoreBE/src/util"
	"github.com/rs/zerolog/log"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GRPCGetRCM(ctx context.Context, bookRequest *BookRequest) (*BookResponse, error) {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	conn, err := grpc.Dial(config.GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Logger.Error().Msgf("did not connect: %v", err)
	}
	log.Info().Msg("Connect GRPC recommend success")
	defer conn.Close()

	client := NewBookRecommendClient(conn)

	books, err := client.GetBookRecommend(ctx, bookRequest)
	if err != nil {
		return nil, err
	}

	return books, nil
}
