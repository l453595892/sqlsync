package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/l453595892/raft"
	"github.com/l453595892/sqlsync/command"
	"github.com/l453595892/sqlsync/db"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

// The raftd server is a combination of the Raft server and an HTTP
// server which acts as the transport.
type Server struct {
	name       string
	host       string
	port       int
	path       string
	router     *mux.Router
	raftServer raft.Server
	httpServer *http.Server
	conf       db.Config
	mutex      sync.RWMutex
}

// Creates a new server.
func New(path string, host string, port int, config db.Config) *Server {
	s := &Server{
		host:   host,
		port:   port,
		path:   path,
		conf:   config,
		router: mux.NewRouter(),
	}

	// Read existing name or generate a new one.
	if b, err := ioutil.ReadFile(filepath.Join(path, "name")); err == nil {
		s.name = string(b)
	} else {
		s.name = fmt.Sprintf("%07x", rand.Int())[0:7]
		if err = ioutil.WriteFile(filepath.Join(path, "name"), []byte(s.name), 0644); err != nil {
			panic(err)
		}
	}

	return s
}

// Returns the connection string.
func (s *Server) connectionString() string {
	return fmt.Sprintf("http://%s:%d", s.host, s.port)
}

// Starts the server.
func (s *Server) ListenAndServe(leader string) error {
	var err error

	log.Printf("Initializing Raft Server: %s", s.path)

	// Initialize and start Raft server.
	transporter := raft.NewHTTPTransporter("/raft", 200*time.Millisecond)
	s.raftServer, err = raft.NewServer(s.name, s.path, transporter, nil, s.conf, "")
	if err != nil {
		log.Fatal(err)
	}
	transporter.Install(s.raftServer, s)
	s.raftServer.Start()

	if leader != "" {
		// Join to leader if specified.

		log.Println("Attempting to join leader:", leader)

		if !s.raftServer.IsLogEmpty() {
			log.Fatal("Cannot join with an existing log")
		}
		if err := s.Join(leader); err != nil {
			log.Fatal(err)
		}

	} else if s.raftServer.IsLogEmpty() {
		// Initialize the server by joining itself.

		log.Println("Initializing new cluster")

		_, err := s.raftServer.Do(&raft.DefaultJoinCommand{
			Name:             s.raftServer.Name(),
			ConnectionString: s.connectionString(),
		})
		if err != nil {
			log.Fatal(err)
		}

	} else {
		log.Println("Recovered from log")
	}

	log.Println("Initializing HTTP server")

	// Initialize and start HTTP server.
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}

	s.router.HandleFunc("/create", s.createHandler).Methods("POST")
	s.router.HandleFunc("/insert", s.insertHandler).Methods("POST")
	s.router.HandleFunc("/update", s.updateHandler).Methods("POST")
	s.router.HandleFunc("/delete", s.deleteHandler).Methods("POST")

	s.router.HandleFunc("/join", s.joinHandler).Methods("POST")

	log.Println("Listening at:", s.connectionString())
	return s.httpServer.ListenAndServe()
}

// This is a hack around Gorilla mux not providing the correct net/http
// HandleFunc() interface.
func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.router.HandleFunc(pattern, handler)
}

// Joins to the leader of an existing cluster.
func (s *Server) Join(leader string) error {
	command := &raft.DefaultJoinCommand{
		Name:             s.raftServer.Name(),
		ConnectionString: s.connectionString(),
	}

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(command)
	resp, err := http.Post(fmt.Sprintf("http://%s/join", leader), "application/json", &b)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func (s *Server) joinHandler(w http.ResponseWriter, req *http.Request) {
	command := &raft.DefaultJoinCommand{}

	if err := json.NewDecoder(req.Body).Decode(&command); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := s.raftServer.Do(command); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Query struct {
	Sql []string `json:"sql"`
}

func (s *Server) createHandler(w http.ResponseWriter, req *http.Request) {

	// Read the value from the POST body.
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := Query{}
	err = json.Unmarshal(b, &query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Execute the command against the Raft server.
	_, err = s.raftServer.Do(command.NewWriteCommand(query.Sql))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) insertHandler(w http.ResponseWriter, req *http.Request) {

	// Read the value from the POST body.
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := Query{}
	err = json.Unmarshal(b, &query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Execute the command against the Raft server.
	_, err = s.raftServer.Do(command.NewWriteCommand(query.Sql))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) updateHandler(w http.ResponseWriter, req *http.Request) {

	// Read the value from the POST body.
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := Query{}
	err = json.Unmarshal(b, &query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Execute the command against the Raft server.
	_, err = s.raftServer.Do(command.NewWriteCommand(query.Sql))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) deleteHandler(w http.ResponseWriter, req *http.Request) {

	// Read the value from the POST body.
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := Query{}
	err = json.Unmarshal(b, &query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Execute the command against the Raft server.
	_, err = s.raftServer.Do(command.NewWriteCommand(query.Sql))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
