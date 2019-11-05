package main

import (
	"time"
)

func NewSecondsTimer(t time.Duration) *countdownTimer {
	return &countdownTimer{time.NewTimer(t), time.Now().Add(t)}
}

func (s *countdownTimer) Reset(t time.Duration) {
	s.timer.Reset(t)
	s.end = time.Now().Add(t)
}

func (s *countdownTimer) Stop() {
	s.timer.Stop()
}

func (s *countdownTimer) TimeRemaining() time.Duration {
	return s.end.Sub(time.Now()) / 1000000000
} 
   