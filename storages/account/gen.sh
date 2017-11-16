#!/usr/bin/env bash

mysql-orm-gen -sql_file=./account.sql -orm_file=./account-gen.go -package_name="account"