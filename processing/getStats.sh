#!/bin/bash
./processing getStats ~/School/data/indexs/td-original/titles-descriptions.xml &
./processing getStats ~/School/data/indexs/td-after-bp/article.xml &
./processing getStats ~/School/data/indexs/td-after-bp/default.xml &
./processing getStats ~/School/data/indexs/td-after-bp/largest.xml &
./processing getStats ~/School/data/further-bp/bte/bte.out &
./processing getStats ~/School/data/further-bp/goose/goose.out &
./processing getStats ~/School/data/further-bp/justext/justext.out &