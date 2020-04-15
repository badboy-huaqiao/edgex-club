var dateToString = function(num) {
    var date = new Date();
    date.setTime(num);
    var y = date.getFullYear();
    var M = date.getMonth() + 1;
    M = (M < 10) ? ('0' + M) : M;
    var d = date.getDate();
    d = (d < 10) ? ('0' + d) : d;
    var hh = date.getHours();
    hh = (hh < 10) ? ('0' + hh) : hh;
    var mm = date.getMinutes();
    mm = (mm < 10) ? ('0' + mm) : mm;
    var ss = date.getSeconds();
    ss = (ss < 10) ? ('0' + ss) : ss;

    var str = y + '-' + M + '-' + d + ' ' + hh + ':' + mm + ':' + ss
    return str;
}


// 时间转化
function edgexFmtDate(inputTime) {

    var YEAR = 1000 * 60 * 60 * 24 * 365;
    var MONTH = 1000 * 60 * 60 * 24 * 30;
    var DAY = 1000 * 60 * 60 * 24;
    var HOUR = 1000 * 60 * 60;
    var MINUTE = 1000 * 60;

    var date = new Date(inputTime);
    var now = new Date();
    var between = now.getTime() - date.getTime();
    if (between > YEAR) {
        return parseInt((between - YEAR) / YEAR + 1) + "年前";
    }
    if (between > MONTH) {
        return parseInt((between - MONTH) / MONTH + 1) + "月前";
    }
    if (between > DAY) {
        return parseInt((between - DAY) / DAY + 1) + "天前";
    }
    if (between > HOUR) {
        return parseInt((between - HOUR) / HOUR + 1) + "小时前";
    }
    if (between > MINUTE) {
        return parseInt((between - MINUTE) / MINUTE + 1) + "分钟前";
    }
    return "刚刚";
}