package main

import (
	"context"
	"testing"
	"time"

	"github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/githubtasks/proto"
)

func InitTestServer() *Server {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	s.GoServer.KSclient.Save(context.Background(), KEY, &pb.Config{LastUpdate: time.Now().Unix()})
	return s
}

func TestEmptyConfig(t *testing.T) {
	s := InitTestServer()
	err := s.validateIntegrity(context.Background())

	if err != nil {
		t.Errorf("Error in validation: %v", err)
	}
}

func TestEmptyMilestones(t *testing.T) {
	s := InitTestServer()

	_, err := s.AddProject(context.Background(), &pb.AddProjectRequest{Add: &pb.Project{Name: "Hello", Milestones: []*pb.Milestone{&pb.Milestone{Name: "teting"}}}})
	if err != nil {
		t.Errorf("Error adding project: %v", err)
	}

	err = s.validateIntegrity(context.Background())

	if err != nil {
		t.Errorf("Error in validation: %v", err)
	}
}

func TestActiveMilestone(t *testing.T) {
	s := InitTestServer()

	_, err := s.AddProject(context.Background(), &pb.AddProjectRequest{Add: &pb.Project{Name: "Hello", Milestones: []*pb.Milestone{&pb.Milestone{Name: "teting", State: pb.Milestone_ACTIVE}}}})
	if err != nil {
		t.Errorf("Error adding project: %v", err)
	}

	err = s.validateIntegrity(context.Background())

	if err != nil {
		t.Errorf("Error in validation: %v", err)
	}
}
