CREATE DATABASE login;
USE login;
CREATE TABLE USER(id INT auto_increment PRIMARY KEY, username TEXT, password TEXT);
INSERT INTO USER(username, password)
VALUES ('Frendy', '12345678');