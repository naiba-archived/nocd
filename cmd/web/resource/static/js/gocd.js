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

function addServer() {
    $.post("/server/", $("#formAddServer").serialize(), function () {
        alert("添加成功");
        window.location.reload()
    }).fail(function (jq) {
        alert("错误[" + jq.status + "，" + jq.responseText + "]请重试")
    });
    return false
}