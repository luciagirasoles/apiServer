CREATE DATABASE IF NOT EXISTS apiServer;
CREATE TABLE apiServer.domains (id SERIAL PRIMARY KEY, domain STRING);
CREATE USER IF NOT EXISTS maxroach;
GRANT ALL ON apiserver.domains TO maxroach;