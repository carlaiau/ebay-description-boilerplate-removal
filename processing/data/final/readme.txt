These files are the complete Ebay Dumps as described in the struct:
https://github.com/carlaiau/ebay-description-boilerplate-removal/blob/master/processing/main.go#L33-#L146

a - r.xml contain the scrape results from the intial API Endpoint after removal of empty results. 
s.xml contains the results from the the secondary API Endpoint after removal of empty results.
t.xml contains the documents not retrieved via API, but recreated from the original documents.tsv file.

Empty was based on a missing title, not a missing description.

Filename	Files	Title	Description
a.xml		35901	35901	35900
b.xml		35375	35375	35373
c.xml		35303	35303	35302
d.xml		34759	34759	34755
e.xml		34681	34681	34679
f.xml		34580	34580	34580
g.xml		34367	34367	34367
h.xml		34068	34068	34067
i.xml		34059	34059	34057
j.xml		33791	33791	33789
k.xml		33654	33654	33653
l.xml		33584	33584	33578
m.xml		33591	33591	33590
n.xml		33266	33266	33265
o.xml		33185	33185	33184
p.xml		33142	33142	33135
q.xml		32923	32923	32922
r.xml		32711	32711	32710
s.xml		65805	65805	65800
t.xml		220937	0	0

Totals
Files	Title	Description
899682	678745	678706
