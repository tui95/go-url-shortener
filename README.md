# URL Shortener

An implementation of URL shortener by using Golang. This is intended for learning the Golang.

## Algorithm

This uses the algorithm from https://stackoverflow.com/a/742047/18571631

## API Reference

### **POST** `/`

Create short URL

#### Request Body

```json
{
  "url": "https://example.com"
}
```

#### Response Body

```json
{
  "url": "http://localhost:8080/abc"
}
```

### **GET** `/{id}`

Get short URL

## Development

### Running server with live reload

```shell
make watch
```

### Running tests

```shell
go test -v ./...
```