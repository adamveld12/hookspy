package main

import (
	b64 "encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dancannon/gorethink"

	"gopkg.in/satori/go.uuid.v1"
)

type HookSession struct {
	ID          string      `gorethink:"hook_id" json:"hookId"`
	Created     time.Time   `gorethink:"created" json:"created"`
	HookEntries []HookEntry `gorethink:"entries" json:"entries"`
}
type HookDiff struct {
	Old HookSession `gorethink:"old_val"`
	New HookSession `gorethink:"new_val"`
}

func (h HookDiff) NewEntries() []HookEntry {
	oldCount := len(h.Old.HookEntries)
	return h.New.HookEntries[oldCount:]
}

type HookEntry struct {
	ID         string            `gorethink:"hook_id" json:"hookId"`
	Created    time.Time         `gorethink:"created" json:"created"`
	Protocol   string            `gorethink:"proto" json:"proto"`
	RemoteAddr string            `gorethink:"remote_addr" json:"remoteAddr"`
	Host       string            `gorethink:"host" json:"host"`
	URL        string            `gorethink:"url" json:"url"`
	Method     string            `gorethink:"method" json:"method"`
	Headers    map[string]string `gorethink:"headers" json:"headers"`
	Body       string            `gorethink:"body,omitempty" json:"body"`
}

type sessionManager struct {
	session *gorethink.Session
	debug   bool
}

func (s *sessionManager) NewSession() (*HookSession, error) {
	hookID := uuid.NewV4().String()

	hs := HookSession{ID: hookID, Created: time.Now().UTC()}

	if err := gorethink.DB("hookspy").Table("hook_sessions").Insert(hs).Exec(s.session); err != nil {
		log.Println("Failed to make a new session", err.Error())
		return nil, err
	}

	return &hs, nil
}

func (s *sessionManager) UpdateSession(hookID string, req *http.Request) error {
	b, _ := ioutil.ReadAll(req.Body)
	sEnc := b64.StdEncoding.EncodeToString(b)
	log.Printf("%s", b)
	he := HookEntry{
		ID:         "",
		Created:    time.Now().UTC(),
		Protocol:   req.Proto,
		RemoteAddr: req.RemoteAddr,
		Host:       req.Host,
		URL:        req.URL.String(),
		Method:     req.Method,
		Headers:    map[string]string{},
		Body:       sEnc,
	}

	h := req.Header
	log.Printf("%+v", h)
	for k, v := range h {
		he.Headers[k] = strings.Join(v, ";")
	}

	// p, _ := json.MarshalIndent(he, "", " ")
	// log.Printf("%s", p)

	hs, err := s.LookupSession(hookID)
	if err != nil {
		return err
	}

	hs.HookEntries = append(hs.HookEntries, he)

	return gorethink.DB("hookspy").Table("hook_sessions").Update(hs, gorethink.UpdateOpts{}).Exec(s.session)
}

func (s *sessionManager) Changes(hookID string) <-chan HookEntry {
	entries := make(chan HookEntry)
	go func() {
		defer close(entries)
		for {
			res, err := gorethink.DB("hookspy").Table("hook_sessions").Get(hookID).Changes().Run(s.session)
			log.Println("trying to find changes")
			if err != nil {
				log.Println(err)
				return
			}

			hd := HookDiff{}
			for res.Next(&hd) {
				for _, e := range hd.NewEntries() {
					entries <- e
				}
			}
		}
	}()

	return entries
}

func (s *sessionManager) LookupSession(hookID string) (*HookSession, error) {
	_, err := uuid.FromString(hookID)
	if err != nil {
		return s.NewSession()
	}

	res, err := gorethink.DB("hookspy").Table("hook_sessions").Get(hookID).Run(s.session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	h := &HookSession{}
	if err := res.One(h); err != nil {
		return nil, errors.New("No result found")
	}

	return h, nil
}
