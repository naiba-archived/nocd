/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

function login() {
    $("#loginForm").submit()
}

function logout() {
    $.removeCookie("uid", {path: '/'});
    $.removeCookie("token", {path: '/'});
    window.location.href = "/"
}