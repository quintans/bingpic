# bingpic
Bing picture of the day hack

This project has two utilities to download Picture of the Day from Bing. One uses the API and the other scrapes the Bing site.

The API has a limitation on the number of files that can be downloaded. You can only get 8 pictures per call and you can only request until 8 days prior today. Essentially you can only get a maximum of 16 images (well, 15 because there is a bug the pagination on the API side?)

Scrapping does not have that limitation. At least I was able to scrape up to 100 images.

## Downloading 
Arguments:
- n: Number of picture we want to download (default is 1)
- dest: Directory where we store the downloaded files (defaults to current dir)
- roll: number of files we keep on the dest directory (default is 100)
- v: verbose

Eg:
```sh
go run ./cmd/bingscraper/main.go -dest=./output -n=50
```

## Cross platform build
```sh
GOOS=linux GOARCH=amd64 go build -o ./dist  ./cmd/bingscraper
```

```sh
GOOS=linux GOARCH=amd64 go build -o ./dist  ./cmd/bingapi
```
