#!/usr/bin/env bash

mysql-orm-gen -sql_file=./account_db.sql -orm_file=./account_db-gen.go -package_name="account_db"