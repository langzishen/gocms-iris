// 首页加载完后，取消加载中状态
$(window).load(function () {
    $("#loading").fadeOut();
});
var isFullScreen = false;

var App = function () {
    var setFullScreen = function () {
        var docEle = document.documentElement;
        if (docEle.requestFullscreen) {
            //W3C
            docEle.requestFullscreen();
        } else if (docEle.mozRequestFullScreen) {
            //FireFox
            docEle.mozRequestFullScreen();
        } else if (docEle.webkitRequestFullScreen) {
            //Chrome等
            docEle.webkitRequestFullScreen();
        } else if (docEle.msRequestFullscreen) {
            //IE11
            docEle.msRequestFullscreen();
        } else {
            $.iMessager.alert('温馨提示', '该浏览器不支持全屏', 'messager-warning');
        }
    };

    //退出全屏 判断浏览器种类
    var exitFullScreen = function () {
        // 判断各种浏览器，找到正确的方法
        var exitMethod = document.exitFullscreen || //W3C
            document.mozCancelFullScreen ||    //Chrome等
            document.webkitExitFullscreen || //FireFox
            document.msExitFullscreen; //IE11
        if (exitMethod) {
            exitMethod.call(document);
        }
        else if (typeof window.ActiveXObject !== "undefined") {//for Internet Explorer
            var wscript = new ActiveXObject("WScript.Shell");
            if (wscript !== null) {
                wscript.SendKeys("{F11}");
            }
        }
    };

    return {
        init: function () {

        },
        handleFullScreen: function () {
            if (isFullScreen) {
                exitFullScreen();
                isFullScreen = false;
            } else {
                setFullScreen();
                isFullScreen = true;
            }
        }
    };
};

$(function () {
    $(".collapseMenu").on("click", function () {
        var p = $("#index_layout").iLayout("panel", "west")[0].clientWidth;
        if (p > 0) {
            $('#index_layout').iLayout('collapse', 'west');
            $(this).children('span').removeClass('fa-chevron-circle-left').addClass('fa-chevron-circle-right');
        } else {
            $('#index_layout').iLayout('expand', 'west');
            $(this).children('span').removeClass('fa-chevron-circle-right').addClass('fa-chevron-circle-left');
        }
    });

    // 首页tabs选项卡
    var index_tabs = $('#index_tabs').iTabs({
        fit: true,
        tools: [{
            iconCls: 'fa fa-home',
            handler: function () {
                $('#index_tabs').iTabs('select', 0);
            }
        }, {
            iconCls: 'fa fa-refresh',
            handler: function () {
                var refresh_tab = $('#index_tabs').iTabs('getSelected');
                var refresh_iframe = refresh_tab.find('iframe')[0];
                refresh_iframe.contentWindow.location.href = refresh_iframe.src;
                //$('#index_tabs').trigger(TOPJUI.eventType.initUI.base);

                //var index = $('#index_tabs').iTabs('getTabIndex', $('#index_tabs').iTabs('getSelected'));
                //console.log(index);
                //$('#index_tabs').iTabs('getTab', index).iPanel('refresh');
            }
        }, {
            iconCls: 'fa fa-close',
            handler: function () {
                var index = $('#index_tabs').iTabs('getTabIndex', $('#index_tabs').iTabs('getSelected'));
                var tab = $('#index_tabs').iTabs('getTab', index);
                if (tab.iPanel('options').closable) {
                    $('#index_tabs').iTabs('close', index);
                }
            }
        }, {
            iconCls: 'fa fa-arrows-alt',
            handler: function () {
                App().handleFullScreen();
            }
        }],
        //监听右键事件，创建右键菜单
        onContextMenu: function (e, title, index) {
            e.preventDefault();
            if (index >= 0) {
                $('#mm').iMenu('show', {
                    left: e.pageX,
                    top: e.pageY
                }).data("tabTitle", title);
            }
        }
    });

    //tab右键菜单
    $("#mm").iMenu({
        onClick: function (item) {
            tabMenuOprate(this, item.name);
        }
    });

    // 初始化accordion
    $("#RightAccordion").iAccordion({
        fit: true,
        border: false
    });

    // 绑定横向导航菜单点击事件
    $(".systemName").on("click", function (e) {
        //generateMenu(e.currentTarget.dataset.menuid, e.target.textContent); //IE9及以下不兼容data-menuid属性
        //generateMenu(e.target.getAttribute('data-menuid'), e.target.textContent);
        generateMenu($(this).attr("id"), $(this).attr("title"));
        $(".systemName").removeClass("selected");
        $(this).addClass("selected");
    });

    $("#index_layout").css("width",$(document).width()+"px");
    $("#index_layout").css("height",$(document).height()+"px");
    // 导航栏
    hideNavBtn();
    $(window).resize(function(){
        $("#index_layout").css("width",$(document).width()+"px");
        $("#index_layout").css("height",$(document).height()+"px");
        setTimeout(function () {
            hideNavBtn();
        },400)
    });
    $(".nav-right-btn").on('click',function () {
        navsScrollTo("next");
    });
    $(".nav-left-btn").on('click',function () {
        navsScrollTo("previous");
    });
    $(".nav-group").on("click",'li',function(e){
        //添加选择样式
        $(".nav-group .selected").removeClass("selected");
        $(this).addClass("selected");

        navsScrollTo("first",this);
        // 刷新左侧导航菜单
        var id = $(this).attr('id');
        var title =  $(this).attr('title');
        if(id && title){
            generateMenu(id, title);
        }
    });
    var changing = false; //是否正在滚动中
    function navsScrollTo(step,target){
        if(changing){
            return
        }else{
            changing = !changing;
            setTimeout(function () {
                changing = !changing;
            },400);
        }

        var contentWidth,itemWidth,areaWidth,currentLeft,times,$navContentBox;
        $navContentBox = $(".nav-group"); // 拿到导航项的盒子

        //1.得到滚动内容的总长度 contentWidth；
        contentWidth = parseInt($navContentBox.width());

        //2.得到滚动元素的宽度 itemWidth；
        itemWidth = 137;

        //3.得到当前currentLeft值
        currentLeft = parseInt($navContentBox.position().left);

        //4.得到滚动区域的总长度 areaWidth;
        areaWidth = parseInt( $navContentBox.parent().width());

        //5.按钮是否生效  并左右滚动
        if(step == 'next'){ //滚动到下一页
            if( areaWidth + Math.abs(currentLeft) < contentWidth){ // 不需要向右移动
                //5.1 移动多少个
                times = Math.floor(areaWidth / itemWidth);
                $navContentBox.css('left',currentLeft - itemWidth * times)
            }
        }else if(step == 'previous'){//滚动到上一页
            if (currentLeft != 0){
                if(Math.abs(currentLeft) <= areaWidth){
                    $navContentBox.css('left',0)
                }else{
                    times = Math.floor(areaWidth / itemWidth);
                    $navContentBox.css('left',-(Math.abs(currentLeft) - itemWidth * times))
                }
            }
        }else if(step == 'first'){//将点击的元素滚动到第一个
            liLeft = $(target).position().left;
            //获取点击项距离浏览器右侧的距离
            var liRight  =  $(".nav-bar").width() - $(target).offset().left - itemWidth - $('.nav-right').width();
            if(liRight < itemWidth){
                $navContentBox.css('left',-liLeft)
            }
        }
    }
    function hideNavBtn() {
        var contentWidth,areaWidth;
        var $navContentBox = $(".nav-group");
        //1.得到滚动内容的总长度 contentWidth；
        contentWidth = parseInt($navContentBox.width());

        //2.得到滚动区域的总长度 areaWidth;
        areaWidth = parseInt( $navContentBox.parent().width());

        if(contentWidth < areaWidth){
            $(".nav-btn").addClass('hide');
            $navContentBox.css('left',0)
        }else{
            $(".nav-btn").removeClass('hide')
        }
    }

    // 主页打开初始化时显示第一个系统的菜单
    $('.nav-group li').eq('0').trigger('click');


    //generateMenu(1325, "系统配置");

    // 显示系统首页
    /*setTimeout(function () {
     var indexTab = [];
     indexTab.iconCls = "icon-house";
     indexTab.text = "系统门户";
     var portal = $.getUrlParam("portal");
     if (portal == "system" || portal == null) portal = "system";
     indexTab.url = "html/portal/index.html";
     indexTab.closable = false;
     indexTab.border = false;
     addTab(indexTab);
     }, 1);*/

});

// Tab菜单操作
function tabMenuOprate(menu, type) {
    var allTabs = $('#index_tabs').iTabs('tabs');
    var allTabtitle = [];
    $.each(allTabs, function (i, n) {
        var opt = $(n).iPanel('options');
        if (opt.closable)
            allTabtitle.push(opt.title);
    });
    var curTabTitle = $(menu).data("tabTitle");
    var curTabIndex = $('#index_tabs').iTabs("getTabIndex", $('#index_tabs').iTabs("getTab", curTabTitle));
    switch (type) {
        case "1"://关闭当前
            if (curTabIndex > 0) {
                $('#index_tabs').iTabs("close", curTabTitle);
                return false;
                break;
            } else {
                $.iMessager.show({
                    title: '操作提示',
                    msg: '首页不允许关闭！'
                });
                break;
            }
        case "2"://全部关闭
            for (var i = 0; i < allTabtitle.length; i++) {
                $('#index_tabs').iTabs('close', allTabtitle[i]);
            }
            break;
        case "3"://除此之外全部关闭
            for (var i = 0; i < allTabtitle.length; i++) {
                if (curTabTitle != allTabtitle[i])
                    $('#index_tabs').iTabs('close', allTabtitle[i]);
            }
            $('#index_tabs').iTabs('scrollBy', 10);
            $('#index_tabs').iTabs('select', curTabTitle);
            break;
        case "4"://当前侧面右边
            for (var i = curTabIndex; i < allTabtitle.length; i++) {
                $('#index_tabs').iTabs('close', allTabtitle[i]);
            }
            $('#index_tabs').iTabs('select', curTabTitle);
            break;
        case "5": //当前侧面左边
            for (var i = 0; i < curTabIndex - 1; i++) {
                $('#index_tabs').iTabs('close', allTabtitle[i]);
            }
            $('#index_tabs').iTabs('select', curTabTitle);
            break;
        case "6": //刷新
            var currentTab = $('#index_tabs').iTabs('getSelected');
            var opts = $.data(currentTab[0], 'panel').options;
            if (opts.iframe) {
                var currentIframe = currentTab.find('iframe')[0];
                currentIframe.contentWindow.location.href = currentIframe.src;
            } else {
                $(currentTab[0]).panel('refresh');
            }
            break;
        case "7": //在新窗口打开
            var refresh_tab = $('#index_tabs').iTabs('getSelected');
            var refresh_iframe = refresh_tab.find('iframe')[0];
            window.open(refresh_iframe.src);
            break;
    }

}

//选项卡右键菜单
function findTabElement(target) {
    var $ele = $(target);
    if (!$ele.is("a")) {
        $ele = $ele.parents("a.menu_tab");
    }
    return $ele;
}

//保存页面id的field
var pageIdField = "data-pageId";

function getPageId(element) {
    if (element instanceof jQuery) {
        return element.attr(pageIdField);
    } else {
        return $(element).attr(pageIdField);
    }
}

function findTabTitle(pageId) {
    var $ele = null;
    $(".page-tabs-content").find("a.menu_tab").each(function () {
        var $a = $(this);
        if ($a.attr(pageIdField) == pageId) {
            $ele = $a;
            return false;//退出循环
        }
    });
    return $ele;
}

function getTabUrlById(pageId) {
    var $iframe = findIframeById(pageId);
    return $iframe[0].contentWindow.location.href;
}

function getTabUrl(element) {
    var pageId = getPageId(element);
    getTabUrlById(pageId);
}

function findTabPanel(pageId) {
    var $ele = null;
    $("#index_tabs").find("div.tab-pane").each(function () {
        var $div = $(this);
        if ($div.attr(pageIdField) == pageId) {
            $ele = $div;
            return false;//退出循环
        }
    });
    return $ele;
}

function findIframeById(pageId) {
    return findTabPanel(pageId).children("iframe");
}

function getActivePageId() {
    var $a = $('.page-tabs-content').find('.active');
    return getPageId($a);
}

//激活Tab,通过id
function activeTabByPageId(pageId) {
    $(".menu_tab").removeClass("active");
    $("#index_tabs").find(".active").removeClass("active");
    //激活TAB
    var $title = findTabTitle(pageId).addClass('active');
    findTabPanel(pageId).addClass("active");
    // scrollToTab($('.menu_tab.active'));
    scrollToTab($title[0]);
}

/**
 * 更换页面风格
 * @param topjuiThemeName
 */
function changeTheme(themeName) {/* 更换主题 */
    var $dynamicTheme = $('#dynamicTheme');
    var themeHref = $dynamicTheme.attr('href');
    var themeHrefNew = themeHref.substring(0, themeHref.indexOf('themes')) + 'themes/default/topjui.' + themeName + '.css';
    // 更换主页面风格
    $dynamicTheme.attr('href', themeHrefNew);

    // 更换iframe页面风格
    var $iframe = $('iframe');
    if ($iframe.length > 0) {
        for (var i = 1; i < $iframe.length; i++) {
            var ifr = $iframe[i];
            var $iframeDynamicTheme = $(ifr).contents().find('#dynamicTheme');
            var iframeThemeHref = $iframeDynamicTheme.attr('href');
            if (iframeThemeHref != undefined) {
                var iframeThemeHrefNew = iframeThemeHref.substring(0, iframeThemeHref.indexOf('themes')) + 'themes/default/topjui.' + themeName + '.css';
                $iframeDynamicTheme.attr('href', iframeThemeHrefNew);
            }
        }
    }

    $.cookie('topjuiThemeName', themeName, {
        expires: 7,
        path: '/'
    });
};
if ($.cookie('topjuiThemeName')) {
    changeTheme($.cookie('topjuiThemeName'));
}

// 退出系统
function logout() {
    $.iMessager.confirm('提示', '确定要退出吗?', function (r) {
        if (r) {
            $.iMessager.progress({
                text: '正在退出中...'
            });
            window.location.href = './login.html' + location.search;
        }
    });
}

// 生成左侧导航菜单
function generateMenu(menuId, systemName) {
    $(".panel-header .panel-title:first").html(systemName);
    var allPanel = $("#RightAccordion").iAccordion('panels');
    var size = allPanel.length;
    if (size > 0) {
        for (i = 0; i < size; i++) {
            var index = $("#RightAccordion").iAccordion('getPanelIndex', allPanel[i]);
            $("#RightAccordion").iAccordion('remove', 0);
        }
    }

    var url = "/admin/index/menu?menu_id=" + menuId;
    if (topJUI.config.appendRefererParam) url = appendRefererParam(url);
    $.get(
        url, {"levelId": "2"}, // 获取第一层目录
        function (data) {
            if (data == "0") {
                window.location = "/Account";
            }
            $.each(data, function (i, e) {// 循环创建手风琴的项
                var pid = e.pid;
                var isSelected = i == 0 ? true : false;
                $('#RightAccordion').iAccordion('add', {
                    fit: false,
                    title: e.title,
                    content: "<ul id='tree" + e.id + "' ></ul>",
                    border: false,
                    selected: isSelected,
                    iconCls: e.icon
                });
                //$.parser.parse();
                var secondUrl = "/admin/index/menu_zree?pid=" + e.id;
                if (topJUI.config.appendRefererParam) secondUrl = appendRefererParam(secondUrl);
                $.get(secondUrl, function (data) {// 循环创建树的项
                    $("#tree" + e.id).tree({
                        data: data,
                        lines: false,
                        animate: true,
                        onBeforeExpand: function (node, param) {
                            //var thirdUrl = "./json/menu/menu_" + node.id + ".json";
                            //if (topJUI.config.appendRefererParam) thirdUrl = appendRefererParam(thirdUrl);
                            //$("#tree" + e.id).tree('options').url = thirdUrl;
                        },
                        onClick: function (node) {
                            if (node.url) {
                                /*if(typeof node.attributes != "object") {
                                 node.attributes = $.parseJSON(node.attributes);
                                 }*/
                                addTab(node);
                            } else {
                                if (node.state == "closed") {
                                    $("#tree" + e.id).tree('expand', node.target);
                                } else if (node.state == 'open') {
                                    $("#tree" + e.id).tree('collapse', node.target);
                                }
                            }
                        }
                    });
                }, 'json');
            });
        }, "json"
    );
}

//打开Tab窗口
function addTab(params) {
    var t = $('#index_tabs');
    var $selectedTab = t.iTabs('getSelected');
    var selectedTabOpts = $selectedTab.iPanel('options');
    if (topJUI.config.appendRefererParam && params.url) params.url = appendRefererParam(params.url);
    var iframe = '<iframe src="' + params.url + '" scrolling="auto" frameborder="0" style="width:100%;height:100%;"></iframe>';

    var defaults = {
        id: getRandomNumByDef(),
        refererTab: {},
        title: params.text,
        iframe: topJUI.config.iframe,
        onlyInitParse: true,
        iconCls: 'fa fa-file-text-o',
        border: true,
        fit: true,
        closable: true
        //cls: 'leftBottomBorder'
    };
    var opts = $.extend(defaults, params);
    var ifOpts = opts.iframe ? {content: iframe} : {href: params.url};
    opts = $.extend(opts, ifOpts);

    if (t.iTabs('exists', opts.title)) {
        t.iTabs('select', opts.title);
    } else {
        var lastMenuClickTime = $.cookie("menuClickTime");
        var nowTime = new Date().getTime();
        if ((nowTime - lastMenuClickTime) >= 500) {
            $.cookie("menuClickTime", new Date().getTime());
            t.iTabs('add', opts);
        } else {
            $.iMessager.show({
                title: '温馨提示',
                msg: '操作过快，请稍后重试！'
            });
        }
    }
}

function addParentTab(options) {
    var src, title;
    src = options.href;
    title = options.title;

    var iframe = '<iframe src="' + src + '" frameborder="0" style="border:0;width:100%;height:100%;"></iframe>';
    parent.$('#index_tabs').iTabs("add", {
        title: title,
        content: iframe,
        closable: true,
        iconCls: 'fa fa-th',
        border: true
    });
}

function modifyPwd() {
    var opts = {
        id: 'pwdDialog',
        title: '修改密码',
        width: 400,
        height: 300,
        iconCls: 'fa fa-key',
        href: '/html/user/modifyPassword.html',
        buttons: [{
            text: '确定',
            iconCls: 'fa fa-save',
            btnCls: 'topjui-btn-green',
            handler: function () {
                if ($('#pwdDialog').form('validate')) {
                    if ($("#password").val().length < 6) {
                        $.iMessager.alert('警告', '密码长度不能小于6位', 'messager-warning');
                    } else {
                        var formData = $("#pwdDialog").serialize();
                        $.ajax({
                            url: './json/response/success.json',
                            type: 'post',
                            cache: false,
                            data: formData,
                            beforeSend: function () {
                                $.iMessager.progress({
                                    text: '正在操作...'
                                });
                            },
                            success: function (data, response, status) {
                                $.iMessager.progress('close');
                                if (data.statusCode == 200) {
                                    $.iMessager.show({
                                        title: '提示',
                                        msg: '操作成功'
                                    });
                                    $("#pwdDialog").iDialog('close').form('reset');

                                } else {
                                    $.iMessager.alert('操作失败！', '未知错误或没有任何修改，请重试！', 'messager-error');
                                }
                            }
                        });
                    }
                }
            }
        }, {
            text: '关闭',
            iconCls: 'fa fa-close',
            btnCls: 'topjui-btn-red',
            handler: function () {
                $("#pwdDialog").iDialog('close');
            }
        }]
    };
    $('#' + opts.id).iDialog('openDialog', opts);
};
