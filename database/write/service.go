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
	ErrChan       chan error
	Path          string
	mu            sync.Mutex
}

func NewWriteService(path string, jobs chan WriteData, errs chan error) *WriteService {
	time := time.Now()
	exists := helpers.FileExists(path)
	if !exists {
		writeable := os.WriteFile(path, []byte("{}"), 0777)
		if writeable != nil {
			log.Fatalln("file not writeable")
		}
	}
	return &WriteService{
		LastWrite: time,
		JobsChan:  jobs,
		ErrChan:   errs,
		Path:      path,
	}
}

func (s *WriteService) write(job *WriteData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !job.Time.After(s.LastWrite) {
		return nil
	}
	if len(s.JobsChan) > 1 && s.writesSkipped < 5 {
		s.writesSkipped++
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

	s.writesSkipped = 0
	s.LastWrite = time.Now()
	return nil
}

func (s *WriteService) Serve() error {
	for {
		select {
		case job, ok := (<-s.JobsChan):
			go func() {
				if err := s.write(&job); err != nil {
					s.ErrChan <- err
				}
			}()

			if !ok {
				close(s.ErrChan)
				return errors.New("job channel closed")
			}

			time.Sleep(time.Second * 10)

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
