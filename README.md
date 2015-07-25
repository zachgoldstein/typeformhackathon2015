# A simple app to create custom user-tagged forms from slack with the typeform I/O API

This app utilises the typeform i/o api, slack, firebase, and reactjs.

## Running it locally:

* run the server to connect to slack's real-time API and receive messages from the webhook `go run main.go`
* run the file server so the images endpoints can get images off the local filesystem `go run ./fileServer.go`
* run webpack to continuously update the page that displays the saved data from firebase `webpack-dev-server --hot --inline --port 3001`

You'll need to expose ports 3001, 3000, and 8080 via `ngrok` for the webhooks, realtime api and file serving to work
