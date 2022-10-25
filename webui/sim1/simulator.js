/*jshint esversion: 6 */

class Simulator {
    constructor() {
        this.dtStart = null;
        this.dtStop = null;
        this.dtPredictStart = null;
        this.dtPredictStop = null;
        this.simulationStartTime = null;
        this.investors = [];
        this.investorCount = 1;
    }

    init() {

        for (let i = 0; i < this.investorCount; i++) {
            let inv = new Investor();
            inv.makeInfluencers();
            this.investors.push( inv );
        }

        this.dtPredictStart = app.dtPredictStart;
        this.dtPredictStop = app.dtPredictStop;

        for (let i = 0; i < this.investors.length; i++) {
            this.investors[i].init();
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
