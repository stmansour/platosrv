/*jshint esversion: 6 */

let app = {
    width: 800,
    height: 500,
    config: null,
    loggedIn: false,
    uid: null,
    name: null,
    imageurl: null,
};

function setup() {
    let c = createCanvas(app.width, app.height);
    c.parent('theCanvas');
    setConfig();
    login();
}

function draw() {
    background(220);
    text("info = " + app.config.user,200,200);

    if (!app.loggedIn) {
        text("logging in...", 200, 220);
    } else {
        text("successfully logged in " + app.name + "  uid = " + app.uid, 200, 220);
    }
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
