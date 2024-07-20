package grpc

import (
	"assignment3/entity"
	pb "assignment3/proto/shortlink/v1"
	"assignment3/service"
	"context"
	"fmt"
	"log"
)

type ShortlinkHandler struct {
	pb.UnimplementedShortlinkServiceServer
	shortlinkService service.IShortlinkService
}

func NewShortlinkHandler(shortlinkService service.IShortlinkService) *ShortlinkHandler {
	return &ShortlinkHandler{
		shortlinkService: shortlinkService,
	}
}

func (u *ShortlinkHandler) CreateShortlink(ctx context.Context, req *pb.CreateShortlinkRequest) (*pb.MutationResponse, error) {
	createdShortlink, err := u.shortlinkService.CreateShortlink(ctx, &entity.Shortlink{
		Url: req.GetUrl(),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success created shortlink. Your shortlink: http://localhost:8080/%v", createdShortlink.Shortlink),
	}, nil
}

func (u *ShortlinkHandler) GetLongUrl(ctx context.Context, req *pb.GetLongUrlRequest) (*pb.GetLongUrlResponse, error) {
	longURL, err := u.shortlinkService.GetLongURL(ctx, req.GetShortlink())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.GetLongUrlResponse{LongUrl: longURL}, nil
}
