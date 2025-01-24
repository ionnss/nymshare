# Use official Golang image
FROM golang:1.23.4

# Define work directory
WORKDIR /nymshare

# Copy files and dependencies
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Download the wait-for-it script
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh

# Make script executable
RUN chmod +x /wait-for-it.sh

# Copy code and static files
COPY . .

# Ensure static directory exists and has proper permissions
RUN mkdir -p /nymshare/static && chmod -R 755 /nymshare/static

# Compile application with optimizations for production
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/app_nymshare -ldflags="-s -w" && \
    chmod +x /usr/local/bin/app_nymshare

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["/wait-for-it.sh", "db:5432", "--", "/usr/local/bin/app_nymshare"]