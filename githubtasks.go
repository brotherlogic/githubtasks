package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	ghcpb "github.com/brotherlogic/githubcard/proto"
	pb "github.com/brotherlogic/githubtasks/proto"
	pbg "github.com/brotherlogic/goserver/proto"
)

type github interface {
	createMilestone(ctx context.Context, m *pb.Milestone) (int32, error)
}

type prodGithub struct {
	dial func(server string) (*grpc.ClientConn, error)
}

func (p *prodGithub) createMilestone(ctx context.Context, m *pb.Milestone) (int32, error) {
	conn, err := p.dial("githubcard")
	if err != nil {
		return -1, err
	}
	defer conn.Close()

	client := ghcpb.NewGithubClient(conn)
	resp, err := client.AddMilestone(ctx, &ghcpb.AddMilestoneRequest{Title: m.GetName(), Repo: m.GetGithubProject()})
	if err != nil {
		return -1, err
	}
	return resp.GetNumber(), err
}

const (
	// KEY where the config is stored
	KEY = "/github.com/brotherlogic/githubtasks/config"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	config *pb.Config
	github github
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
		config:   &pb.Config{},
	}
	s.github = &prodGithub{dial: s.DialMaster}
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterTasksServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "no", Value: int64(233)},
	}
}

func (s *Server) save(ctx context.Context) error {
	return s.KSclient.Save(ctx, KEY, s.config)
}

func (s *Server) load(ctx context.Context) error {
	data, _, err := s.KSclient.Read(ctx, KEY, s.config)
	if err != nil {
		return err
	}

	s.config = data.(*pb.Config)
	return nil
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server
	err := server.RegisterServerV2("githubtasks", false, true)
	if err != nil {
		return
	}

	server.RegisterRepeatingTask(server.validateIntegrity, "validate_integrity", time.Hour)
	server.RegisterLockingTask(server.processProjects, "process_projects")

	fmt.Printf("%v", server.Serve())
}
