/*jshint esversion: 6 */

// MaxExch is an influencer that simply finds the maximum value for its
// currency echange rate for a specific day.
class MaxExch extends Influencer {
    constructor(ticker) {
        super();                // generates the id
        this.ticker = ticker;
        this.records = [];      // place to store exch data used
        this.predictions = 0;
        this.correctPredictions = 0;
        this.archive = [];
        this.dt = null;         // what date am I currently processing
        this.dataCollected = false;
    }

    init() {
        this.dt = app.simulator.dtStart;
    }

    go() {
        if (!this.dataCollected) {
            if (this.records.length > 0) {
                let v = this.processRecords();
                let info = {
                    dt: this.dt,
                    low: v[0],
                    high: v[1],
                };
                this.archive.push(info);
                this.records = [];
                this.dt.setDate(this.dt.getDate() + 1); // move to the next day
                if (this.dt > app.simulator.dtStop) {
                    this.dataCollected = true;
                }
            } else {
                this.fetchExchData([this.ticker],this.dt);
            }
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
        var dat = JSON.stringify(params);
        $.post('http://localhost:8277/v1/exch/', dat, null, "json")
        .done( (data) => {
            if (data.total == 0) {
                console.log("NO RECORDS FOUND for " +s);
                return;
            }
            if (data.status === "success") {
                console.log("Successfully retrieved " + data.records.length + " records from plato server.");
                this.records = data.records;
            } else {
                console.log("Login service returned unexpected status: " + data.status);
            }
        })
        .fail(function(/*data*/){
            console.log("Request failed");
        });
    }

    // print out relevant information about this MaxExch object
    // show(id) {
    // }

    // processRecords returns an array of 2 numbers: the low and the high of the day
    //------------------------------------------------------------------------------
    processRecords() {
        let high = 0;
        let low = Infinity;
        for (var i = 0; i < this.records.length; i++) {
            if (this.records[i].Close > high) {
                high = this.records[i].High;
            }
            if (low > this.records[i].Low) {
                low = this.records[i].Low;
            }
        }
        return [low,high];
    }

}
