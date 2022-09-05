/*jshint esversion: 6 */

function initUI() {
    let s = formatDate(app.dtStart) + "  -  " + formatDate(app.dtStop);
    setInnerHTML("simDates",s);
    setInnerHTML("currentCycle",app.currentCycle);
    setInnerHTML("populationSize",app.populationSize);
    setTDBG("toprow14","red");
}

function formatDate(dt) {
    const m = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    return m[dt.getMonth()] + " " + dt.getDate() + ", " + dt.getFullYear();
}

// add the slash for readability in the database ticker string
function formatTicker(s) {
    if (s.length < 4) {
        return s;
    }
    return s.slice(0, 3) + "/" + s.slice(3);
}
