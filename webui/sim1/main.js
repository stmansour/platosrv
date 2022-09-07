/*jshint esversion: 6 */

let app = {
    width: 800,
    height: 500,
    config: null,
    uid: null,
    name: null,
    imageurl: null,
    simulationCycles: 5,
    currentCycle: 0,
    initialFunding: 100.00,
    populationSize: 10,
    dtStart: "8/22/2022",       // simulation start date
    dtStop: "8/31/2022",        // simulation stop date
    dt: null,                   // current date
    loggedIn: false,            // wait til we log in before starting the simulation
    platoReqActive: false,      // is a response to a request pending?
    records: [],                // temporary result set
    GUID: 0,                    // not for production, but perfect and highly efficient for our little simulator
    fontSize: 12,               // point size for default font
    simulator: null,
    grid: null,
};


function setup() {
    let c = createCanvas(app.width, app.height);
    c.parent('theCanvas');
    setConfig();
    login();
    let d = new Date(app.dtStart);
    app.dtStart = d;
    app.dt = d;
    d = new Date(app.dtStop);
    app.dtStop = d;
    app.simulator = new Simulator();
    app.simulator.dtStart = app.dtStart;
    app.simulator.dtStop = app.dtStop;
    app.simulator.init();
    app.grid = new Grid(50,50,app.width-100,app.height-100);

    app.grid.setXRange(0,1000);  // debug only
    app.grid.xLabelDecimals = 0;
    app.grid.setYRange(80,90);   // debug only

    initUI();
}

function draw() {
    background(39,40,61);
    fill(201,204,212);
    textSize(app.fontSize);

    if (!app.loggedIn) {
        text("info = " + app.config.user,200,200);
        text("logging in...", 200, 210);
    }

    if (app.simulator.infs[0].archive.length > 0) {
        textSize(2*app.fontSize);
        text(formatTicker(app.simulator.tmpGuy.ticker), 50,25);
        textSize(1.5 * app.fontSize);
        fill(201,204,212);
        let y = 250;
        let lows = [];
        for (let i = 0; i < app.simulator.tmpGuy.archive.length; i++) {
            // let s = formatDateSlash(app.simulator.tmpGuy.archive[i].dtLow);
            // if (!app.simulator.tmpGuy.archive[i].found) {
            //     s += "  NO DATA";
            // } else {
            //     s +=  "   low = " + number_format(app.simulator.tmpGuy.archive[i].low,2) + " @ " + formatTime(app.simulator.tmpGuy.archive[i].dtLow) +
            //           "   high = " + number_format(app.simulator.tmpGuy.archive[i].high,2) + " @ " + formatTime(app.simulator.tmpGuy.archive[i].dtHigh);
            // }
            // text(s,50,y);
            // y += 15;

            // plot low...
            lows.push(app.simulator.tmpGuy.archive[i].low);
        }
        app.grid.data = [];
        app.grid.data.push(lows);
        app.grid.xTicks = lows.length;
    }
    app.grid.show();

    app.simulator.go(); // let the simulation proceed

}

function login() {
    var params = {user: app.config.user, pass: app.config.pass };
    var dat = JSON.stringify(params);
    $.post('http://localhost:8277/v1/authn/', dat, null, "json")
    .done(function(data) {
        if (data.status === "success") {
            app.uid = data.uid;
            app.name = data.Name;
            app.imageurl = data.ImageURL;
            app.loggedIn = true;
            setTDBG("toprow14","green");
            setInnerHTML("loginName","Hi, " + app.name);
            setImageSrc("userImage",app.imageurl);
            app.simulator.begin();
        } else {
            console.log("Login service returned unexpected status: " + data.status);
        }
    })
    .fail(function(/*data*/){
        console.log("Login failed");
    });
}
