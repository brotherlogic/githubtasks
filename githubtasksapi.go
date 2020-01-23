package main

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	task := &pb.Task{Title: req.GetTitle(), Body: req.GetBody()}

	for _, p := range s.config.GetProjects() {
		for _, m := range p.GetMilestones() {
			if m.GetName() == req.GetMilestoneName() && m.GetNumber() == req.GetMilestoneNumber() && m.GetGithubProject() == req.GetGithubProject() {
				for _, t := range m.GetTasks() {
					if t.GetTitle() == task.GetTitle() {
						return &pb.AddTaskResponse{Task: t}, status.Errorf(codes.AlreadyExists, "Task exists")
					}
				}
				m.Tasks = append(m.Tasks, task)
				return &pb.AddTaskResponse{Task: task}, s.save(ctx)
			}
		}
	}

	return nil, fmt.Errorf("Could not locate milestone %v/%v", req.GetMilestoneName(), req.GetMilestoneNumber())
}

// DeleteTask to the system
func (s *Server) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {

	for _, p := range s.config.GetProjects() {
		for _, m := range p.GetMilestones() {
			for i, t := range m.GetTasks() {
				if t.GetUid() == req.GetUid() {
					m.Tasks = append(m.Tasks[:i], m.Tasks[i+1:]...)
					return &pb.DeleteTaskResponse{Task: t}, s.save(ctx)
				}
			}
		}
	}

	return &pb.DeleteTaskResponse{}, nil
}

// GetProjects for the system
func (s *Server) GetProjects(ctx context.Context, req *pb.GetProjectsRequest) (*pb.GetProjectsResponse, error) {
	resp := &pb.GetProjectsResponse{Projects: []*pb.Project{}}
	for _, p := range s.config.GetProjects() {
		resp.Projects = append(resp.Projects, p)
	}

	return resp, nil
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
