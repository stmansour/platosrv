/*jshint esversion: 6 */

class Grid {
    constructor(x,y,w,h) {
        this.x = x;
        this.y = y;
        this.width = w;
        this.height = h;
        this.xTicks = 10;
        this.yTicks = 10;
        this.xLabels = [];      // could be dates
        this.showXLabels = true;
        this.showYLabels = true;
        this.xMin = 0;
        this.xMax = 0;
        this.yMin = 0;
        this.yMax = 0;
        this.xLabelDecimals = 2;
        this.yLabelDecimals = 2;
        this.data = []; // an array of arrays with user data.  Each elemnt is a single data set
        this.dataColors = [
            [210,193,134],
            [201,128,94],
            [127,164,185],
            [169,205,103],
            [132,128,192],
        ];
    }

    show () {
        push();
        this.drawBasicGrid();
        this.drawLabels();
        this.drawData();
        pop();
    }

    drawData() {
        if (this.data.length == 0) {
            return;
        }
        let y;
        let offset = 0;
        let width = 0;
        for (let i = 0; i < this.data.length; i++) {
            let dx = this.width/this.data[i].length;
            width = dx / this.data.length;  // width to accommodate this many data sets
            offset = width * i;             // each dataset gets an offset
            fill(this.dataColors[i][0],this.dataColors[i][1],this.dataColors[i][2]);
            for (let j = 0; j < this.data[i].length; j++) {
                if (this.data[i][j] == 0) {
                    continue;
                }
                y = map(this.data[i][j],
                        this.yMin, this.yMax,
                        this.y + this.height, this.y);
                rect(offset + this.x + j * dx, y,  width - 1,  this.y + this.height - y);
            }
        }
    }

    drawLabels() {
        let useLabels = (this.xLabels.length > 0);
        let labelsAreDates =(useLabels && typeof this.xLabels[0] == 'object' && this.xLabels[0].constructor.name == 'Date');
        if (labelsAreDates && this.xTicks > 10) {
            this.xTicks = 10; // they overlap if there are any more than 10.
        }
        let dx = this.width/this.xTicks;
        let dy = this.height/this.yTicks;
        let x1, y1, xldx, yldy, xl, yl;
        noStroke();
        fill(201,204,212);
        textSize(12);

        if (this.xMin == this.xMax ) {
            console.log("X Range has not been set");
            return;
        } else {
            textAlign(CENTER, CENTER);
            y1 = this.y + this.height + 15;
            xldx = (this.xMax - this.xMin)/this.xTicks;
            let inc = 1;
            let iMax = this.xTicks;
            if (useLabels) {
                inc = this.xLabels.length/this.xTicks;
                if (this.xLabels.length > this.xTicks) {
                    iMax = this.xLabels.length;
                }
            }
            let k = 0;
            for (let i = 0; i <= iMax; i += inc) {
                if (useLabels) {
                    x1 = this.x + k * dx;
                    k++;
                    let idx = floor(i + 0.5);
                    if (idx < this.xLabels.length) {
                        let t = this.xLabels[idx];
                        if (labelsAreDates) {
                            let s = formatDateSlash(t);
                            t = s;
                        }
                        text( t, x1, y1);
                    }
                } else {
                    xl = this.xMin + i * xldx;
                    text( number_format(xl,this.xLabelDecimals), x1,y1);
                }
            }
        }

        if (this.yMin == this.yMax ) {
            console.log("Y Range has not been set");
            return;
        } else {
            textAlign(CENTER, CENTER);
            x1 = this.x - textWidth(''+this.yMax) - 8;
            y1 = this.y + this.height;
            yldy = (this.yMax - this.yMin)/this.yTicks;
            for (let i = 0; i <= this.yTicks; i++) {
                y1 = this.y + this.height - i * dy;
                yl = this.yMin + i * yldy;
                text( number_format(yl,this.yLabelDecimals), x1,y1);
            }
        }
    }

    setXRange(a,b) {
        this.xMin = a;
        this.xMax = b;
    }

    setYRange(a,b) {
        this.yMin = a;
        this.yMax = b;
    }

    drawBasicGrid() {
        stroke(66,67,87);
        strokeWeight(3);
        noFill();
        rect(this.x,this.y, this.width, this.height);
        let dx = this.width/this.xTicks;
        let dy = this.height/this.yTicks;
        let x1 = this.x + this.width;
        let y1 = this.y + this.height;
        stroke(255);
        line(this.x,this.y+this.height, this.x + this.width, this.y + this.height);
        strokeWeight(1);
        for (let i = 1; i < this.xTicks; i++) {
            let x = this.x + dx*i;
            stroke(66,67,87);
            line(x, this.y, x, y1); // the line
            stroke(255);
            line(x, y1, x, y1+6);   // the tick
        }
        for (let i = 1; i < this.yTicks; i++) {
            let y = this.y + dy*i;
            stroke(66,67,87);
            line(this.x, y, x1, y );    // the line
            stroke(255);
            line(this.x-4, y, this.x, y );  // the tick
        }
    }
}
