#!/bin/sh
JS=/usr/share/nginx/html/assets/index-SDZGhHK5.js
for s in "Table Defaults" "Advanced Spark" "SSE-KMS" "rest-catalog" "glue-catalog" "nessie" "pyiceberg" "Objects are not valid" "formatApiDetail"; do
  c=$(grep -o "$s" "$JS" 2>/dev/null | head -1)
  echo "$s => '$c'"
done
