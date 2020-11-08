package scheduler

import (
	"reflect"
	"testing"
)

func initScheduler() *Service {
	sites := []string{"site1.com"}
	s := New(sites, 3)
	return s
}

func TestNew(t *testing.T) {
	s := initScheduler()

	if (len(s.History) != 1) || (s.ExpDays != 3) {
		t.Fatalf("incorrect init values")
	}
}

func TestScheduler_UpdateHistory(t *testing.T) {
	s := initScheduler()

	oldTime := s.History["site1.com"]
	s.UpdateHistory([]string{"site1.com"})
	newTime := s.History["site1.com"]
	if newTime.Sub(oldTime) <= 0 {
		t.Fatalf("time not updated")
	}

	s.UpdateHistory([]string{"site2.com"})
	newTime = s.History["site2.com"]
	if (newTime.Sub(oldTime) <= 0) || (len(s.History) != 2) {
		t.Fatalf("new site not added")
	}
}

func TestScheduler_ExpiredSites(t *testing.T) {
	s := initScheduler()

	s.UpdateHistory([]string{"site2.com"})
	expired := s.ExpiredSites()
	if len(expired) != 1 {
		t.Fatalf("initial site not returned")
	}
}

func TestScheduler_Save_LoadFromFile(t *testing.T) {
	s1 := initScheduler()

	var file = "scheduler.bin"
	err := s1.SaveTo(file)
	if err != nil {
		t.Fatalf("incorrect encoding")
	}
	s2, err := LoadFrom(file)
	if err != nil {
		t.Fatalf("incorrect decoding")
	}
	if !reflect.DeepEqual(s1.History, s2.History) {
		t.Fatalf("history changed")
	}
	if s1.ExpDays != s2.ExpDays {
		t.Fatalf("expiration days changed")
	}
}
