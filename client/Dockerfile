FROM alpine:latest

# Copy the pre-built server binary
COPY main . 

# Copy the certificate files
COPY certs/. . 

# Command to run the server
CMD ["/main", "8081"]