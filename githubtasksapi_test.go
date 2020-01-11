package main

import (
	"context"
	"testing"

	pb "github.com/brotherlogic/githubtasks/proto"
)

func TestAddProject(t *testing.T) {
	s := InitTestServer()

	s.AddProject(context.Background(),
		&pb.AddProjectRequest{Add: &pb.Project{Name: "test project"}})

	if len(s.config.GetProjects()) != 1 {
		t.Errorf("Project was not added")
	}
}
