package main

import (
	"context"
	"testing"

	pb "github.com/brotherlogic/githubtasks/proto"
)

func TestDeleteTask(t *testing.T) {
	s := InitTestServer()

	tr, err := s.DeleteTask(context.Background(), &pb.DeleteTaskRequest{})
	if err != nil {
		t.Errorf("Failure to delete empty task: %v", err)
	}

	if tr.GetTask() != nil {
		t.Errorf("Stars have clashed: %v", tr)
	}
}

func TestBadAddTask(t *testing.T) {
	s := InitTestServer()
	s.GoServer.KSclient.Fail = true

	_, err := s.AddTask(context.Background(), &pb.AddTaskRequest{})

	if err == nil {
		t.Errorf("Failing add did not fail")
	}
}

func TestAddProject(t *testing.T) {
	s := InitTestServer()

	s.AddProject(context.Background(),
		&pb.AddProjectRequest{Add: &pb.Project{Name: "test project", Milestones: []*pb.Milestone{&pb.Milestone{Name: "Testing", Number: 1, GithubProject: "madeup"}}}})

	if len(s.config.GetProjects()) != 1 {
		t.Errorf("Project was not added")
	}

	presp, err := s.GetProjects(context.Background(), &pb.GetProjectsRequest{})
	if err != nil {
		t.Fatalf("Error getting projects: %v", err)
	}
	if len(presp.GetProjects()) != 1 {
		t.Fatalf("Project not created correctly: %v", presp)
	}

	resp, err := s.GetMilestones(context.Background(), &pb.GetMilestonesRequest{GithubProject: "madeup"})
	if err != nil {
		t.Fatalf("Error getting milestones: %v", err)
	}
	if len(resp.GetMilestones()) != 1 {
		t.Fatalf("Milestone not created correctly: %v", resp)
	}

	_, err = s.AddTask(context.Background(),
		&pb.AddTaskRequest{MilestoneName: "Testing", MilestoneNumber: 1, Title: "Add stuff", Body: "Do Stuff", GithubProject: "madeup"})

	if err != nil {
		t.Errorf("Task add failed: %v", err)
	}

	_, err = s.AddTask(context.Background(),
		&pb.AddTaskRequest{MilestoneName: "Testing", MilestoneNumber: 1, Title: "Add stuff", Body: "Do Stuff", GithubProject: "madeup"})

	if err == nil {
		t.Errorf("Double Task add failed: %v", err)
	}

	_, err = s.AddTask(context.Background(),
		&pb.AddTaskRequest{MilestoneName: "Testing_No", MilestoneNumber: 10, Title: "Add stuff", Body: "Do Stuff"})

	if err == nil {
		t.Errorf("Task add did not fail")
	}

	resp, err = s.GetMilestones(context.Background(), &pb.GetMilestonesRequest{})
	if err != nil {
		t.Fatalf("Cannot get milestones")
	}

	for _, m := range resp.GetMilestones() {
		for _, tsk := range m.GetTasks() {
			td, err := s.DeleteTask(context.Background(), &pb.DeleteTaskRequest{Uid: tsk.GetUid()})
			if err != nil {
				t.Errorf("Bad task delete")
			}

			if td.GetTask() == nil {
				t.Errorf("Bad task find")
			}
		}
	}

	dresp, err := s.DeleteProject(context.Background(), &pb.DeleteProjectRequest{Name: "test projectsssss"})
	if err != nil {
		t.Errorf("Bad project delete: %v", err)
	}
	if dresp.GetDeleted() != int32(0) {
		t.Errorf("Hmm - no projects deleted")
	}

	dresp, err = s.DeleteProject(context.Background(), &pb.DeleteProjectRequest{Name: "test project"})
	if err != nil {
		t.Errorf("Bad project delete: %v", err)
	}
	if dresp.GetDeleted() != int32(1) {
		t.Errorf("Hmm - no projects deleted")
	}
}
