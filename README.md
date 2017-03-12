
# What's my IP

Simple Go service to detect the remote IP address of the caller.

The application is designed to run on [Heroku](https://www.heroku.com).

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```sh
$ go get -u github.com/ericjesse/whatsmyip
$ cd $GOPATH/src/github.com/ericjesse/whatsmyip
$ heroku local web
```

Your app should now be running on [localhost:5000](http://localhost:5000/ip).

Note: if you get the erreur "/bin/sh: whatsmyip: command not found", make sure the current directory is in your $PATH.

You can now test it with the cURL command:
```sh
curl -X GET -H "Cache-Control: no-cache" "http://localhost:5000/ip"
```

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
$ heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)


## Documentation

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)
