/*jshint esversion: 6 */

class Investor {
    constructor(dna) {
        this.dna = (typeof dna == "undefined") ? [] : dna;
        this.influencers = [];
        this.amount = 100;      // we'll stake the investor with $100
        this.dtStart = app.simulator.dtStart;
        this.dtStop = app.simulator.dtStop;
        if (this.dna.length > 0) {
            for (var i = 0; i < this.dna.length; i++) {
                // Extract the class name and create them:
                // {
                //
                //    "ticker": "NZDJPY"},
                //    ...
                // }
            }
        } else {
            // how many influencers should we have?
            // one for now...
            this.dna.push('{"Influencer":"MaxExch","ticker":"NZDJPY"}');
        }
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
            this.influencers.push( inf );
        }
    }

    go() {
        for (var i = 0; i < this.influencers.length; i++) {
            if (!this.influencers[i].ready) {
                this.influencers[i].go();
            } else {
                // do your thing
            }
        }
    }
}
