package main

import (
	"fmt"
	"sort"
	"time"

	"golang.org/x/net/context"

	ghcpb "github.com/brotherlogic/githubcard/proto"
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
					} else {
						activeTask := false
						for _, task := range milestone.GetTasks() {
							if task.State == pb.Task_ACTIVE {
								activeTask = true
							}
						}

						if !activeTask {
							s.RaiseIssue(ctx, "Task Issue", fmt.Sprintf("%v of %v has no active tasks", project.GetName(), milestone.GetName()), false)
						}
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

func (s *Server) updateProjects(ctx context.Context) (time.Time, error) {
	err := s.load(ctx)

	if err == nil {
		for _, project := range s.config.GetProjects() {
			for _, milestone := range project.GetMilestones() {
				if milestone.GetState() == pb.Milestone_ACTIVE {
					for _, task := range milestone.GetTasks() {
						if task.GetState() == pb.Task_ACTIVE {
							issue, err := s.github.getIssue(ctx, milestone.GetGithubProject(), task.GetNumber())
							if err != nil {
								return time.Now().Add(time.Minute * 5), err
							}

							if issue.GetState() == ghcpb.Issue_CLOSED {
								task.State = pb.Task_COMPLETE
							}
							err = s.save(ctx)
						}
					}
				}
			}
		}
	}

	return time.Now().Add(time.Minute * 5), err
}

func (s *Server) processProjects(ctx context.Context) (time.Time, error) {
	err := s.load(ctx)

	if err == nil {
		for _, project := range s.config.GetProjects() {
			for _, milestone := range project.GetMilestones() {
				if milestone.GetState() == pb.Milestone_ACTIVE {

					// Sort tasks by the UID
					sort.SliceStable(milestone.GetTasks(), func(i, j int) bool {
						return milestone.GetTasks()[i].GetUid() < milestone.GetTasks()[j].GetUid()
					})

					for _, task := range milestone.GetTasks() {
						if task.GetState() == pb.Task_ACTIVE {
							break
						}

						if task.GetState() == pb.Task_CREATED {
							num, err := s.github.createTask(ctx, task, milestone.GetGithubProject(), milestone.GetNumber())
							s.Log(fmt.Sprintf("Added task %v -> %v,%v", task.GetTitle(), num, err))
							if err != nil {
								return time.Now().Add(time.Minute * 5), err
							}

							task.Number = num
							task.State = pb.Task_ACTIVE
							break
						}
					}

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
