package main

import (
	"context"
	"testing"

	"github.com/brotherlogic/keystore/client"
)

func InitTestServer() *Server {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	return s
}

func TestEmptyConfig(t *testing.T) {
	s := InitTestServer()
	err := s.validateIntegrity(context.Background())

	if err != nil {
		t.Errorf("Error in validation")
	}
}
