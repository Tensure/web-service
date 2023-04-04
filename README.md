# web-service

## Developer Setup

### Before you start

You will need [Git](https://git-scm.com/), [Google Cloud CLI](https://cloud.google.com/sdk/gcloud), and [Golang 1.18 or later](https://golang.org/dl/). It is also recommended that you use [GoLand](https://www.jetbrains.com/go/) or [Visual Studio Code](https://code.visualstudio.com) for development.

Once you have those tools installed you should clone the repo [here](Github URL).

### Setting up web-service

To begin, open up the IDE of your choice (preferably Visual Studio Code, as GoLand requires a module import in the settings when you clone from GitHub) and use the terminal to pull the code using:

``` git
git pull
```

This ensures you have the latest code available.

After pulling run `go mod init`.  This will set up a couple different tools we've configured for this repository.


If you followed the directions to the letter, you should now direct your terminal instance to your working directory if you haven't already using:

``` bash
cd path/to/your/web-service's/root
```

When you're in this directory you can use the command:

``` golang
go build .
```

This builds the web-service so you can test routes.
The command to run web-service after build is:

``` bash
go run .
```
When you're ready to push to GCP, visit this [page](https://cloud.google.com/run/docs/quickstarts/build-and-deploy/deploy-go-service).