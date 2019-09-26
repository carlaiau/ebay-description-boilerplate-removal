
for file in goose/*
do
    curl -s -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/goose/_bulk" --data-binary "@$file"
    echo "@$file indexed!"
done 