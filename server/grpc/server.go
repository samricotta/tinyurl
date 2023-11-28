package server

import (
	"context"
	"net"

	"github.com/samricotta/tinyurl/store"
	"google.golang.org/grpc"
)

var _ TinyURLServer = (*Server)(nil)

type Server struct {
	store *store.Store
	UnimplementedTinyURLServer
}

func (s *Server) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	longUrl, err := s.store.Get(req.TinyUrl)
	stlurl := string(longUrl)
	rs := &GetResponse{LongUrl: stlurl}
	return rs, err
}

func (s *Server) Post(ctx context.Context, req *PostRequest) (*PostResponse, error) {

	longUrl := req.LongUrl
	tinyUrl, err := s.store.Set([]byte(longUrl))
	if err != nil {
		return nil, err
	}
	rs := &PostResponse{TinyUrl: tinyUrl}

	return rs, nil
}

func Serve(path string) error {
	store, err := store.New(path)
	if err != nil {
		return err
	}
	s := &Server{store: store}
	grpcServer := grpc.NewServer()
	RegisterTinyURLServer(grpcServer, s)
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	defer listener.Close()

	return grpcServer.Serve(listener)
}
