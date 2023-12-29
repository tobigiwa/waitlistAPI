![Logo](/templ/BlockRideLogo.png)

# Waitlist Backend

## Version

#### _version: 0.9.0_ - first release

#### _version: 0.9.1_

> added healthcheck handler

#### _version: 1.0.0_ - official release

> added swagger 2.0 OpenAPI documentation

> added containerisation with Docker

#### _version: 1.0.1_ - minor release

## START SERVER

### DEV MODE:

> $ make dev

Defaults .env variables are enabled

### PROD MODE:

> hosting services matters

## CODE:

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
