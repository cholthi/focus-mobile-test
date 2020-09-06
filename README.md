# focus-mobile-test
This repo contains code for my interview with Focus Mobile

### Installation
 `go get github.com/cholthi/focus-mobile-test`

 Please be sure your `GOPATH` env variable is set. This project does not use go modules and `go get` requires $GOPATH to work.

 The above command will fetch the project to your `GOPATH/src` and install the bin to `GOPATH/bin`. Be sure those directory exit. Which they should if you're using go.

 ### Running the command line app.

 `focus-mobile-test --help` // Will out the usage info

 `focus-mobile-test supported KES` // Will find out if ISO 4217 currency code is supported. The MVP feature

 The `supported` is a sub command for CLI app. This is so the app can be easily extended.


 ## Attention.
  The app can take two global options file paths. The --cacheFile is used to cache contents of the currency file locally to avoid hitting the network if the file has not changed.
  The --versionFile is used to track changes to remote currency file.

  These files must be readable and writable to the app and have no special format and are just text.

  if not supplied, they default to `./modified.lock` and `./currency.cache` for --versionFile and --cacheFile respectively. Make sure those files exist in the current dir you're running the CLI app.

  ### Running Test
   This project comes with some tests. To run, do the following.

   `cd $GOPATH/src/github.com/cholthi/focus-mobile-test`
   Then run
   `go test`
