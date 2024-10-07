# Sleeper CLI

Tool built in [Golang](https://github.com/golang/go) to access your Sleeper leagues data using Cobra CLI.

## Installation

To install the Sleeper CLI, you need to have Go installed on your machine.

Afterwards, copy `env-example.yaml` and create a new file called `env.yaml`. Then, replace the example LEAGUE_ID value with your actual Sleeper LEAGUE_ID, found in your league settings on the app.

Once you have Go installed, you can clone the repository and build the CLI tool.

```sh
git clone git@github.com:miriarte33/sleeper-cli.git
cd sleeper-cli
go build
```

## Usage

```sh
./sleeper -h
```
