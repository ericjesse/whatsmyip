
# ![logo](.README_images/youarehere.png)  What's my IP


Simple Go service to detect the remote IP address of the caller.

The application is originally designed to run on [Heroku](https://www.heroku.com) and in a Docker container.  
However, you can run it locally by simply add the port to listen as first parameter or as the environment variable PORT.

## Getting the code

Make sure you have [Go](http://golang.org/doc/install) installed.

```sh
$ go get -u github.com/ericjesse/whatsmyip
$ cd $GOPATH/src/github.com/ericjesse/whatsmyip
```

## Running locally

```sh
$ cd $GOPATH/src/github.com/ericjesse/whatsmyip
$ go build
```

If some dependencies are missing on your local $GOPATH directory, you can install them using [godep](https://github.com/tools/godep):
```sh
$ godep go install
```

Then, run the app:
```sh
$ ./whatsmyip -port 5000 -dbType postgres -dbUrl <url to the db> -dataRetentionDuration <data retention duration>
```
Your app should now be running on [localhost:5000/ip](http://localhost:5000/ip).

## Parameters
- port: the port to listen for the HTTP communications. It can also be set by the environment variable ```PORT```, which has the priority.
- dbType: the type of the database to use to save usage data. Only ```postgres``` is supported, which is the default value if the parameter is omitted.
- dbUrl: the URL to the database, e.g: ```postgres://whatsmyip:whatsmyip@hostname/whatsmyip?sslmode=disable```. It can also be set by the environment variable ```DATABASE_URL```, which has the priority.
- dataRetentionDuration: the duration for the data retention, e.g: "300ms", "-1.5h" or "2h45m". Default is ```6 weeks``` if omitted. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h". It can also be set by the environment variable ```DATA_RETENTION```, which has the priority.

## Running locally using Heroku

Make sure you have the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

Create the file .env containing the environment variables in your working directory:
```
PORT=5000
DATABASE_URL=postgres://whatsmyip:whatsmyip@localhost/whatsmyip?sslmode=disable
DATA_RETENTION=1008h
LOG_DEBUG=false
```
Then run the following commands:

```sh
$ cd $GOPATH/src/github.com/ericjesse/whatsmyip
$ go install .
$ heroku local web
```

Your app should now be running on [localhost:5000/ip](http://localhost:5000/ip).

Note: if you get the error "/bin/sh: whatsmyip: command not found", make sure the current directory is in your $PATH.

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
```

## Create and run the Docker image

```sh
$ make
$ docker run -d -p 5000:5000 -e DATABASE_URL=postgres://whatsmyip:whatsmyip@ipOfTheDb/whatsmyip?sslmode=disable ericjesse/whatsmyip
```
You can now use the commands below to use the service.

## Using the service with cURL

To get the result as JSON (default if the HTTP header _"Accept"_ is omitted):
```sh
curl -H "Accept-Encoding: gzip" -H "Accept: application/javascript" "http://localhost:5000/ip"
```
To get the result as XML:
```sh
curl -H "Accept-Encoding: gzip" -H "Accept: application/xml" "http://localhost:5000/ip"
```

If you deployed the service on Heroku, replace _http://localhost:5000_ by the URL of your Heroku's app.  

You can also omit the HTTP header _"Accept-Encoding: gzip"_ to receive a plain response.
