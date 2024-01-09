![Logo](/templ/Screenshot%20from%202024-01-09%2005-33-38.png)

# Waitlist Backend with Email Confirmation (using Go)

Welcome to the "Waitlist Backend with Email Confirmation" project! This repository serves as a comprehensive tutorial and implementation guide for building a robust waitlist system with email confirmation using the Go programming language.

## Project Overview

The goal of this project is to provide developers with a step-by-step guide on creating a backend system that manages a waitlist and ensures user sign-ups are confirmed through email verification. Leveraging the power and simplicity of Go, we'll explore the essential components of building a scalable and secure waitlist application.
As you would suspect, company "companyXYZ" was made up, and the Ethereum Foundation logo was just a random pick. Thank you.

## Key Features

- **Waitlist Management**: Implement a reliable system to handle user registrations saved to a MongoDB database.

- **Email Confirmation**: Securely confirm user sign-ups through a email verification process with expiry link.

- **Swagger Documentation of Endpoints**: Document the API endpoints using Swagger to enhance the developer experience and facilitate integration with the backend.

- **Scalable Design**: Utilize best practices in Go to ensure scalability wit no need for extra services.

## Prerequisites

Before diving into the tutorial, make sure you have the following prerequisites installed:

- **Go** (any version should do, but I developed with 1.21).

- Database system (MongoDB for this tutorial).

- SMTP server for sending confirmation emails (I used [zohomail](https://www.zoho.com/mail/), it allows quick setup of app passwords).

## Getting Started

### DEV MODE:

> $ make setup

> add your .env to the bin/ created folder

> $ make dev

Defaults .env variables are enabled

### PROD MODE:

> hosting services matters

## Endpoints:

### Routes:

**POST:** _"/joinwaitlist"_ **sends user registration email**

**GET:** _"/confirmuser"_ **confirms user registration from email link**

**GET:** _"/healthcheck"_ **Report application status**

**GET:** _"/swagger/\*"_ **OpenAPI Schema documentation**

### Frontend endpoint integration:

#### **POST:** _"/joinwaitlist"_

> This endpoint takes a url-encoded form data, see **GET:** _"/swagger/\*"_.

> No response body is sent, headers are set instead.

- 200 StatusOK - redirect to a page that tells user to check their mail.
- 400 StatusBadRequest - formData/validation error.
- 409 StatusConflict - user with email already exist, redirects to homepage/login.

#### **GET:** _"/confirmuser"_

> This endpoint is called from the sent mail.

> No response body is sent, headers are set instead.

- 200 StatusOK - redirects to a page that says "user added to waitlist".
- 400 StatusBadRequest - Url is not authentic.
- 404 StatusNotFound - The url has expired.
- 409 StatusConflict - user with email already exist, redirects to homepage/login.

## Contributing

Contributions are welcome! If you find issues or have ideas for improvements, please feel free to open an issue or submit a pull request. Let's build this waitlist backend together!

Happy coding!
