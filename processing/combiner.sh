#!/bin/bash
# Run within folder of combined XML files
for filename in *.xml; do
    sed '1d;$d' $filename > "pro-$filename"
done

echo "<Superset>" > merged.xml
for filename in pro-*.xml; do
    cat $filename >> merged.xml
done
echo "</Superset>" >> merged.xml    
rm pro-*.xml