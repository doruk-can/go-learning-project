# Foover

**Foover** is a backend service designed to handle product votes. It allows users to generate sessions, submit votes, and retrieve aggregated product scores for products in a machine.

## Features

- Generate user sessions
- Submit votes for products
- Retrieve aggregated product scores

## Prerequisites

To run this service locally, you need the following:

- **Go**: Version 1.16 or higher
- **Docker**: For MongoDB setup

## Installation

Follow these steps to install and run the service:

1. **Run Docker Compose**:
   The service uses MongoDB, which runs at port `27017`. You can start MongoDB by running the following command in the deployments directory:

      ```bash
      docker-compose up
      ```

2. **Run the Application**:
   After ensuring MongoDB is running, you can start the Go application:

   ```bash
   go run cmd/main.go
    ```

## API Documentation

You can access the Swagger UI by navigating to:

- **Swagger UI URL**: http://localhost:8080/docs

## Usage

When the application starts, it retrieves products from an external API during initialization. You can use the product IDs from this initial product fetch to send votes.

The following are product IDs you can use for testing, assuming the response from the external API hasn't changed:

- "3aba3a59-fd44-45e8-80db-7d4771b8f822"
- "4d9c21d9-26e2-4191-b971-0fb6587d3cb9"
- "7aeacdb5-9a99-4210-a6ec-adce3d0f4bfc"
- "83c3e42e-1147-481c-bdd8-3ccfa1e54c96"
- "9c13c0b2-a4c3-4a62-aaf9-c5b40dd4394d"
- "bfe734c2-e51a-4524-82b5-e5a20c36a94b"
- "ddb7dcef-cedb-4c13-b218-4cf0a7f64e11"
- "e2ba10ec-3c82-4dd0-bdb1-86a418d54a87"