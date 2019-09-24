package main

import (
	"flag"
	"github.com/l453595892/raft"
	"github.com/l453595892/sqlsync/command"
	"github.com/l453595892/sqlsync/db"
	"github.com/l453595892/sqlsync/env"
	"github.com/l453595892/sqlsync/server"
	"log"
	"math/rand"
	"os"
	"time"
)

type config struct {
	Trace    bool
	Debug    bool
	DataPath string `default:"~/data"`

	Config db.Config
}

var host string
var port int
var join string

func init() {
	flag.StringVar(&join, "join", "", "host:port of leader to join")
	flag.StringVar(&host, "h", "localhost", "hostname")
	flag.IntVar(&port, "p", 4001, "port")
}

func main() {
	log.SetFlags(0)
	// flag
	flag.Parse()
	// env
	conf := new(config)
	env.IgnorePrefix()
	err := env.Fill(conf)
	if err != nil {
		panic(err)
	}

	// trace and debug
	if conf.Trace {
		raft.SetLogLevel(raft.Trace)
		log.Print("Raft trace debugging enabled.")
	} else if conf.Debug {
		raft.SetLogLevel(raft.Debug)
		log.Print("Raft debugging enabled.")
	}

	rand.Seed(time.Now().UnixNano())
	raft.RegisterCommand(&command.WriteCommand{})

	// set store path
	if err := os.MkdirAll(conf.DataPath, 0744); err != nil {
		log.Fatalf("Unable to create path: %v", err)
	}

	log.SetFlags(log.LstdFlags)
	s := server.New(conf.DataPath, host, port, conf.Config)
	log.Fatal(s.ListenAndServe(join))
}
