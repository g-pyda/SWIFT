#!/bin/sh
echo "Running tests in xlsxParser..."
go test ./src/xlsxParser 
echo "Running tests in databaseControl..." 
go test ./src/databaseControl
echo "Running tests in REST..." 
go test ./src/REST
