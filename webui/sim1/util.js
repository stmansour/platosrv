/*jshint esversion: 6 */

function setImageSrc(id,val) {
    let x = document.getElementById(id);
    if (x != null) {
        x.src = val;
    }
}

function formatDateSlash(dt) {
    return (1 +dt.getMonth()) + "/" + dt.getDate() + "/" + dt.getFullYear();
}

function setInnerHTML(id,s) {
    let el = document.getElementById(id);
    if (el != null) {
        el.innerHTML = s;
    }
}

function setTDBG(id,s) {
    let el = document.getElementById(id);
    if (el != null) {
        el.style.backgroundColor = s;
    }
}
