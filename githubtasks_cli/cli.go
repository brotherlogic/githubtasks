package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/githubtasks/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(&utils.DiscoveryClientResolverBuilder{})
}

func main() {
	conn, err := grpc.Dial("discovery:///githubtasks", grpc.WithInsecure(), grpc.WithBalancerName("my_pick_first"))
	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewTasksServiceClient(conn)
	ctx, cancel := utils.BuildContext("githubtasks-cli", "githubtasks")
	defer cancel()

	switch os.Args[1] {
	case "milestones":
		milestoneFlags := flag.NewFlagSet("Milestones", flag.ExitOnError)
		var ghp = milestoneFlags.String("github_project", "", "Project file to add")

		if err := milestoneFlags.Parse(os.Args[2:]); err == nil {
			resp, err := client.GetMilestones(ctx, &pb.GetMilestonesRequest{GithubProject: *ghp})
			if err != nil {
				log.Fatalf("Error getting milestones: %v", err)
			}

			for _, m := range resp.GetMilestones() {
				fmt.Printf("%v. %v (%v) [%v]\n", m.GetNumber(), m.GetName(), m.GetTasks(), m.GetGithubProject())
			}
		}
	case "projects":
		resp, err := client.GetProjects(ctx, &pb.GetProjectsRequest{})
		if err != nil {
			log.Fatalf("Error getting milestones: %v", err)
		}

		for i, p := range resp.GetProjects() {
			fmt.Printf("%v. %v\n", i, p.GetName())
		}
	case "project":
		projectFlags := flag.NewFlagSet("Project", flag.ExitOnError)
		var file = projectFlags.String("file", "", "Project file to add")

		if err := projectFlags.Parse(os.Args[2:]); err == nil {
			file, err := os.Open(*file)
			if err != nil {
				log.Fatalf("Error reading file: %v", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Scan()
			pname := scanner.Text()
			project := &pb.Project{Name: pname}
			for scanner.Scan() {
				elems := strings.Split(scanner.Text(), "~")
				project.Milestones = append(project.Milestones, &pb.Milestone{Name: elems[0], GithubProject: elems[1]})
			}

			_, err = client.AddProject(ctx, &pb.AddProjectRequest{Add: project})
			if err != nil {
				log.Fatalf("Error adding project: %v", err)
			}
		}
	case "milestone_tasks":
		projectFlags := flag.NewFlagSet("Project", flag.ExitOnError)
		var file = projectFlags.String("file", "", "Project file to add")

		if err := projectFlags.Parse(os.Args[2:]); err == nil {
			file, err := os.Open(*file)
			if err != nil {
				log.Fatalf("Error reading file: %v", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Scan()
			ms := scanner.Text()
			elems := strings.Split(ms, "~")
			number, err := strconv.Atoi(elems[1])
			if err != nil {
				log.Fatalf("Pah: %v", err)
			}
			for scanner.Scan() {
				task := scanner.Text()
				_, err = client.AddTask(ctx, &pb.AddTaskRequest{MilestoneName: elems[0], MilestoneNumber: int32(number), Title: task, Body: "Auto added"})
			}

		}

	}
}
