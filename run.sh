#!/bin/bash
#for letter in a b c d e f g h i j
#do
#   (unzip merged-50ks/$letter.xml.zip
#    ./Ebay-Collection filter $letter.xml filtered/$letter.xml
#    rm $letter.xml) &
#done

for letter in k l m n o 
do 
    (unzip merged-50ks/$letter.xml.zip
    ./Ebay-Collection filter $letter.xml filtered/$letter.xml
    rm $letter.xml) &
done

for letter in p q r missing;
do 
    (unzip merged-50ks/$letter.xml.zip
    ./Ebay-Collection filter $letter.xml filtered/$letter.xml
    rm $letter.xml) &
done