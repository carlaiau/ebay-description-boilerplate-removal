
for file in article/*
do
    curl -s -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/article_bm25/_bulk" --data-binary "@$file"
    echo "@$file indexed!"
done 