package main

import (
	"fmt"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/githubtasks/proto"
)

func (s *Server) validateIntegrity(ctx context.Context) error {
	err := s.load(ctx)

	if err == nil {
		if len(s.config.GetProjects()) == 0 {
			s.RaiseIssue(ctx, "Task Issue", fmt.Sprintf("There are no projects listed"), false)
		}

		for _, project := range s.config.GetProjects() {
			activeMilestone := false
			for _, milestone := range project.GetMilestones() {
				if milestone.GetState() == pb.Milestone_ACTIVE {
					activeMilestone = true
				}
			}

			if !activeMilestone {
				s.RaiseIssue(ctx, "Task Issue", fmt.Sprintf("%v has no active milestones", project.GetName()), false)
			}
		}
	}

	return err
}
