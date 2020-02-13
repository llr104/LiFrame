#!/bin/sh
go build server/main/createtables.go
./createtables orm syncdb db