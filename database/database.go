package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// DBCrud is implemented by DB
type DBCrud interface {
	Connect()
	Create()
	Read()
	Update()
	Delete()
	ReadMany()
	Disconnect()
}

// New initializes a database to a given location and sets it's internal DB to an empty map
func New(location string, memory bool) *DB {
	return &DB{
		location: location,
		database: make(map[string]interface{}),
		memory:   memory,
	}
}

// DB represents the database struct
type DB struct {
	location string
	database map[string]interface{}
	memory   bool
}

// Connect connects to file and saves it's contents to database field
func (d *DB) Connect() error {
	if d.database == nil {
		return errors.New("must Init() the database first")
	}

	if d.location == "" {
		return nil
	}

	// try to read file
	_, err := ioutil.ReadFile(d.location)

	// create new if doesn't exist
	if err != nil {
		f, err := os.Create(d.location)

		// return in case of error
		if err != nil {
			return errors.New(err.Error())
		}
		// write empty valid json to file
		// TODO: error check
		_, _ = f.WriteString("{}")
		f.Close()
	}

	// Read newly created file
	content, err := ioutil.ReadFile(d.location)
	if err != nil {
		return errors.New(err.Error())
	}

	// Unmarshal it's contents into in-memory database
	err = json.Unmarshal(content, &d.database)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// Disconnect encodes database with json and saves it to location if provided
func (d *DB) Disconnect() error {
	if len(d.database) == 0 || d.location == "" || d.memory {
		return nil
	}

	jsonBody, err := json.Marshal(d.database)
	if err != nil {
		return errors.New(err.Error())
	}

	f, err := os.OpenFile(d.location, os.O_WRONLY, 0664)
	if err != nil {
		return errors.New(err.Error())
	}

	defer f.Close()

	_, err = f.Write(jsonBody)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// Create creates a new record
func (d *DB) Create(key string, value interface{}) error {
	if _, ok := d.database[key]; ok {
		return errors.New("key already exists")
	}
	d.database[key] = value
	return nil
}

// Read reads from a single key
func (d *DB) Read(key string) (interface{}, error) {
	if _, ok := d.database[key]; !ok {
		return nil, errors.New("key doesn't exist")
	}
	return d.database[key], nil
}

// ReadMany returns multiple keys
func (d *DB) ReadMany(keys ...string) map[string]interface{} {
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
	str := ""
	for k, v := range d.database {
		str += fmt.Sprintf("%v => %v\n", k, v)
	}
	return str
}

// Update updates a single entry
func (d *DB) Update(key string, value interface{}) error {
	if _, ok := d.database[key]; !ok {
		return errors.New("key doesn't exist")
	}

	d.database[key] = value
	return nil
}

// Delete deletes a single entry
func (d *DB) Delete(key string) error {
	if _, ok := d.database[key]; !ok {
		return errors.New("key doesn't exist")
	}

	delete(d.database, key)
	return nil
}
