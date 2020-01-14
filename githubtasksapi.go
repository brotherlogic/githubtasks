package main

import (
	"fmt"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/githubtasks/proto"
)

// AddProject to the system
func (s *Server) AddProject(ctx context.Context, req *pb.AddProjectRequest) (*pb.AddProjectResponse, error) {
	err := s.load(ctx)
	if err == nil {
		s.config.Projects = append(s.config.Projects, req.GetAdd())
		err = s.save(ctx)
	}
	return &pb.AddProjectResponse{}, err
}

// AddTask to the system
func (s *Server) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	task := &pb.Task{Title: req.GetTitle(), Body: req.GetBody()}

	for _, p := range s.config.GetProjects() {
		for _, m := range p.GetMilestones() {
			if m.GetName() == req.GetMilestoneName() && m.GetNumber() == req.GetMilestoneNumber() {
				m.Tasks = append(m.Tasks, task)
				return &pb.AddTaskResponse{Task: task}, nil
			}
		}
	}

	return nil, fmt.Errorf("Could not locate milestone %v/%v", req.GetMilestoneName(), req.GetMilestoneNumber())
}

// GetMilestones for the system
func (s *Server) GetMilestones(ctx context.Context, req *pb.GetMilestonesRequest) (*pb.GetMilestonesResponse, error) {
	resp := &pb.GetMilestonesResponse{Milestones: []*pb.Milestone{}}
	for _, p := range s.config.GetProjects() {
		for _, m := range p.GetMilestones() {
			if len(req.GetGithubProject()) == 0 || m.GetGithubProject() == req.GetGithubProject() {
				resp.Milestones = append(resp.Milestones, m)
			}
		}
	}

	return resp, nil
}
