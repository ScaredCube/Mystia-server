package service

import (
	"context"
	"time"

	"mystia-voice-backend/internal/db"
	"mystia-voice-backend/proto"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	proto.UnimplementedAuthServiceServer
	DB        *db.DB
	JWTSecret []byte
}

func (s *AuthService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if req.Username == "" || req.Password == "" || req.Nickname == "" {
		return nil, status.Error(codes.InvalidArgument, "missing fields")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	count, err := s.DB.GetUserCount()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check user count")
	}

	role := "USER"
	if count == 0 {
		role = "SUPER_ADMIN"
	}

	user := &db.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		PasswordHash: string(hash),
		Nickname:     req.Nickname,
		Role:         role,
		CreatedAt:    time.Now(),
	}

	if err := s.DB.CreateUser(user); err != nil {
		return nil, status.Error(codes.AlreadyExists, "username already taken")
	}

	msg := "Registration successful"
	if role == "SUPER_ADMIN" {
		msg = "Registration successful. You are the first user and have been granted SUPER_ADMIN privileges."
	}

	return &proto.RegisterResponse{
		Success: true,
		Message: msg,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	user, err := s.DB.GetUserByUsername(req.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}

	if user == nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token, err := GenerateJWT(s.JWTSecret, user.ID, user.Nickname, user.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &proto.LoginResponse{
		Token: token,
		User: &proto.User{
			Id:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Role:     user.Role,
		},
	}, nil
}
