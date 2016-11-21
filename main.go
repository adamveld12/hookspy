package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/kelseyhightower/envconfig.v1"

	"github.com/adamveld12/muxwrap"
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
)

type Config struct {
	Addr  string
	DB    string
	Debug bool
}

func main() {
	log.SetFlags(log.Ltime | log.LUTC | log.Llongfile)

	c := Config{}
	envconfig.MustProcess("hookspy", &c)

	if c.Debug {
		log.Printf("Running w/ %+v", c)
	}

	dbSession, err := r.Connect(r.ConnectOpts{
		Address:    c.DB,
		InitialCap: 10,
		MaxOpen:    10,
	})

	if err != nil {
		log.Fatal("Could not connect to database", err)
		return
	}

	// for now, just generate the database and table at start up and ignore errs
	r.DBCreate("hookspy").RunWrite(dbSession)
	r.DB("hookspy").TableCreate("hook_sessions", r.TableCreateOpts{PrimaryKey: "hook_id"}).RunWrite(dbSession)

	s := &sessionManager{session: dbSession}

	mux := muxwrap.New(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			header := res.Header()
			if c.Debug {
				header.Set("Access-Control-Allow-Origin", "*")
			} else {
				header.Set("Access-Control-Allow-Origin", "https://hookspy.veldhousen.ninja")
				header.Set("Strict-Transport-Security", "max-age=31536000")
			}

			next.ServeHTTP(res, req)
		})
	})

	fs := assetFS()

	fs.Prefix = "client/build"
	mux.Embed("/", http.FileServer(fs))

	mux.Handle("/hook/", logHook(s))

	mux.Get("/session/", startSession(s))

	log.Println("Serving on 3001")
	log.Fatal(http.ListenAndServe(c.Addr, mux))
}

func socketSession(sock *websocket.Conn, s *sessionManager, hookID string) {
	defer sock.Close()

	hookSession, err := s.LookupSession(hookID)
	if err != nil {
		log.Println("Can't create a new hook session")
		return
	}

	hookID = hookSession.ID

	log.Println("HookID", hookID)
	p, _ := json.Marshal(hookSession)
	handshake := []byte(fmt.Sprintf(`{ "type": "CREATE_SESSION", "payload" : %s }`, p))
	if err := sock.WriteMessage(websocket.TextMessage, handshake); err != nil {
		log.Println("Could not write to websocket", err.Error())
	}

	// open up change update thingy with matching hook ID
	changes := s.Changes(hookID)
	timer := time.NewTimer(time.Second * 10)
	defer timer.Stop()
	defer log.Println("Disconnecting from client")

	for {
		select {
		case <-timer.C:
			if err := sock.WriteMessage(websocket.TextMessage, []byte(`{ "type": "HEARTBEAT" }`)); err != nil {
				return
			}
			timer.Reset(time.Second * 10)
		case change := <-changes:
			p, _ := json.Marshal(change)
			payload := []byte(fmt.Sprintf(`{"type":"NEW_REQUEST", "payload": %s }`, p))
			if err := sock.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		}
	}

}

func startSession(s *sessionManager) http.HandlerFunc {
	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		conn, err := wsUpgrader.Upgrade(res, req, nil)
		if err != nil {
			log.Println("Could not upgrade connection to websocket: ", err.Error())
			return
		}

		hookID := strings.Trim(req.URL.Path, "/session/")
		socketSession(conn, s, hookID)
	})
}

func logHook(s *sessionManager) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			header := res.Header()
			header.Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS, PUT, HEAD")
			header.Set("Access-Control-Max-Age", "86400")
			header.Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}

		hookID := strings.TrimPrefix(req.URL.Path, "/hook/")
		if hookID == "" || s.UpdateSession(hookID, req) != nil {
			http.Error(res, "404 not found", http.StatusNotFound)
		}
	})
}
