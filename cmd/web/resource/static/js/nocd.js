/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
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
            let extPrefix = '';
            if (e.getAttribute('data-ext') == 'protocol_prefix') {
                extPrefix = window.location.protocol + '//';
            }
            return extPrefix + e.getAttribute('data-text')
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
        const mth = button.data('method');
        $(this).find('form').attr('method', mth)
        $(this).find('#inputHDRepoID').val(button.data('repo'));
        var eventCheckBox = $(this).find('#listEventCheckbox');
        eventCheckBox.empty();
        for (var k in platform) {
            var container = $('<div/>', { class: 'form-check form-check-inline' });
            $('<input/>', {
                name: 'events[]',
                class: 'form-check-input',
                type: 'checkbox',
                id: 'ic' + k,
                value: k
            }).appendTo(container);
            $('<label/>', { class: 'form-check-label', for: 'ic' + k, text: platform[k] }).appendTo(container);
            container.appendTo(eventCheckBox);
        }
        const modal = $(this)
        // 是否是添加
        var toggle = function (md, togg, clazzold, clazznew, btn) {
            modal.find('#btnEditRepo').removeClass(clazzold);
            modal.find('#btnEditRepo').addClass(clazznew);
            $(md).find('input').attr('disabled', togg);
            $(md).find('select').attr('disabled', togg);
            $(md).find('textarea').attr('disabled', togg);
            modal.find('#btnEditRepo').text(btn);
        };
        var setForm = function (button, isAdd) {
            modal.find('#iMEditPlName').val(isAdd ? '' : button.data('name'));
            modal.find('#inputID').val(isAdd ? '' : button.data('id'));
            modal.find('#inputHDRepoID').val(button.data('repo'));
            modal.find('#inputBranch').val(isAdd ? 'master' : button.data('branch'));
            modal.find('#inputShell').text(isAdd ? './nb-deploy' : button.data('shell'));
            if (!isAdd) {
                button.data('events').forEach(function (value) {
                    modal.find('#ic' + value).attr('checked', true)
                });
                modal.find('#slMPlatform').find('>option[value=' + button.data('server') + ']').attr('selected', true);
            } else {
                modal.find('#slMPlatform').find('>option:selected').attr('selected', false);
            }
        };
        if (mth === 'PATCH') {
            toggle(this, false, 'btn-danger', 'btn-primary', "修改")
            setForm(button, false);
        } else if (mth === 'DELETE') {
            toggle(this, true, 'btn-primary', 'btn-danger', "删除")
            setForm(button, false);
        } else {
            // 填充表单
            setForm(button, true);
        }
    });

    $('#modalServer').on('show.bs.modal', function (event) {
        const button = $(event.relatedTarget)
        const server = button.parent().data('server')
        const method = button.data('method')
        $(this).find('form').attr('method', method)
        switch (method) {
            case "POST":
                $(this).find("input,select,textarea").each(function () {
                    $(this).attr('disabled', false)
                });
                $(this).find("#btnSubmit").removeClass("btn-danger");
                $(this).find("#btnSubmit").removeClass("btn-primary");
                $(this).find("#btnSubmit").addClass("btn-secondary");
                $(this).find("#btnSubmit").text("添加");
                $(this).find("#inputID").val('');
                $(this).find("#inputName").val('');
                $(this).find("#inputAddress").val('');
                $(this).find("#inputPort").val('');
                $(this).find("#inputLogin").val('');
                $(this).find("#inputLoginType").val('');
                $(this).find("#inputPassword").text('');
                break;

            case "DELETE":
                $(this).find("input,select,textarea").each(function () {
                    $(this).attr('disabled', true)
                });
                $(this).find("#btnSubmit").removeClass("btn-secondary");
                $(this).find("#btnSubmit").removeClass("btn-primary");
                $(this).find("#btnSubmit").addClass("btn-danger");
                $(this).find("#btnSubmit").text("删除");
                $(this).find("#inputID").val(server.id);
                $(this).find("#inputName").val(server.name);
                $(this).find("#inputAddress").val(server.address);
                $(this).find("#inputPort").val(server.port);
                $(this).find("#inputLogin").val(server.login);
                $(this).find("#inputLoginType").val(server.login_type);
                $(this).find("#inputPassword").text(server.password);
                break;

            case "PATCH":
                $(this).find("input,select,textarea").each(function () {
                    $(this).attr('disabled', false)
                });
                $(this).find("#btnSubmit").removeClass("btn-secondary");
                $(this).find("#btnSubmit").removeClass("btn-danger");
                $(this).find("#btnSubmit").addClass("btn-primary");
                $(this).find("#btnSubmit").text("修改");
                $(this).find("#inputID").val(server.id);
                $(this).find("#inputName").val(server.name);
                $(this).find("#inputAddress").val(server.address);
                $(this).find("#inputPort").val(server.port);
                $(this).find("#inputLogin").val(server.login);
                $(this).find("#inputLoginType").val(server.login_type);
                $(this).find("#inputPassword").text(server.password);
                break

            default:
                break;
        }
    });

    $('#modalWebhook').on('show.bs.modal', function (event) {
        const button = $(event.relatedTarget)
        const webhook = button.parent().data('webhook')
        const method = button.data('method')
        $(this).find('form').attr('method', method)
        $(this).find("#checkVerifySSL").removeAttr('checked');
        $(this).find("#checkPushSuccess").removeAttr('checked');
        $(this).find("#checkEnable").removeAttr('checked');
        switch (method) {
            case "POST":
                $(this).find("input,select,textarea").each(function () {
                    $(this).attr('disabled', false)
                });
                $(this).find("#btnSubmit").removeClass("btn-danger");
                $(this).find("#btnSubmit").removeClass("btn-primary");
                $(this).find("#btnSubmit").addClass("btn-secondary");
                $(this).find("#btnSubmit").text("添加");

                $(this).find("#inputID").val('');
                $(this).find("#inputPipelineID").val(button.data('pipeline'));
                $(this).find("#inputURL").val('');
                $(this).find("#selectMethod").val('');
                $(this).find("#selectType").val('');
                $(this).find("#textareaBody").text('');
                break;

            case "DELETE":
                $(this).find("input,select,textarea").each(function () {
                    $(this).attr('disabled', true)
                });
                $(this).find("#btnSubmit").removeClass("btn-secondary");
                $(this).find("#btnSubmit").removeClass("btn-primary");
                $(this).find("#btnSubmit").addClass("btn-danger");
                $(this).find("#btnSubmit").text("删除");

                $(this).find("#inputID").val(webhook.id);
                $(this).find("#inputPipelineID").val(webhook.pipeline_id);
                $(this).find("#inputURL").val(webhook.url);
                $(this).find("#selectMethod").val(webhook.request_method);
                $(this).find("#selectType").val(webhook.request_type);
                $(this).find("#textareaBody").text(webhook.request_body);
                if (webhook.verify_ssl !== null) $(this).find("#checkVerifySSL").attr('checked', 'checked');
                if (webhook.push_success !== null) $(this).find("#checkPushSuccess").attr('checked', 'checked');
                if (webhook.enable !== null) $(this).find("#checkEnable").attr('checked', 'checked');
                break;

            case "PATCH":
                $(this).find("input,select,textarea").each(function () {
                    $(this).attr('disabled', false)
                });
                $(this).find("#btnSubmit").removeClass("btn-secondary");
                $(this).find("#btnSubmit").removeClass("btn-danger");
                $(this).find("#btnSubmit").addClass("btn-primary");
                $(this).find("#btnSubmit").text("修改");

                $(this).find("#inputID").val(webhook.id);
                $(this).find("#inputPipelineID").val(webhook.pipeline_id);
                $(this).find("#inputURL").val(webhook.url);
                $(this).find("#selectMethod").val(webhook.request_method);
                $(this).find("#selectType").val(webhook.request_type);
                $(this).find("#textareaBody").text(webhook.request_body);
                if (webhook.verify_ssl) $(this).find("#checkVerifySSL").attr('checked', 'checked');
                if (webhook.push_success) $(this).find("#checkPushSuccess").attr('checked', 'checked');
                if (webhook.enable) $(this).find("#checkEnable").attr('checked', 'checked');
                break

            default:
                break;
        }
    });

    $('#modalRepo').on('show.bs.modal', function (event) {
        var button = $(event.relatedTarget);
        $(this).find("#inputID").val(button.data('id'));
        $(this).find("#inputName").val(button.data('name'));
        $(this).find("#selectPlatform>option[value=" + button.data('platform') + "]").attr('selected', true);
    });
    // 部署日志展示
    var running = $('#console').attr("running");
    var consoleText = $('#console').attr("text");
    var logLineNumber = 0;
    if (consoleText && consoleText.length > 0) {
        parseLog(consoleText, 0);
    }
    if (running === "true") {
        var timerHandler = setInterval(function () {
            $.get('?ajax=1&act=view&line=' + logLineNumber, function (res) {
                if (res.log && res.log.length > 0) {
                    parseLog(res.log, logLineNumber);
                }
                if (res.end === "false") {
                    logLineNumber = res.line;
                } else {
                    clearInterval(timerHandler)
                }
            })
        }, 3000)
    }
};

function saveForm(form) {
    $(form).submit()
}

function stopDeploy(fromAdmin) {
    if (confirm("确定停止部署吗？")) {
        $.get('?ajax=1&act=stop', function (res) {
            if (res === "success") {
                alert("成功停止部署");
                window.location.href = fromAdmin ? "/admin/running/" : "/pipelog/";
            } else {
                alert(res)
            }
        })
    }
}

function logout() {
    $.removeCookie("uid", { path: '/' });
    $.removeCookie("token", { path: '/' });
    window.location.href = "/"
}

function parseLog(str, number) {
    var text = "";
    str.split("\n").forEach(function (value, index) {
        text += '<div class="row"><span class="line-number col-1  text-center">' + (index + parseInt(number) + 1) + '.</span><span class="code col-11"\n' +
            ' data-toggle="tooltip"\n' +
            ' data-placement="top"\n' +
            ' title="' + value.substr(0, 8) + '">' + $(document.createElement('code')).text(value.substr(9)).prop("outerHTML") + '</span>\n' +
            '</div>';
    });
    $('#console').append(text)
}

function addRepoHandler(mth) {
    return ajaxUtil("/repository/", "#formAddRepo", mth)
}

function editRepoHandler(mth) {
    return ajaxUtil("/repository/", "#formEditRepo", mth)
}

function accountTransferHandler() {
    return ajaxUtil("/user/transfer", "#formAccountTransfer", "POST")
}

function pipelineHandler() {
    return ajaxUtil("/pipeline/", "#formAddPipeline", $('#formAddPipeline').attr('method'))
}

function serverHandler() {
    return ajaxUtil("/server/", "#formEditServer", $('#formEditServer').attr('method'))
}

function webhookHandler() {
    return ajaxUtil("/pipeline/webhook", "#formWebhook", $('#formWebhook').attr('method'))
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
