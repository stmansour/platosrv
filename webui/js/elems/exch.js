/*global
    $, window, w2ui, document, console,
    parseInt, w2ui, app, getExchInitRecord, updateBUDFormList, renderMoney,
*/

"use strict";

var exchMeta = {
    Dt: "",
    Tickers: ["AUDUSD","NZDJPY","GBPAUD"],
};

function exchInit() {
    var dt = new Date();
    dt.setTime( dt.getTime() - 1*24*60*60*1000); // days * hours/day * mins/hour * sec/min * msec/sec
    exchMeta.Dt = dt.toISOString().substring(0,10);
}

function getExchDtId() {
    var dt = document.getElementById("exchDt");
    if (typeof dt == "undefined") {
        return null;
    }
    return dt;
}

function updateExchDate() {
    var dt = getExchDtId();
    if (dt == null) {
        return;
    }
    if (dt.value == "") {
        return;
    }
    var ds = dt.value;
    var dd = new Date(ds);
    if (dd.getFullYear() < 2011) {
        return;
    }
    exchMeta.Dt = dt.value;
    w2ui.exchGrid.postData = {
        Dt: exchMeta.Dt,
        Tickers: exchMeta.Tickers,
    };
    w2ui.exchGrid.reload();
}

window.getExchInitRecord = function(){
    return {
        recid: 0,
        XID: 0,
        Dt: new Date(),
        Ticker: 0,
        Open: 0,
        High: 0,
        Low: 0,
        Close: 0,
    };
};

window.buildExchMethodElements = function () {
    //------------------------------------------------------------------------
    //          Deposit Methods Grid
    //------------------------------------------------------------------------
    $().w2grid({
        name: 'exchGrid',
        url: '/v1/exch',
        multiSelect: false,
        limit: 360,
        show: {
            toolbar        : true,
            footer         : true,
            toolbarAdd     : false,   // indicates if toolbar add new button is visible
            toolbarDelete  : false,   // indicates if toolbar delete button is visible
            toolbarSave    : false,   // indicates if toolbar save button is visible
            selectColumn   : false,
            expandColumn   : false,
            toolbarEdit    : false,
            toolbarSearch  : false,
            toolbarInput   : false,
            searchAll      : false,
            toolbarReload  : true,
            toolbarColumns : true,
        },
        columns: [
            {field: 'recid',      hidden: true,  caption: 'recid',       size: '40px', sortable: true},
            {field: 'XID',        hidden: true,  caption: 'XID',         size: '40px', sortable: true},
            {field: 'Dt',         hidden: false, caption: 'Dt',          size: '150px', sotrable: true,
                render: function(record) {
                    return record.Dt.substring(0,19);
                },
            },
            {field: 'Ticker',     hidden: false, caption: 'Ticker',      size: '70px', sotrable: true, style: 'text-align: center'},
            {field: 'Open',       hidden: false, caption: 'Open',        size: '50px', sotrable: true, style: 'text-align: right',
            render: function(record) {
                    return renderMoney(record.Open,2);
                },
            },
            {field: 'High',       hidden: false, caption: 'High',        size: '50px', sotrable: true, style: 'text-align: right',
            render: function(record) {
                    return renderMoney(record.Open,2);
                },
            },
            {field: 'Low',        hidden: false, caption: 'Low',         size: '50px', sotrable: true, style: 'text-align: right',
            render: function(record) {
                    return renderMoney(record.Open,2);
                },
            },
            {field: 'Close',      hidden: false, caption: 'Close',       size: '50px', sotrable: true, style: 'text-align: right',
            render: function(record) {
                    return renderMoney(record.Open,2);
                },
            },
            {field: 'LastModTime',hidden: true,  caption: 'LastModTime', size: '20%',  sortable: true},
            {field: 'LastModBy',  hidden: true,  caption: 'LastModBy',   size: '150px',sortable: true},
            {field: 'CreateTS',   hidden: true,  caption: 'CreateTS',    size: '20%',  sortable: true},
            {field: 'CreateBy',   hidden: true,  caption: 'CreateBy',    size: '150px',sortable: true},
        ],
        onRefresh: function(event) {
            event.onComplete = function() {
                var sel_recid = parseInt(this.last.sel_recid);
                if (app.active_grid == this.name && sel_recid > -1) {
                    if (app.new_form_rec) {
                        this.selectNone();
                    }
                    else{
                        this.select(app.last.grid_sel_recid);
                    }
                }
            };
        },
        onRequest: function(event) {
            this.postData = {
                Dt: exchMeta.Dt,
                Tickers: exchMeta.Tickers,
            };
            var dt = getExchDtId();
            if (dt != null) {
                dt.value = exchMeta.Dt;
            }
        },
        onLoad: function(event) {
            exchMeta.Tickers = this.toolbar.get("tickers").selected;
            w2ui.exchGrid.postData = {
                Dt: exchMeta.Dt,
                Tickers: exchMeta.Tickers,
            };
            event.onComplete = function() {
                var txt = document.getElementById("exchDt");
                txt.onchange = updateExchDate;

            };
        },
    });
    var html1 = '<div class="w2ui-field" style="padding: 0px 5px;">Exchange Date: <input type="date" id="exchDt" name="exchDt"></div>';
    w2ui.exchGrid.toolbar.add( [
        { type: 'html',  id: 'D1',  html: function(){
                var html1 = '<div class="w2ui-field" style="padding: 0px 5px;">Exchange Date: <input type="date" id="exchDt" value="';
                if (exchMeta.Dt == "") {
                    exchInit();
                }
                var dt = new Date(exchMeta.Dt);
                html1 += dt.toISOString().substring(0,10);
                html1 += '"></div>';
                return html1;
            },
        },
        { type: 'break' },
        { type: 'menu-check', id: 'tickers', text: 'Tickers', icon: 'fa fa-exchange',
            selected: ['AUDUSD', 'NZDJPY', 'GBPAUD'],
            onRefresh: function (event) {
                event.item.count = event.item.selected.length;
            },
            items: [
                { id: 'AUDUSD', text: 'AUDUSD', icon: 'fa fa-usd' },
                { id: 'NZDJPY', text: 'NZDJPY', icon: 'fa fa-yen' },
                { id: 'GBPAUD', text: 'GBPAUD', icon: 'fa fa-gbp' },
            ],
        },
    ]);

};
