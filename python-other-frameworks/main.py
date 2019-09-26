'''
This file needs to be given the ready xml files containing the fields:
DOCNO
ORIGTITLE
CSDESCRIPTION

It will then run the description through the requested boilerplate removal technique
and restructure it back into it's original format. 
This document then needs processed using the processing/main.go JSON creator so that 
they can be loaded into ElasticSearch
'''

import io
import sys
from bte import html2text
import justext

technique = sys.argv[1]
input_file = io.open(sys.argv[2], mode="r",  encoding="utf-8")
output_file = io.open(sys.argv[3], "w+",  encoding="utf-8")

contents = input_file.read()
individual_docs = contents.split("<DOC>")
counter = 0
stringBuilder = []
count = 0
for doc in individual_docs[1:]:
    sections = doc.split("<CSDESCRIPTION>")
    header = "<DOC>" + sections[0] + "<CSDESCRIPTION>"
    print(sections[0])
    remaining = sections[1]
    sections = remaining.split("</CSDESCRIPTION>")
    extracted = ""
    if technique == 'bte' and len(sections[0]) > 0:
        extracted = html2text(sections[0], False, False)
    
    elif technique == 'justext' and len(sections[0]) > 1:
        paragraphs = justext.justext(sections[0], justext.get_stoplist("English"))
        for paragraph in paragraphs:
            if not paragraph.is_boilerplate:
                extracted += paragraph.text
        
    footer = "</CSDESCRIPTION>" + sections[1]

    stringBuilder.append((header + extracted + footer)) 
    count += 1
    if count % 1000 == 0:
        print(str(count) + " completed!")

output_file.write( ''.join( stringBuilder ) )
    
input_file.close()
output_file.close()