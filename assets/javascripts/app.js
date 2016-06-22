(function($) {
  var AppView, JSON_CODE;

  JSON_CODE = {
    noLogin: -1
  };

  AppView = Backbone.View.extend({
    el: "body",
    repliesPerPage: 50,
    events: {
      "click .topic-detail .panel-footer a.watch": "toggleWatch",
      "click .topic-detail .panel-footer a.star": "toggleStar",
      "click .md-dropdown .dropdown-menu li": "toggleDropdown",
      "click #replies .reply .btn-reply": "reply",
      "click #replies a.mention-floor": "mentionFloor",
      "click .button-captcha": "refreshCaptcha",
      "click .header .form-search .btn-search": "openHeaderSearchBox",
      "click .header .form-search .btn-close": "closeHeaderSearchBox",
      "click .social-share-button a": "sharePage"
    },
    initialize: function() {
      //this.initWebSocket();
      this.initShareButtonPopover();
      this.initHighlight();
      this.setupAjaxCommonResult();
      return $.notifier.checkOrRequirePermission();
    },
//    initWebSocket: function() {
//      //this.ws = new WebSocket("ws://" + location.host + "/msg");
//      //return this.ws.onmessage = this.onWebSocketMessage;
//    },
    initHighlight: function() {
      return $("pre code").each(function(i, block) {
        return hljs.highlightBlock(block);
      });
    },
    setupAjaxCommonResult: function() {
      return $.ajaxSetup({
        success: function(res) {
          if (res.code === JSON_CODE.noLogin) {
            return location.href = "/signin";
          }
        }
      });
    },
//    onWebSocketMessage: function(res) {
//      var badge, counter, notify;
//      notify = JSON.parse(res.data);
//      badge = $(".notification-count a");
//      counter = badge.find(".count");
//      if (notify.unread_count > 0) {
//        badge.addClass("new");
//        counter.text(notify.unread_count);
//        if (notify.is_new) {
//          return $.notifier.notify(notify.avatar, "回帖通知", notify.title, notify.path);
//        }
//      } else {
//        badge.removeClass("new");
//        return counter.text(0);
//      }
//    },
    toggleDropdown: function(e) {
      var $target;
      $target = $(e.currentTarget);
      $target.closest('.input-group-btn').find('[data-bind="value"]').val($target.data("id")).end().find('[data-bind="label"]').text($target.text()).end().children('.dropdown-toggle').dropdown('toggle');
      return false;
    },
    toggleStar: function(e) {
      var a, count, topicId;
      a = $(e.target);
      topicId = a.data("id");
      count = parseInt(a.data("count"));
      if (a.hasClass("followed")) {
        $.post("/topics/" + topicId + "/unstar").done(function(res) {
          var labelText, newCount;
          newCount = count - 1;
          labelText = "" + newCount + " 人收藏";
          return a.removeClass("followed").data("count", newCount).html('<i class="fa fa-star-o"></i> ' + labelText);
        });
      } else {
        $.post("/topics/" + topicId + "/star").done(function(res) {
          var labelText, newCount;
          newCount = count + 1;
          labelText = "" + newCount + " 人收藏";
          return a.addClass("followed").data("count", newCount).html('<i class="fa fa-star"></i> ' + labelText);
        });
      }
      return false;
    },
    toggleWatch: function(e) {
      var a, topicId;
      a = $(e.target);
      topicId = a.data("id");
      if (a.hasClass("followed")) {
        $.post("/topics/" + topicId + "/unwatch").done(function(res) {
          return a.removeClass("followed").attr("title", "关注此话题，当有新回帖的时候会收到通知").html('<i class="fa fa-eye"></i> 关注');
        });
      } else {
        $.post("/topics/" + topicId + "/watch").done(function(res) {
          return a.addClass("followed").attr("title", "点击取消关注").html('<i class="fa fa-eye"></i> 已关注');
        });
      }
      return false;
    },
    reply: function(e) {
      var floor, login, new_text, reply_body, _el;
      _el = $(e.target);
      floor = _el.data("floor");
      login = _el.data("login");
      reply_body = $(".reply-form textarea");
      new_text = "#" + floor + "楼 @" + login + " ";
      if (reply_body.val().trim().length === 0) {
        new_text += '';
      } else {
        new_text = "\n" + new_text;
      }
      reply_body.focus().val(reply_body.val() + new_text);
      return false;
    },
    mentionFloor: function(e) {
      var floor, page, replyEl, url, _el;
      _el = $(e.target);
      floor = _el.data('floor');
      replyEl = $("#reply" + floor);
      if (replyEl.length > 0) {
        this.highlightReply(replyEl);
      } else {
        page = this.pageOfFloor(floor);
        url = window.location.pathname + ("?page=" + page) + ("#reply" + floor);
        this.gotoUrl(url);
      }
      return replyEl;
    },
    highlightReply: function(replyEl) {
      $("#replies .reply").removeClass("light");
      return replyEl.addClass("light");
    },
    pageOfFloor: function(floor) {
      return Math.floor((floor - 1) / Topics.repliesPerPage) + 1;
    },
    gotoUrl: function(url) {
      return location.href = url;
    },
    refreshCaptcha: function(e) {
      var img;
      img = $(e.target);
      img.attr("src", "/captcha?t=" + ((new Date).getTime()));
      return false;
    },
    initShareButtonPopover: function(e) {
      var btn, sharePanelHTML;
      btn = $(".share-button");
      if (btn.size() > 0) {
        btn.on("click", function() {
          return false;
        });
        sharePanelHTML = $(".social-share-button")[0].outerHTML;
        btn.data("html", true).data("trigger", "click").data("placement", "top").data("content", sharePanelHTML);
        return btn.popover();
      }
    },
    sharePage: function(e) {
      var link;
      link = $(e.currentTarget);
      $(".share-button").popover("hide");
      return SocialShareButton.share(link);
    },
    openHeaderSearchBox: function(e) {
      $(".header .form-search").addClass("active");
      $(".header .form-search input").focus();
      return false;
    },
    closeHeaderSearchBox: function(e) {
      $(".header .form-search input").val("");
      $(".header .form-search").removeClass("active");
      return false;
    }
  });

  $(document).on("ready page:load", function() {
    return new AppView();
  });

  window.SocialShareButton = {
    openUrl: function(url) {
      window.open(url);
      return false;
    },
    share: function(el) {
      var $parent, appkey, desc, get_tumblr_extra, img, site, title, tumblr_params, url, via;
      site = $(el).data('site');
      appkey = $(el).data('appkey') || '';
      $parent = $(el).parent();
      title = encodeURIComponent($parent.data('title') || '');
      img = encodeURIComponent($parent.data("img") || '');
      url = encodeURIComponent($parent.data("url") || '');
      via = encodeURIComponent($parent.data("via") || '');
      desc = encodeURIComponent($parent.data("desc") || ' ');
      if (url.length === 0) {
        url = encodeURIComponent(location.href);
      }
      switch (site) {
        case "email":
          location.href = "mailto:?to=&subject=" + title + "&body=" + url;
          break;
        case "weibo":
          SocialShareButton.openUrl("http://service.weibo.com/share/share.php?url=" + url + "&type=3&pic=" + img + "&title=" + title + "&appkey=" + appkey);
          break;
        case "twitter":
          SocialShareButton.openUrl("https://twitter.com/intent/tweet?url=" + url + "&text=" + title + "&via=" + via);
          break;
        case "douban":
          SocialShareButton.openUrl("http://shuo.douban.com/!service/share?href=" + url + "&name=" + title + "&image=" + img + "&sel=" + desc);
          break;
        case "facebook":
          SocialShareButton.openUrl("http://www.facebook.com/sharer.php?u=" + url);
          break;
        case "qq":
          SocialShareButton.openUrl("http://sns.qzone.qq.com/cgi-bin/qzshare/cgi_qzshare_onekey?url=" + url + "&title=" + title + "&pics=" + img + "&summary=" + desc + "&site=" + appkey);
          break;
        case "tqq":
          SocialShareButton.openUrl("http://share.v.t.qq.com/index.php?c=share&a=index&url=" + url + "&title=" + title + "&pic=" + img + "&appkey=" + appkey);
          break;
        case "baidu":
          SocialShareButton.openUrl("http://hi.baidu.com/pub/show/share?url=" + url + "&title=" + title + "&content=" + desc);
          break;
        case "kaixin001":
          SocialShareButton.openUrl("http://www.kaixin001.com/rest/records.php?url=" + url + "&content=" + title + "&style=11&pic=" + img + "&aid=" + appkey);
          break;
        case "renren":
          SocialShareButton.openUrl("http://widget.renren.com/dialog/share?resourceUrl=" + url + "&srcUrl=" + url + "&title=" + title + "&pic=" + img + "&description=" + desc);
          break;
        case "google_plus":
          SocialShareButton.openUrl("https://plus.google.com/share?url=" + url);
          break;
        case "google_bookmark":
          SocialShareButton.openUrl("https://www.google.com/bookmarks/mark?op=edit&output=popup&bkmk=" + url + "&title=" + title);
          break;
        case "delicious":
          SocialShareButton.openUrl("http://www.delicious.com/save?url=" + url + "&title=" + title + "&jump=yes&pic=" + img);
          break;
        case "plurk":
          SocialShareButton.openUrl("http://www.plurk.com/?status=" + title + ": " + url + "&qualifier=shares");
          break;
        case "pinterest":
          SocialShareButton.openUrl("http://www.pinterest.com/pin/create/button/?url=" + url + "&media=" + img + "&description=" + title);
          break;
        case "tumblr":
          get_tumblr_extra = function(param) {
            var cutom_data;
            cutom_data = $(el).attr("data-" + param);
            if (cutom_data) {
              return encodeURIComponent(cutom_data);
            }
          };
          tumblr_params = function() {
            var params, path, quote, source;
            path = get_tumblr_extra('type') || 'link';
            params = (function() {
              switch (path) {
                case 'text':
                  title = get_tumblr_extra('title') || title;
                  return "title=" + title;
                case 'photo':
                  title = get_tumblr_extra('caption') || title;
                  source = get_tumblr_extra('source') || img;
                  return "caption=" + title + "&source=" + source;
                case 'quote':
                  quote = get_tumblr_extra('quote') || title;
                  source = get_tumblr_extra('source') || '';
                  return "quote=" + quote + "&source=" + source;
                default:
                  title = get_tumblr_extra('title') || title;
                  url = get_tumblr_extra('url') || url;
                  return "name=" + title + "&url=" + url;
              }
            })();
            return "/" + path + "?" + params;
          };
          SocialShareButton.openUrl("http://www.tumblr.com/share" + (tumblr_params()));
      }
      return false;
    }
  };
})(jQuery);
