var news_click = 0;
; (function ($) {

    $.index = $.index || {};
    $.index = {

        //初始化函数
        inits: function () {
            this.loadBanner();
        },
        //小的交互效果
        otherFun: function () {
            $('.citySlider').flexslider({
                animation: "slide",
                controlNav: false,
                animationLoop: false,
                pauseOnHover: true,
                directionNav: true,
                slideshow: false,
                prevText: ' ',
                nextText: ' ',
                itemWidth: 1280,
                before: function () {
                    $("#video_box").hide().empty();
                }
            });
            // citySlider
            $('.citySlider em').click(function () {
                $("#video_box").show();
                $('.citySliderWrapper .close').show();
                var playerSwf = "/Asserts/Global/js/frame/ckplayer/ckplayer.swf";
                var playerName = "video_player";
                var containerId = "video_box";
                var videosrc = $(this).next().attr("data-video");
                var videosrch5 = videosrc + "->video/mp4";

                var videoData = { flashaddr: videosrc, html5addr: videosrch5 };

                var flashvars = {
                    f: videoData.flashaddr,//使用flash播放器时视频地址的
                    c: 0,//风格配置参数
                    p: 1 //1默认播放0暂停
                    //i: preImg,//预览图
                };
                var params = { bgcolor: '#FFF', allowFullScreen: true, allowScriptAccess: 'always', wmode: 'transparent' };
                var html5video = [];
                html5video.push(videoData.html5addr);
                CKobject.embed(playerSwf, containerId, playerName, '100%', '100%', false, flashvars, html5video, params);
            });

            function isDefine(value) {
                if (value == null || value == "" || value == "undefined" || value == undefined || value == "null" || value == "(null)" || value == 'NULL' || typeof (value) == 'undefined') {
                    return false;
                }
                else {
                    value = value + "";
                    value = value.replace(/\s/g, "");
                    if (value == "") {
                        return false;
                    }
                    return true;
                }
            }

            $('.citySliderWrapper .close').click(function () {
                $('.citySliderWrapper .close').hide();
                $("#video_box").hide();
                $('.citySliderWrapper .close').hide();
                $.main.videoFun.clearVideo($(this).attr('data-pid'));
            });

            $('.prodSlider').flexslider({
                animation: "slide",
                controlNav: false,
                animationLoop: false,
                pauseOnHover: true,
                directionNav: true,
                slideshow: false,
                prevText: ' ',
                nextText: ' ',
                itemWidth: 1280
            });

            var maxHeight = 0;
            $(".groupAction .gaContent").each(function(){
                if ($(this).height() > maxHeight) { maxHeight = $(this).height(); }
            });

            $(".groupAction .gaContent").height(maxHeight);

            //$(".citySlider em").click(function () {

            //---视频start---
            //	var playerId = "video_player";
            //	var containerId = $.main.videoFun.showPop(playerId);//弹出层播放

            //实际上线，用$(this).attr("data-video...")存放视频地址
            //	var videoData = { flashaddr: "../video/mov_bbb.mp4", html5addr: "../video/mov_bbb.mp4->video/mp4" };
            //	$.main.videoFun.playVideo(videoData.flashaddr, videoData.html5addr, containerId, playerId, null, null, null);
            //--视频end---

            //	var winH = $(window).height();
            //	var playH = $('.play_box').height();
            //	var topVal = 0;
            //	if (winH > playH) {
            //	    topVal = (winH - playH) / 2;
            //	}
            //	else {
            //	    $('.play_box').height(winH);//撑满高度
            //	}
            //	$('.videoPop .popWrap').css({ top: topVal });
            //	$('body').css({ overflow: 'hidden' });
            //});

            //新闻资讯点击轮番
            $('.news_p_n').click(function () {
                var t = $(this)
                if (t.hasClass('active')) {
                    return false;
                }
                var bdobj = t.parent().siblings('ul');
                var lin = bdobj.find('li').length
                //              if (news_click == 0) {
                //                  var n_l_l = bdobj.find('li').length - 1;
                //                  var n_l_w = bdobj.find('li').innerWidth();
                //                  var htmlf = '<li class="m_hide">' + bdobj.find('li').eq('0').html() + '</li>';
                //                  var htmle = '<li class="m_hide">' + bdobj.find('li').eq(+n_l_l).html() + '</li>';
                //                  bdobj.append(htmlf).prepend(htmle).css('left', -1 * n_l_w + 'px')
                //                  news_click++;
                //              }
                if (isNaN(bdobj.attr('lf_index'))) {
                    var n = 1;
                    bdobj.attr('lf_index', n)
                } else {
                    n = bdobj.attr('lf_index')
                }
                var n_l_l_s = bdobj.find('li').length + 1;
                var n_l_w_s = bdobj.find('li').innerWidth();
                if (t.hasClass('news_p')) {
                    n--;
                    bdobj.attr('lf_index', n)
                } else {
                    n++;
                    bdobj.attr('lf_index', n)
                }
                if(n==1 && t.hasClass('news_p')){
                    t.addClass('visibility_hidden');
                }else{
                    t.siblings('li').removeClass('visibility_hidden');
                }
                if (n == (n_l_l_s - 1) && t.hasClass('news_n')) {
                    t.addClass('visibility_hidden');
                } else {
                    t.siblings('li').removeClass('visibility_hidden');
                }
                //              if(n==1 && t.hasClass('news_p')){
                //              	$('.news_n').addClass('visibility_hidden')
                //              	return false;
                //              }
                //              if(n==0){
                //              	$('.news_p').addClass('visibility_hidden')
                //              	return false;
                //              }
                bdobj.stop().animate({ 'left': n_l_w_s * (n - 1) * -1 + 'px' }, 1000)

            })
        },
        //加载banner图
        loadBanner: function () {

            //初始化banner控件
            var initBannerControl = function (sliderJq, isMobile) {
                var startBeforeFun = function (curSlide, isMb) {
                    var attrStr = isMb ? "data-mobile" : "data-src";
                    if (curSlide.closest(".flexslider").hasClass("resizeMb")) {
                        attrStr = "data-mobile";
                    }

                    if (curSlide.closest(".flexslider").hasClass("resizePc")) {
                        attrStr = "data-src";
                    }

                    var prevSlide = curSlide.prev("li");
                    var nextSlide = curSlide.next("li");

                    var curImg = curSlide.find("img.flexLazy");
                    var prevImg = prevSlide.find("img.flexLazy");
                    var nextImg = nextSlide.find("img.flexLazy");

                    if (curImg.length > 0) {
                        curImg.attr('src', curImg.attr(attrStr)).removeClass("flexLazy");
                    }

                    if (prevImg.length > 0) {
                        prevImg.attr('src', prevImg.attr(attrStr)).removeClass("flexLazy");
                    }

                    if (nextImg.length > 0) {
                        nextImg.attr('src', nextImg.attr(attrStr)).removeClass("flexLazy");
                    }
                };

                if (isMobile) {
                    var firstImg = sliderJq.find("ul.slides img:eq(0)");
                    firstImg.attr('src', firstImg.attr('data-mobile'));
                }
                sliderJq.flexslider({
                    animation: "slide",
                    controlNav: true,
                    animationLoop: false,
                    pauseOnHover: true,
                    directionNav: true,
                    slideshow: true,//自动播放
                    prevText: ' ',
                    nextText: ' ',
                    start: function (slider) {
                        var $slide = $(slider).find("li.flex-active-slide");
                        startBeforeFun($slide, isMobile);
                    },
                    before: function (slider) {
                        var slides = slider.slides,
                            index = slider.animatingTo,
                            $slide = $(slides[index]);
                        startBeforeFun($slide, isMobile);
                    }
                });
            }

            //调用
            if ($(window).width() < 769) {
                initBannerControl($(".homeslider"), true);
                initBannerControl($(".citySlider"), true);
                initBannerControl($(".prodSlider"), true);
            }
            else {
                initBannerControl($(".homeslider"), false);
            }

            $(window).resize(function () {

                if ($(window).width() >= 768) {
                    if ($(".homeslider").hasClass("resizePc") == false) {
                        $(".homeslider img").each(function () {
                            var url = $(this).attr("data-src");
                            $(this).attr("src", url);
                        });
                        $(".homeslider").removeClass("resizeMb").addClass("resizePc");
                    }
                }
                else {
                    if ($(".homeslider").hasClass("resizeMb") == false) {
                        $(".homeslider img").each(function () {
                            var url = $(this).attr("data-mobile");
                            $(this).attr("src", url);
                        });
                        $(".homeslider").removeClass("resizePc").addClass("resizeMb");
                    }
                }


            });
        }
    }


    //$.extend($.fn, $.index);



    //DOM载入就绪执行
    $(function () {
        //初始化
        $.index.inits();
        $.index.otherFun();
    });

})(jQuery);
