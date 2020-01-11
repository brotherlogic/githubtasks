package main

import (
	"bufio"
	"flag"
	"log"
	"os"

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
				milestone := scanner.Text()
				project.Milestones = append(project.Milestones, &pb.Milestone{Name: milestone})
			}

			_, err = client.AddProject(ctx, &pb.AddProjectRequest{Add: project})
			if err != nil {
				log.Fatalf("Error adding project: %v", err)
			}
		}
	}
}
