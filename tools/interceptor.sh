#!/bin/sh

### Define or variables
cmdLine=$*
checklicBin="/Volume/home/admin/checklic"
now=`date`
parent="$(ps -o comm= $PPID)"
logFile="/tmp/sflicCMDLog.log"
logLine="[${now}] (Parent: $parent) $cmdLine"
evLine="${checklicBin} ${cmdLine}"

### Save the command to or logfile
echo $logLine >> $logFile

### Eval the original command calling or checklic
eval "${checklicBin} ${cmdLine}"
