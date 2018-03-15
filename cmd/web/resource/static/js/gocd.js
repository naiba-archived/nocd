/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

document.ready = function () {
    $(function () {
        $('[data-toggle="tooltip"]').tooltip()
    })
};

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

function addRepo() {
    $.post("/repository/", $("#formAddRepo").serialize(), function () {
        alert("添加成功");
        window.location.reload()
    }).fail(function (jq) {
        alert("错误[" + jq.status + "，" + jq.responseText + "]请重试")
    });
    return false
}

function clipSecret(btn) {
    var clipboard = new ClipboardJS(btn, {
        text: function (e) {
            return e.getAttribute("data-original-title").substr(10)
        }
    });
    clipboard.on('success', function () {
        alert('成功复制到剪贴版');
    });
    clipboard.on('error', function (e) {
        alert('复制失败，浏览器不支持' + e);
    })
}