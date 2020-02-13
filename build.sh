#!/bin/sh
go build server/main/loginserver.go
go build server/main/gateserver.go
go build server/main/masterserver.go
go build server/main/worldserver.go
go build server/main/gameserver.go