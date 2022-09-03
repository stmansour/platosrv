/*jshint esversion: 6 */

let app = {
    width: 800,
    height: 500,
    config: null,
    loggedIn: false,
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
};

let tmpGuy = new MaxExch("NZDJPY");

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
    initUI();
}

function draw() {
    background(39,40,61);
    fill(201,204,212);
    text("info = " + app.config.user,200,200);

    if (!app.loggedIn) {
        text("logging in...", 200, 220);
    } else {
        text("successfully logged in " + app.name + "  uid = " + app.uid, 200, 220);
    }

    let on = color(255,0,0);
    let off = color(120,0,0);
    let c = app.platoReqActive ? on : off;
    fill(c);
    circle(200,235,10);
    fill(201,204,212);
    text("plato server", 215, 240);

    if (tmpGuy.records.length > 0) {
        processRecords();
    }

}

function processRecords() {
    let high = 0;
    let low = Infinity;
    for (var i = 0; i < tmpGuy.records.length; i++) {
        if (tmpGuy.records[i].Close > high) {
            high = tmpGuy.records[i].High;
        }
        if (low > tmpGuy.records[i].Low) {
            low = tmpGuy.records[i].Low;
        }
    }
    fill(201,204,212);
    let s = formatTicker(tmpGuy.records[0].Ticker) + ":  low = " + low + "   high = " + high;
    text(s,200,260);
}

function startSimulation() {
    console.log("Starting simulation...");
    tmpGuy.fetch(app.dt);
}


function login() {
    var params = {user: app.config.user, pass: app.config.pass };
    var dat = JSON.stringify(params);
    $.post('http://localhost:8277/v1/authn/', dat, null, "json")
    .done(function(data) {
        if (data.status === "error") {
            console.log(data);
        }
        else if (data.status === "success") {
            app.uid = data.uid;
            app.name = data.Name;
            app.imageurl = data.ImageURL;
            app.loggedIn = true;
            setImageSrc("userImage",app.imageurl);
            startSimulation();
        } else {
            console.log("Login service returned unexpected status: " + data.status);
        }
        return;
    })
    .fail(function(/*data*/){
        console.log("Login failed");
        return;
    });
}
