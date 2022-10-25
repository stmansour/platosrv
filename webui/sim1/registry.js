/*jshint esversion: 6 */

class Registry {
    constructor() {
        this.reg = [
            "MaxExch",
        ];
    }

    generate(c) {
        switch (c) {
            case "MaxExch": return new MaxExch();
            case "TimeTrader": return new TimeTrader();
            default:
                console.log("unrecognized class name: " + c);
        }
    }
}
