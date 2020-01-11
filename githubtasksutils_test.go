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
