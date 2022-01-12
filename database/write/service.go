package write

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/maracko/go-store/database/helpers"
)

type WriteService struct {
	LastWrite  time.Time
	JobsChan   chan *WriteData
	WritesDone chan bool
	ErrChan    chan error
	Path       string
	mu         sync.Mutex
}

func NewWriteService(path string, jobs chan *WriteData, errs chan error, wd chan bool) *WriteService {
	time := time.Now()
	if path != "" {
		exists := helpers.FileExists(path)
		if !exists {
			writeable := os.WriteFile(path, []byte("{}"), 0777)
			if writeable != nil {
				log.Fatalln("file not writeable:", writeable)
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
	if job == nil || job.Data == nil {
		return errors.New("received nil pointer")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Path == "" || !job.Time.After(s.LastWrite) {
		return nil
	}

	data, err := json.Marshal(job.Data)
	if err != nil {
		return errors.New("marshal error: " + err.Error())
	}
	err = os.WriteFile(s.Path, data, 0777)
	if err != nil {
		return errors.New("write error: " + err.Error())
	}
	s.LastWrite = time.Now()
	return nil
}

func (s *WriteService) Serve() {
	for {
		select {
		case <-s.WritesDone:
			log.Println("Exiting...")
			hasError := false
		jobLoop:
			for {
				select {
				case job := <-s.JobsChan:
					if err := s.write(job); err != nil {
						hasError = true
						s.ErrChan <- err
					}
				case <-time.After(time.Millisecond * 500):
					break jobLoop
				}
			}

			if !hasError {
				log.Println("Clean exit")
			}
			close(s.ErrChan)
			s.WritesDone <- true
		case job := (<-s.JobsChan):
			if err := s.write(job); err != nil {
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
