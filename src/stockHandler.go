package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"strings"
)

// assumes the existence of a json file and retrieves the data
// Returns Portfolio pointer
func retrieveData(directory string) *Portfolio {
	data, err := ioutil.ReadFile(directory)
	Check(err)
	var retreivedPorfolio Portfolio
	err2 := json.Unmarshal(data, &retreivedPorfolio)
	Check(err2)
	return &retreivedPorfolio
}

// Check the existence of a file before using
// NOTE: it depends on the directory that the program is running from
func fileExists(directory string) bool {
	info, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Parse command line arguments and redirect the flow of the program
func parseArgs() error {
	// NOTE: all these flags are pointers and must be dereferenced before usage
	// Flag for the file name, has string value
	file := flag.String("f", "", "name of portfolio file")
	// Flag for the operation, has string value
	operation := flag.String("o", "", "name of operation")
	// Flag for the ticker, has string value
	ticker := flag.String("t", "", "ticker of company")
	// Flag for the shares, float64 value
	shares := flag.Float64("s", 0.0, "number of shares")
	// Flag for the price, float64 value
	price := flag.Float64("p", 0.0, "price")
	// Flag that triggers a price refresh
	// if -r is present in arguments, the value of refresh becomes true
	refresh := flag.Bool("r", false, "refresh")

	flag.Parse()

	// get the file directory
	fileName := *file
	if fileName == "" || !strings.Contains(fileName, ".json") {
		return errors.New("Invalid file")
	}
	// handle creation or retrieval of file
	port := loadFile(fileName)

	// Convert the ticker to uppercase
	newTicker := strings.ToUpper(*ticker)
	switch op := strings.ToLower(*operation); op {
	case "buy":
		// Handle 0 or negative arguments
		if (*shares <= 0) || (*price <= 0) {
			return errors.New("Invalid shares/price value")
		}
		port.Buy(newTicker, *shares, *price)
		port.StoreData()

	case "sell":
		// Handle 0 or negative arguments
		if (*shares <= 0) || (*price <= 0) {
			return errors.New("Invalid shares/price value")
		}
		port.Sell(newTicker, *shares, *price)
		port.StoreData()

	case "deposit":
		// NOTE: deposit only considers the shares argument
		if *shares <= 0 {
			return errors.New("Invalid deposit amount")
		}
		port.Deposit(*shares)
		port.StoreData()

	case "withdraw":
		// NOTE: witdraw only considers the shares argument
		if *shares <= 0 {
			return errors.New("Invalid withdraw amount")
		}
		port.Withdraw(*shares)
		port.StoreData()

	// Allow blank operation in case of simple refresh
	case "":
		break
	default:
		return errors.New("Unknown operation")
	}
	// Check refresh arg
	if *refresh {
		port.RefreshData()
		port.StoreData()
	}
	return nil
}

func loadFile(directory string) *Portfolio {
	location := "./portfolio/" + directory
	// NOTE: port is a pointer
	var port *Portfolio
	if fileExists(location) {
		// fmt.Println("File already exists")
		port = retrieveData(location)
	} else {
		// fmt.Println("Creating new Portfolio...")
		// Portfolio has no data
		tempPort := Portfolio{}
		port = &tempPort
		port.Name = directory
		port.Directory = directory
		port.Positions = map[string]Stock{}
		port.History = map[string][]Entry{}
	}
	return port
}

func main() {
	// TODO: add operation history storage

	err := parseArgs()
	Check(err)

}
