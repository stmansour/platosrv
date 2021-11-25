/*global
    $, window, w2ui, document, console,
    parseInt, w2ui, app, getitemInitRecord, updateBUDFormList, renderMoney,
*/

"use strict";

var itemMeta = {
    PubDt: "",
};

function itemInit() {
}

function getItemDateId() {
    var dt = document.getElementById("ItemDt");
    if (typeof dt == "undefined") {
        return null;
    }
    return dt;
}

function updateItemPubDate() {
    var dt = getItemDateId();
    if (dt == null) {
        return;
    }
    var ds = dt.value;
    var dd = new Date(ds);
    if (dd.getFullYear() < 2011) {
        return;
    }
    itemMeta.PubDt = ds;
    w2ui.itemGrid.postData = {
        PubDt: itemMeta.PubDt,
    };
    w2ui.itemGrid.reload();
}

window.getitemInitRecord = function(){
    return {
        recid: 0,
        IID: 0,
        PubDt: new Date(),
        Title: 0,
        Description: 0,
        High: 0,
        Low: 0,
        Close: 0,
    };
};

window.buildItemMethodElements = function () {
    //------------------------------------------------------------------------
    //          Deposit Methods Grid
    //------------------------------------------------------------------------
    $().w2grid({
        name: 'itemGrid',
        url: '/v1/item',
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
            toolbarSearch  : true,
            toolbarInput   : true,
            searchAll      : false,
            toolbarReload  : true,
            toolbarColumns : true,
        },
        columns: [
            {field: 'recid',      hidden: true,  caption: 'recid',       size: '40px', sortable: true},
            {field: 'XID',        hidden: true,  caption: 'XID',         size: '40px', sortable: true},
            {field: 'PuDt',       hidden: false, caption: 'PubDt',       size: '150px', sotrable: true,
                render: function(record) {
                    return record.PubDt.substring(0,19);
                },
            },
            {field: 'Link',       hidden: false, caption: '',        size: '20px', sotrable: true,
                render: function(record) {
                    if (typeof record != "undefined") {
                        var s = '<a href="' + record.Link + '" alt="' + record.Link + '" target="_blank"><i class="fa fa-external-link" aria-hidden="true"></i></a>';
                        return s;
                    }
                },
            },
            {field: 'Title',      hidden: false, caption: 'Title',       size: '500px', sotrable: true },
            {field: 'Description',hidden: false, caption: 'Description', size: '100px', sotrable: true },
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
                Dt: itemMeta.PubDt,
            };
            var dt = getItemDateId();
            if (dt != null) {
                dt.value = itemMeta.PubDt;
            }
        },
        onLoad: function(event) {
            event.onComplete = function() {
                var txt = document.getElementById("ItemDt");
                txt.onchange = updateItemPubDate;
            };
        },
    });
    var html1 = '<div class="w2ui-field" style="padding: 0px 5px;">Item Date: <input type="date" id="ItemDt" name="ItemDt"></div>';
    w2ui.itemGrid.toolbar.add( [
            { type: 'html',  id: 'D1',  html: function(){
                    var html1 = '<div class="w2ui-field" style="padding: 0px 5px;">Item Date: <input type="date" id="ItemDt" value="';
                    if (itemMeta.PubDt.length > 0) {
                        var dt = new Date(itemMeta.PubDt);
                        html1 += dt.toISOString().substring(0,10);
                    }
                    html1 += '"></div>';
                    return html1;
                },
            },
            { type: 'button', id: 'ClearItemDt', caption: 'Clear Date'},
        ],
    );
    w2ui.itemGrid.toolbar.onClick = function(event) {
        console.log('event.target = ' + event.target);
        switch (event.target) {
            case 'ClearItemDt':
                var id = getItemDateId();
                if (null != id) {
                    id.value = "";
                    updateItemPubDate();
                }
                break;
            default:
        }
    };
};
