#!/usr/bin/env bash
set -e
input="$1"
output="$2"
format="$3"
opts=()
opts+=(--standalone)
opts+=(--resource-path "$PWD")
opts+=(--filter pandoc-crossref)
opts+=(--filter pandoc-citeproc)
opts+=(--csl "$HOME/.csl/science.csl")
#opts+=(--bibliography "$HOME/.local/share/pandoc/Bibliography.bib")
# table_captions
opts+=(--from markdown+yaml_metadata_block+implicit_figures+fenced_divs+citations+table_captions)
opts+=(--to "$format")
#opts+=(--pdf-engine xelatex)
#opts+=(--gladtex)
#opts+=(--katex)
#opts+=(--mathjax)
opts+=(--webtex)
#opts+=(--webtex "https://latex.codecogs.com/svg.latex?")
opts+=(--output "$output")
opts+=("$input")
echo "Running pandoc with ${opts[*]}"
#PANDOC="$HOME/.local/bin/pandoc"
#PANDOC="/usr/local/bin/pandoc"
PANDOC="pandoc"
$PANDOC "${opts[@]}"
