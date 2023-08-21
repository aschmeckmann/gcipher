# gcipher: Go Certification & Integrity Platform for Hosted Encryption Resources
## Unfinished Project

gcipher is an open-source Go-based public key infrastructure (PKI) that empowers you to easily manage and secure public-key certificates while maintaining scalability. It provides a comprehensive solution for handling certificate signing requests (CSRs), certificate signing, storage, retrieval, and more.

## Features

- **Certificate Request Handling:** Process incoming certificate signing requests (CSRs) from clients and generate signed certificates using the [x509 package](https://pkg.go.dev/crypto/x509) in Go.
- **Certificate Storage:** Save signed certificates securely to a MongoDB for persistent storage, ensuring certificates remain available even across container restarts.
- **Certificate Revocation List (CRL):** Manage and maintain a CRL to keep track of revoked certificates, enhancing security by preventing the use of compromised certificates.
- **User Authentication:** Authenticate users' requests to ensure secure and authorized access to certificate-related operations.
- **Flexibility:** Modify, extend, and tailor the infrastructure to your organization's unique security requirements.

## Roadmap and Future Enhancements

As we continue to enhance and develop the **gcipher** project, here are some of the key improvements and features we're considering for the future:

- **Optimized Serial Number Storage:** Explore optimizing the storage of serial numbers by directly using big integers, reducing unnecessary type conversions.

- **CA Endpoint:** Implement an endpoint to serve the CA certificate, making it easily accessible for verification purposes.

- **Logging Infrastructure:** Develop a robust and configurable logging system that allows users to choose between logging to standard output, log files, or even MongoDB.

- **Access Control List (ACL):** Consider implementing an Access Control List mechanism to further enhance security and manage user access to various resources and operations.

- **HTTP Method Checking:** Evaluate the implementation of HTTP method checking to ensure that API endpoints are accessed using the appropriate HTTP methods, enhancing endpoint security.

- **Advanced Routing:** Explore utilizing more advanced routing mechanisms to enhance the efficiency and organization of the routing system, potentially replacing the current http/mux with a more robust routing solution.

- **Continuous Improvement:** Continuously review and improve the project, addressing bugs, enhancing documentation, and incorporating user feedback to ensure a seamless and secure experience.

These are just a few of the exciting enhancements we have on our roadmap. We're committed to making **gcipher** a powerful and flexible solution for managing certificates and encryption resources. Your feedback and contributions are invaluable as we work towards these goals.

## Getting Started

1. Clone the repository:

    ```bash
    git clone https://github.com/aschmeckmann/gcipher.git
    cd gcipher
    ```

2. Build the Docker container:

    ```bash
    docker build -t gcipher .
    ```

3. Run the server command:

    ```bash
    docker run -p 8080:8080 --volume certificates:/certificates gcipher server
    ```

4. Access the API at `http://localhost:8080`.

With the updated command system, you can now use different commands to manage the application. The "server" command starts the server, allowing you to access the API. You can explore additional commands as they are implemented in your application.

## Usage

- **Certificate Request:** POST a CSR to `/api/v1/certificate/request` to generate signed certificates.
- **Certificate Retrieval:** POST a serial number to get a certificate using `/api/v1/certificate/retrieve`.
- **Certificate Revocation:** POST a serial number to revoke a certificate using `/api/v1/certificate/revoke`.
- **CRL Retrieval:** GET the latest CRL using `/public/ca/intermediate/crl`

### API Request Structure

```json
{
  "data": {
    "applicant": "John Doe",         // Optional: Name of the certificate applicant
    "csr": "BASE64_CSR_DATA",       // Optional: Certificate signing request in BASE64 format
    "lifetime": 365,                // Optional: Lifetime of the certificate in days
    "type": "client",               // Optional: Type of certificate (client or server)
    "state": "active",              // Optional: State of the certificate (active, revoked, etc.)
    "serialnumber": "1234567890"    // Optional: Serial number of the certificate
  },
  "auth": {
    "username": "your_username",    // Username for authentication
    "password": "your_password"     // Password for authentication
  }
}
```

### API Response Structure

```json
{
  "success": true,                 // Indicates the success of the request
  "errors": [
    {
      "code": 400,                 // HTTP status code
      "message": "Error message"   // Error description
    }
  ],
  "data": {
    "certificate_pem": "PEM_DATA"  // Certificate data in PEM format (for CertificateResponseData)
  }
}
```

### Configuration

The `gcipher` application can be configured using a YAML configuration file. By default, the application looks for `config.yml` or `config.yaml` in the same directory. You can also provide configuration options through environment variables.

#### Configuration File

An example of a `config.yml` configuration file:

```yaml
port: 8080
database_url: "mongodb://localhost:27017"
certificate_lifetime_default: 365
ca_cert_path: "/path/to/ca_cert.pem"
ca_key_path: "/path/to/ca_key.pem"
```

#### Environment Variables

You can override configuration options by setting environment variables. The following environment variables are available:

- `GCIPHER_PORT`: Port on which the application should run.
- `GCIPHER_DATABASE_URL`: Database URL for connecting to MongoDB.
- `GCIPHER_CERTIFICATE_LIFETIME_DEFAULT`: Default lifetime of certificates in days.
- `GCIPHER_CA_CERT_PATH`: Path to the CA certificate file.
- `GCIPHER_CA_KEY_PATH`: Path to the CA private key file.

Please note that environment variables take precedence over configuration file options.

To ensure proper configuration, it's recommended to provide the necessary certificate and key paths, especially when dealing with certificate management operations. Make sure to adjust these values according to your environment.

For more details on the configuration options, please refer to the source code in the `config` package.

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
