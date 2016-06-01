(function() {
  var Notifier,
    __bind = function(fn, me){ return function(){ return fn.apply(me, arguments); }; };

  Notifier = (function() {
    function Notifier() {
      this.checkOrRequirePermission = __bind(this.checkOrRequirePermission, this);
      this.setPermission = __bind(this.setPermission, this);
      this.enableNotification = false;
    }

    Notifier.prototype.hasSupport = function() {
      return window.Notification != null;
    };

    Notifier.prototype.requestPermission = function(cb) {
      return window.Notification.requestPermission(cb);
    };

    Notifier.prototype.setPermission = function() {
      if (this.hasPermission()) {
        $('#notification-alert a.close').click();
        return this.enableNotification = true;
      } else if (window.Notification.permission === "granted") {
        return $('#notification-alert a.close').click();
      }
    };

    Notifier.prototype.hasPermission = function() {
      if (window.Notification.permission === "granted") {
        return true;
      } else {
        return false;
      }
    };

    Notifier.prototype.checkOrRequirePermission = function() {
      if (this.hasSupport()) {
        if (this.hasPermission()) {
          return this.enableNotification = true;
        } else {
          if (window.Notification.permission === "default") {
            return this.showTooltip();
          }
        }
      } else {
        return console.log("Desktop notifications are not supported for this Browser/OS version yet.");
      }
    };

    Notifier.prototype.showTooltip = function() {
      console.log("show notifications tip");
      $('.main-container').prepend("<div class='alert alert-info' id='notification-alert'><a href='#' id='link_enable_notifications'>点击这里</a> 开启桌面提醒通知功能。 <a class='close fa fa-close' data-dismiss='alert' href='#'></a></div>");
      $("#notification-alert").alert();
      return $('#notification-alert').on('click', 'a#link_enable_notifications', (function(_this) {
        return function(e) {
          e.preventDefault();
          return _this.requestPermission(_this.setPermission);
        };
      })(this));
    };

    Notifier.prototype.visitUrl = function(url) {
      return window.location.href = url;
    };

    Notifier.prototype.notify = function(avatar, title, content, url) {
      var opts, popup;
      if (url == null) {
        url = null;
      }
      if (this.enableNotification) {
        opts = {
          icon: avatar,
          body: content
        };
        popup = new window.Notification(title, opts);
        return popup.onclick = function() {
          window.parent.focus();
          return $.notifier.visitUrl(url);
        };
      }
    };

    return Notifier;

  })();

  jQuery.notifier = new Notifier;

}).call(this);
