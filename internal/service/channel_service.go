package service

import (
	"context"
	"time"

	"mystia-voice-backend/internal/db"
	"mystia-voice-backend/internal/livekit"
	"mystia-voice-backend/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChannelService struct {
	proto.UnimplementedChannelServiceServer
	DB      *db.DB
	LiveKit *livekit.LiveKitProvider
}

func (s *ChannelService) ListChannels(ctx context.Context, req *proto.ListChannelsRequest) (*proto.ListChannelsResponse, error) {
	channels, err := s.DB.ListChannels()
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}

	var pbChannels []*proto.Channel
	for _, c := range channels {
		pbChannels = append(pbChannels, &proto.Channel{
			Id:   c.ID,
			Name: c.Name,
		})
	}

	return &proto.ListChannelsResponse{
		Channels: pbChannels,
	}, nil
}

func (s *ChannelService) JoinChannel(ctx context.Context, req *proto.JoinChannelRequest) (*proto.JoinChannelResponse, error) {
	claims, ok := ctx.Value("claims").(*Claims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	token, err := s.LiveKit.GenerateToken(req.ChannelId, claims.UserID, claims.Nickname)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate livekit token")
	}

	return &proto.JoinChannelResponse{
		Token: token,
		Url:   s.LiveKit.Host,
	}, nil
}

func (s *ChannelService) CreateChannel(ctx context.Context, req *proto.CreateChannelRequest) (*proto.Channel, error) {
	claims, ok := ctx.Value("claims").(*Claims)
	if !ok || (claims.Role != "ADMIN" && claims.Role != "SUPER_ADMIN") {
		return nil, status.Error(codes.PermissionDenied, "admin only")
	}

	channel := &db.Channel{
		ID:        uuid.New().String(),
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if err := s.DB.CreateChannel(channel); err != nil {
		return nil, status.Error(codes.Internal, "failed to create channel")
	}

	return &proto.Channel{
		Id:   channel.ID,
		Name: channel.Name,
	}, nil
}
