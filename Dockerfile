# Gets golang version 1.17 to use in the container
FROM golang:1.17

# Sets a maintainer
LABEL maintainer="mmstene@stud.ntnu.no"

# Setting the work directory
WORKDIR /go/assignment-2/main

# Copies the directories needed.
COPY ./main /go/assignment-2/main
COPY ./handler /go/assignment-2/handler
COPY ./client /go/assignment-2/client
COPY ./constants /go/assignment-2/constants
COPY ./json_coder /go/assignment-2/json_coder
COPY ./structs /go/assignment-2/structs
COPY ./util /go/assignment-2/util
COPY ./CredentialsFirestore.json /go/assignment-2/CredentialsFirestore.json
COPY ./go.mod /go/assignment-2/go.mod
COPY ./go.sum /go/assignment-2/go.sum

# Sets up the folder in which the container will run from
RUN CGO_ENABLE=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main

# Gives the command to execute the main folder
CMD ["./main"]