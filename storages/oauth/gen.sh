#!/usr/bin/env bash

mysql-orm-gen -sql_file=./oauth.sql -orm_file=./oautht-gen.go -package_name="oauth"