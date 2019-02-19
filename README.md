
## MatchUps
The purpose of the application is to make arranging casual sports games easier by making 
organisation within your own group of friends easier and also allow people to find other teams
to play against.

Application originally designed and developed for a 2nd year group project at Imperial alongside 
Thomas Yung, Shashwat Dalal and Andy Li.

I have since continued working on it in order to produce a fully working application and also a 
cleaner codebase.

## Implementation
This project was bootstrapped with [Create React App](https://github.com/facebookincubator/create-react-app). 
The front end react implementation can be found in `src/`.

Golang is then used to serve static files generated from the React code by `yarn build`. GoLang serves the files using 
the `gorilla/mux` package. For this see `./main.go`.

`Golang` is also the server side language which handles requests from the client. To do this it communicates with
either our `PostgreSQL` database (which is hosted on AWS) or Google Maps APIs. For this see `./handlers/`.

`Pusher` is also used within the React components to enable real time interaction between several users.