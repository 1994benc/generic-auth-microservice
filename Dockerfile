FROM golang:1.18

# Set the Current Working Directory inside the container
WORKDIR /app
# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . /app/

# Download all the dependencies
RUN go mod download

# This container exposes port 8080 to the outside world
EXPOSE 8080

RUN go build -o /generic-auth ./

# Run the executable
CMD ["/generic-auth"]