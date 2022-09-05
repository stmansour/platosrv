/*jshint esversion: 6 */

function initUI() {
    let s = formatDate(app.dtStart) + "  -  " + formatDate(app.dtStop);
    setInnerHTML("simDates",s);
    setInnerHTML("currentCycle",app.currentCycle);
    setInnerHTML("populationSize",app.populationSize);
    setTDBG("toprow14","red");
}
