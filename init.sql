CREATE DATABASE IF NOT EXISTS swiftdb;
CREATE DATABASE IF NOT EXISTS testswiftdb;

GRANT ALL PRIVILEGES ON swiftdb.* TO 'swiftuser'@'%';
GRANT ALL PRIVILEGES ON testdb.* TO 'swiftuser'@'%';