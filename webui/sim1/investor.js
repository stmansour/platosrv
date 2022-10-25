/*jshint esversion: 6 */

class Investor {
    constructor(dna) {
        this.dna = (typeof dna == "undefined") ? [] : dna;
        this.influencers = [];
        this.amountA = 20000;   // amount of A currency (NZD)
        this.amountB = 200000;  // amount of B currency (JPY)
        this.investA = 5000;    // investment amount for A currency
        this.investB = 50000;   // investment amount for B currency
        this.dtStart = app.simulator.dtStart;
        this.dtStop = app.simulator.dtStop;
        this.dtPredictStart = app.simulator.dtPredictStart;
        this.dtPredictStop = app.simulator.dtPredictStop;
        this.dt = this.dtStart;         // date being processed
        this.allDataCollected = false;  // we'll collect all data before predicting
        this.state = 0;                 // 0 = collect data, 1 = make prediction, 2 = invest, 3 = analyze results
        this.probabilityTrigger = 0.4;  // probability that triggers investment.
        this.simulationComplete = false;
        if (this.dna.length > 0) {
            for (var i = 0; i < this.dna.length; i++) {
                // use supplied dna
            }
        } else {
            // build random dna
            // one for now...
            this.dna.push('{"Influencer":"TimeTrader","ticker":"NZDJPY"}');
        }
    }

    init() {
        this.dtPredictStart = app.simulator.dtPredictStart;
        this.dtPredictStop = app.simulator.dtPredictStop;
    }

    makeInfluencers() {
        for (let i = 0; i < this.dna.length; i++) {
            let gene = JSON.parse(this.dna[i]);
            let inf = null;
            //----------------------------------------
            // first pass... generate the influencer
            //----------------------------------------
            for (const p in gene) {
                if (p == "Influencer") {
                    let name = gene[p];
                    inf = app.registry.generate(name);
                }
            }
            //----------------------------------------
            // next pass... initialize it
            //----------------------------------------
            for (const p in gene) {
                if (p != "Influencer") {
                    inf[p] = gene[p];
                }
            }
            inf.init();
            this.influencers.push( inf );
        }
    }

    go() {
        //---------------------------------------------------
        // have all influencers gotten the data they need?
        //---------------------------------------------------
        if (this.simulationComplete) {
            return;
        }
        if (!this.allDataCollected) {
            //-------------------------------------
            //  NO: continue to collect data...
            //-------------------------------------
            let allData = 0;
            for (let i = 0; i < this.influencers.length; i++) {
                //--------------------------------------------------
                //  has this influencer collected all its data?
                //--------------------------------------------------
                if (this.influencers[i].allDataCollected) {
                    allData++;
                } else {
                    if (this.influencers[i].ready) {
                        this.influencers[i].go();
                    }
                }
            }
            if (allData == this.influencers.length) {
                this.allDataCollected = true;
            }
        } else {
            //------------------------------------------------------------
            //  YES: start asking for predictions about the dates in the
            //  the simulation range...
            //------------------------------------------------------------
            console.log("Investor starts with  A: " + this.amountA + "    B: " + this.amountB);
            let dtEnd = new Date(this.dtPredictStop);

            let predictions = [];
            let sum = 0;
            let n = this.influencers.length;
            for (let dt = new Date(this.dtPredictStart); dt <= dtEnd; dt.setDate(dt.getDate() + 1)) {
                if (dt.getDay() == 0) {
                    continue; // we're going to skip sundays
                }
                for (let i = 0; i < n; i++) {
                    let a = this.influencers[i].predict(dt);  // returns info struct
                    predictions.push(a);
                }

                //-------------------------------------------------------------------------------
                // go through the predictions and summarize preditions for the different tickers
                //-------------------------------------------------------------------------------
                let tickerDict = {};
                for (let i = 0; i < predictions.length; i++) {
                    let r = null;
                    if (predictions[i].ticker in tickerDict) {
                        r = tickerDict[predictions[i].ticker];
                    } else {
                        r = newRecommendation();
                        r.ticker = predictions[i].ticker;
                        tickerDict[r.ticker] = r;
                    }
                    r.timeLow += predictions[i].timeLow;
                    r.timeHigh += predictions[i].timeHigh;
                    r.probA += predictions[i].probA;
                    r.probB += predictions[i].probB;
                    r.count++;
                }

                //------------------------------------------------------------------
                // normalize all the predictions, then decide if we want to invest
                //------------------------------------------------------------------
                for (const [key, r] of Object.entries(tickerDict)) {
                    r.timeLow = floor((r.timeLow + 0.5) / r.count);
                    r.timeHigh = floor((r.timeHigh + 0.5) / r.count);
                    r.probA /= r.count;
                    r.probB /= r.count;

                    if (r.probA >= this.probabilityTrigger) {
                        console.log("Investor recommends buying " + r.ticker.substr(0,3) + " on " + formatDateYMD(dt) + " @ " + formatMinsToHHMM(r.timeLow) + " and selling @ " + formatMinsToHHMM(r.timeHigh));
                        this.buyA(dt,r);
                    } else if (r.probB >= this.probabilityTrigger) {
                        console.log("Investor recommends buying " + r.ticker.substr(3) + " on " + formatDateYMD(dt) + " @ " + formatMinsToHHMM(r.timeLow) + " and selling @ " + formatMinsToHHMM(r.timeHigh));
                        this.buyB(dt,r);
                    } else {
                        console.log("Investor has no recommendation for " + formatDateYMD(dt));
                    }
                }
            }  // date dt
            console.log("Investor finishes with  A: " + this.amountA + "    B: " + this.amountB);
            this.simulationComplete = true;
        }
    }

    buyA(dt,r) {
        //------------------------------------------------------
        // use currency B to buy A
        // determine the price at time low and time high...
        //------------------------------------------------------
        let d = app.cache.fetchByMinute(r.ticker,dt);
        let buyPrice = d[r.timeLow].Close;
        let sellPrice = d[r.timeHigh].Close;

        this.amountB -= this.investB;
        let buy = this.investB / buyPrice;
        let sell = buy * sellPrice;
        this.amountB = sell - buy;
    }

    buyB(dt,r) {
        //------------------------------------------------------
        // use currency A to buy B
        // determine the price at time low and time high...
        // also note that the times will be reversed.
        //------------------------------------------------------
        let d = app.cache.fetchByMinute(r.ticker,dt);
        if (typeof d == "undefined") {
            console.log("list is undefined for " + dt);
            return;
        }
        let buyPrice = d[r.timeHigh].Close;
        let sellPrice = d[r.timeLow].Close;

        this.amountA -= this.investA;
        let buy = this.investA * buyPrice;
        let sell = buy / sellPrice;
        this.amountA = sell - buy;
    }
}

function newRecommendation() {
    return {
        rateLow: 0,
        timeLow: 0,
        rateHigh: 0,
        timeHigh: 0,
        buyA: false,
        buyB: false,
        probA: 0,
        probB: 0,
        ticker: "",
        count: 0,
    };

}
