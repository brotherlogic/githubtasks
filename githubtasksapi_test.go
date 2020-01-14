package main

import (
	"context"
	"testing"

	pb "github.com/brotherlogic/githubtasks/proto"
)

func TestAddProject(t *testing.T) {
	s := InitTestServer()

	s.AddProject(context.Background(),
		&pb.AddProjectRequest{Add: &pb.Project{Name: "test project", Milestones: []*pb.Milestone{&pb.Milestone{Name: "Testing", Number: 1, GithubProject: "madeup"}}}})

	if len(s.config.GetProjects()) != 1 {
		t.Errorf("Project was not added")
	}

	resp, err := s.GetMilestones(context.Background(), &pb.GetMilestonesRequest{GithubProject: "madeup"})
	if err != nil {
		t.Fatalf("Error getting milestones: %v", err)
	}
	if len(resp.GetMilestones()) != 1 {
		t.Fatalf("Milestone not created correctly: %v", resp)
	}

	_, err = s.AddTask(context.Background(),
		&pb.AddTaskRequest{MilestoneName: "Testing", MilestoneNumber: 1, Title: "Add stuff", Body: "Do Stuff"})

	if err != nil {
		t.Errorf("Task add failed: %v", err)
	}

	_, err = s.AddTask(context.Background(),
		&pb.AddTaskRequest{MilestoneName: "Testing_No", MilestoneNumber: 10, Title: "Add stuff", Body: "Do Stuff"})

	if err == nil {
		t.Errorf("Task add did not fail")
	}

}
