#!/usr/bin/env zsh

# rm -rf rsstmp
# mkdir rsstmp
# waybackpack https://feeds.a.dj.com/rss/RSSOpinion.xml -d rsstmp > w.log 2>&1 &

KillIt () {
    echo "KillIt - looking for ${KILLME}"
    I=0
    ps -ef | grep waybackpack | grep -v grep | while read line ;
    do
        a=(${(s/ /)line})
        (( I++ ))
        # echo "line = ${line}"
        echo -n "Instance ${I}   -->  kill $a[2]"
        kill "$a[2]"
        echo "kill return code: $?"
        echo
    done;

    echo "Done."
    echo "Killed ${I} instance(s) of ${KILLME}"
}

KILLME="wd.sh" ; KillIt
KILLME="wbp.sh" ; KillIt
KILLME="waybackpack" ; KillIt
