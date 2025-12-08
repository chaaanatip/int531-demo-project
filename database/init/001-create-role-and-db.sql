-- create role and database; change password if needed
CREATE ROLE sreuser WITH LOGIN PASSWORD 'srepass';
CREATE DATABASE appdb OWNER sreuser;
