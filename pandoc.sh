#!/usr/bin/env bash
# Transforms Markdown format to html for use with md-publisher.
# Requires pandoc, pandoc-crossref, and pandoc-citeproc
# # See https://pandoc.org/MANUAL.html#extension-yaml_metadata_block
set -e
input="$1"
output="$2"
format="$3"
opts=()
opts+=(--standalone)
opts+=(--resource-path "$PWD")
opts+=(--filter pandoc-crossref)
opts+=(--filter pandoc-citeproc)

## Add here your citation style
opts+=(--csl "$HOME/.csl/science.csl")

## Add here your global BibTex library if you like
#opts+=(--bibliography "$HOME/.local/share/pandoc/Bibliography.bib")
opts+=(--from markdown+yaml_metadata_block+implicit_figures+fenced_divs+citations+table_captions)
opts+=(--to "$format")
opts+=(--webtex)
opts+=(--output "$output")
opts+=("$input")
echo "Running pandoc with ${opts[*]}"
pandoc "${opts[@]}"
