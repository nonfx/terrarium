#!/bin/bash
# Copyright (c) Ollion
# SPDX-License-Identifier: Apache-2.0


# Specify the path to the database dump file
dump_file="/docker-entrypoint-initdb.d/cc_terrarium.psql"

# Check if the dump file exists
if [ ! -f "$dump_file" ]; then
  pwd
  echo "File $dump_file not found. Skipping database restore."
  exit 0
fi

# Restore the database using pg_restore command
pg_restore -U $POSTGRES_USER -d $POSTGRES_DB -W $dump_file

# Check the exit status of the pg_restore command
if [ $? -ne 0 ]; then
  echo "Error: Database restore failed."
  exit 1
else
  echo "Database restore completed successfully."
fi
