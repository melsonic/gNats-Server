package data

import (
	"fmt"
	"sync"
)

// SUB is map <subject, sid>
// maintain a hashmap with <sid, []subs clients>
type SubjectSID struct {
	mu             sync.Mutex
	subSid         map[string]int
	subSubscribers map[int][]chan string
}

func (s *SubjectSID) Add(subject string, sid int, channel chan string) bool {
	s.mu.Lock()
	// if already exist return
	subjectSid, ok := s.subSid[subject]
	if ok {
		for _, ch := range s.subSubscribers[subjectSid] {
			if ch == channel {
				return false
			}
		}
	} else {
		subjectSid = sid
	}
	s.subSid[subject] = subjectSid
	s.subSubscribers[subjectSid] = append(s.subSubscribers[subjectSid], channel)
	s.mu.Unlock()
	return true
}

func (s *SubjectSID) Publish(subject string, msg []byte) {
	s.mu.Lock()
	subjectSid := s.subSid[subject]
	subsList := s.subSubscribers[subjectSid]
	s.mu.Unlock()
	var response string = fmt.Sprintf("MSG %s %d %d\r\n%s\r\n", subject, subjectSid, len(msg), string(msg))
	for _, subscriber := range subsList {
		subscriber <- response
	}
}

var GSubjectSIDs SubjectSID = SubjectSID{subSid: make(map[string]int), subSubscribers: make(map[int][]chan string)}
