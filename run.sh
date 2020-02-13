#!/bin/sh
nohup ./loginserver &
nohup ./gateserver &
nohup ./masterserver &
nohup ./worldserver &
nohup ./gameserver &
