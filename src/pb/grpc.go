package pb

import (
	"context"

	"github.com/rs/zerolog/log"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GRPCGetRCM(ctx context.Context, bookRequest *BookRequest) (*BookResponse, error) {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Logger.Error().Msgf("did not connect: %v", err)
	}
	defer conn.Close()

	client := NewBookRecommendClient(conn)

	books, err := client.GetBookRecommend(ctx, bookRequest)
	if err != nil {
		return nil, err
	}

	return books, nil
}
