package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/maracko/go-store/database/helpers"
	"github.com/maracko/go-store/database/write"
)

// DB represents the database struct
type DB struct {
	location       string
	database       map[string]interface{}
	memory         bool
	continousWrite bool
	errRcv         chan error
	jobSender      chan write.WriteData
	writeService   *write.WriteService
	mu             sync.Mutex
}

// New initializes a database to a given location and sets it's internal DB to an empty map or reads from file first
func New(location string, memory bool, continousWrite bool, ec chan error) *DB {

	jc := make(chan write.WriteData, 10)
	ws := write.NewWriteService(location, jc, ec)
	return &DB{
		location:       location,
		database:       make(map[string]interface{}),
		errRcv:         ec,
		jobSender:      jc,
		memory:         memory,
		writeService:   ws,
		continousWrite: continousWrite,
	}
}

// Connect connects to file and saves it's contents to database field
func (d *DB) Connect() error {
	if d.database == nil {
		return errors.New("db not initialized")
	}

	if d.location != "" {
		if !helpers.FileExists(d.location) {
			f, err := os.Create(d.location)
			if err != nil {
				return err
			}
			if _, err = f.WriteString("{}"); err != nil {
				return err
			}
		}
		db, err := helpers.ReadJsonToMap(d.location)
		if err != nil {
			return errors.New("cannot read file: " + err.Error())
		}
		d.database = db
	}
	if d.memory {
		return nil
	}

	go func() {
		log.Println("Starting write service")
		err := d.writeService.Serve()
		if err != nil {
			log.Println(err)
		}
	}()
	return nil
}

func (d *DB) NewWrite() {
	data := write.NewWriteData(d.database)
	d.jobSender <- data
}

// Disconnect encodes database with json and saves it to location if provided
func (d *DB) Disconnect() error {
	if len(d.database) == 0 || d.location == "" || d.memory {
		return nil
	}

	writeFinished := make(chan bool)
	go func() {
		d.NewWrite()
		writeFinished <- true
	}()
	<-writeFinished
	time.Sleep(time.Millisecond * 500)
	close(d.jobSender)
	return nil
}

// Create creates a new record
func (d *DB) Create(key string, value interface{}) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.database[key]; ok {
		return errors.New("key already exists")
	}
	d.database[key] = value
	d.NewWrite()
	return nil
}

// Read reads from a single key
func (d *DB) Read(key string) (interface{}, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.database[key]; !ok {
		return nil, errors.New("key doesn't exist")
	}
	return d.database[key], nil
}

// ReadMany returns multiple keys
func (d *DB) ReadMany(keys ...string) map[string]interface{} {
	d.mu.Lock()
	defer d.mu.Unlock()
	results := make(map[string]interface{})
	for _, k := range keys {
		if v, ok := d.database[k]; !ok {
			results[k] = nil
		} else {
			results[k] = v
		}
	}
	return results
}

// ReadAll returns all entries from DB
func (d *DB) ReadAll() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	str := ""
	for k, v := range d.database {
		str += fmt.Sprintf("%v => %v\n", k, v)
	}
	return str
}

// Update updates a single entry
func (d *DB) Update(key string, value interface{}) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.database[key]; !ok {
		return errors.New("key doesn't exist")
	}

	d.database[key] = value
	d.NewWrite()
	return nil
}

// Delete deletes a single entry
func (d *DB) Delete(key string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.database[key]; !ok {
		return errors.New("key doesn't exist")
	}

	delete(d.database, key)
	d.NewWrite()
	return nil
}

func (d *DB) DeleteMany(keys ...string) map[string]interface{} {
	d.mu.Lock()
	defer d.mu.Unlock()

	res := make(map[string]interface{})

	del := make(map[string]bool, 1)
	del["deleted"] = true

	err := make(map[string]string, 1)
	err["error"] = "key doesn't exist"

	for _, key := range keys {
		if _, ok := d.database[key]; !ok {
			res[key] = err
		} else {
			delete(d.database, key)
			res[key] = del
		}

	}
	d.NewWrite()
	return res
}
