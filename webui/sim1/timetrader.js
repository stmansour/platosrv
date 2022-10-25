/*jshint esversion: 6 */

// Time trader is a test influencer that was written to help build the infrastructure for real
// single strategy.
// Call go() to have it collect the data it needs.  When it completes data
// collecting the baseClass variable allDataCollected will be set to true.
//-----------------------------------------------------------------------------

class TimeTrader extends Influencer {
    constructor() {
        super();
        this.ticker = "";
        this.dt  = this.dtStart;
        this.n = 3;     // number of days prior to analyze
    }

    setTicker(s) {
        this.ticker = s;
    }

    go() {
        this.getDataBG();
    }

    init() {
        super.init();
        this.ready = true;
    }

    getDataBG() {
        if (!this.ready) {
            return -1;  // can't do it now
        }
        this.ready = false;
        if (this.dt <= this.dtStop) {
            this.fetchExchData(this.dt);
        }

        return 0;
    }

    fetchExchData(dt) {
        let s = formatDateSlash(dt);
        var params = {
            cmd: "get",
            limit: 1440,
            tickers: [this.ticker],
            dt: s,
        };
        super.fetchExchData(params,'http://localhost:8277/v1/exch/');
    }

    predict(dt) {
        //--------------------------------------------------------------
        // analyze the data from this.n days back. predict next day
        //--------------------------------------------------------------
        let dt1 = dt;
        let dp = [];
        dt1.setDate(dt.getDate() - this.n);
        for (let x = 0; x < this.n; x++) {
            let d = app.cache.fetchByMinute(this.ticker,dt1);
            let dmin = Infinity;
            let minmins = 0;
            let dmax = 0;
            let maxmins = 0;
            if (typeof d != "undefined") {
                for (let i = 0; i < d.length; i++) {
                    if (d[i].Low < dmin) {
                        dmin = d[i].Low;
                        minmins = i;
                    }
                    if (d[i].High > dmax) {
                        dmax = d[i].High;
                        maxmins = i;
                    }
                }
            }
            dp.push( {
                low: dmin,
                lowt: minmins,
                high: dmax,
                hight: maxmins,
            });
            dt1.setDate(dt.getDate()+1);
        }

        let plow = 0;
        let tlow = 0;
        let phigh = 0;
        let thigh = 0;
        let proba = 0;
        let probb = 0;
        let dxp = 1/(dp.length+2);
        for (let i = 0; i < dp.length; i++) {
            plow += dp[i].low;
            tlow += dp[i].lowt;
            phigh += dp[i].high;
            thigh = dp[i].hight;
            if (tlow < thigh) {
                proba += dxp;
            } else if (tlow > thigh) {
                probb += dxp;
            }
        }
        let r = {
            rateLow: plow / this.n,
            timeLow: tlow / this.n,
            rateHigh: phigh / this.n,
            timeHigh: thigh / this.n,
            buyA: false,
            buyB: false,
            probA: proba,
            probB: probb,
            ticker: this.ticker,
        };
        r.buyA = r.timeLow < r.timeHigh;
        r.buyB = r.timeLow > r.timeHigh;
        // console.log('plow = ' + plow.toFixed(2) + '  @ ' + formatMinsToHrsMins(tlow) );
        // console.log('phigh = ' + phigh.toFixed(2) + '  @ ' + formatMinsToHrsMins(thigh) );
        return r;
    }
}
