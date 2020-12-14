FROM python:3.8

WORKDIR /usr/src/app

# copy all the files to the container
COPY . .

# Install go
RUN wget -c https://golang.org/dl/go1.15.6.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

RUN go version

# Install go dependency
RUN go get github.com/piquette/finance-go

# Intsall pipenv
RUN pip install pipenv

# # Install dependencies using pipenv
# RUN pipenv install

# Expose the streamlit port
EXPOSE 8501

# Run the streamlit command
CMD pipenv run streamlit run src/web_interface.py