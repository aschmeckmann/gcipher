FROM golang:1.17

WORKDIR /app
COPY . .

# Create a directory for certificates and set it as the volume
RUN mkdir /certificates
VOLUME ["/certificates"]

# Set environment variables for configuration
ENV GCIPHER_PORT=8080
ENV GCIPHER_DATABASE_URL="mongodb://localhost:27017"
ENV GCIPHER_CERTIFICATE_LIFETIME_DEFAULT=365
ENV GCIPHER_CA_CERT_PATH="/app/ca.crt"
ENV GCIPHER_CA_KEY_PATH="/app/ca.key"

CMD ["./gcipher"]
