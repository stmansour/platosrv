/*jshint esversion: 6 */

function setImageSrc(id,val) {
    let x = document.getElementById(id);
    if (x != null) {
        x.src = val;
    }
}

//  Jan 25, 2022
//-----------------------------------------------------------------------------
function formatDate(dt) {
    const m = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    return m[dt.getMonth()] + " " + dt.getDate() + ", " + dt.getFullYear();
}

//  01/25/2022  or with sep:  01-25-2022
//-----------------------------------------------------------------------------
function formatDateSlash(dt,sep) {
    if (typeof sep == "undefined") {
        sep = "/";
    }
    if (typeof dt == "string") {
        let x = new Date(dt);
        dt = x;
    }
    return (""+(1 +dt.getMonth())).padStart(2,'0') + sep + ("" + dt.getDate()).padStart(2,'0') + sep + dt.getFullYear();
}

//  2022-01-25 , or use sep
//-----------------------------------------------------------------------------
function formatDateYMD(dt,sep) {
    if (typeof sep == "undefined") {
        sep = "-";
    }
    if (typeof dt == "string") {
        let x = new Date(dt);
        dt = x;
    }
    return '' + dt.getFullYear() + sep + (''+(1 +dt.getMonth())).padStart(2,'0') + sep + ("" + dt.getDate()).padStart(2,'0');
}

function formatMinsToHrsMins(x) {
    return '' + floor((x / 60)) + ":" + ('' + floor((x % 60) + 0.5)).padStart(2,'0');
}

// function formatDateDash(dt,sep) {
//     if (typeof sep == "undefined") {
//         sep = "-";
//     }
//     if (typeof dt == "string") {
//         let x = new Date(dt);
//         dt = x;
//     }
//     return ("" + dt.getFullYear() + sep + (1 +dt.getMonth())).padStart(2,'0') + sep + ("" + dt.getDate()).padStart(2,'0');
// }

function formatTime(dt) {
    if (typeof dt == "string") {
        let x = new Date(dt);
        dt = x;
    }
    let h = "" + dt.getHours();
    let m = "" + dt.getMinutes();
    return h.padStart(2,'0') + ":" + m.padStart(2,'0');
}

function formatMinsToHHMM(m) {
    let hrs = '' + floor(m/60);
    let mins = '' + (m % 60);
    return "" + hrs.padStart(2,'0') + ":" + mins.padStart(2,'0');
}

function formatDT(dt) {
    return formatDateSlash(dt) + " " + formatTime(dt);
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

// add the slash for readability in the database ticker string
function formatTicker(s) {
    if (s.length < 4) {
        return s;
    }
    return s.slice(0, 3) + "/" + s.slice(3);
}

function number_format(number, decimals, dec_point, thousands_sep) {
    // *     example 1: number_format(1234.56);
    // *     returns 1: '1,235'
    // *     example 2: number_format(1234.56, 2, ',', ' ');
    // *     returns 2: '1 234,56'
    // *     example 3: number_format(1234.5678, 2, '.', '');
    // *     returns 3: '1234.57'
    // *     example 4: number_format(67, 2, ',', '.');
    // *     returns 4: '67,00'
    // *     example 5: number_format(1000);
    // *     returns 5: '1,000'
    // *     example 6: number_format(67.311, 2);
    // *     returns 6: '67.31'
    // *     example 7: number_format(1000.55, 1);
    // *     returns 7: '1,000.6'
    // *     example 8: number_format(67000, 5, ',', '.');
    // *     returns 8: '67.000,00000'
    // *     example 9: number_format(0.9, 0);
    // *     returns 9: '1'
    // *    example 10: number_format('1.20', 2);
    // *    returns 10: '1.20'
    // *    example 11: number_format('1.20', 4);
    // *    returns 11: '1.2000'
    // *    example 12: number_format('1.2000', 3);
    // *    returns 12: '1.200'
    var n = !isFinite(+number) ? 0 : +number,
        prec = !isFinite(+decimals) ? 0 : Math.abs(decimals),
        sep = (typeof thousands_sep === 'undefined') ? ',' : thousands_sep,
        dec = (typeof dec_point === 'undefined') ? '.' : dec_point,
        toFixedFix = function (n, prec) {
            // Fix for IE parseFloat(0.55).toFixed(0) = 0;
            var k = Math.pow(10, prec);
            return Math.round(n * k) / k;
        },
        s = (prec ? toFixedFix(n, prec) : Math.round(n)).toString().split('.');
    if (s[0].length > 3) {
        s[0] = s[0].replace(/\B(?=(?:\d{3})+(?!\d))/g, sep);
    }
    if ((s[1] || '').length < prec) {
        s[1] = s[1] || '';
        s[1] += new Array(prec - s[1].length + 1).join('0');
    }
    return s.join(dec);
}
