/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

document.ready = function () {
    // tooltip
    $(function () {
        $('[data-toggle="tooltip"]').tooltip()
    });
    // clipboard
    var clipboard = new ClipboardJS('.copy-btn', {
        text: function (e) {
            return e.getAttribute('data-text')
        }
    });
    clipboard.on('success', function () {
        alert('成功复制到剪贴版');
    });
    clipboard.on('error', function (e) {
        alert('复制失败，浏览器不支持' + e);
    });
    // modal
    $('#modalAddPipeline').on('show.bs.modal', function (event) {
        var button = $(event.relatedTarget);
        var platform = events[button.data('whatever')];
        $("#inputHDRepoID").val(button.data('repo'));
        var modal = $("#listEventCheckbox");
        modal.empty();
        for (var k in platform) {
            var container = $('<div/>', {class: 'form-check form-check-inline'});
            $('<input/>', {
                name: 'events[]',
                class: 'form-check-input',
                type: 'checkbox',
                id: 'ic' + k,
                value: k
            }).appendTo(container);
            $('<label/>', {class: 'form-check-label', for: 'ic' + k, text: platform[k]}).appendTo(container);
            container.appendTo(modal);
        }
    });
    $('#modalEditServer').on('show.bs.modal', function (event) {
        var button = $(event.relatedTarget);
        var isd = button.data('delete');
        isDelete = isd ? "DELETE" : "PATCH";
        $("#formMethod").val(isDelete);
        if (isd) {
            $("#modalEditServer").find("input").each(function () {
                $(this).attr('disabled', true)
            });
            $("#btnSubmit").text("删除");
            $("#btnSubmit").removeClass("btn-primary");
            $("#btnSubmit").addClass("btn-danger");
        } else {
            $("input").each(function () {
                $(this).attr('disabled', false)
            });
            $("#btnSubmit").text("修改");
            $("#btnSubmit").removeClass("btn-danger");
            $("#btnSubmit").addClass("btn-primary");
        }
        $("#inputID").val(button.data('id'));
        $("#inputName").val(button.data('name'));
        $("#inputAddress").val(button.data('address'));
        $("#inputPort").val(button.data('port'));
        $("#inputLogin").val(button.data('login'));
    });
};

function saveForm(form) {
    $(form).submit()
}

function logout() {
    $.removeCookie("uid", {path: '/'});
    $.removeCookie("token", {path: '/'});
    window.location.href = "/"
}

function addServer() {
    return postForm("/server/", "#formAddServer")
}

function addRepo() {
    return postForm("/repository/", "#formAddRepo")
}

function addPipeline() {
    return postForm("/pipeline/", "#formAddPipeline")
}

function editServer(method) {
    return postForm("/server/", "#formEditServer", method)
}

function postForm(url, form, method) {
    if (method === undefined) {
        method = "POST"
    }
    if (method === "DELETE") {
        $(form).find(':input:disabled').removeAttr('disabled')
    }
    var ajaxUrl = method === "DELETE" ? url + "?" + $(form).serialize() : url;
    var ajaxData = method === "DELETE" ? {} : $(form).serialize();
    $.ajax({
        url: ajaxUrl,
        type: method,    // 此处发送的是POST请求
        data: ajaxData,
        success: function () {
            alert("成功");
            window.location.reload()
        },
        error: function (jq) {
            alert("错误[" + jq.status + "，" + jq.responseText + "]请重试")
        }
    });
    return false
}