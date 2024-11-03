package data

import (
	"net"
	"sync"

	set "github.com/golang-collections/collections/set"
)

// SUB is map <subject, sid>
// maintain a hashmap with <sid, []subs clients>
type SubjectSID struct {
	mu    sync.Mutex
	ssid  map[string]int
	ssubs map[string]*set.Set
}

func (s *SubjectSID) Add(conn net.Conn, subject string, sid int) {
	s.mu.Lock()
	s.ssid[subject] = sid
	if s.ssubs[subject] == nil {
		s.ssubs[subject] = set.New()
	}
	s.ssubs[subject].Insert(conn)
	s.mu.Unlock()
}

var GSubjectSIDs SubjectSID = SubjectSID{ssid: make(map[string]int), ssubs: make(map[string]*set.Set)}
