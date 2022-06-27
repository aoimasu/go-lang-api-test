## Setup/installation instructions
- On Mac machine, download the package from https://go.dev/doc/install and install it
- Go to root folder which contain the `main.go` file and run `go build`
- Run the application with `./api` command
- To run test, use `go test -v` command

## Solution
As Go lang is the language required for this position so I tried my best to create the API in go lang. 

The approach I take is the TDD so I developed the test cases to make the API failed at the beginning, then I wrote the API to make the test cases passed. 

Folder structure:
```
.
├── main.go                 # Entry point for this api module
├── app.go                  # Application entry point, used to  
|                           # initialize services and router
├── models                  # Contains type definitions for 
|                           # articles, tags, error...
├── services                # Contains the services, i.e. articles ...
├── utils                   # Utilities and helpers
├── main_test.go            # Unit tests
└── README.md               # this file
```

The error message structure example:
```
{
    "code": 400,
    "message" : "Invalid input"
  }
```
## List of assumptions
- JSON input for POST request to `/articles` endpoint is an array, for example:
```
[
  {
    "title": "title 1",
    "date" : "2016-09-22",
    "body" : "body 1",
    "tags" : ["health", "test", "test2"]
  },
  {
    "id": "2",
    "title": "title 2",
    "date" : "2016-09-22",
    "body" : "body 2",
    "tags" : ["health", "fitness", "science"]
  },
  {
    "title": "title 3",
    "date" : "2016-09-23",
    "body" : "body 3",
    "tags" : ["health", "test3", "test4"]
  }
]
```
- Data submitted to `/articles` endpoint will be upserted into local storage. Entry without id will be insert and entry with id will be updated if it exists in local storage.
- GET request for single article will return `{}` object if there is no available article
- I assumed this is an open api so there is no rate limit and no authentication
- I assumed the service store data in memory or local storage so data will not be persistent and will be wiped out upon restarting the service.


## [Optional] Tell us what you thought of the test and how long it took you to complete
The test is a basic API implementation and I spent 8 hours to work on this test, most of the time I spent on reading Go documentation
## [Optional] Tell us what else you would have added to the code if you had more time
- I would move the local storage openration into a separated lib so then if I want to swap in a DB repo, it would be easier. 
- I would add request data validator, I'm not familiar with Go lang so I decided to leave it out. I only did some basic validation at the handlers.