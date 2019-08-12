# Ebay Parser
Basic scraper that takes a tsv file of titles and categories as input.
The scraper Searches Ebay's Finder API using each product title.

Upon recieveing a response, we concatenate the product title and category coming out of the API and compare that to the title and category from the TSV file. 

If these match, then we conclude that the product is correctly identified. 13
After we have 20 correctly identified product IDs, we then query the Shopping API to get the actual details about each listing

### Prerequisities
Go
Ebay Production API ID
TSV File of products

### Usage
```
go get github.com/valyala/tsvreader
go build .
./Ebay-Parser
```

### CLI arguments
app_id required
Usage of ./ebay-play for scraping
```
./Ebay-Parser scrape api_key input_file output_file
```

For Testing, this will load a scraped file and give counts of found numbers, as well as dump items where the
Title != OrigTitle
```
/.Ebay-Parser test input_file
```


