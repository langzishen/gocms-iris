$(function () {
    // 等待页面加载完成，链接websocket
    var socket = new SockJS('/'+request_app+'/message');
    stompClient = Stomp.over(socket);
    stompClient.connect({}, function (frame) {
        stompClient.subscribe('/'+request_app+'/webnotice/'+userNameId, function (result) {
            var message = JSON.parse(result.body)
            num = 1;
            if(message && (message instanceof Array)){
                num = message.length;
            }
            $.iMessager.show({
                title:"提示",
                msg:"您有"+num+"条新的消息提醒，请注意查收！<a style=\"text-decoration: underline;color: #1E9FFF\" href='javascript:addNewMessageTab()'>立即查看</a>" ,
                timeout:5000,
                showType:'slide'
            });
        });
    });
    // 获取未读消息的数量
    setTimeout(function () {
        $.post("/"+request_app+"/message/noreadnum",{},function (res) {
            var num = res ;
            if(parseInt(num) == 0){
                return;
            }
            try {
                $.iMessager.show({
                    title:"提示",
                    msg:"您有"+num+"条新的消息提醒，请注意查收！<a style=\"text-decoration: underline;color: #1E9FFF\" href='javascript:addNewMessageTab()'>立即查看</a>" ,
                    timeout:10000,
                    showType:'slide'
                });
            }catch (e) {
                console.log('未读消息获取异常！')
            }

        })
    },2000);
});
// 打开一个新的tab，如果已存在则刷新原来的tab
function addNewMessageTab() {
    var url = '/'+request_app+'/message/list', title ='我的消息', iconCls='fa fa-commenting';
    var iframe = '<iframe src="' + url + '" frameborder="0" style="border:0;width:100%;height:98%;"></iframe>';
    var t = $('#index_tabs');
    var opts = {
        title: title,
        closable: true,
        iconCls: iconCls,
        content: iframe,
        iframe: true,//加上后才能刷新
        border: false,
        fit: true,
        cls: 'leftBottomBorder'
    };
    if (t.tabs('exists', opts.title)) {
        t.tabs('select', opts.title);
        var tab = t.iTabs('getSelected');
        tab.iPanel('refresh',url);
    } else {
        t.tabs('add', opts);
    }
}
