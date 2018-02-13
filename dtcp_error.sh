#!/bin/bash

commands="dtcp -mean_field_recursion"
labels="0"
 n=4
 while [  $n -lt 100 ]; do
     commands="$commands, dtcp -full_complete -n=$n"
     labels="$labels, n=$n"
     let n=n+4
 done

python main.py \
-commands="$commands" \
-labels="$labels" \
-shared="-T=4 -steps=200000" \
-show_plot -type=error