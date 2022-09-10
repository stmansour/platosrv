/*jshint esversion: 6 */

class Simulator {
    constructor() {
        this.dtStart = null;
        this.dtStop = null;
        this.simulationStartTime = null;
        this.investors = [];
        this.investorCount = 1;
    }

    init() {

        // this.infs.push(new MaxExch("NZDJPY"));
        // and push a lot more things here in the future...

        for (let i = 0; i < this.investorCount; i++) {
            let inv = new Investor();
            inv.makeInfluencers();
            this.investors.push( inv );
        }

    }

    begin() {
        this.simulationStartTime = new Date();
        console.log("Starting simulation...");
        this.go();
    }

    go() {
        for (let i = 0; i < this.investors.length; i++) {
            this.investors[i].go();
        }
    }
}
