FROM golang:1.21-alpine3.18

# Set the woeking directory to /app
WORKDIR /app

# Copy the go.mod and go.sum fules to the working directory
COPY go.mod go.sub ./

# Download and indtall requested dependencies
RUN go mod Download

# Copy the entire source code to hte working dir
COPY . .

# Build the Go application
RUN go build -o main

# Expose the port, specified by PORT env variable
EXPOSE 3000

#Set the entry point
CMD [ "./main" ]