package write

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/maracko/go-store/database/helpers"
)

type WriteService struct {
	LastWrite time.Time
	JobsChan  chan WriteData
	// writesSkipped int
	WritesDone chan bool
	ErrChan    chan error
	Path       string
	mu         sync.Mutex
}

func NewWriteService(path string, jobs chan WriteData, errs chan error, wd chan bool) *WriteService {
	time := time.Now()
	if path != "" {
		exists := helpers.FileExists(path)
		if !exists {
			writeable := os.WriteFile(path, []byte("{}"), 0777)
			if writeable != nil {
				log.Fatalln("file not writeable")
			}
		}
	}
	return &WriteService{
		LastWrite:  time,
		JobsChan:   jobs,
		ErrChan:    errs,
		Path:       path,
		WritesDone: wd,
	}
}

func (s *WriteService) write(job *WriteData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Path == "" || !job.Time.After(s.LastWrite) {
		return nil
	}

	data, err := json.Marshal(job.Data)
	if err != nil {
		return errors.New("marshall error: " + err.Error())
	}
	err = os.WriteFile(s.Path, data, 0777)
	if err != nil {
		return errors.New("write error: " + err.Error())
	}
	s.LastWrite = time.Now()
	fmt.Println("WROTE")
	return nil
}

func (s *WriteService) Serve() {
	if s.Path == "" {
		s.WritesDone <- true
		return
	}
	for {
		select {
		case <-s.WritesDone:
			fmt.Println("Shutting down writing service")
			if len(s.JobsChan) > 0 {
				lastJob := &WriteData{}
				for data := range s.JobsChan {
					lastJob = &data
				}
				if err := s.write(lastJob); err != nil {
					s.ErrChan <- err
				}
			}
			close(s.ErrChan)
			s.WritesDone <- true
		case job := (<-s.JobsChan):
			if err := s.write(&job); err != nil {
				s.ErrChan <- err
			}
		}
	}
}

type WriteData struct {
	Time time.Time
	Data map[string]interface{}
}

func NewWriteData(data map[string]interface{}) WriteData {
	return WriteData{time.Now(), data}
}
