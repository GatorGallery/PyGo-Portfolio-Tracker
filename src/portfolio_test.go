package main

// Command to run test cases
// go test -v -cover ./src
// MUST BE IN REPO ROOT

import (
	"fmt"
	"testing"
	"time"
)

func getEmptyPort() *Portfolio {
	tempPort := Portfolio{}
	port := &tempPort
	port.Name = "EmptyPort"
	port.Directory = "EmptyPort.json"
	port.Positions = map[string]Stock{}
	port.History = map[string][]Entry{}
	return port
}

func equalEntry(entry1 Entry, entry2 Entry) bool {
	if entry1.Type != entry2.Type {
		return false
	}
	if entry1.Ticker != entry2.Ticker {
		return false
	}
	if entry1.Shares != entry2.Shares {
		return false
	}
	if entry1.UnitPrice != entry2.UnitPrice {
		return false
	}
	if entry1.OrderTotal != entry2.OrderTotal {
		return false
	}
	return true
}

// Does not check for latest price
func equalPosition(stock1 Stock, stock2 Stock) bool {
	if stock1.Shares != stock2.Shares {
		return false
	}
	if stock1.AvgPrice != stock2.AvgPrice {
		return false
	}
	if stock1.TotalCost != stock2.TotalCost {
		return false
	}
	return true
}

func TestGetStockData(t *testing.T) {
	stockList := []string{"tsla", "TsLa", "MSFT", "msft", "Zm", "AMZN"}
	for _, stock := range stockList {
		stockPrice := getStockData(stock)
		if stockPrice <= 0 {
			t.Errorf("invalid stock price %v", stockPrice)
		}
	}
}

func TestFileExists(t *testing.T) {
	var tests = []struct {
		directory string
		expected  bool
	}{
		{"web_interface.py", true},
		{"stockHandler.go", true},
		{"portfolio.go", true},
		{"randomFile.json", false},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("Test %v", tt.directory)
		t.Run(testname, func(t *testing.T) {
			ans := fileExists(tt.directory)
			if ans != tt.expected {
				t.Errorf("Directory %v not found", tt.directory)
			}
		})
	}
}

func TestPortDeposit(t *testing.T) {
	port := getEmptyPort()
	// Operation done on portfolio
	port.Deposit(450.33)
	// Store actual values/results
	cashResult := port.Cash
	valueResult := port.Value
	// Check cash change
	if cashResult != 450.33 {
		t.Errorf("Found Cash value %v , expecting %v", cashResult, 450.33)
	}
	// Check portfolio value change
	if valueResult != 450.33 {
		t.Errorf("Found Cash value %v , expecting %v", cashResult, 450.33)
	}
	today := time.Now().Format("01-02-2006")
	entry := port.History[today][0]
	expectedEntry := Entry{"Deposit", "N/A", 450.33, 1.0, 450.33}
	if !equalEntry(entry, expectedEntry) {
		t.Errorf("Mismatched entries")
	}
}

func TestPortWithdraw(t *testing.T) {
	port := getEmptyPort()
	// Operation done on portfolio
	port.Deposit(500)
	port.Withdraw(250)
	// Store actual values/results
	cashResult := port.Cash
	valueResult := port.Value
	// Check cash change
	if cashResult != 250 {
		t.Errorf("Found Cash value %v , expecting %v", cashResult, 250)
	}
	// Check portfolio value change
	if valueResult != 250 {
		t.Errorf("Found Cash value %v , expecting %v", cashResult, 250)
	}
	today := time.Now().Format("01-02-2006")
	entry1 := port.History[today][0]
	expectedEntry1 := Entry{"Deposit", "N/A", 500, 1.0, 500}
	if !equalEntry(entry1, expectedEntry1) {
		t.Errorf("Mismatched entries")
	}
	entry2 := port.History[today][1]
	expectedEntry2 := Entry{"Withdraw", "N/A", 250, 1.0, 250}
	if !equalEntry(entry2, expectedEntry2) {
		t.Errorf("Mismatched entries")
	}
}

func TestPortBuy(t *testing.T) {
	port := getEmptyPort()
	// Operation done on portfolio
	port.Deposit(1000)
	port.Buy("tsla", 3, 100)
	port.RefreshData()
	// Store actual values/results
	cashResult := port.Cash
	stockResult := port.Positions["tsla"]
	costResult := port.Cost
	// Check cash change
	if cashResult != 700 {
		t.Errorf("Found Cash value %v , expecting %v", cashResult, 700)
	}
	// Check portfolio value change
	if costResult != 300 {
		t.Errorf("Found Cost value %v , expecting %v", costResult, 300)
	}
	today := time.Now().Format("01-02-2006")
	// get the last entry in the list
	entry := port.History[today][len(port.History[today])-1]
	expectedEntry := Entry{"Buy", "tsla", 3, 100, 300}
	if !equalEntry(entry, expectedEntry) {
		t.Errorf("Mismatched entries")
		fmt.Println("Expected entry", expectedEntry)
		fmt.Println("Found entry", entry)
	}

	// Check the position validity
	expectedPosition := Stock{3, 100, 0, 300, 0, 0, 0, 0}
	if !equalPosition(stockResult, expectedPosition) {
		t.Errorf("Mismatched positions")
		fmt.Println("Expected entry", expectedPosition)
		fmt.Println("Found entry", stockResult)
	}
}

func TestPortDoubleBuy(t *testing.T) {
	port := getEmptyPort()
	// Operation done on portfolio
	port.Deposit(1000)
	port.Buy("tsla", 3, 100)
	port.Buy("tsla", 3, 200)
	port.RefreshData()
	// Store actual values/results
	cashResult := port.Cash
	stockResult := port.Positions["tsla"]
	costResult := port.Cost
	// Check cash change
	if cashResult != 100 {
		t.Errorf("Found Cash value %v , expecting %v", cashResult, 100)
	}
	// Check portfolio value change
	if costResult != 900 {
		t.Errorf("Found Cost value %v , expecting %v", costResult, 900)
	}
	today := time.Now().Format("01-02-2006")
	// get the last entry in the list
	entry := port.History[today][len(port.History[today])-1]
	expectedEntry := Entry{"Buy", "tsla", 3, 200, 600}
	if !equalEntry(entry, expectedEntry) {
		t.Errorf("Mismatched entries")
		fmt.Println("Expected entry", expectedEntry)
		fmt.Println("Found entry", entry)
	}

	// Check the position validity
	expectedPosition := Stock{6, 150, 0, 900, 0, 0, 0, 0}
	if !equalPosition(stockResult, expectedPosition) {
		t.Errorf("Mismatched positions")
		fmt.Println("Expected entry", expectedPosition)
		fmt.Println("Found entry", stockResult)
	}
}

func TestPortSell(t *testing.T) {
	port := getEmptyPort()
	// Operation done on portfolio
	port.Deposit(1000)
	port.Buy("tsla", 3, 100)
	port.Buy("tsla", 3, 200)
	port.Sell("tsla", 3, 500)
	port.RefreshData()
	// Store actual values/results
	cashResult := port.Cash
	stockResult := port.getStock("tsla")
	costResult := port.Cost
	// Check cash change
	expectedCash := 1600.0
	expectedCost := 450.0
	if cashResult != expectedCash {
		t.Errorf("Found Cash value %v , expecting %v", cashResult, expectedCash)
	}
	// Check portfolio value change
	if costResult != expectedCost {
		t.Errorf("Found Cost value %v , expecting %v", costResult, expectedCost)
	}
	today := time.Now().Format("01-02-2006")
	// get the last entry in the list
	entry := port.History[today][len(port.History[today])-1]
	expectedEntry := Entry{"Sell", "tsla", 3, 500, 1500}
	if !equalEntry(entry, expectedEntry) {
		t.Errorf("Mismatched entries")
		fmt.Println("Expected entry", expectedEntry)
		fmt.Println("Found entry", entry)
	}

	// Check the position validity
	expectedPosition := Stock{3, 150, 0, 450, 0, 0, 0, 0}
	if !equalPosition(*stockResult, expectedPosition) {
		t.Errorf("Mismatched positions")
		fmt.Println("Expected entry", expectedPosition)
		fmt.Println("Found entry", stockResult)
	}
}
