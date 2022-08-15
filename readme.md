# Generic Authentication Microservice 
Clone this and adapt it to your needs.

## Description
This is a generic microservice that handles generating tokens for accessing your resource services. See comments above endpoints for more details for each endpoint.

## Instructions

### Removing "generic" prefix from all the file names and adding your own prefix
Search for "generic" and replace it with your project name

### Environment variables

To get started quickly, add a file called .env.local and add the following environment variables. In production, the program will read environment variables from the host environment

```
JWT_SECRET=Your JWT Secret
MONGODB_URI=Your database uri. Should start with mongodb.
GOOGLE_CLIENT_ID=Your Google Client ID
GOOGLE_CLIENT_SECRET=Your Google Client Secret
DATABASE_NAME=The main database where your users should be saved
```

### Dependencies

```
go mod download
```

### Running the program

```bash
go run main.go
```

## Notes
Currently only Google sign in is supported. I'm planning to add more providers soon.