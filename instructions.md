# PyGo Stock Portfolio Builder

## Prerequisites

### Python Environment

Python makes up all front-end aspects of this project. Specifically, web application and data retrieval are done using Python. To ensure that all Python components run smoothly, you will need to install:

- A version of Python through Pyenv
- Pip, a Python package manager
- Pipenv, to create a Python virtualenv and retrieve needed packages. This package can be installed using:

```bash
pip install pipenv
```

### Golang Environment

Back-end operations and calculations, as well as data storage are made possible through Go. The following prerequisites are needed:

- A new version of Go (1.14 or 1.15)
- `finance-go` package, can be installed using:

```bash
go get github.com/piquette/finance-go
```

## Installation

After all needed packages and tools are installed, the Python virtual environment needs to be created in the repository's root using:

```bash
pipenv install
```

This might take some time to locate all dependencies and create `Pipfile.lock`. However, once this is completed, you can run the tool using:

```bash
pipenv run streamlit run src/web_interface.py
```

The test suite for the tool can be run using the following command:

```bash
go test -v -cover ./src
```

## Using Back-end Independently

- build all modules: `go build portfolio.go stockHandler.go`
- run stockHandler using client line argument flags, example: `go run stockHandler -o sell -t TSLA -s 54.345 -p 43`
