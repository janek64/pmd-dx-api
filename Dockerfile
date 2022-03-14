FROM golang:1.17.8

# add a user to run the pmd-dx-api
RUN addgroup api && adduser --ingroup api --gecos "" --disabled-password api

# set the working directory to the home of the api user
WORKDIR /home/api

# switch to the api user
USER api

# copy go.mod, then download and verify the dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# copy all source files
COPY main.go main.go
COPY api api

# build the pmd-dx-api
RUN go build -v

# expose port 3000 since it is the default port
EXPOSE 3000

# start the API
CMD ["./pmd-dx-api"]