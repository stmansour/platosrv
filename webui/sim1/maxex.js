/*jshint esversion: 6 */

// MaxExch is an influencer that simply finds the maximum value for its
// currency echange rate for a specific day.
class MaxExch extends Influencer {
    constructor(ticker) {
        super();                // generates the id
        this.ticker = ticker;
        this.records = [];      // place to store exch data used
        this.dtRecords = null,  // date associated with this.records (which may be empty, that's why we need the date)
        this.predictions = 0;
        this.correctPredictions = 0;
        this.archive = [];
        this.dt = null;         // what date am I currently processing
        this.dataCollected = false;
        this.responseProcessed = true;
        this.awaitingServerReply = false;
    }

    init() {
        this.dt = app.simulator.dtStart;
    }

    go() {
        if (this.dataCollected || this.awaitingServerReply) {
            return;
        }
        if (!this.responseProcessed) {
            //-------------------------------------
            // Process the server response now...
            //-------------------------------------
            let info = this.getNewInfo();
            if (this.records.length >= 0) {
                if (this.records.length > 1) { // we added one record when the result was 0 records.  The expected value is 1440
                    info.found = true;
                    let r = this.processRecords();
                    info.low = r.low;
                    info.high = r.high;
                    info.dtLow = r.dtLow;
                    info.dtHigh = r.dtHigh;
                } else {
                    info = this.records[0];
                }
                this.archive.push(info);
                this.records = [];
            }
            this.responseProcessed = true;
            //---------------------------------------
            // Set up future requests...
            //---------------------------------------
            this.dt.setDate(this.dt.getDate() + 1); // move to the next day
            if (this.dt > app.simulator.dtStop) {
                this.dataCollected = true;
            }
        } else {
            //----------------------------------------------
            // We have not collected all the data we need.
            // Fetch the next bit.
            //----------------------------------------------
            this.fetchExchData([this.ticker],this.dt);
        }
    }

    fetchExchData(tickers,dt) {
        let s = formatDateSlash(dt);
        var params = {
            cmd: "get",
            limit: 1440,
            tickers: tickers,
            dt: s,
        };
        let saveDt = s;
        var dat = JSON.stringify(params);
        this.awaitingServerReply = true;
        this.responseProcessed = false;
        $.post('http://localhost:8277/v1/exch/', dat, null, "json")
        .done( (data) => {
            this.awaitingServerReply = false;
            if (data.total == 0) {
                console.log("NO RECORDS FOUND for " +s);
                let info = this.getNewInfo();
                info.dtLow = saveDt;
                info.dtHigh = saveDt;
                this.records.push(info);  // save this so we have the date
                return;
            }
            if (data.status === "success") {
                console.log("Successfully retrieved " + data.records.length + " records from plato server.");
                this.records = data.records;
                this.dtRecords = saveDt;
            } else {
                console.log("Login service returned unexpected status: " + data.status);
            }
        })
        .fail(function(/*data*/){
            this.awaitingServerReply = false;
            console.log("Request failed");
        });
    }

    getNewInfo() {
        return {
            dtLow: null,
            dtHigh: null,
            low: 0,
            high: 0,
            found: false, // assume nothing found on this date
        };
    }

    // processRecords returns an array of 2 numbers: the low and the high of the day
    //------------------------------------------------------------------------------
    processRecords() {
        let r = {
            high: 0,
            dtHigh: null,
            low: Infinity,
            dtLow: null,
        };
        for (var i = 0; i < this.records.length; i++) {
            if (this.records[i].Close > r.high) {
                r.high = this.records[i].High;
                r.dtHigh = this.records[i].Dt;
            }
            if (r.low > this.records[i].Low) {
                r.low = this.records[i].Low;
                r.dtLow = this.records[i].Dt;
            }
        }
        return r;
    }

}
