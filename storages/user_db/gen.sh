#!/usr/bin/env bash

mysql-orm-gen -sql_file=./neuron-user_db.sql -orm_file=./neuron-user_db-gen.go -package_name="user_db"