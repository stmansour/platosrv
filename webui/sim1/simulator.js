/*jshint esversion: 6 */

class Simulator {
    constructor() {
        this.dtStart = null;
        this.dtStop = null;
        this.simulationStartTime = null;
        this.infs = [];
    }

    init() {
        this.tmpGuy = new MaxExch("NZDJPY");
        this.infs.push(this.tmpGuy);
        // and push a lot more things here in the future...

        for (let i = 0; i < this.infs.length; i++) {
            this.infs[i].init();
        }

    }

    begin() {
        this.simulationStartTime = new Date();
        console.log("Starting simulation...");
        this.go();
    }

    go() {
        for (let i = 0; i < this.infs.length; i++) {
            this.infs[i].go();
        }
    }
}
