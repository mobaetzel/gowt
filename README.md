# GoWT

GoWT is a small cli and web service to fill word templates with data from a json file.

## How to build and run the CLI

Clone this repo and build it with

```bash
go build -o gowt src/main.go
```

You can then run the cli with either

```bash
./gowt process --template ./examples/example.odt --data ./examples/data.json --output output.odt
```

or

```bash
./gowt serve --host 0.0.0.0 --port 8080 --webui
```

## How to build and run the docker image

Run

```bash
docker build -t gowt .
docker run -p 3000:3000 gowt
```