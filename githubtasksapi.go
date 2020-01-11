package main

import "golang.org/x/net/context"

import pb "github.com/brotherlogic/githubtasks/proto"

// AddProject to the system
func (s *Server) AddProject(ctx context.Context, req *pb.AddProjectRequest) (*pb.AddProjectResponse, error) {
	err := s.load(ctx)
	if err == nil {
		s.config.Projects = append(s.config.Projects, req.GetAdd())
		err = s.save(ctx)
	}
	return &pb.AddProjectResponse{}, err
}
