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
	config, err := s.load(ctx)
	if err != nil {
		return err
	}

	if err == nil {
		if len(config.GetProjects()) == 0 {
			s.RaiseIssue("Task Issue", fmt.Sprintf("There are no projects listed"))
		}

		for _, project := range config.GetProjects() {
			activeMilestone := false
			noComplete := true
			mstone := &pb.Milestone{}
			for _, milestone := range project.GetMilestones() {
				if milestone.GetState() == pb.Milestone_ACTIVE {
					activeMilestone = true
					if len(milestone.GetTasks()) == 0 {
						s.RaiseIssue("Task Issue", fmt.Sprintf("%v of %v has no tasks", project.GetName(), milestone.GetName()))
					} else {
						activeTask := false
						for _, task := range milestone.GetTasks() {
							if task.State == pb.Task_ACTIVE {
								activeTask = true
							}
						}

						if !activeTask && len(milestone.GetTasks()) > 0 {
							s.RaiseIssue("Task Issue", fmt.Sprintf("%v of %v for %v has no active tasks", project.GetName(), milestone.GetName(), milestone.GetGithubProject()))
						}
					}
				} else if milestone.GetState() != pb.Milestone_COMPLETE {
					noComplete = false
					mstone = milestone
				}

			}

			if (!activeMilestone && len(project.GetMilestones()) > 0) && !noComplete {
				s.RaiseIssue("Task Issue", fmt.Sprintf("%v has no active milestones (%v is not complete)", project.GetName(), mstone))
			}
		}

	}

	return err
}

func (s *Server) updateProjects(ctx context.Context) (time.Time, error) {
	config, err := s.load(ctx)

	if err == nil {
		for _, project := range config.GetProjects() {
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
	config, err := s.load(ctx)

	if err == nil {
		for _, project := range config.GetProjects() {
			for _, milestone := range project.GetMilestones() {
				s.CtxLog(ctx, fmt.Sprintf("Process: %v -> %v", milestone.GetName(), milestone.GetState()))
				time.Sleep(time.Second * 5)
				if milestone.GetState() == pb.Milestone_ACTIVE {
					countActive := 0

					// Sort tasks by the UID
					sort.SliceStable(milestone.GetTasks(), func(i, j int) bool {
						return milestone.GetTasks()[i].GetUid() < milestone.GetTasks()[j].GetUid()
					})

					for _, task := range milestone.GetTasks() {
						if task.GetState() == pb.Task_ACTIVE {
							countActive++
							break
						}

						if task.GetState() == pb.Task_CREATED {
							countActive++
							num, err := s.github.createTask(ctx, task, milestone.GetGithubProject(), milestone.GetNumber())
							s.CtxLog(ctx, fmt.Sprintf("Added task %v -> %v,%v", task.GetTitle(), num, err))
							if err != nil {
								return time.Now().Add(time.Minute * 5), err
							}

							task.Number = num
							task.State = pb.Task_ACTIVE
							break
						}
					}

					// If we've reached this point, we have no active taskss and no tasks that need creation
					if len(milestone.GetTasks()) > 0 && countActive == 0 {
						//Close out the milestone
						milestone.State = pb.Milestone_COMPLETE
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
