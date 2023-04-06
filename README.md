# go-micro-urlshortner

You can create tables directly by importing users.sql file to db

OR

we can use soda migration to create everthing as code.

Step 1: create database.yml file
Step 2: soda generate fizz CreateUrlMappingTable
Step 3: write fizz query to create table in up.fizz file - https://gobuffalo.io/documentation/database/fizz/
Step 4: write fizz query to drop table in down.fizz file
Step 4: soda migrate
