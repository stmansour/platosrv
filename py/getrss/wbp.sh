#!/usr/bin/env zsh

declare -a urls=(
  "https://feeds.a.dj.com/rss/RSSOpinion.xml"
  "https://feeds.a.dj.com/rss/RSSWorldNews.xml"
  "https://feeds.a.dj.com/rss/WSJcomUSBusiness.xml"
  "https://feeds.a.dj.com/rss/RSSMarketsMain.xml"
  "https://feeds.a.dj.com/rss/RSSWSJD.xml"
  "https://feeds.a.dj.com/rss/RSSLifestyle.xml"

  "https://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/World.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Africa.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Americas.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/AsiaPacific.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Europe.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/MiddleEast.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/US.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Education.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Politics.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Upshot.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/NYRegion.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Business.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/EnergyEnvironment.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/SmallBusiness.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Economy.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Dealbook.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/MediaandAdvertising.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/YourMoney.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/PersonalTech.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Sports.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Baseball.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/CollegeBasketball.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/CollegeFootball.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Golf.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Hockey.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/ProBasketball.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/ProFootball.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Soccer.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Tennis.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Science.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Climate.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Space.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Well.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Arts.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/ArtandDesign.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Books.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/SundayBookReview.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Dance.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Movies.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Music.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Television.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Theater.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/FashionandStyle.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/DiningandWine.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/tmagazine.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Jobs.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/RealEstate.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Automobiles.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Lens.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/Obituaries.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/MostEmailed.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/MostShared.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/MostViewed.xml"
  "https://rss.nytimes.com/services/xml/rss/nyt/sunday-review.xml"
)

#-------------------------------------------------------------------------------
#  usage  -  Display information explaining how to use this script
#-------------------------------------------------------------------------------
usage() {
    PROGNAME="wbp.sh"
    cat <<ZZEOF

WayBackPack RSS Feed Collector

    Usage:   ${PROGNAME}

    This command updates a mysql database named "plato" based on the OPTIONS
    provided in the run command.  If no OPTIONS are specified, then it will
    update every RSS Feed from its internal list based on the internal DTSTART
    value (currently set to ${DTSTART}).

CMD:
    -d  <startdt>
        Begin pulling feed information from this datetime value when feedurl
        is reached.  All subsequent reads will start from ${DTSTART}. The
        format is YYYYMMDDHHMMSS.  Hours, minutes, seconds are optional.

    -f  <feedurl>
        Skip feed URLs from the internal URL list until this url is reached.

    -h  Print this help message

    -p  <path>
        By default, ${PROGNAME} stores its temporary information in a volume
        containing "plato" and a directory named "rss".  You can override
        this default by specifying the full path using the -p option.

    -u  Bring the database up-to-date by filling in all the data starting from
        the latest date found in the Items table and ending on today's date.

EXAMPLES:
    Command to update plato database with all articles since ${DTSTART}:

        bash$  ./${PROGNAME}

    Command to start ${PROGNAME} and update all RSS feeds with articles that
    have been released since the last time this script was run:

    	bash$  ./${PROGNAME} update

    Command to start ${PROGNAME} and get all RSS feeds with articles that
    have been released since Aug 20, 2022:

    	bash$  ./${PROGNAME} update -d20220820

ZZEOF
}

#-------------------------------------------------------------------------------
#  SetUpdateStartDate-query the database for the latest date in the Exch table
#                     then set the range to cover from that date to the current
#                     date
#-------------------------------------------------------------------------------
SetUpdateStartDate () {
    DTSTART=$(echo "SELECT PubDt FROM Item ORDER BY PubDt DESC limit 1;" | mysql plato | grep -v PubDt | sed 's/ .*//' | sed 's/-//g')
}
#-------------------------------------------------------------------------------
#  SetDEST - determine the plato volume and set the temporary directory for
#            the RSS feed downloads. By default, it assumes there will be a
#            directory on the Plato volume named rss.  Failing that, it will
#            set the name of the directory to "./rsstmp"
#
#            Note that DEST can also be set from a runtime option -p
#-------------------------------------------------------------------------------
GetDEST () {
    echo "Trace: GetDEST:  DEST = ${DEST}"
    if [ "${DEST}x" = "x" ]; then
        # DEST=$(ls -l /Volumes | sed 's/........................................................//' | sed 's/\///' | grep -i plato)
        # LC=$(echo ${DEST} | wc -l)
        # if (( LC > 1 )); then
        #     echo "There are multiple volumes containing 'plato'"
        #     exit 1
        # fi
        # if [ "${DEST}" = "" ]; then
        #     DEST="./rsstmp"
        # else
        #     DEST="/Volumes/${DEST}/rss"
        # fi
        DEST="./rsstmp"
    fi
    mkdir -p "${DEST}"
    if (( $? != 0 )); then
        echo "Could not detect or create ${DEST}"
        exit 1
    fi
    echo "RSS temp directory is: ${DEST}"
}

#-------------------------------------------------------------------------------
#  KillWBP - kill all instances of waybackpack running on the system
#-------------------------------------------------------------------------------
KillWBP () {
    ps -ef | grep waybackpack | grep -v grep | while read line ;
    do
        a=(${(s/ /)line})
        kill "$a[2]"
        echo "killed $a{2}" >> "${RESTARTS}"
    done;
}

#-------------------------------------------------------------------------------
#  INIIALIZE...
#-------------------------------------------------------------------------------
DTSTART="20110202"  # Use this date to start from scratch -- ** don't remove this line
# DTSTART="20211223"  # Use this to set to a different start date
DOWNLOADED="completed.txt"
RESTARTS="restarted.txt"
FINISHED="finished.txt"
WLOG="w.log"

# RESTARTME="restartme.txt"
# rm -f "${DOWNLOADED}" "${RESTARTS}" "${WLOG}" "${FINISHED}" "${RESTARTME}"

echo "URLS downloaded to disk during this run:" >> "${DOWNLOADED}"
echo "Log of killed and restarted waybackpack processes" >> "${RESTARTS}"
MYSQL=$(which mysql)
if [ "x" = "${MYSQL}x" ]; then
    echo "mysql command not found. Ensure that mysql is installed an in your PATH then try again."
    exit 1
fi
MYSQL="${MYSQL} --no-defaults"

#-------------------------------------------------------------------------------
#  Handle command line args...
#-------------------------------------------------------------------------------
retryURL=""   # assume no retry for now
retryDT=""
SKIPTORETRY=0

while getopts "d:f:Hhp:u" o; do
	echo "o = ${o}"
	case "${o}" in
		h | H)
			usage
			exit 0
			;;
        u)  SetUpdateStartDate
            ;;
        f)  retryURL="${OPTARG}"
            ;;
        d)  retryDT="${OPTARG}"
            ;;
        p)  DEST="${OPTARG}"
            echo "DEST set to ${DEST}"
            ;;
		*) 	usage
			exit 1
			;;
	esac
done
shift $((OPTIND-1))

GetDEST  # do this after processing options

if [ "${retryURL}x" != "x" ]; then
    SKIPTORETRY=1
    if [ "${retryDT}x" != "x" ]; then
        DT="${retryDT}"
    else
        DT="${DTSTART}"
    fi
    echo "Skipping to ${retryURL} @ ${DT}..."
fi

#-------------------------------------------------------------------------------
#  On with it!
#-------------------------------------------------------------------------------
for url in "${urls[@]}"; do
    #---------------------------------------------------------------------------
    # If started in retry mode, skip all URLs until we hit ${retruURL}
    #---------------------------------------------------------------------------
    if [ "${SKIPTORETRY}" = "1" -a "${retryURL}" = "${url}" ]; then
        SKIPTORETRY=0
    fi
    if (( SKIPTORETRY == 0 )); then
        rm -rf "${DEST}"
        mkdir -p "${DEST}"
        if [ "${retryURL}" = "${url}" -a "${retryDT}x" != "x" ]; then
            DT="${retryDT}"
        else
            DT="${DTSTART}"
        fi

        echo "Starting WayBackPack for ${url}" >> "${WLOG}"
        date >> "${WLOG}"
        echo "waybackpack ${url} --max-retries 3 --from-date ${DT} -d ${DEST}"
        waybackpack "${url}" --max-retries 3 --from-date "${DT}" -d "${DEST}" >> "${WLOG}" 2>&1

        echo "${url} " >> ${DOWNLOADED}
        echo "Calling unpack.sh \"${url}\" \"${DEST}\"" | tee -a "${WLOG}"
        ./unpack.sh "${url}" "${DEST}"
        echo "unpack.sh completed" | tee -a "${WLOG}"
        echo "finished" >> ${DOWNLOADED}
    fi
done
echo "wbp.sh finished" > "${FINISHED}"
