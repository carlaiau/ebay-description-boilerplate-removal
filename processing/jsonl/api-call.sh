#!/bin/bash

for file in filtered/original/*
do
    curl -s -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/original_filtered/_bulk" --data-binary "@$file"
    echo "@$file indexed!"
done 

for file in filtered/bp-a/*
do
    curl -s -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/bp_article_filtered/_bulk" --data-binary "@$file"
    echo "@$file indexed!"
done 

for file in filtered/bp-d/*
do
    curl -s -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/bp_default_filtered/_bulk" --data-binary "@$file"
    echo "@$file indexed!"
done 

for file in filtered/bp-l/*
do
    curl -s -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/bp_largest_filtered/_bulk" --data-binary "@$file"
    echo "@$file indexed!"
done 

for file in filtered/goose/*
do
    curl -s -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/goose_filtered/_bulk" --data-binary "@$file"
    echo "@$file indexed!"
done 

for file in filtered/justext/*
do
    curl -s -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/justext_filtered/_bulk" --data-binary "@$file"
    echo "@$file indexed!"
done 
# curl -s -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/original/_bulk" --data-binary "@deletes.jsonl"
# echo "Completed!"
# curl -s -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/goose/_bulk" --data-binary "@deletes.jsonl"
# echo "Completed!"
# curl -s -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/bte/_bulk" --data-binary "@deletes.jsonl"
# echo "Completed!"
# curl -s -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/default/_bulk" --data-binary "@deletes.jsonl"
# echo "Completed!"
# curl -s -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/justext/_bulk" --data-binary "@deletes.jsonl"
# echo "Completed!"
