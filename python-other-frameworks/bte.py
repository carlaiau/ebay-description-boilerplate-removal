"""
This module implements Finn's BTE (Body Text Extraction) algorithm for
extracting the main body text of a web page and avoiding the surrounding
irrelevant information. The description of the algorithm can be found
in A. Finn, N. Kushmerick, and B. Smyth. Fact or Fiction: Content
classification for digital libraries. In DELOS Workshop: Personalisation
and Recommender System in Digital Libraries, 2001.

Python implementation by Jan Pomikalek <xpomikal@fi.muni.cz>
"""

import re

def html2text(html_text, preserve_par=False, preserve_head_list_par=False):
    """
    Converts HTML to plain text with boilerplate removed.
    If preserve_par is set to True, paragraph mark-up will be preserverd.
    If preserve_head_list_par is set to True, paragraph mark-up will be
    preserverd and headers and list items marked with <h> and <l> tags
    respectively.
    """

    cleaned_text = preclean(html_text)
    tokens = tokenise(cleaned_text)
    (start, end) = bte(tokens)
    main_body = tokens[start:end+1]
    cleaned_body = find_paragraphs(main_body, preserve_head_list_par)

    # separate paragraphs with newlines
    blocks = []
    block = []
    for token in cleaned_body:
        if token_value(token) > 0:
            block.append(token)
        else:
            if len(block) > 0:
                blocks.append(" ".join(block))
                block = []
            if preserve_par or preserve_head_list_par:
                block.append(token)
    if len(block) > 0:
        blocks.append(" ".join(block))

    return "\n".join(blocks)


def preclean(html_text):
    """
    HTML preprocessing -- striping headers, scripts, styles; replacing HTML
    entities.
    """

    

    # strip all but body
    cleaned_text = re.compile('^.*<body(\s+[^>]*)?>', re.S | re.I
            ).sub('', html_text)
    cleaned_text = re.compile('</body>.*$', re.S | re.I
            ).sub('', cleaned_text)

    # strip scripts
    cleaned_text = re.compile('<script\\b[^>]*>(.*?)<\\/script>',
            re.S | re.I).sub('', cleaned_text)

    # strip styles
    cleaned_text = re.compile('<style\\b[^>]*>(.*?)<\\style>',
            re.S | re.I).sub('<style></style>', cleaned_text)

    # html entities
    cleaned_text = html_entities(cleaned_text)
    
    return cleaned_text


def html_entities(html_text):
    "Substitution of the most commonly used HTML entities."
    html_text = re.sub('&quot;', '"', html_text)
    html_text = re.sub('&nbsp;', ' ', html_text)
    html_text = re.sub('&#39;', "'", html_text)
    return html_text


def tokenise(html_text):
    """
    Tokenises HTML document to a sequence of HTML tags and strings of
    non-whitespace characters (words).
    """
    return [g1 for (g1, g2) in re.findall('(<([^>]|\s)+>|[^\s<]+)', html_text)]


def bte(tokens):
    """
    BTE algorithm. Expects a sequence of HTML tags and words as input parameter.
    Outputs a pair of indices which indicate the beginning and end of the main
    body.
    """

    # find breakpoints
    breakpoints = []
    prev_value = None
    sum_value = 0
    for i in range(len(tokens)):
        cur_value = token_value(tokens[i])
        if prev_value and cur_value != prev_value:
            breakpoints.append((i-1, sum_value))
            sum_value = 0
        sum_value+= cur_value
        prev_value = cur_value
    breakpoints.append((len(tokens)-1, sum_value))

    # find breakpoints range which maximises the score
    max_score = 0
    max_start = 0
    max_end   = 0
    for i in range(len(breakpoints)):
        score = breakpoints[i][1]
        if score > max_score:
            max_score = score
            if i > 0: max_start = breakpoints[i-1][0]+1
            else:     max_start = 0
            max_end   = breakpoints[i][0]
        for j in range(i+1, len(breakpoints)):
            score+= breakpoints[j][1]
            if score > max_score:
                max_score = score
                if i > 0: max_start = breakpoints[i-1][0]+1
                else:     max_start = 0
                max_end   = breakpoints[j][0]

    return (max_start, max_end)


def token_value(token):
    "Returns -1 if the token is HTML tag, 1 otherwise (if word)."
    if token.startswith('<'):
        return -1
    else:
        return 1


def find_paragraphs(tokens, tag_h_l=False):
    """
    Marks paragraph blocks with <p>. If tag_h_l set to True, headers and
    list items are also detected and marked with <h> and <l> respectively.
    """

    PAR_FIND_TAGS = ['p', 'div', 'hr', 'blockquote', 'table']
    PAR_REPLACE_TAG = '<p>'
    HEADER_FIND_TAGS = ['h1', 'h2', 'h3']
    HEADER_REPLACE_TAG = '<h>'
    LIST_FIND_TAGS = ['li']
    LIST_REPLACE_TAG = '<l>'
    result = [PAR_REPLACE_TAG]

    in_paragraph = False
    for token in tokens:
        if token_value(token) > 0:
            result.append(token)
            in_paragraph = True
        else:
            if not in_paragraph:
                continue
            m = re.search('^<([^\s>]+)', token)
            if not m:
                continue
            tag = m.group(1).lower()
            if tag in PAR_FIND_TAGS:
                result.append(PAR_REPLACE_TAG)
                in_paragraph = False
            if tag in HEADER_FIND_TAGS:
                if tag_h_l:
                    result.append(HEADER_REPLACE_TAG)
                else:
                    result.append(PAR_REPLACE_TAG)
                in_paragraph = False
            if tag in LIST_FIND_TAGS:
                if tag_h_l:
                    result.append(LIST_REPLACE_TAG)
                else:
                    result.append(PAR_REPLACE_TAG)
                in_paragraph = False

    return result

