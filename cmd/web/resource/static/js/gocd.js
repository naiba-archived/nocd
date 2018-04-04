/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

window.onload = function () {
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
        // 是否是添加
        mth = button.data('method');
        var toggle = function (md, togg, clazzold, clazznew, btn) {
            $('#btnEditRepo').removeClass(clazzold);
            $('#btnEditRepo').addClass(clazznew);
            $(md).find('input').attr('disabled', togg);
            $(md).find('select').attr('disabled', togg);
            $(md).find('textarea').attr('disabled', togg);
            $('#btnEditRepo').text(btn);
        };
        if (mth === 'POST') {
            toggle(this, false, 'btn-danger', 'btn-primary', "添加")
        } else {
            if (mth === 'PATCH') {
                toggle(this, false, 'btn-danger', 'btn-primary', "修改")
            } else if (mth === 'DELETE') {
                toggle(this, true, 'btn-primary', 'btn-danger', "删除")
            }
            // 填充表单
            $('#iMEditPlName').val(button.data('name'));
            $('#inputID').val(button.data('id'));
            $('#inputHDRepoID').val(button.data('repo'));
            $('#inputBranch').val(button.data('branch'));
            $('#inputShell').text(button.data('shell'));
            button.data('events').forEach(function (value) {
                $('#ic' + value).attr('checked', true)
            });
            $('#slMPlatform').find('>option[value=' + button.data('server') + ']').attr('selected', true);
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
    $('#modalRepo').on('show.bs.modal', function (event) {
        var button = $(event.relatedTarget);
        $(this).find("#inputID").val(button.data('id'));
        $(this).find("#inputName").val(button.data('name'));
        $(this).find("#selectPlatform>option[value=" + button.data('platform') + "]").attr('selected', true);
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

function addServerHandler() {
    return ajaxUtil("/server/", "#formAddServer", 'POST')
}

function setNotification() {
    return ajaxUtil("/settings/", "#formNotify", 'POST')
}

function addRepoHandler(mth) {
    return ajaxUtil("/repository/", "#formAddRepo", mth)
}

function editRepoHandler(mth) {
    return ajaxUtil("/repository/", "#formEditRepo", mth)
}

function pipelineHandler(mth) {
    return ajaxUtil("/pipeline/", "#formAddPipeline", mth)
}

function editServerHandler(mth) {
    return ajaxUtil("/server/", "#formEditServer", mth)
}

function toggleUser(uid, toggle, col) {
    $.get("/admin/user/" + uid + "/" + col + "/" + (toggle ? "off" : "on"), function () {
        alert("操作成功");
        window.location.reload()
    }).fail(function (x) {
        alert(x.responseText)
    })
}

function ajaxUtil(url, form, mth) {
    if (mth === "DELETE") {
        $(form).find(':input:disabled').removeAttr('disabled')
    }
    var ajaxUrl = mth === "DELETE" ? url + "?" + $(form).serialize() : url;
    var ajaxData = mth === "DELETE" ? {} : $(form).serialize();
    $.ajax({
        url: ajaxUrl,
        type: mth,
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