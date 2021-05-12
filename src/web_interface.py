"""Web Interface for building a stock portfolio."""

import streamlit as st
import pandas as pd
import matplotlib.pyplot as plt
import altair as alt
import os
import subprocess
from os import path
import json
import shlex
import yfinance as yf
import datetime as dt


def web_interface():
    """Execute the web interface."""
    st.set_page_config(page_title="PyGo Portfolio Builder",
                       layout='wide', initial_sidebar_state='auto')
    st.title("Welcome to PyGo Portfolio Builder")
    # Build the golang files
    stdout, stderr = run_command(
        "go build src/stockHandler.go src/portfolio.go")
    if not(stdout == ""):
        st.warning(stdout)
    if not(stderr == ""):
        st.warning(stderr)

    if not path.isdir("./portfolio/"):
        os.mkdir("portfolio")

    # get the json files in the portfolio folder
    portfolio_list = []
    for file in os.listdir("portfolio"):
        if file.endswith(".json"):
            portfolio_list.append(file)

    port_path = st.sidebar.selectbox(
        "Please select a portfolio to open, Note: json files must be placed in portfolio folder", options=portfolio_list)

    # Check if data exists
    if portfolio_list == []:
        st.warning("No portfolio found, create new one")
        # redirect program if no file found
        create_portfolio()
        if st.button("refresh data", key="First creation"):
            refresh(port_path)
            st._RerunData()
        st.stop()

    # Open and load data
    with open(f"portfolio/{port_path}") as f:
        data = json.load(f)

    portfolio_name = data["Name"]
    st.markdown(f"## {portfolio_name}:")

    page_selection = st.sidebar.selectbox("Please select a page", options=[
                                          "Chart", "Graphs", "History", "Company Search", "Run Tests"])

    if st.button("refresh data"):
        refresh(port_path)
        st._RerunData()

    if page_selection == "Chart":
        show_chart(data)

    if page_selection == "Graphs":
        show_graph(data)

    if page_selection == "History":
        show_history(data)

    if page_selection == "Run Tests":
        run_tests()

    if page_selection == "Company Search":
        company_search()

    show_sidebar(port_path, data)


def company_search():
    company = st.text_input('Please input the stock ticker you would like to look for:')
    tick = yf.Ticker(company)
    st.header("Company summary:")
    st.write(tick.info["longBusinessSummary"])
    st.header("Company stock information:")
    st.write("Company Sector: " + tick.info["sector"])
    st.write("Company Industry Category: ", tick.info["industry"])
    st.write("Company's Shares Regular Opening Price: ", tick.info["regularMarketOpen"])
    st.write("Company's Shares Previous Closing Price: ", tick.info["previousClose"])
    st.write(tick.recommendations)

    # tick_df = tick.history(period="max")
    # fig = tick_df['Close'].plot(title="Company's stock price")
    # st.show(fig)
    # fig, ax = plt.subplots() #solved by add this line
    # ax.lineplot(tick_df, x="Closing", y="Year")
    # st.pyplot(fig)

    start = dt.datetime.today()-dt.timedelta(2 * 365)
    end = dt.datetime.today()
    df = yf.download(company,start,end)
    df = df.reset_index()
    fig = go.Figure(
            data=go.Scatter(x=df['Date'], y=df['Adj Close'])
        )
    fig.update_layout(
        title={
            'text': "Stock Prices Over Past Two Years",
            'y':0.9,
            'x':0.5,
            'xanchor': 'center',
            'yanchor': 'top'})
    st.plotly_chart(fig, use_container_width=True)


def show_history(data):
    """Create and format a chart with operations history"""
    history_dict = data["History"]
    if len(history_dict) == 0:
        st.warning(
            "Error: no stock data found, please purchase a stock or make a deposit first")
        return
    data_list = []
    for date in history_dict.keys():
        for entry in history_dict[date]:
            new_format = [date, entry["Type"], entry["Ticker"], str(
                entry["Shares"]), str(entry["UnitPrice"]), str(entry["OrderTotal"])]
            data_list.append(new_format)

    frame = pd.DataFrame(data_list, columns=[
                         'Date', 'Type', 'Ticker', 'Shares', 'UnitPrice', 'OrderTotal'])
    st.dataframe(frame)


def show_chart(data):
    # ============ Portfolio Values =====================
    col1, col2 = st.beta_columns(2)

    with col1:
        value = '{:.2f}'.format(data["Value"])
        st.markdown(f"### Value: `${value}`")
        cost = '{:.2f}'.format(data["Cost"])
        st.markdown(f"### Cost: `${cost}`")
        total_gain_loss = '{:+.2f}'.format(data["TotalGainLoss"])
        st.markdown(f"### Total Gain/Loss: `${total_gain_loss}`")

    with col2:
        equity = '{:.2f}'.format(data["Equity"])
        st.markdown(f"### Equity: `${equity}`")
        cash = '{:.3f}'.format(data["Cash"])
        st.markdown(f"### Cash: `${cash}`")
        total_gain_loss_prcnt = '{:+.2%}'.format(data["TotalGainLossPrcnt"])
        st.markdown(f"### Gain/Loss Percentage: `{total_gain_loss_prcnt}`")

        # ============ End Portfolio Values =====================

    # ============= Portfolio Table===========================
    stock_data_frame = pd.DataFrame.from_dict(data["Positions"]).T
    if len(stock_data_frame) == 0:
        st.warning("No stocks found to display")
    else:
        st.markdown("### Positions:")
        # sort data frame by equity column
        stock_data_frame = stock_data_frame.sort_values(
            by=["Equity"], ascending=False)
        # rearrange the columns differently
        stock_data_frame = stock_data_frame[['Shares', 'AvgPrice', 'LatestPrice',
                                             'GainLoss', 'GainLossPrcnt', "Equity",
                                             "PrcntOfPort", "TotalCost"]]
        # add dataframe styling to decimal points and conditional formatting
        data_frame_style = (stock_data_frame.style
                            .applymap(color_negative_red, subset=[
                                'GainLoss', 'GainLossPrcnt'])
                            .format(
                                {'AvgPrice': "${:.2f}", 'LatestPrice': '${:.2f}',
                                    'GainLoss': '${:+.2f}', 'GainLossPrcnt': '{:+.2%}',
                                    'Equity': '${:.2f}', 'PrcntOfPort': '{:.2%}',
                                    'TotalCost': '${:.2f}'}))
        st.dataframe(data_frame_style)

    # ============= End Portfolio Table ===========================


def color_negative_red(val):
    """
    Takes a scalar and returns a string with
    the css property `'color: red'` for negative
    strings, black otherwise.
    """
    color = 'red' if val < 0 else 'green'
    return 'color: %s' % color


def show_sidebar(port_path, data):
    selection = st.sidebar.selectbox("Please select an operation", options=[
        "Buy", "Sell", "Deposit", "Withdraw", "Create New Portfolio"])
    # Add a placeholder for textboxes and buttons
    placeholder = st.sidebar.empty()

    if selection == "Buy":
        with placeholder.beta_container():
            buy(port_path, data)

    if selection == "Sell":
        with placeholder.beta_container():
            sell(port_path, data)

    if selection == "Deposit":
        with placeholder.beta_container():
            deposit(port_path)

    if selection == "Withdraw":
        with placeholder.beta_container():
            withdraw(port_path, data)

    if selection == "Create New Portfolio":
        with placeholder.beta_container():
            create_portfolio()


def refresh(port_path):
    # run refresh command and show output
    stdout, stderr = run_command("./stockHandler -f "+port_path + " -r")
    if not(stdout == ""):
        st.warning(stdout)
    if not(stderr == ""):
        st.warning(stderr)


def buy(port_path, data):
    """Execute buy operation"""
    st.header('Buy')
    ticker = st.text_input('Input stock ticker:', key="buy")
    shares = st.number_input('Input number of shares:', key="buy")
    price = st.number_input('Input average price:', key="buy")
    st.write("Order Total: ", round(shares*price, 2))
    if st.button('buy', key="buy"):
        # Check for ticker existence and empty ticker
        if not(valid_ticker(ticker)) or ticker == "":
            st.warning("Invalid ticker")
        # Check positive numbers
        elif (shares <= 0) or (price <= 0):
            st.warning("Invalid share/price")
        # check sufficient cash
        elif data["Cash"] < (shares * price):
            st.warning("Not enough cash")
        else:
            # Run the buy command and show results
            stdout, stderr = run_command("./stockHandler -f "+port_path + " -o buy" + " -t " + str(ticker) +
                                         " -s " + str(shares) + " -p " + str(price) + " - r ")
            if not(stdout == ""):
                st.warning(stdout)
            if not(stderr == ""):
                st.warning(stderr)
            st.success("Shares bought successfully! Please refresh few times.")


def sell(port_path, data):
    """Execute sell operation"""
    st.header('Sell')
    ticker = st.text_input('Input stock ticker:', key="sell")
    shares = st.number_input('Input number of shares:', key="sell")
    price = st.number_input('Input average price:', key="sell")
    st.write("Order Total: ", round(shares*price, 2))
    if st.button("sell", key="sell"):
        # Check positive numbers
        if (shares <= 0) or (price <= 0):
            st.warning("Invalid share/price")
        # check ticker existence in positions
        elif (ticker == "") or not(ticker.upper() in data["Positions"].keys()):
            st.warning("Ticker not found")
        # Check if enough shares are available to sell
        elif data["Positions"][ticker.upper()]["Shares"] < shares:
            st.warning("Sold shares can't exceed owned shares")
        else:
            # Run the sell command and show results
            stdout, stderr = run_command("./stockHandler -f "+port_path + " -o sell" + " -t " + str(ticker) +
                                         " -s " + str(shares) + " -p " + str(price) + " -r ")
            if not(stdout == ""):
                st.warning(stdout)
            if not(stderr == ""):
                st.warning(stderr)
            st.success("Shares sold successfully! Please refresh few times.")


def deposit(port_path):
    """Execute deposit operation"""
    st.header('Deposit')
    cash = st.number_input('Input deposit amount', key="deposit")
    if st.button("deposit", key="deposit"):
        if cash <= 0:
            st.warning("Invalid zero or negative deposit amount")
        else:
            # run deposit command and show results
            command = "./stockHandler -f "+port_path + \
                " -o deposit" + " -s " + str(cash)
            stdout, stderr = run_command(command)
            if not(stdout == ""):
                st.warning(stdout)
            if not(stderr == ""):
                st.warning(stderr)
            st.success("Deposit successfull! Please refresh few times.")


def withdraw(port_path, data):
    """Execute withdraw operation"""
    st.header('Withdraw')
    cash = st.number_input('Input witdraw amount', key="withdraw")
    if st.button("withdraw", key="withdraw"):
        if cash > data["Cash"]:
            st.warning("Insufficient funds to withdraw")
        elif cash <= 0:
            st.warning("Invalid amount")
        else:
            # run withdraw command and show results
            stdout, stderr = run_command(
                "./stockHandler -f "+port_path + " -o withdraw" + " -s " + str(cash))
            if not(stdout == ""):
                st.warning(stdout)
            if not(stderr == ""):
                st.warning(stderr)
            st.success("Withdraw successfull! Please refresh few times.")


def run_command(command_string):
    """Use shlex and subprocess to run a command, return tuple of strings (stdout, stderr)"""
    command = shlex.split(command_string)
    process = subprocess.Popen(
        command, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    stdout, stderr = process.communicate()
    return (stdout.decode('utf-8'), stderr.decode('utf-8'))


def create_portfolio():
    """Execute creation operation"""
    st.header('Create New')
    name = st.text_input('Input new portfolio name')
    cash = st.number_input('Input deposit amount')
    if st.button("create"):
        if cash <= 0:
            st.warning("Invalid zero or negative deposit amount")
        elif " " in name or name == "":
            st.warning("Invalid portfolio name")
        else:
            # run deposit command and show output
            port_path = name + ".json"
            stdout, stderr = run_command(
                "./stockHandler -f "+port_path + " -o deposit" + " -s " + str(cash))
            if not(stdout == ""):
                st.warning(stdout)
            if not(stderr == ""):
                st.warning(stderr)
            st.success("New portfolio create! Please refresh few times.")


def show_graph(data):
    col1, col2 = st.beta_columns(2)
    stock_data = data["Positions"]
    if len(stock_data) == 0:
        st.warning("Error: no stock data found, please purchase a stock first")
        return
    # ========== Pie Chart ======
    with col1:
        st.write("### Postitions Distribution")
        labels_list = list(stock_data.keys())
        labels_list.append("Cash")
        labels_tuple = tuple(labels_list)
        sizes = []
        for stock in stock_data.keys():
            sizes.append(stock_data[stock]["PrcntOfPort"])
        cashPercentage = (data["Cash"] / data["Value"])
        sizes.append(cashPercentage)

        fig1, ax1 = plt.subplots()
        ax1.pie(sizes, labels=labels_tuple, autopct='%1.1f%%',
                shadow=False, startangle=90, normalize=True)
        st.pyplot(fig1)
    # ========== End Pie Chart ======

    # ========== End Stacked Bar Chart ====
    with col2:
        st.write("### Cost-Gain Comparison")
        profit_dict = {}
        for stock in stock_data.keys():
            profit_dict[stock] = {"Profit": stock_data[stock]
                                  ["GainLoss"], "Cost": stock_data[stock]["TotalCost"], }

        profit_dict_frame = pd.DataFrame.from_dict(profit_dict).T
        profit_dict_frame = profit_dict_frame.sort_values(
            by=["Cost"], ascending=True)
        st.bar_chart(profit_dict_frame, height=500)
    # ========== Stacked Bar Chart ====

    # ========== Percentages Chart ====
    with col1:
        st.write("### Stock Profitability")
        stock_list = []
        prcnt_list = []
        for stock in stock_data.keys():
            stock_list.append(stock)
            prcnt_list.append(
                stock_data[stock]["GainLossPrcnt"] * 100)

        percentage_dataframe = pd.DataFrame(
            {"Stocks": stock_list, "Profit Percentage": prcnt_list})
        bars = alt.Chart(percentage_dataframe).mark_bar().encode(
            x="Profit Percentage",
            y="Stocks"
        )
        st.altair_chart(bars, use_container_width=True)

    # ========== End Percentages Chart ====

     # ========== Activity Bar Chart ====
    with col2:
        st.write("### Trading Activity")
        history = data["History"]
        graph_dict = {}
        for date in sorted(history.keys()):
            graph_dict[date] = len(history[date])

        line_frame = pd.DataFrame(graph_dict, index=["# of Operations"]).T
        st.bar_chart(line_frame)

    # ========== End Activity Bar Chart ====


def valid_ticker(name):
    """Checks tickers data returns true if data exists"""
    # create the ticker with uppercase letters
    ticker = yf.Ticker(name.upper())
    # check if there is data on this stock for the past week
    if len(ticker.history(period="1w")) == 0:
        return False
    # Data found return true
    return True


def run_tests():
    """Runs golang test command and shows output"""
    command = "go test -v -cover ./src"
    st.write(
        f"### Test cases can be ran separately using `{command}` or you can click the button below")
    if st.button("Run Tests"):
        stdout, stderr = run_command(command)
        if not(stderr == ""):
            st.warning(stderr)
        if "FAIL" in stdout:
            st.write("### Tests Failed!")
            st.warning(stdout)
        else:
            st.write("### Tests Passed!")
            st.success(stdout)


if __name__ == "__main__":
    web_interface()
