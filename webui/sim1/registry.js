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
            default:
                console.log("unrecognized class name: " + c);
        }
    }
}
