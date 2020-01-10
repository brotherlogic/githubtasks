package main

import (
	"fmt"

	"golang.org/x/net/context"
)

func (s *Server) validateIntegrity(ctx context.Context) error {
	if len(s.config.GetProjects()) == 0 {
		s.RaiseIssue(ctx, "Task Issue", fmt.Sprintf("There are no projects listed"), false)
	}

	return nil
}
