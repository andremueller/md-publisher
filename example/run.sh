#!/bin/bash
set -e
echo "Converting demo_article markdown to html"
bash ../pandoc-medium --output demo_article.html demo_article.md
echo "Done. Now you could check demo_article.html before uploading"
echo ""
read -p "PRESS ENTER TO UPLOAD OR CTRL-C TO ABORT"

echo "Uploading demo_article to Medium as DRAFT"
md-publisher --log-level 4 publish demo_article.html
echo "Done."

