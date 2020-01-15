package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/githubtasks/proto"
)

type testGithub struct {
	fail bool
}

func (t *testGithub) createMilestone(ctx context.Context, m *pb.Milestone) (int32, error) {
	if t.fail {
		return -1, fmt.Errorf("Built to fail")
	}
	return 10, nil
}

func (t *testGithub) createTask(ctx context.Context, m *pb.Task, service string, num int32) (int32, error) {
	if t.fail {
		return -1, fmt.Errorf("Built to fail")
	}
	return 10, nil
}

func InitTestServer() *Server {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	s.GoServer.KSclient.Save(context.Background(), KEY, &pb.Config{LastUpdate: time.Now().Unix()})
	s.github = &testGithub{}
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

	_, err = s.processProjects(context.Background())

	if err != nil {
		t.Errorf("Error when processing: %v", err)
	}
}

func TestEmptyMilestonesWithAddFail(t *testing.T) {
	s := InitTestServer()
	s.github = &testGithub{fail: true}

	_, err := s.AddProject(context.Background(), &pb.AddProjectRequest{Add: &pb.Project{Name: "Hello", Milestones: []*pb.Milestone{&pb.Milestone{Name: "teting"}}}})
	if err != nil {
		t.Errorf("Error adding project: %v", err)
	}

	err = s.validateIntegrity(context.Background())

	if err != nil {
		t.Errorf("Error in validation: %v", err)
	}

	_, err = s.processProjects(context.Background())

	if err == nil {
		t.Errorf("Processing did not fail")
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

	_, err = s.processProjects(context.Background())

	if err != nil {
		t.Errorf("Error when processing: %v", err)
	}

	err = s.validateIntegrity(context.Background())
	if err != nil {
		t.Errorf("Error in validation: %v", err)
	}

}

func TestNoActiveTasks(t *testing.T) {
	s := InitTestServer()

	_, err := s.AddProject(context.Background(), &pb.AddProjectRequest{Add: &pb.Project{Name: "Hello", Milestones: []*pb.Milestone{&pb.Milestone{Name: "teting", State: pb.Milestone_ACTIVE, Tasks: []*pb.Task{&pb.Task{Title: "Hello"}}}}}})
	if err != nil {
		t.Errorf("Error adding project: %v", err)
	}

	err = s.validateIntegrity(context.Background())

	if err != nil {
		t.Errorf("Error in validation: %v", err)
	}

	_, err = s.processProjects(context.Background())

	if err != nil {
		t.Errorf("Error when processing: %v", err)
	}

	err = s.validateIntegrity(context.Background())
	if err != nil {
		t.Errorf("Error in validation: %v", err)
	}

}

func TestNoActiveTasksAddFail(t *testing.T) {
	s := InitTestServer()
	s.github = &testGithub{fail: true}

	_, err := s.AddProject(context.Background(), &pb.AddProjectRequest{Add: &pb.Project{Name: "Hello", Milestones: []*pb.Milestone{&pb.Milestone{Name: "teting", State: pb.Milestone_ACTIVE, Tasks: []*pb.Task{&pb.Task{Title: "Hello"}}}}}})
	if err != nil {
		t.Errorf("Error adding project: %v", err)
	}

	err = s.validateIntegrity(context.Background())

	if err != nil {
		t.Errorf("Error in validation: %v", err)
	}

	_, err = s.processProjects(context.Background())

	if err == nil {
		t.Errorf("Error when processing: %v", err)
	}

}

func TestActiveTasks(t *testing.T) {
	s := InitTestServer()

	_, err := s.AddProject(context.Background(), &pb.AddProjectRequest{Add: &pb.Project{Name: "Hello", Milestones: []*pb.Milestone{&pb.Milestone{Name: "teting", State: pb.Milestone_ACTIVE, Tasks: []*pb.Task{&pb.Task{Title: "Hello", State: pb.Task_ACTIVE}}}}}})
	if err != nil {
		t.Errorf("Error adding project: %v", err)
	}

	err = s.validateIntegrity(context.Background())

	if err != nil {
		t.Errorf("Error in validation: %v", err)
	}

	_, err = s.processProjects(context.Background())

	if err != nil {
		t.Errorf("Error when processing: %v", err)
	}

	err = s.validateIntegrity(context.Background())
	if err != nil {
		t.Errorf("Error in validation: %v", err)
	}

}
