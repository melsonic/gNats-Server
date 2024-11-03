package data

import (
	"sync"
)

// SUB is map <subject, sid>
// maintain a hashmap with <sid, []subs clients>
type SubjectSID struct {
	mu   sync.Mutex
	ssid map[string]int
}

func (s *SubjectSID) Add(subject string, sid int) {
	s.mu.Lock()
	s.ssid[subject] = sid
	s.mu.Unlock()
}

var GSubjectSIDs SubjectSID = SubjectSID{ssid: make(map[string]int)}
