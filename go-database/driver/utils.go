package driver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jcelliott/lumber"
)

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

// to export the function always make sure that function name starts with  caps
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

func Stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

func (d *Driver) GetOrCreateMutex(collection string) *sync.Mutex {

	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}
	return m
}

func AddNewCollection(db *Driver, reader *bufio.Reader) {

	fmt.Print("Enter new collection name: ")
	collectionName, _ := reader.ReadString('\n')
	collectionName = strings.TrimSpace(collectionName)

	mutex := db.GetOrCreateMutex(collectionName)
	mutex.Lock()
	defer mutex.Unlock()
	collectionData := make(map[string]string)

	for {
		fmt.Print("Enter a key (or 'exit' to finish): ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		if key == "exit" {
			break
		}

		fmt.Print("Enter value: ")
		value, _ := reader.ReadString('\n')
		value = strings.TrimSpace(value)

		collectionData[key] = value
	}

	dataJSON, err := json.Marshal(collectionData)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	//err = db.Write(collectionName, dataJSON)
	if err := os.WriteFile(".//Collections//"+collectionName+".json", dataJSON, 0644); err != nil {
		fmt.Println("Error writing to database: ", err)
	}

}

func ListCollections(db *Driver) {
	// Read the contents of the directory
	files, err := os.ReadDir(".//Collections//")
	if err != nil {
		fmt.Println("Error reading collections directory:", err)
		return
	}

	fmt.Println("Available Collections:")
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			// Print the collection name without the '.json' extension
			fmt.Println(strings.TrimSuffix(file.Name(), ".json"))
		}
	}
}

func ViewCollection(db *Driver, reader *bufio.Reader) {
	fmt.Print("Enter the name of the collection to view : ")
	collectionName, _ := reader.ReadString('\n')
	collectionName = strings.TrimSpace(collectionName)

	mutex := db.GetOrCreateMutex(collectionName)
	mutex.Lock()
	defer mutex.Unlock()

	b, err := os.ReadFile(".//Collections//" + collectionName + ".json")
	if err != nil {
		fmt.Println("Error reading collection:", err)
		return
	}

	//fmt.Println(b)
	// Unmarshal JSON data into a map
	collectionData := make(map[string]interface{})
	er := json.Unmarshal(b, &collectionData)
	if er != nil {
		fmt.Println("Error unmarshalling data:", err)
		return
	}
	updatedDataJSON, err := json.MarshalIndent(collectionData, "", "\t")
	fmt.Println(string(updatedDataJSON))
}

func UpdateCollection(db *Driver, reader *bufio.Reader) {
	fmt.Print("Enter the name of the collection to update: ")
	collectionName, _ := reader.ReadString('\n')
	collectionName = strings.TrimSpace(collectionName)

	mutex := db.GetOrCreateMutex(collectionName)
	mutex.Lock()
	defer mutex.Unlock()

	b, err := os.ReadFile(".//Collections//" + collectionName + ".json")
	if err != nil {
		fmt.Println("Error reading collection:", err)
		return
	}

	// Unmarshal JSON data into a map
	collectionData := make(map[string]interface{})
	er := json.Unmarshal(b, &collectionData)
	if er != nil {
		fmt.Println("Error unmarshalling data:", err)
		return
	}

	// Updating the collection
	for {
		fmt.Print("Enter a key to update (or 'exit' to finish): ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		if key == "exit" {
			break
		}

		fmt.Print("Enter new value: ")
		value, _ := reader.ReadString('\n')
		value = strings.TrimSpace(value)

		collectionData[key] = value // Update the value for the key
	}

	// Marshal the updated map back to JSON
	updatedDataJSON, err := json.Marshal(collectionData)
	if err != nil {
		fmt.Println("Error marshalling updated data:", err)
		return
	}

	// Write the updated JSON data back to the database
	if err := os.WriteFile(".//Collections//"+collectionName+".json", updatedDataJSON, 0644); err != nil {
		fmt.Println("Error writing updated data to database:", err)
	}
}

func DeleteCollection(db *Driver, reader *bufio.Reader) {
	fmt.Print("Enter the name of the collection to delete: ")
	collectionName, _ := reader.ReadString('\n')
	collectionName = strings.TrimSpace(collectionName)
	mutex := db.GetOrCreateMutex(collectionName)
	mutex.Lock()
	defer mutex.Unlock()

	// Construct the file path for the collection
	filePath := ".//Collections//" + collectionName + ".json"

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Collection does not exist.")
		return
	}

	// Confirm deletion
	fmt.Printf("Are you sure you want to delete the collection '%s'? (yes/no): ", collectionName)
	confirmation, _ := reader.ReadString('\n')
	confirmation = strings.TrimSpace(confirmation)

	if confirmation != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}

	// Delete the file
	err := os.Remove(filePath)
	if err != nil {
		fmt.Println("Error deleting collection:", err)
	} else {
		fmt.Println("Collection deleted successfully.")
	}
}
