// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0 

$(document).ready(function() {
    $.ajaxSetup({
        cache: false, //prevent browser cache result to redirect  failed.
        headers: { "edgex-club-token": window.localStorage.getItem("edgex-club-token") },
        statusCode: {
            302: function() {
                window.location.href = '/login.html?ran=' + Math.random(); //prevent browser cache result to redirect  failed.
            },
            429: function() {
                //alert("当前操作过于频繁，请稍后再试！");
            },
            401: function() {
                window.location.href = '/login.html?ran=' + Math.random();
            },
            400: function() {
                window.location.href = '/error.html?ran=' + Math.random();
            }

        }
    });
    edgexClubMainModule.loginIsVaild = edgexClubMainModule.checkLogin();
    //edgexClubMainModule.loadArticleList();
    $("div.header_login").on("click", function() {
        mainModuleBtnGroup.login();
    });
    $("div.header_user").on("click", function() {
        mainModuleBtnGroup.user();
    });
    $("div.header_post").on("click", function() {
        mainModuleBtnGroup.post();
    });
    // $("div.header_user").mouseover(function(){
    //
    //   $("div.user_hide_info").show()
    // });
    // $("div.header_user").mouseleave(function(){
    //   $("div.user_hide_info").hide()
    // });

    if (edgexClubMainModule.loginIsVaild) {
        var userInfo = JSON.parse(window.localStorage.getItem("edgex-club-userInfo"));
        $.ajax({
            url: "/api/v1/auth/message/" + userInfo.name + "/count",
            type: "GET",
            success: function(data) {
                if (data == 0) {
                    return;
                } else if (data > 100) {
                    $("div.header_user span.badge").text(data + "+");
                } else {
                    $("div.header_user span.badge").text(data);
                }
                $("div.header_user span.badge").show();
            }
        });
    }
});

var edgexClubMainModule = {
    loginIsVaild: false,
    checkLogin: function() {
        var token = window.localStorage.getItem("edgex-club-token");
        var user = JSON.parse(window.localStorage.getItem("edgex-club-userInfo"));
        var isVaild = false;
        if (token) {
            $.ajax({
                url: "/api/v1/isvalid/" + token,
                type: "GET",
                async: false,
                success: function(data) {
                    if (data == 1) {
                        isVaild = true;
                        $("div.header_login").hide();
                        $("div.header_user img").prop("src", user["avatarUrl"])
                        $("div.header_user").show();
                    } else {
                        isVaild = false;
                    }
                }
            });
            return isVaild;
        } else {
            $("div.header_user").hide();
            return isVaild;
        }
    }
}

var mainModuleBtnGroup = {
    login: function() {
        window.location.href = "/login.html"
    },
    user: function() {
        if (edgexClubMainModule.loginIsVaild) {
            var user = JSON.parse(window.localStorage.getItem("edgex-club-userInfo"));
            window.location.href = "/user/" + user["name"];
        } else {
            window.location.href = "/login.html"
        }
    },
    post: function() {
        if (edgexClubMainModule.loginIsVaild) {
            var user = JSON.parse(window.localStorage.getItem("edgex-club-userInfo"));
            window.location.href = "/api/v1/auth/article/add";
        } else {
            window.location.href = "/login.html"
        }
    }
}