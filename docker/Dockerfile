FROM golang:1.17

WORKDIR /app
COPY . .

# Create a directory for certificates and set it as the volume
RUN mkdir /certificates
VOLUME ["/certificates"]

CMD ["./gcipher"]
