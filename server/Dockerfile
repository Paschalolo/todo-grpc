FROM alpine:latest

# Copy the pre-built server binary
COPY main . 

# Copy the certificate files
COPY certs/. . 

# Expose the port your server listens on (e.g., 8081)
EXPOSE 8081

# Command to run the server
CMD ["/main", "8081"]