# Define the base image with desired Go version
FROM golang:latest

# Copy your Go source code
COPY . .

# Install dependencies (if using Go modules)
RUN go mod download

# Build your Go application
RUN go build -o main ./main

# Expose port (if applicable)
EXPOSE 7777

# Command to run the application
CMD ["./main/main"]