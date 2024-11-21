package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const Version = "1.0.0"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	opts := Options{}

	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		dir:     dir,
		log:     opts.Logger,
		mutexes: make(map[string]*sync.Mutex),
	}

	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("using '%s' (database already exists)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating the database at '%s'\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("collection is required (can't save record!)")
	}

	if resource == "" {
		return fmt.Errorf("resource is required (can't save record!)")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, fnlPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("collection is required (can't read record!)")
	}

	if resource == "" {
		return fmt.Errorf("resource is required (can't read record!)")
	}

	record := filepath.Join(d.dir, collection, resource+".json")

	if _, err := stat(record); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(record)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)
}

func (d *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection is required (can't read all records!)")
	}

	dir := filepath.Join(d.dir, collection)

	if _, err := stat(dir); err != nil {
		return nil, err
	}

	files, _ := ioutil.ReadDir(dir)

	var records []string

	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return records, err
		}

		records = append(records, string(b))
	}
	return records, nil
}

func (d *Driver) Delete() error {

}
func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {

	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func main() {
	dir := "./"

	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	friends = []User{
		{"Yash", "21", "0987654321", "Google", Address{"Bangalore", "Karnataka", "India", "560001"}},
		{"Siddharth", "22", "0987654444", "Microsoft", Address{"Bangalore", "Karnataka", "India", "560001"}},
		{"Manav", "20", "0987653333", "Adobe", Address{"Bangalore", "Karnataka", "India", "560001"}},
		{"Shailendra", "23", "0987652222", "Netflix", Address{"Bangalore", "Karnataka", "India", "560001"}},
		{"Aman", "21", "0987651111", "Apple", Address{"Gwalior", "Madhya Pradesh", "India", "474011"}},
		{"Priyanshu", "21", "0987650000", "Amazon", Address{"Bangalore", "Karnataka", "India", "560001"}},
	}

	for _, value := range friends {
		db.Write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}

	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println(records)

	allusers := []User{}

	for _, f := range records {
		userFound := User{}
		if err := json.Unmarshal([]byte(f), userFound); err != nil {
			fmt.Println("Error", err)
		}
		allusers = append(allusers, userFound)
	}
	fmt.Println((allusers))

	// if err := db.Delete("users", "Yash"); err != nil {
	// 	fmt.Println("Error", err)
	// }

	// if err := db.Delete("users", ""); err != nil {
	// 	fmt.Println("Error", err)
	// }

}
