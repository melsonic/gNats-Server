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

func (s *SubjectSID) Add(subject string, sid int, channel chan string) {
	s.mu.Lock()
	// if already exist return
	subjectSid, ok := s.subSid[subject]
	if ok {
		for _, ch := range s.subSubscribers[subjectSid] {
			if ch == channel {
				return
			}
		}
	} else {
		subjectSid = sid
	}
	s.subSid[subject] = subjectSid
	s.subSubscribers[subjectSid] = append(s.subSubscribers[subjectSid], channel)
	s.mu.Unlock()
}

func (s *SubjectSID) Publish(subject string, msgLen int, msg []byte) {
	s.mu.Lock()
	subjectSid := s.subSid[subject]
	subsList := s.subSubscribers[subjectSid]
	s.mu.Unlock()
	var response string = fmt.Sprintf("MSG %s %d %d\r\n%s\r\n", subject, subjectSid, msgLen, string(msg))
	for _, subscriber := range subsList {
		subscriber <- response
	}
}

func (s *SubjectSID) Unsub(sid int, channel chan string) {
	s.mu.Lock()
	subsList := s.subSubscribers[sid]
	s.subSubscribers[sid] = nil
	for _, subscriber := range subsList {
		if subscriber == channel {
			continue
		}
		s.subSubscribers[sid] = append(s.subSubscribers[sid], subscriber)
	}
	s.mu.Unlock()
}

var GSubjectSIDs SubjectSID = SubjectSID{subSid: make(map[string]int), subSubscribers: make(map[int][]chan string)}
