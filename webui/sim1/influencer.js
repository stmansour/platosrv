/*jshint esversion: 6 */

// Influencer is the base class for an investment strategy. It implements a
// single strategy.
//-----------------------------------------------------------------------------

class Influencer {
    constructor(id) {
        if (typeof id == "undefined") {
            this.id = app.GUID;
            app.GUID++;
        }  else {
            this.id = id;
        }
        this.fitness = 0;               // how well this strategy did in the simulation
        this.nrmFit = 0;                // a normalized fitness value. It is set by a Population object.
        this.ready = false;             // true when the influencer is available to execute another command
        this.dtStart = app.simulator.dtStart;
        this.dtStop = app.simulator.dtStop;
        this.dt = null;                 // date for which we're collecting data
        this.allDataCollected = false;  // true when all the data needed by this influencer has been collected and cached
    }

    init() {
        // nothing yet.  Subclasses can override with any initialization needs they have.
    }

    fetchExchData(params,url) {
        let s = params.dt;
        let saveDt = s;
        var dat = JSON.stringify(params);
        this.awaitingServerReply = true;
        this.responseProcessed = false;
        $.post(url, dat, null, "json")
        .done( (data) => {
            this.awaitingServerReply = false;
            if (data.total == 0) {
                console.log("NO RECORDS FOUND for " +s);
            } else if (data.status === "success") {
                app.cache.store(data.records);
            } else {
                console.log("Login service returned unexpected status: " + data.status);
            }
            this.dt.setDate(this.dt.getDate() + 1); // move to the next day
            if (this.dt > this.dtStop) {
                this.allDataCollected = true;   // mark that we have all the data we need
            }
            this.ready = true;  // ready for next command
        })
        .fail(function(/*data*/){
            this.awaitingServerReply = false;
            console.log("Request failed");
        });
    }

    predict() {
        // must be implemented by the subclass
        console.log("Influencer.predict was called. This method must be implemented by the subclass.");
    }

}
