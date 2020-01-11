package main

import (
	"fmt"
	"time"

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
					if len(milestone.GetTasks()) == 0 {
						s.RaiseIssue(ctx, "Task Issue", fmt.Sprintf("%v of %v has no tasks", project.GetName(), milestone.GetName()), false)
					}
				}
			}

			if !activeMilestone {
				s.RaiseIssue(ctx, "Task Issue", fmt.Sprintf("%v has no active milestones", project.GetName()), false)
			}
		}

	}

	return err
}

func (s *Server) processProjects(ctx context.Context) (time.Time, error) {
	err := s.load(ctx)

	if err == nil {
		for _, project := range s.config.GetProjects() {
			for _, milestone := range project.GetMilestones() {
				if milestone.GetState() == pb.Milestone_ACTIVE {
					break
				}

				if milestone.GetState() == pb.Milestone_CREATED {
					num, err := s.github.createMilestone(ctx, milestone)

					if err != nil {
						return time.Now().Add(time.Minute * 5), err
					}

					milestone.Number = num
					milestone.State = pb.Milestone_ACTIVE
					break
				}
			}
		}

		err = s.save(ctx)
	}

	return time.Now().Add(time.Minute * 5), err
}
