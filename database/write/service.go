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
	LastWrite     time.Time
	JobsChan      chan WriteData
	writesSkipped int
	isLastWrite   bool
	writesDone    chan bool
	ErrChan       chan error
	Path          string
	mu            sync.Mutex
}

func NewWriteService(path string, jobs chan WriteData, errs chan error, wd chan bool) *WriteService {
	time := time.Now()
	exists := helpers.FileExists(path)
	if !exists {
		writeable := os.WriteFile(path, []byte("{}"), 0777)
		if writeable != nil {
			log.Fatalln("file not writeable")
		}
	}
	return &WriteService{
		LastWrite:  time,
		JobsChan:   jobs,
		ErrChan:    errs,
		Path:       path,
		writesDone: wd,
	}
}

func (s *WriteService) write(job *WriteData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Println(s.isLastWrite)
	if !job.Time.After(s.LastWrite) {
		return nil
	}
	if !s.isLastWrite {
		if len(s.JobsChan) > 1 && s.writesSkipped < 5 {
			s.writesSkipped++
			return nil
		}
	}
	data, err := json.Marshal(job.Data)
	if err != nil {
		return errors.New("marshall error: " + err.Error())
	}
	err = os.WriteFile(s.Path, data, 0777)
	if err != nil {
		return errors.New("write error: " + err.Error())
	}

	if s.isLastWrite {
		log.Println("last write finished")
		s.writesDone <- true
		return nil
	}

	s.writesSkipped = 0
	s.LastWrite = time.Now()
	return nil
}

func (s *WriteService) Serve() error {
	for {
		select {
		case job, ok := (<-s.JobsChan):
			if !ok {
				s.isLastWrite = true
				if err := s.write(&job); err != nil {
					s.ErrChan <- err
				}
				close(s.ErrChan)
				return nil
			}

			if err := s.write(&job); err != nil {
				s.ErrChan <- err
			}

		default:
			continue
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
