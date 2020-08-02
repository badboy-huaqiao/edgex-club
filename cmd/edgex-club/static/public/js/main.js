// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0 

$(document).ready(function() {
    $.ajaxSetup({
        cache: false, //prevent browser cache result to redirect  failed.
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
});