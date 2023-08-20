# gcipher: Go Certification & Integrity Platform for Hosted Encryption Resources

# Unfinished

gcipher is an open-source Go-based public key infrastructure (PKI) that empowers you to easily manage and secure public-key certificates. It provides a comprehensive solution for handling certificate signing requests (CSRs), certificate signing, storage, retrieval, and more.

## Features

- **Certificate Request Handling:** Process incoming certificate signing requests (CSRs) from clients and generate signed certificates using the x509 package in Go.
- **Certificate Storage:** Save signed certificates securely to a Mongoose DB for persistent storage, ensuring certificates remain available even across container restarts.
- **User Authentication:** Authenticate users' requests to ensure secure and authorized access to certificate-related operations.
- **Flexibility:** Modify, extend, and tailor the infrastructure to your organization's unique security requirements.

## Getting Started

1. Clone the repository:

    ```bash
    git clone https://github.com/aschmeckmann/gcipher.git
    cd gcipher
    ```

2. Build and run the Docker container:

    ```bash
    docker build -t gcipher .
    docker run -p 8080:8080 --volume certificates:/certificates gcipher
    ```

3. Access the API at `http://localhost:8080`.

## Usage

- **Certificate Request:** POST a CSR to `/api/v1/certificate/request` to generate signed certificates.
- **Certificate Retrieval:** GET a certificate by serial number using `/api/v1/certificate/retrieve?serialNumber=SERIAL_NUMBER`.

## Dependencies

- Go (1.17)
- Docker

## Contributing

Contributions to gcipher are welcome! Feel free to open issues for bug reports, feature requests, or questions. If you'd like to contribute code, please fork the repository and create a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

---

For more information, documentation, and updates, visit the [gcipher GitHub repository](https://github.com/aschmeckmann/gcipher).

**Disclaimer:** gcipher is provided as-is and without warranty. Always follow best practices for security and encryption when using public-key infrastructure.
