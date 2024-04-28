# K-Tax Application

K-Tax is a backend API service designed to calculate personal income tax based on the provided income, deductions, and allowances. It simplifies the process of determining the tax payable or refundable for a given year, taking into account various factors such as withholding tax (WHT) and allowances consisting of personal deduction, donation deduction, and k-receipt (Shop to Reduce Tax) deduction. The service is developed using the Go programming language, providing a fast and efficient solution for tax calculations.

This project is developed as part of the Go KBTG Bootcamp, which aims to enhance skills in Go programming and build practical applications.

## Progressive Tax Bracket Calculations

Tax is calculated based on the following progressive income brackets:

- Income from 0 - 150,000: Exempt from tax
- Income from 150,001 - 500,000: Tax rate of 10%
- Income from 500,001 - 1,000,000: Tax rate of 15%
- Income from 1,000,001 - 2,000,000: Tax rate of 20%
- Income over 2,000,000: Tax rate of 35%

## Features

- Calculate personal income tax based on the provided `total income` and deductions
- Support for different types of allowances, including `personal allowance`, `donation`, and `k-receipt`
- Handle withholding tax (WHT) and calculate the tax `refund when applicable`
- Provide a detailed breakdown of tax calculations for each progressive tax brackets
- Allow `admin users` to configure personal allowance and k-receipt deduction limits
- Swagger documentation for API exploration and testing
- Containerization using Docker for easy deployment and scalability

## Getting Started

This section provides instructions on how to set up and run the application.

### Prerequisites

- Go programming language (version 1.21 or higher)
- Docker (for containerization)
- Docker Compose (for running the PostgreSQL database)

### Installation

Clone the project using:

```
git clone https://github.com/thitiphum-bluesage/assessment-tax.git .
```

Start the PostgreSQL database environment using Docker:

```
docker-compose up -d --build
```

### Running the Application

#### Method 1: Using Go Directly

Set up all environment variables:

For Windows:

```
$env:DATABASE_URL = "postgres://godou:1111@localhost:2000/ktaxes?sslmode=disable"
$env:PORT = "8080"
$env:ADMIN_USERNAME = "adminTax"
$env:ADMIN_PASSWORD = "admin!"
```

For macOS/Linux:

```
export DATABASE_URL=postgres://godou:1111@localhost:2000/ktaxes?sslmode=disable
export PORT=8080
export ADMIN_USERNAME=adminTax
export ADMIN_PASSWORD=admin!
```

Start the application:

```
go run main.go
```

#### Method 2: Using Docker

Build the Go project image:

```
docker build . -t ktax_ou
```

Since the environment variables are set within the Dockerfile, you can start the container directly from the built image with network host settings:

```
docker run -p 8080:8080 --network host ktax_ou
```

## Default Configuration

Upon initial setup, the K-Tax Application is configured with default deduction limits to help you get started quickly. These settings are intended to provide a baseline for tax calculations.

### Initial Deduction Limits

- **Personal Deduction**: 60,000
- **K-Receipt Deduction Maximum**: 50,000
- **Donation Deduction Maximum**: 100,000 (fixed and cannot be adjusted)

These values are the starting points for tax calculations. The personal and k-receipt deduction limits can be adjusted by authorized admin users as needed. The donation deduction limit is fixed and cannot be changed.

### How to Adjust Deduction Configuration

Admin users can update the settings for personal and k-receipt deductions by authenticating and sending requests to the respective admin endpoints. Here are the endpoints available for configuration adjustments:

- **POST /admin/deductions/personal**: To update the personal deduction.
- **POST /admin/deductions/k-receipt**: To update the k-receipt deduction limit.

For more details on how to authenticate and modify these settings, refer to the descriptions provided under each relevant API endpoint.

## API Endpoints

### POST /tax/calculations

Calculates the total tax based on total income, withholding tax (WHT), and specified allowances. Returns the total tax and a breakdown by tax brackets, along with any applicable tax refund.

#### Request Example

Calculate tax with no refund:

```json
{
  "totalIncome": 750000.0,
  "wht": 20000.0,
  "allowances": [
    {
      "allowanceType": "k-receipt",
      "amount": 70000.0
    },
    {
      "allowanceType": "donation",
      "amount": 50000.0
    }
  ]
}
```

#### Response Example

```json
{
  "tax": 28500,
  "taxLevel": [
    {
      "level": "0-150,000",
      "tax": 0
    },
    {
      "level": "150,001-500,000",
      "tax": 35000
    },
    {
      "level": "500,001-1,000,000",
      "tax": 13500
    },
    {
      "level": "1,000,001-2,000,000",
      "tax": 0
    },
    {
      "level": "2,000,001 ขึ้นไป",
      "tax": 0
    }
  ]
}
```

<details>
<summary>Calculation Logic</summary>

**Initial Total Income:** 750,000

**Allowances:**

- **K-receipt:** 70,000 (max allowed 50,000)
- **Donation:** 50,000 (max allowed 10,0000)
- **Personal Deduction:** 60,000 (standard for all)

**Taxable Income Calculation:**

- 750,000 (Total Income) - 50,000 (K-receipt) - 50,000 (Donation) - 60,000 (Personal Deduction) = 590,000

**Tax Calculation:**

- **First 150,000:** Tax-Free
- **Next 350,000:** 10% = 35,000
- **Remaining 90,000:** 15% = 13,500
- **Total Tax Due:** 35,000 + 13,500 = 48,500

| Tax Level           | Tax    |
| ------------------- | ------ |
| 0-150,000           | 0      |
| 150,001-500,000     | 35,000 |
| 500,001-1,000,000   | 13,500 |
| 1,000,001-2,000,000 | 0      |
| 2,000,001 ขึ้นไป    | 0      |

**Withholding Tax and Final Tax:**

- **Tax Due:** 48,500
- **Less:** Withholding Tax (WHT) of 20,000
- **Net Tax Payable:** 48,500 - 20,000 = 28,500
</details>

#### Request Example for refund response

```json
{
  "totalIncome": 750000.0,
  "wht": 50000.0,
  "allowances": [
    {
      "allowanceType": "k-receipt",
      "amount": 70000.0
    },
    {
      "allowanceType": "donation",
      "amount": 50000.0
    }
  ]
}
```

#### Refund Response Example

```json
{
  "taxRefund": 1500,
  "taxLevel": [
    {
      "level": "0-150,000",
      "tax": 0
    },
    {
      "level": "150,001-500,000",
      "tax": 35000
    },
    {
      "level": "500,001-1,000,000",
      "tax": 13500
    },
    {
      "level": "1,000,001-2,000,000",
      "tax": 0
    },
    {
      "level": "2,000,001 ขึ้นไป",
      "tax": 0
    }
  ]
}
```

<details>
<summary>Calculation Logic</summary>

**Initial Total Income:** 750,000

**Allowances:**

- **K-receipt:** 70,000 (max allowed 50,000)
- **Donation:** 50,000 (max allowed 10,0000)
- **Personal Deduction:** 60,000 (standard for all)

**Taxable Income Calculation:**

- 750,000 (Total Income) - 50,000 (K-receipt) - 50,000 (Donation) - 60,000 (Personal Deduction) = 590,000

**Tax Calculation:**

- **First 150,000:** Tax-Free
- **Next 350,000:** 10% = 35,000
- **Remaining 90,000:** 15% = 13,500
- **Total Tax Due:** 35,000 + 13,500 = 48,500

| Tax Level           | Tax    |
| ------------------- | ------ |
| 0-150,000           | 0      |
| 150,001-500,000     | 35,000 |
| 500,001-1,000,000   | 13,500 |
| 1,000,001-2,000,000 | 0      |
| 2,000,001 ขึ้นไป    | 0      |

**Withholding Tax and Final Tax:**

- **Tax Due:** 48,500
- **Less:** Withholding Tax (WHT) of 50,000
- **Net Tax Payable:** 48,500 - 50,000 = -1,500
</details>

### POST /tax/calculations/upload-csv

Allows batch processing of tax calculations by uploading a CSV file containing columns for `totalIncome`, `wht`, and `donation`. This endpoint is useful for calculating taxes for multiple entries at once and returns the tax calculated or tax refund for each row in the CSV.

#### CSV Format

The CSV file should contain the following columns, with only `totalIncome` being mandatory:

- `totalIncome`: The total income of the individual.
- `wht` (optional): Withholding tax already paid.
- `donation` (optional): Donation deductions.
- `k-receipt` (optional): k-receipt deductions.

#### Example CSV Content

```
totalIncome,wht,donation,k-receipt
500000,0,0,0
600000,40000,20000,0
750000,50000,15000,30000
```

#### Request

The request involves uploading the CSV file through a form or API client that supports file uploads. Ensure that the file is attached with the key `taxFile` for the request to be processed correctly.

#### Response Example

The response returns an array of tax calculations for each row provided in the CSV file:

```json
{
  "taxes": [
    {
      "totalIncome": 500000,
      "tax": 29000
    },
    {
      "totalIncome": 600000,
      "taxRefund": 2000
    },
    {
      "totalIncome": 750000,
      "tax": 6750
    }
  ]
}
```

### POST /admin/deductions/personal

Allows admin users to configure the personal allowance deduction limits. This endpoint requires basic authentication with admin credentials to ensure only authorized users can make changes.

#### Authentication

To access this endpoint, you must provide basic authentication credentials:

- **Username**: `adminTax`
- **Password**: `admin!`

#### Request Example

To update the personal deduction limit:

```json
{
  "amount": 70000
}
```

#### Response Example

```json
{
  "personalDeduction": 70000
}
```

#### Unauthorized Access

If the authentication credentials are incorrect or not provided, the response will be:

```json
{
  "message": "Unauthorized: Incorrect credentials"
}
```

### POST /admin/deductions/k-receipt

Allows admin users to configure the k-receipt (Shop to Reduce Tax) deduction limits. This endpoint requires basic authentication with admin credentials to ensure only authorized users can make adjustments.

#### Authentication

To access this endpoint, you must provide basic authentication credentials:

- **Username**: `adminTax`
- **Password**: `admin!`

#### Request Example

To update the k-receipt deduction limit:

```json
{
  "amount": 40000
}
```

#### Response Example

If the credentials are correct and the request is successful, the response will confirm the updated k-receipt deduction limit:

```json
{
  "kReceipt": 40000
}
```

#### Unauthorized Access

If the authentication credentials are incorrect or not provided, the response will be:

```json
{
  "message": "Unauthorized: Incorrect credentials"
}
```

### Documentation and API Exploration

You can explore the API documentation and interact with the endpoints using Swagger UI. This provides a user-friendly web interface where you can see all available endpoints, their expected parameters, and even test them in real-time.

#### Accessing the Documentation

The Swagger UI can be accessed via the following paths:

- `/docs`
- `/swagger/index.html`

Both URLs map to the same Swagger interface, so you can use either URL to access the documentation.

## Contact

- Email: cthitipum@gmail.com
