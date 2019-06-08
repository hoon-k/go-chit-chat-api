#!/bin/bash

psql -U postgres -c "CREATE DATABASE chitchat_users"
psql -U postgres -c "CREATE DATABASE chitchat_discussions"
psql -U postgres -d chitchat_users -f /docker-entrypoint-initdb.d/setup-users-db.txt
psql -U postgres -d chitchat_discussions -f /docker-entrypoint-initdb.d/setup-discussions-db.txt