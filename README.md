# Project AlfredEventsOranizer

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

All the backend code lives in the `backend/` folder.

### Setup

```bash
cd backend
cp .env.example .env
go mod tidy
```

Place your Firebase Admin service account file in `backend/FirebaseAdminCredential.json` and fill in the `.env` values (PostgreSQL `DSN`, Firebase UIDs).

Run migrations from `backend/` with [Goose](https://github.com/pressly/goose):

```bash
goose up
```

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
