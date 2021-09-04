# go-cloudrun-boilerplate
This repository is a boilerplate for go project in STUDIO.

# Pre Condition
- Go 1.17.0 >=
- Docker 3.6.0 >=

# Quickstart
1. Initialize project 
   ```
   go mod init go-cloudrun-boilerplate
   go mod tidy
   ```
   
# Tips
## How to run tests
```
APP_ENV=test PROJECT_UUID=<Project UUID> PROJECT_ID=<Project ID here> GOOGLE_APPLICATION_CREDENTIALS=<Service Account file path here> go test -v -race -run=. -bench=. ./...
```
## How to format all go files
```
go fmt ./...
```
## Clean up go.mod
```
go mod tidy
```