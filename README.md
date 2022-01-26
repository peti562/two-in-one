# Golang Coding and Logical challenge

Thank you for taking the time and checking my solution

## Notes
- I do not currently have the private infrastructure to run the api on the cloud as I normally would so I just made a locally runnable application
- I received two tasks from you, a fibonacci logical task for big numbers and a CRUD api. 
- I've done both inside this one to prevent unnecessary overcomplicating and to save your time and effort checking them
- I created a basic free MySQL instance for the comments, in real life I would not leave the env file here with credentials in it, I just wanted to keep it as easy and simple as possible

## Setup

1. After you installed the dependencies using `go get .`,
2. You can run linter normally with `golangci-lint run`
3. You can run the unit tests as usual with `go test ./...`
4. You can run the application itself using `go run .`

## Requests
### *** Please note that I only used GET everywhere (and dummy data on update and create) to save your time and effort for testing. In real life I would obviously use POST/PUT/PATCH/DELETE etc. ***
- GET `localhost:3000/fibonacci/<n>` replacing `n` here with a number you wish to use and you will receive the fibonacci value back
- GET `localhost:3000/comments/<userId>` gets all comments related to the user
- GET `localhost:3000/comment/<commentId>` gets a single comment by Id
- GET `localhost:3000/comment/<commentId>` (soft) deletes a comment
- GET `localhost:3000/comment/<commentId>` updates a comment body
- GET `localhost:3000/comment/create` creates a comment - i used dummy data to prevent unnecessary overcomplicating 
