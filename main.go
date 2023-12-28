package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

type Address struct {
	City    string
	Pincode json.Number
	state   string
	country string
}

type User struct {
	Name    string
	Age     json.Number
	Company string
	Address Address
}

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

	/*
		This function creates a new driver (that interacts with db). It attaches a logger(if it doesn't exist)
	*/
	dir = filepath.Clean(dir)
	opts := Options{}

	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}
	// Create a driver , as per the parameters used in the driver structure
	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}
	// Check if the database exists or not using os.stat
	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil // driver exits , error is nil
	}

	opts.Logger.Debug("Creating the database at '%s' \n ", dir)

	return &driver, os.MkdirAll(dir, 0755) //0755 is access permission
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
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

// Collection represents the directory where json files (records) would be stored
func (d *Driver) Write(collection, resource string, v interface{}) error {

	if collection == "" {
		fmt.Println("Collection not found.No place to save data ")
	}
	if resource == "" {
		fmt.Println("Resource not found. No such record available")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock() // frees the mutex lock after entire function has been executed.

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "\t") // format the json
	if err != nil {
		return err
	}
	b = append(b, byte('\n'))
	if err := os.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}
	return os.Rename(tmpPath, fnlPath)
}

func (d *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("Missing Collection, Unable to read")
	}
	dir := filepath.Join(d.dir, collection)

	if _, err := stat(dir); err != nil { // stat function checks for collection of directory's existence
		return nil, err
	}
	files, _ := os.ReadDir(dir)

	var records []string
	for _, file := range files {
		b, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		records = append(records, string(b))
	}
	return records, nil
}

func (d *Driver) Delete(collection, reource string) error {

	path := filepath.Join(collection, reource)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := stat(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("unable to find file or directory named %v\n", path)
	case fi.Mode().IsDir(): // to remove entire folder
		return os.RemoveAll(dir)
	case fi.Mode().IsRegular(): // remove all files inside the folder
		return os.RemoveAll(dir + ".json")
	}

	return nil
}

func (d *Driver) Read(collection, resource string, v interface{}) error {

	if collection == "" {
		fmt.Println("Collection not found. No place to save data ")
	}
	if resource == "" {
		fmt.Println("Resource not found. No such record available")
	}
	record := filepath.Join(d.dir, collection)

	if _, err := stat(record); err != nil {
		return err
	}
	b, err := os.ReadFile(record + ".json")
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &v)
}

const Version = "1.0.0" //version of db

func main() {
	dir := "./"
	db, err := New(dir, nil)

	if err != nil {
		fmt.Println("Error ", err)
	}
	employees := []User{
		{"Mayank", "27", "Google", Address{"SF", "30020", "CA", "USA"}},
		{"Raj", "27", "Google", Address{"LA", "30021", "CA", "USA"}},
		{"John", "24", "MSF", Address{"NYC", "30022", "NY", "IND"}},
		{"Doe", "22", "FBK", Address{"CH", "30023", "IL", "IND"}},
		{"Paul", "21", "NKE", Address{"MI", "30024", "FL", "CHN"}},
		{"Mike", "29", "APL", Address{"AS", "30025", "TX", "CHN"}},
	}

	for _, value := range employees {
		db.Write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Company: value.Company,
			Address: value.Address,
		})
	}
	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error ", err)
	}
	fmt.Println(records)

	allrecords := []User{}

	for _, f := range records {
		employee := User{}
		if err := json.Unmarshal([]byte(f), &employee); err != nil {
			fmt.Println("error ", err)
		}
		allrecords = append(allrecords, employee)
	}
	fmt.Println(allrecords) // fetch all records , add to a slice and print

	// if err := db.Delete("users", "John"); err != nil {
	// 	fmt.Println("Error ", err)
	// }

	// if err := db.Delete("users", ""); err != nil {
	// 	fmt.Println("Error ", err)
	// }
}
