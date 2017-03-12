
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
$ ./whatsmyip 5000
```

Your app should now be running on [localhost:5000/ip](http://localhost:5000/ip).

## Running locally using Heroku

Make sure you have the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```sh
$ cd $GOPATH/src/github.com/ericjesse/whatsmyip
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
$ docker run -d -p 5000:5000 ericjesse/whatsmyip
```
You can now use the commands below to use the service.

## Using the service with cURL

To get the result as JSON (default if the HTTP header _"Accept"_ is omitted):
```sh
curl -H "Accept-Encoding: application/gzip" -H "Accept: application/javascript" "http://localhost:5000/ip"
```
To get the result as XML:
```sh
curl -H "Accept-Encoding: application/gzip" -H "Accept: application/xml" "http://localhost:5000/ip"
```

If you deployed the service on Heroku, replace _http://localhost:5000_ by the URL of your Heroku's app.  

You can also omit the HTTP header _"Accept-Encoding: application/gzip"_ to receive a plain response.

