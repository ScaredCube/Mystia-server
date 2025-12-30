package service

import (
	"context"

	"mystia-voice-backend/internal/db"
	"mystia-voice-backend/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminService struct {
	proto.UnimplementedAdminServiceServer
	DB *db.DB
}

func (s *AdminService) SetAdminStatus(ctx context.Context, req *proto.SetAdminStatusRequest) (*proto.SetAdminStatusResponse, error) {
	claims, ok := ctx.Value("claims").(*Claims)
	if !ok || (claims.Role != "ADMIN" && claims.Role != "SUPER_ADMIN") {
		return nil, status.Error(codes.PermissionDenied, "admin only")
	}

	targetUser, err := s.DB.GetUserByID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch target user")
	}
	if targetUser == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Protect super admin
	if targetUser.Role == "SUPER_ADMIN" {
		return nil, status.Error(codes.PermissionDenied, "cannot modify SUPER_ADMIN status")
	}

	role := "USER"
	if req.IsAdmin {
		role = "ADMIN"
	}

	if err := s.DB.UpdateUserRole(req.UserId, role); err != nil {
		return nil, status.Error(codes.Internal, "failed to update user role")
	}

	return &proto.SetAdminStatusResponse{Success: true}, nil
}

func (s *AdminService) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	claims, ok := ctx.Value("claims").(*Claims)
	if !ok || (claims.Role != "ADMIN" && claims.Role != "SUPER_ADMIN") {
		return nil, status.Error(codes.PermissionDenied, "admin only")
	}

	users, err := s.DB.ListUsers()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	var pbUsers []*proto.User
	for _, u := range users {
		pbUsers = append(pbUsers, &proto.User{
			Id:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
			Role:     u.Role,
		})
	}

	return &proto.ListUsersResponse{Users: pbUsers}, nil
}
