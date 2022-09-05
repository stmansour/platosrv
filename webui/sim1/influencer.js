/*jshint esversion: 6 */

// Influencer is the base class for an investment strategy. It implements a
// single strategy.
//-----------------------------------------------------------------------------

class Influencer {
    constructor(id) {
        if (typeof id == "undefined") {
            this.id == app.GUID;
            app.GUID++;
        }  else {
            this.id = id;
        }
        this.fitness = 0;                           // how well this strategy did in the simulation
        this.nrmFit = 0;                            // a normalized fitness value. It is set by a Population object.
    }

}
