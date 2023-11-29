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
