package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"github.com/piquette/finance-go/quote"
)

// Portfolio instance variables must begin with capital letter to be exported by json
type Portfolio struct {
	Name string
	// Must contain `.json`
	Directory string
	// the sum Equity and Cash
	Value float64
	// The sum of all stocks cost
	Cost float64
	// The sum of all stocks equity
	Equity float64
	// TotalGainLoss equals Equity - Cost
	TotalGainLoss float64
	// TotalGainLossPrcnt equals Equity - Cost
	TotalGainLossPrcnt float64
	// amount of uninvested cash
	Cash float64
	// stock ticker is key, Stock object is value
	Positions map[string]Stock
	// date is key, list of entry object is value
	// date format is DD-MM-YYY
	History map[string][]Entry
}

// Stock instance variables must begin with capital letter to be exported by json
type Stock struct {
	// CompanyName string
	Shares float64
	// Average price each share was bought at
	AvgPrice float64
	// retrieved from GoogleFinance/yahoo finance
	LatestPrice float64
	// Amount of Cash invested in stock
	TotalCost float64
	// latest price times the number of shares
	Equity float64
	// GainLoss is the Equity - TotalCost
	GainLoss float64
	// GainLossPrcnt is the percentage of earning/loss per share
	// GainLossPrcnt = [(LatestPrice / AvgPrice) -1] * 100
	GainLossPrcnt float64
	// PrcntOfPort is how much stock equity makes up from total portfolio value
	PrcntOfPort float64
}

// Entry instance variables must begin with capital letter to be exported by json
type Entry struct {
	// Type of operation (buy, sell, deposit, withdraw)
	Type       string
	Ticker     string
	Shares     float64
	UnitPrice  float64
	OrderTotal float64
}

// getStock is not exported
func (port *Portfolio) getStock(ticker string) *Stock {
	myStock := port.Positions[ticker]
	return &myStock
}

// Buy is exported function, retruns no value
func (port *Portfolio) Buy(ticker string, shares float64, price float64) {
	port.RefreshData()
	// fmt.Println("You bought", shares, ticker, "shares at", price, "dollars per share")
	// Calculate order total
	orderTotal := shares * price
	// Check if there is enough cash
	if port.Cash < orderTotal {
		Check(errors.New("Can't execute order: insufficient funds"))
	}
	// get the pointer to the stock
	stock := port.getStock(ticker)
	// add purchase shares
	stock.Shares += shares
	// Add the orderTotal to the stock's cost
	stock.TotalCost += orderTotal
	// Calculate new average price
	stock.AvgPrice = stock.TotalCost / stock.Shares

	// store the modified stock in the portfolio
	port.Positions[ticker] = *stock
	// Remove used cash from portfolio
	port.Cash -= orderTotal

	// Create and add a new entry
	today := time.Now().Format("01-02-2006")
	newEntry := Entry{"Buy", ticker, shares, price, orderTotal}
	port.History[today] = append(port.History[today], newEntry)

	port.RefreshData()
}

// Sell is exported function, retruns no value
func (port *Portfolio) Sell(ticker string, shares float64, price float64) {
	port.RefreshData()
	// fmt.Println("You sold ", shares, ticker, " shares at ", price, " dollars per share")
	// Calculate order total
	orderTotal := shares * price
	// get the pointer to the stock
	stock := port.getStock(ticker)
	// Check if there is enough shares to sell
	if stock.Shares < shares {
		Check(errors.New("Can't sell non-existent shares"))
	}
	// subtract sold shares
	stock.Shares -= shares
	// Adjust the total cost to reflect current shares
	stock.TotalCost = stock.Shares * stock.AvgPrice

	// store the modified stock in the portfolio
	port.Positions[ticker] = *stock
	// Add resulting cash to portfolio
	port.Cash += orderTotal

	// Create and add a new entry
	today := time.Now().Format("01-02-2006")
	newEntry := Entry{"Sell", ticker, shares, price, orderTotal}
	port.History[today] = append(port.History[today], newEntry)

	port.RefreshData()
}

// Deposit is exported function, retruns no value
// Add to the cash instance variable and write to json
func (port *Portfolio) Deposit(cash float64) {
	if cash <= 0 {
		Check(errors.New("Invalid deposit amount, zero or negative"))
	}
	// fmt.Println("You deposited ", cash, "in your portfolio")
	port.Cash += cash
	port.Value += cash
	// Create and add a new entry
	today := time.Now().Format("01-02-2006")
	newEntry := Entry{"Deposit", "N/A", cash, 1.0, cash}
	port.History[today] = append(port.History[today], newEntry)

	port.RefreshData()
}

// Withdraw is exported function, retruns no value
// Subtract from the cash instance variable and write to json
func (port *Portfolio) Withdraw(cash float64) {
	if cash <= 0 {
		Check(errors.New("Invalid witdraw amount, zero or negative"))
	}
	if cash > port.Cash {
		Check(errors.New("Withdraw amount can't exceed portfolio balance"))
	}
	// fmt.Println("You withdrew ", cash, "from your portfolio")
	port.Cash -= cash
	port.Value -= cash

	// Create and add a new entry
	today := time.Now().Format("01-02-2006")
	newEntry := Entry{"Withdraw", "N/A", cash, 1.0, cash}
	port.History[today] = append(port.History[today], newEntry)
	port.RefreshData()
}

// RefreshData is exported function, retruns no value
func (port *Portfolio) RefreshData() {
	// TODO: go over this more
	// Handle running it on empty values
	if len(port.Positions) == 0 {
		return
	}
	totalCost := 0.0
	totalEquity := 0.0
	// Iterate through positions
	// STOCK related changes
	for ticker := range port.Positions {
		// get stock pointer
		currentStock := port.getStock(ticker)
		if currentStock.Shares == 0 {
			// Stock has no shares and it must be deleted
			delete(port.Positions, ticker)
		} else {
			// get latest price from finance package
			currentStock.LatestPrice = getStockData(ticker)
			// calculate equity
			currentStock.Equity = currentStock.Shares * currentStock.LatestPrice
			// Calculate gain or loss
			currentStock.GainLossPrcnt = ((currentStock.LatestPrice / currentStock.AvgPrice) - 1)
			currentStock.GainLoss = currentStock.Equity - currentStock.TotalCost
			port.Positions[ticker] = *currentStock
			totalEquity += currentStock.Equity
			totalCost += currentStock.TotalCost
		}

	}

	// PORTFOLIO RELATED CHANGES
	port.Cost = totalCost
	port.Equity = totalEquity
	port.TotalGainLoss = port.Equity - port.Cost
	port.TotalGainLossPrcnt = ((port.Equity / port.Cost) - 1)
	port.Value = port.Cash + totalEquity

	// STOCK RTELATED CHANGES AGAIN
	// iterate again to find percent of portfolio value
	// TODO: go over this more
	for ticker := range port.Positions {
		currentStock := port.getStock(ticker)
		currentStock.PrcntOfPort = (currentStock.Equity / port.Value)
		port.Positions[ticker] = *currentStock
	}
	// fmt.Println("data refreshed")
}

// StoreData writes the content of the Portfolio to the directory
// It overwrites any previous contents
// assumes that the directory contains `.json`
func (port *Portfolio) StoreData() {
	data, _ := json.Marshal(port)
	err := ioutil.WriteFile("./portfolio/"+port.Directory, data, 0644)
	Check(err)
}

// Check function is exported
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func getStockData(ticker string) float64 {
	ticker = strings.ToUpper(ticker)
	q, err := quote.Get(ticker)
	Check(err)
	// handle non existent stocks
	if q == nil {
		Check(errors.New("Stock does not exist"))
	}
	return q.RegularMarketPrice
}

// TODO: add Clear(): erases everything on json and writes `{}`
