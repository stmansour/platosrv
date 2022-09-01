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
    initUI();
}

function draw() {
    background(220);
    fill(0);
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
    text("plato server", 215, 240);

    if (app.records.length > 0) {
        processRecords();
    }

}

function processRecords() {
    let high = 0;
    let low = Infinity;
    for (var i = 0; i < app.records.length; i++) {
        if (app.records[i].Close > high) {
            high = app.records[i].High;
        }
        if (low > app.records[i].Low) {
            low = app.records[i].Low;
        }
    }
    fill(0);
    let s = app.records[0].Ticker + ":  low = " + low + "   high = " + high;
    text(s,200,260);
}

function startSimulation() {
    fetchExchData(["NZDJPY"],app.dt);
    console.log("We did it!");
}

function fetchExchData(tickers,dt) {
    let s = formatDateSlash(dt);
    var params = {
        cmd: "get",
        limit: 1440,
        Tickers: tickers,
        Dt: s,
    };
    var dat = JSON.stringify(params);
    $.post('http://localhost:8277/v1/exch/', dat, null, "json")
    .done(function(data) {
        if (data.status === "error") {
            console.log(data);
        }
        else if (data.status === "success") {
            console.log("Success!");
            app.records = data.records;
        } else {
            console.log("Login service returned unexpected status: " + data.status);
        }
        return;
    })
    .fail(function(/*data*/){
        console.log("Request failed");
        return;
    });
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
