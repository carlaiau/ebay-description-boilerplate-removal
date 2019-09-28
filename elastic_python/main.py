import sys
from elasticsearch import Elasticsearch
es = Elasticsearch(
    urls=['localhost'], 
    port=9200)

queries = [
"kazar comic",
"dc digest",
"vintage bronze sculpture",
"roman replica",
"deer sculpture",
"usmc",
"multicam",
"new woodland pants",
"jim shore bunny",
"disney traditions",
"heartwood creek angels",
"mexico carved wood",
"oaxacan",
"milagros",
"wooden radio",
"rare zenith",
"stromberg carlson",
"hexagon box",
"hinged trinket box",
"blue glass box",
"whale oil",
"aladdin lamp",
"glass oil lamp",
"attack on titan",
"tokyo ghoul",
"joan crawford",
"barbara feldon",
"coin banks vintage",
"wooden bank",
"building bank",
"cgc ss stan lee",
"spiderman 1",
"marvel comics 9.8",
"gerz stein",
"stein",
"german beer mugs",
"tole lamp",
"brass cherub",
"vintage desk",
"predator hunters",
"the walking dead 20",
"walking dead comic #100",
"kankakee",
"galesburg illinois",
"decatur illinois",
"payphone",
"vintage rotary phone",
"rotary wall phone",
"parker duofold",
"parker 51",
"zippo harley",
"camel zippo",
"beer signs budweiser neon",
"budweiser clydesdale",
"skeleton keys",
"padlock",
"breyerfest",
"breyer roan",
"vintage planter",
"cat planter",
"vintage pic",
"toy steam roller",
"plaid wool blanket",
"pendleton blanket",
"disney pewter",
"disney ceramic figurines",
"gumball machine",
"penny machine",
"soda pop sign",
"fence sign",
"match holder porcelain",
"pig match",
"bullet shot glass",
"espresso shot glass",
"vintage metal chest",
"wooden tool box",
"kennedy memorabilia",
"king george queen mary",
"masonic apron",
"trail of painted ponys",
"united delco",
"christmas carolers",
"calif license plates",
"ny license plate",
"military police patches",
"white flower pitcher vase",
"ride poster",
"greek vase",
"absolut vodka bottle",
"holt howard",
"girl scout books",
"angel ornaments",
"piggy piggy banks",
"queen mary",
"goofy phone",
"tupperware bowls",
"sluice",
"x-men 9.8",
"lenox picture frame",
"precious moments figurines",
"victorian photo album",
"woven basket",
"epaulettes",
"waterman rollerball pen",
"magnavox transistor radio",
"tom clark gnomes",
"pinball machine",
"scales of justice",
"see hear speak no evil",
"railroad book",
"disney bank",
"ronald reagan autographed",
"stop sign",
"darth vader 3",
"wolverine 15",
"fire department historical memorabilia",
"tea tags",
"resin easter",
"batman statue",
"michael myers mask",
"large easter eggs",
"m998 hmmwv",
"star trek cards",
"poker table",
"national cash register",
"sewing box",
"empty green bottles",
"herb grinder",
"barber bottle",
"book storage box",
"rooster creamer",
"beer tap handle",
"pyramid of giza",
"pilot helmet",
"old fitzgerald",
"hess truck",
"alpha kappa alpha",
"slide rule",
"medals ribbons us",
"movie refrigerator magnets",
"batman and robin 1",
"vintage sheaffer fountain pen",
"harley-davidson zippo lighter",
"star wars card trader pack art",
"the movie collection goku figure",
"dragon ball statue resin",
"bell key chain keyfob",
"a charlie brown christmas",
"atlantic city casino chips",
"green lantern green arrow comics"]

boost = [0,1,2,3,4,5,6,7,8,9,10]
index_to_use = sys.argv[1]


# Create a filename based on output
for b in range(len(boost)):
    hits = []
    for q in range(len(queries)):
        if boost[b] == 10:
            query = \
            {  
                "size": 5000,
                "query": { "match": { "Title": queries[q] } }
            }
        elif boost[b] == 0:
            query = \
            {  
                "size": 5000,
                "query": { "match": { "Description": queries[q] } }
            }
        else:
            query = \
            {  
                "size": 5000,
                "query": {
                    "multi_match" : {
                        "query": queries[q],
                        "fields": ["Title^" + str(boost[b]), "Description^" + str(10 - boost[b])]
                    }
                }
            }

        res = es.search(
            index=index_to_use, 
            body=query, 
            request_timeout=120,
            scroll='2m'
            )
        
        # Get the scroll ID
        scroll_id = res['_scroll_id']
        scroll_size = len(res['hits']['hits'])

        while scroll_size > 0:

            for rank, hit in enumerate(res['hits']['hits'], 1):
                hits.append('{}\tQ{}\t{}\t{}\t{}\t{}'.format(q + 1, 0, hit['_id'], rank, hit['_score'], 'title-test'))
            
            res = es.scroll(scroll_id=scroll_id, scroll='2m')
            
            # Update the scroll ID
            scroll_id = res['_scroll_id']

            # Get the number of results that returned in the last scroll
            scroll_size = len(res['hits']['hits'])

        print(str(q + 1) + ": " + queries[q] + " done!")

    file_name = index_to_use  + "-" + str(boost[b]) + ".result"
    text_file = open("filtered-output/" + file_name, "w")
    
    text_file.write('\n'.join(hits))
    
    text_file.close()
    print(file_name +" written!\n")

