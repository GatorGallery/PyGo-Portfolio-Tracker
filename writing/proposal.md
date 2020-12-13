# PyGo Stock Portfolio Builder

## Noor Buchi, Anh Tran

## Description of the Main Idea

This tool allows users to build and keep track of a stock portfolio. It provides ways to visualize positions and give the user the information they need to make investment decisions. One of the main operations that this tool allows is the buying and selling of stock shares. The user can specify the company, the number of shares, and the average price that the stock was traded at. Following that, the tool will store this information using `json` file, taking into consideration what was stored before and adding/subtracting accordingly. In addition to the previously mentioned general operations, the tool will create graphs that show the makeup of the portfolio using a pie chart. It will also use real-time data to calculate and graph the current profit and equity of the portfolio. This program will use two programming languages to accomplish these tasks. In the front-end, it will use Python and packages like Streamlit, Pandas, and many others to create a graphical interface and to graph the data. In the back-end, the tool will use Go to accept command-line arguments, store and update `json`, and retrieve stock data using `GoogleFinance`. The two languages will be connected through system commands where an action in Streamlit will cause the Python front end to run the Go program through its command line interface. Overall, the tool provides an intuitive and simple way to keep track of a stock portfolio.

## Description of the Tasks that You Will Complete

There are two main parts in this project. The first part is data collection and calculation. This aspect of the tool is independent of the second part and it can operate through a command line interface. It will be written in Go and it should:
- Parse basic command line arguments
  - Include tokens for buying, selling, and refreshing data
  - Example structure: `go run myProgram.go -O Buy -T TSLA -S 3.2 -P 480.43`
  - `O` stands for operation, `T` stands for ticker, `S` stands for shares, and `P` stands for price
- Write to and update a `json` file
  - Totals and positions will be stored in the `json` file to facilitate calculations
  - Example structure:

  ```json
  {

    "UNINVESTEDCASH" : 500,
    "INVESTEDCASH": 2000,
    "TOTALEQUITY": 2500,

    "TSLA": {
        "shares": 4,
        "avgPrice": 450,
        "latestPrice": 500,
        "totalEquity" : 2000,
        "percentage" : 65
    },

    "MSFT": {
        "shares": 8,
        "avgPrice": 150,
        "latestPrice": 200,
        "totalEquity" : 2000,
        "percentage" : 35
    }

  }
  ```

- Utilize `GoogleFinance` API to get the latest stock market prices
  - Iterate through the database and calculate profit and equity
  - Implement a refresh functionality that updates the database with the newest prices

The second part of the tool will focus on front-end interface and graphing of data. It will depend on the presence of data collected using the first part. Additionally, it will be responsible for making requests to part 1 through system commands. The second section will be implemented in Python and have these features:
- Runs the Golang program based on user input
- A Streamlit web interface
  - Pages for buy and sell operations with text boxes to input information and buttons to execute orders
  - Page to display the data retrieved from the `json` file
    - Includes a table, stacked bar chart, and pie chart.

## Detailed Plan for Your Team to Complete the Project

The first step of completing this project will be to implement the part1 and to ensure that data retrieval and storage is working as expected. It's also be important to make sure that the back-end part will be able to handle user input and catch any exceptions. The initial work will focus on creating a command line interface and understanding `GoogleFinance` and its API. Once a solid back-end foundation is created, we plan to start creating the web interface and come up with ways to iterate and graph the data. This will require some basic knowledge in Streamlit as well as data frame management tools such as Pandas. This step will also have to take into consideration user input and handle empty or incomplete orders.
