/*jshint esversion: 6 */

//
//  Data Structures
//     byDay[ TICKER + "YYYY-MM-DD"] -> record = { high: 0, dtHigh: null, low: Infinity, dtLow: null, };
//     byMinute[ TICKER + "YYYY-MM-DD"] -> [0 - 1439] = record { Close : 0.6985 Dt : "2022-08-01 00:01:00 UTC" High : 0.6985 Low : 0.6984 Open : 0.6984 }
//
//
//

class ExchCache {
    constructor() {
        this.byDayDict = {};
        this.byMinuteDict = {};
    }

    store(recs) {
        let key = recs[0].Ticker + recs[0].Dt.substr(0,10);
        let byMinute = [];
        let r = {
            high: 0,
            dtHigh: null,
            low: Infinity,
            dtLow: null,
        };
        for (var i = 0; i < recs.length; i++) {
            byMinute.push( {
                Open:   recs[i].Open,
                High:   recs[i].High,
                Low:    recs[i].Low,
                Close:  recs[i].Close,
            } );
            if (recs[i].High > r.high) {
                r.high = recs[i].High;
                r.dtHigh = recs[i].Dt;
            }
            if (r.low > recs[i].Low) {
                r.low = recs[i].Low;
                r.dtLow = recs[i].Dt;
            }
        }
        this.byDayDict[key] = r;
        this.byMinuteDict[key] = byMinute;
    }

    fetchByMinute(ticker,dt) {
        let key = ticker + formatDateYMD(dt);
        let a = this.byMinuteDict[key];
        return a;
    }
}
