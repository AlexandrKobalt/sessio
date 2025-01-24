package grpc

import (
	"context"

	session "github.com/AlexandrKobalt/sessio/api/proto"
	"github.com/AlexandrKobalt/sessio/internal"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	uc internal.UseCase
	session.UnimplementedSessioServer
}

func New(uc internal.UseCase) session.SessioServer {
	return &Server{uc: uc}
}

func (s *Server) CreateSession(ctx context.Context, req *session.CreateSessionRequest) (*session.CreateSessionResponse, error) {
	authKey, err := s.uc.CreateSession(ctx, req.GetUserId(), req.GetUserAgent())
	if err != nil {
		return nil, err
	}

	return &session.CreateSessionResponse{AuthKey: authKey}, nil
}

func (s *Server) Authenticate(ctx context.Context, req *session.AuthenticateRequest) (*session.AuthenticateResponse, error) {
	userID, err := s.uc.Authenticate(ctx, req.GetAuthKey(), req.GetUserAgent())
	if err != nil {
		return nil, err
	}

	return &session.AuthenticateResponse{UserId: userID}, nil
}

func (s *Server) GetAllActiveSessions(ctx context.Context, req *session.GetAllActiveSessionsRequest) (*session.GetAllActiveSessionsResponse, error) {
	activeSessions, err := s.uc.GetAllActiveSessions(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	sessions := make([]*session.Session, len(activeSessions))
	for i, s := range activeSessions {
		sessions[i] = &session.Session{
			AuthKey:   s.AuthKey,
			UserAgent: s.UserAgent,
			CreatedAt: s.CreatedAt.Unix(),
			ExpiresAt: s.ExpiresAt.Unix(),
		}
	}

	return &session.GetAllActiveSessionsResponse{Sessions: sessions}, nil
}

func (s *Server) DeleteSession(ctx context.Context, req *session.DeleteSessionRequest) (*emptypb.Empty, error) {
	err := s.uc.DeleteSession(ctx, req.GetUserId(), req.GetUserAgent())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
