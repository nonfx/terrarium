#!/bin/bash
echo "# Table of contents\n" > TOC.md
rsync -avn . /dev/shm --exclude-from .gitignore --exclude-from .git/info/exclude | grep "\.md$" | while IFS= read -r line; do
  printf "* [%s](./%s)\n"  $line $line >> TOC.md
done;
