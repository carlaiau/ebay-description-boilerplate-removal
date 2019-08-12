# Ebay Parser
Basic scraper that takes a tsv file of titles and categories (from SIGIR 2019 workshop) as input.
The scraper Searches Ebay's Finder API using each product title. We are presently not using the findCompletedItems API endpoint.

Upon recieveing a response, we concatenate the product title and category coming out of the API and compare that to the title and category from the TSV file. 
If these match, then we conclude that the product is correctly identified. 13
After we have 20 correctly identified product IDs, we then query the Shopping API to get the actual details about each listing

### Prerequisities
Go
Ebay Production API ID
SIGIR 2019 TSV File

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

### Timings / Strategy
We're still getting rate limited on the findingAPI (Getting EbayID for products)

If we leave out the findCompletedItems API requests until we have done an initial run through, it makes parallelising easier. I just cut up the document list into bunches of 5000, and daily assign each bunch to one API key. That way we are predictibly maximising API Limits, whilst getting around a 70% hit rate on requests.
We can then split up all the documents that were missed on the initial run and add to the collection after wards. ExpiredListings are still API queriable for up to 90 days after they have expired.

### This results in:
895182 Remaining Documents
5000 documents queried per day per API Key
179 Days
5 keys = 35, 10 keys = 18, 20 keys = 9

After initial run through, I can then re-query all missing documents using the findCompletedItems API. The benefit of this approach is pre-known request number, so that we can optimise for our limited set of API keys.

### Misc
Regarding boilerplate removal, I think I can start working on techniques well before we have all the data. 
Web2Text's main exploration/discussion was on a set of ~1000 documents (from memory).
The slower data retrieval isn't ideal, but it's not a complete deal breaker. Our collection will suffer from the expiring items, but I still expect over 650K items. 
I will try to obtain 10 or more keys. 

