/*jshint esversion: 6 */

// PopBase is the base class for populations. It contains many methods that
// are common to any population used in a Genetic Algorithm.
//

class PopBase {
    constructor() {
        p = [];         // the population of influencers
    }

    getMaxFitness() {
        let maxfit = 0;
        for (let i = 0; i < this.p.length; i++) {
            if (this.p[i].fitness > maxfit) {
                maxfit = this.p[i].fitness;
            }
        }
        return maxfit;
    }

    normalizeFitness(maxfit) {
        if (maxfit == "undefined") {
            maxfit = getMaxFitness();
        }
        for (let i = 0; i < this.p.length; i++) {
            this.p[i].fitness = this.p[i].fitness / maxfit;
        }
    }

}
