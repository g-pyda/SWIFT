# SWIFT Go Application

## Overview

This project is a Go-based application providing an REST API for relational database (MySQL) containing the SWIFT codes data. It is built using Docker and can be deployed in a lightweight Alpine Linux container.

## Prerequisites

Make sure you have the following installed before proceeding:

- [Go](https://go.dev/dl/) (version 1.24.1 or later)
- [Docker](https://www.docker.com/get-started)
- [Git](https://git-scm.com/downloads) (optional, for cloning the repository)

## Setup Instructions

### 1. Clone the Repository

```sh
git clone https://github.com/g-pyda/SWIFT.git
cd SWIFT
```

### 2. Build and Run the Application

#### Running Locally (without Docker)

1. **Go to the *src* directory**

```sh
cd src
```

2. **Install swagger to get the swagger UI**

```sh
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g ./main.go
```

3. **Run the project**

```sh
go run .
```

#### Running with Docker

1. **Build the Docker Image**

```sh
docker-compose build
```

2. **Run the Container**

```sh
docker-compose up
```

The application will now be accessible at `http://localhost:8080`.

## Dockerfile Overview

This project uses a multi-stage Docker build to optimize the container size:

1. **Build Stage**: Uses a Golang container to build the application.
2. **Run Stage**: Uses a lightweight Alpine container to run the compiled binary.

## API Endpoints

- `GET /v1/swift-codes` - Return all SWIFT codes 
- `GET /v1/swift-codes/{swift-code}` - Retrieve details of a single SWIFT code whether for a headquarters or branches
- `GET /v1/swift-codes/country/{countryISO2code}` - Return all SWIFT codes with details for a specific country (both headquarters and branches)
- `POST /v1/swift-codes` - Add new SWIFT code entries to the database for a specific country, given the request stucture is the following:

```json
{
    "address": string,
    "bankName": string,
    "countryISO2": string,
    "countryName": string,
    "isHeadquarter": bool,
    "swiftCode": string,
}
```

- `DELETE /v1/swift-codes/{swift-code}` - Delete swift-code data if swiftCode matches the one in the database

**Endpoint can be further investigated and tested within the *swagger UI* under the following endpoint**

- `GET /swagger/index.html`

## Testing

Testing can be easily performed through *testing* go module in this project as follows:

1. Ensure you're placed in the *src* directory 
2. Run the first test for *xlsxParser* package 

```sh
go test ./xlsxParser
```

3. Run the first test for *databaseControl* package 

```sh
go test ./databaseControl
```

4. Run the first test for *REST* package 

```sh
go test ./REST
```

**REMARK: REST test won't be executed if the databaseControl test has failed**

