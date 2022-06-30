var ctx = "";
var getCodeitemListByPid = ctx + "/system/codeItem/getListByPid";

var LODOP; //声明为全局变量
function printPreview(printAreaId) {
    createOneFormPage(printAreaId);
    LODOP.SET_PRINT_PAGESIZE(1, 0, 0, "A4");
    LODOP.PREVIEW();
}

function createOneFormPage(printAreaId) {
    LODOP = getLodop();
    var strStyleCSS = "<link href='/static/topjui/ext/css/print.css' type='text/css' rel='stylesheet'>";
    var strFormHtml = strStyleCSS + document.getElementById(printAreaId).innerHTML;
    LODOP.PRINT_INIT(printAreaId);
    LODOP.ADD_PRINT_HTM(40, 40, 350, 600, strFormHtml);
}

/**
 * 获得打印时额外的CSS样式
 * @returns {string}
 */
function getSinglePrintCss() {
    var printCSS = "<link href='/static/topjui/ext/css/print.css?version=" + Math.random() + "' type='text/css' rel='stylesheet'>";
    console.log(printCSS);
    return printCSS;
}

function getPrintCss(elementId) {
    var strBodyStyle = "<style>" + document.getElementById("style").innerHTML + "</style>";
    return strBodyStyle + "<body>" + document.getElementById(elementId).innerHTML + "</body>";
}

function addMainTab(url, title, iconCls) {
    var iframe = '<iframe src="' + url + '" frameborder="0" style="border:0;width:100%;height:98%;"></iframe>';
    var t = $('#index_tabs');
    var opts = {
        title: title,
        closable: true,
        iconCls: iconCls,
        content: iframe,
        border: false,
        fit: true,
        cls: 'leftBottomBorder'
    };
    if (t.tabs('exists', opts.title)) {
        t.tabs('select', opts.title);
    } else {
        t.tabs('add', opts);
    }
}

/**
 * 自动完成下拉框中格式化显示数据
 * @param row
 * @returns {string}
 */
function formatUserAutoCompleteItem(row) {
    var s = '<i class="icon-user"></i><span style="font-weight:bold">' + row.userName + '</span><br/>' +
        '<span style="color:#888">帐号：' + row.userNameId + '</span><br/>' +
        '<span style="color:#888">机构：' + row.orgName + '</span>';
    return s;
}

function formatDateTime(inputTime) {
    inputTime = inputTime<9999999999 ? inputTime*1000:inputTime;
    var date = new Date(inputTime);
    var y = date.getFullYear();
    if (y<=1970) {
        return "-";
    }
    var m = date.getMonth() + 1;
    m = m < 10 ? ('0' + m) : m;
    var d = date.getDate();
    d = d < 10 ? ('0' + d) : d;
    var h = date.getHours();
    h = h < 10 ? ('0' + h) : h;
    var minute = date.getMinutes();
    var second = date.getSeconds();
    minute = minute < 10 ? ('0' + minute) : minute;
    second = second < 10 ? ('0' + second) : second;
    return y + '-' + m + '-' + d+' '+h+':'+minute+':'+second;
};

function formatDate(inputTime) {
    inputTime = inputTime<9999999999 ? inputTime*1000:inputTime;
    var date = new Date(inputTime);
    var y = date.getFullYear();
    if (y<=1970) {
        return "-";
    }
    var m = date.getMonth() + 1;
    m = m < 10 ? ('0' + m) : m;
    var d = date.getDate();
    d = d < 10 ? ('0' + d) : d;
    var h = date.getHours();
    h = h < 10 ? ('0' + h) : h;
    var minute = date.getMinutes();
    var second = date.getSeconds();
    minute = minute < 10 ? ('0' + minute) : minute;
    second = second < 10 ? ('0' + second) : second;
    return y + '-' + m + '-' + d;
};

/**
 * 时间戳格式化函数
 * @param {string} format 格式
 * @param {int} timestamp 要格式化的时间 默认为当前时间
 * @return {string}   格式化的时间字符串
 */
function date(format, timestamp){
    format = format == '' ? "Y-m-d H:i:s" : format;
    var a, jsdate=((timestamp) ? new Date(timestamp*1000) : new Date());
    var pad = function(n, c){
        if((n = n + "").length < c){
            return new Array(++c - n.length).join("0") + n;
        } else {
            return n;
        }
    };
    var txt_weekdays = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
    var txt_ordin = {1:"st", 2:"nd", 3:"rd", 21:"st", 22:"nd", 23:"rd", 31:"st"};
    var txt_months = ["", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
    var f = {
        // Day
        d: function(){return pad(f.j(), 2)},
        D: function(){return f.l().substr(0,3)},
        j: function(){return jsdate.getDate()},
        l: function(){return txt_weekdays[f.w()]},
        N: function(){return f.w() + 1},
        S: function(){return txt_ordin[f.j()] ? txt_ordin[f.j()] : 'th'},
        w: function(){return jsdate.getDay()},
        z: function(){return (jsdate - new Date(jsdate.getFullYear() + "/1/1")) / 864e5 >> 0},

        // Week
        W: function(){
            var a = f.z(), b = 364 + f.L() - a;
            var nd2, nd = (new Date(jsdate.getFullYear() + "/1/1").getDay() || 7) - 1;
            if(b <= 2 && ((jsdate.getDay() || 7) - 1) <= 2 - b){
                return 1;
            } else{
                if(a <= 2 && nd >= 4 && a >= (6 - nd)){
                    nd2 = new Date(jsdate.getFullYear() - 1 + "/12/31");
                    return date("W", Math.round(nd2.getTime()/1000));
                } else{
                    return (1 + (nd <= 3 ? ((a + nd) / 7) : (a - (7 - nd)) / 7) >> 0);
                }
            }
        },

        // Month
        F: function(){return txt_months[f.n()]},
        m: function(){return pad(f.n(), 2)},
        M: function(){return f.F().substr(0,3)},
        n: function(){return jsdate.getMonth() + 1},
        t: function(){
            var n;
            if( (n = jsdate.getMonth() + 1) == 2 ){
                return 28 + f.L();
            } else{
                if( n & 1 && n < 8 || !(n & 1) && n > 7 ){
                    return 31;
                } else{
                    return 30;
                }
            }
        },

        // Year
        L: function(){var y = f.Y();return (!(y & 3) && (y % 1e2 || !(y % 4e2))) ? 1 : 0},
        //o not supported yet
        Y: function(){return jsdate.getFullYear()},
        y: function(){return (jsdate.getFullYear() + "").slice(2)},

        // Time
        a: function(){return jsdate.getHours() > 11 ? "pm" : "am"},
        A: function(){return f.a().toUpperCase()},
        B: function(){
            // peter paul koch:
            var off = (jsdate.getTimezoneOffset() + 60)*60;
            var theSeconds = (jsdate.getHours() * 3600) + (jsdate.getMinutes() * 60) + jsdate.getSeconds() + off;
            var beat = Math.floor(theSeconds/86.4);
            if (beat > 1000) beat -= 1000;
            if (beat < 0) beat += 1000;
            if ((String(beat)).length == 1) beat = "00"+beat;
            if ((String(beat)).length == 2) beat = "0"+beat;
            return beat;
        },
        g: function(){return jsdate.getHours() % 12 || 12},
        G: function(){return jsdate.getHours()},
        h: function(){return pad(f.g(), 2)},
        H: function(){return pad(jsdate.getHours(), 2)},
        i: function(){return pad(jsdate.getMinutes(), 2)},
        s: function(){return pad(jsdate.getSeconds(), 2)},
        //u not supported yet

        // Timezone
        //e not supported yet
        //I not supported yet
        O: function(){
            var t = pad(Math.abs(jsdate.getTimezoneOffset()/60*100), 4);
            if (jsdate.getTimezoneOffset() > 0) t = "-" + t; else t = "+" + t;
            return t;
        },
        P: function(){var O = f.O();return (O.substr(0, 3) + ":" + O.substr(3, 2))},
        //T not supported yet
        //Z not supported yet

        // Full Date/Time
        c: function(){return f.Y() + "-" + f.m() + "-" + f.d() + "T" + f.h() + ":" + f.i() + ":" + f.s() + f.P()},
        //r not supported yet
        U: function(){return Math.round(jsdate.getTime()/1000)}
    };

    return format.replace(/[\]?[a-zA-Z]/g, function(t, s){
        if( t!=s ){
            // escaped
            ret = s;
        } else if( f[s] ){
            // a date function exists
            ret = f[s]();
        } else{
            // nothing special
            ret = s;
        }
        return ret;
    });
}

/***
 * 封装一些常用的方法
 */
(function ($) {
    var rootNodeId = "1";

    function appendNodes(dataList, target) {
        // 如果不存在则是根节点
        if (target.length == 0) {
            var arr = findChildByPid(dataList, rootNodeId);
            for (var i = 0; i < arr.length; i++) {
                target.push(arr[i]);
            }
        }
        // 如果不是根节点插入children
        $.each(target, function (i) {
            var resArr = findChildByPid(dataList, target[i].id);
            if (resArr.length > 0) {
                target[i].children = resArr;
                appendNodes(dataList, target[i].children);
            }
        })
    };

    function findChildByPid(data, pid) {
        var resultArr = []; // 保存查询结果
        $.each(data, function (index) {
            if (data[index]["pid"] == pid) {
                resultArr.push(data[index]);
            }
        });
        return resultArr;
    };
    $.extend({
        /***
         * 根据根节点组合生成树形
         * @param dataList 数据列表
         * @param root 根节点ID
         * @returns {[]} 列表
         */
        getTreeList: function (dataList, root) {
            rootNodeId = root || rootNodeId;
            var target = [];
            appendNodes(dataList, target);
            return target;
        },
        /***
         * 将列表数据转换成指定结构树形的数据
         * @param oldData 原始数据
         * @param textField 文本字段名
         * @param idFiled id字段名
         * @param pidFiled pid字段名
         * @returns {[]}
         */
        formateTreeData: function (oldData, textField, idFiled, pidFiled) {
            var newData = [];
            idFiled = idFiled || 'id';
            textField = textField || "text";
            pidFiled = pidFiled || 'pid';
            $.each(oldData, function (i) {
                newData.push($.extend(oldData[i], {
                    id: oldData[i][idFiled],
                    text: oldData[i][textField],
                    pid: oldData[i][pidFiled]
                }))
            });
            return newData;
        },
        /***
         * 通过js触发打开一个弹框
         * @param opts 需要覆盖的属性
         * @param handler 点击提交按钮处理的方法
         * @returns {*|jQuery|HTMLElement}
         */
        openDialog: function (opts, handler) {
            var myDialogId = opts.id || (new Date()).getTime();
            var $myDialog = $("<form id='" + myDialogId + "' style='overflow-x: hidden' ></form>");
            var defaultOptions = {
                id: myDialogId,// 唯一标识id
                title: '操作页面',// 标题
                closed: false, // 关闭状态
                height: 450, // 高度
                content: "空空如也~~",
                buttons: [
                    {
                        text: '提交', iconCls: 'fa fa-save', btnCls: 'topjui-btn-green', handler: function () {
                            // . 验证表单是否填写
                            var flag = $myDialog.iForm("validate");
                            if (!flag) {
                                $.iMessager.alert('提示', '请将必填信息填写完整！', 'messager-info');
                                return;
                            }
                            // 执行自定义方法
                            if (typeof handler == "function") {
                                // this只按钮
                                handler.call($myDialog, this);
                            }
                        }
                    },
                    {
                        text: '关闭', iconCls: 'fa fa-close', btnCls: 'topjui-btn-red', handler: function () {
                            $myDialog.iDialog("destroy")
                        }
                    }
                ]
            };
            $myDialog.iDialog($.extend(true, {}, defaultOptions, opts));
            return $myDialog;
        },
        /**
         * 打开人员选择器弹框
         * @param rows 已选中的行
         * @param handler 处理的方法
         */
        openUserDialog: function (rows, handler) {
            if (rows instanceof Function) {
                handler = rows;
            }
            //userArr == [{userName:"赵侠策",userNameId:'ewsd0001'},{userName:"宋景民",userNameId:'ewsd0002'}];
            this.openDialog({
                content: '页面加载中……',
                href: '/mdata/user/selector',
                height: 580,
                title: "人员选择器",
                width: 900,
                onLoad: function () {
                    if (rows instanceof Array && rows.length > 0) {
                        _Operate.receiverDatagrid.iDatagrid("loadData", rows);
                        _Operate.onRowChange();
                    }
                },
                buttons: [{text: "保存", iconCls: 'fa fa-plus'}]
            }, function () {
                // 获取选中的行记录
                var rows = _Operate.userList;
                if (rows.length > 0) {
                    handler.call(this, rows);
                } else {
                    $.iMessager.alert('提示', '请选择人员！', 'messager-info');
                }

            })
        },
        /***
         * 导出excel
         * @param reportName 报表名称
         * @param excelName Excel名称
         * @param params Array 需要传递到报表里面的参数信息 {name:123,title:123}
         */
        exportExcel: function (reportName, excelName, params) {
            var href = '/report/ureport/excel?_u=mysql-' + reportName+'.ureport.xml&_n=' + excelName;
            $.each(params, function (item) {
                href += params[item] ? ('&' + item + "=" + params[item]) : "";
            });
            window.location.href = href;
        },
        /**
         * 通过uuid 预览pdf插件
         * @param uuid 附件在附件表中的uuid
         */
        pdfViewer:function (uuid) {
            var encodeURL ='/system/attachment/pdfStreamHandeler?uuid='+uuid;
            var str1 = encodeURIComponent(encodeURL);
            return window.open("/plugins/pdfjs/web/viewer.html?file="+str1,"_blank");
        },
        /**
         * 打开一个新的tab页面
         * @param params {titil:'标题',url:'地址',id:'唯一标识',iconCls:'图标'}
         * @returns {*|void}
         */
        addTab:function (params){
            if(!params.url){
                return $.iMessager.alert('我的消息','新的tab页面打开失败，未知URL！','messager-info');
            }
            var iframe = '<iframe src="' + params.url + '" scrolling="auto" frameborder="0" style="width:100%;height:100%;"></iframe>';
            var t = parent.$('#index_tabs');
            var opts = {
                id: getRandomNumByDef(),
                title: params.title,
                closable: typeof (params.closable) != "undefined" ? params.closable : true,
                iconCls: params.iconCls ? params.iconCls : 'fa fa-file-text-o',
                content: iframe,
                //href: params.url,
                border: params.border || true,
                fit: true
                //cls: 'leftBottomBorder'
            };
            var tabRecord = t.iTabs('options').record;
            // 判断面板是否存在
            if (params.id && tabRecord[params.id]) {
                var index = t.iTabs('getTabIndex', tabRecord[params.id]);
                if (index != -1) {
                    t.iTabs('select', index);
                    if(params.refresh){
                        var tab = t.iTabs('getSelected')
                        var refresh_iframe = tab.find('iframe')[0];
                        refresh_iframe.contentWindow.location.href = refresh_iframe.src;
                    }
                    return t;
                }
            }
            t.iTabs('add', opts);
            tabRecord[params.id||opts.id] = t.iTabs("getSelected");
        }
    });
})(jQuery);