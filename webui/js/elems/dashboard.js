
/*global
    w2ui, app, $, console, document,
    setInnerHTML, number_format, monthName,
*/

"use strict";

function fmtDate(d) {
    return monthName(d.getMonth(),3) + " " + d.getDate() + ", " + d.getFullYear();
}

function dateRange(d1,d2) {
    var dt1 = new Date(d1);
    var dt2 = new Date(d2);
    return fmtDate(dt1) + " - " + fmtDate(dt2);
}

function getDashboard() {

    var params = {
        cmd: "get",
    };
    var dat = JSON.stringify(params);
    var url = '/v1/dashboard/';

    $.post(url, dat, null, "json")
    .done(function(data) {
        if (typeof data == 'string') {  // it's weird, a successful data add gets parsed as an object, an error message does not
            var msg = JSON.parse(data);
            console.log('Response to dashboard: ' + msg.status);
            return;
        }
        if (data.status == 'success') {
            var rec = data.record;
            setInnerHTML('pTableName0',    '' + data.record.Tables[0].Name);
            setInnerHTML('pTableCount0',   '' + number_format(data.record.ExchCount,0));
            setInnerHTML('pTableDtRange0', dateRange(rec.ExchDtMin,rec.ExchDtMax));
            setInnerHTML('pTableSize0',    '' + number_format(data.record.Tables[0].Size,0) + " bytes   (" + number_format(data.record.Tables[0].Size/1024/1024,2) + "MB) ");
            setInnerHTML('pTableName1',    '' + data.record.Tables[1].Name);
            setInnerHTML('pTableCount1',   '' + number_format(data.record.ItemCount,0));
            setInnerHTML('pTableDtRange1', dateRange(rec.ItemDtMin,rec.ItemDtMax));
            setInnerHTML('pTableSize1',    '' + number_format(data.record.Tables[1].Size,0) + " bytes   (" + number_format(data.record.Tables[1].Size/1024/1024,2) + "MB) ");
        } else {
            console.log('data.status = ' + data.status);
        }
    })
    .fail(function(data) {
        console.log('data = ' + data);
    });
}
