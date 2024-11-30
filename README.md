# bingpic
Bing picture of the day hack

This has two utilities to download Picture of the Day from Bing usi an API or by scraping the Bing site.

The API has limitations one the number of files that can be downloaded. At the moment limits to 15 with calls returning 8 each time.

Scrapping does not have that limitation.

## Downloading 
Arguments:
- n: Number of picture we want to download (default is 1)
- dest: Directory where we store the downloaded files (defaults to current dir)
- roll: number of files we keep on the dest directory (default is 100)

Eg:
```sh
go run ./cmd/bingscraper/main.go -dest=./output -n=50
```

## Cross platform build
```sh
GOOS=linux GOARCH=amd64 go build -o ./dist  ./cmd/bingscraper;
```

```sh
GOOS=linux GOARCH=amd64 go build -o ./dist  ./cmd/bingapi
```
