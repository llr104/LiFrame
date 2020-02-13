#!/bin/sh
go build createtables.go
./createtables orm syncdb db