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
        $("#inputHDRepoID").val(button.data('repo'));
        var modal = $("#listEventCheckbox");
        modal.empty();
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
        var setForm = function (button, isAdd) {
            $('#iMEditPlName').val(isAdd ? '' : button.data('name'));
            $('#inputID').val(isAdd ? '' : button.data('id'));
            $('#inputHDRepoID').val(button.data('repo'));
            $('#inputBranch').val(isAdd ? 'master' : button.data('branch'));
            $('#inputShell').text(isAdd ? './nb-deploy' : button.data('shell'));
            if (!isAdd) {
                button.data('events').forEach(function (value) {
                    $('#ic' + value).attr('checked', true)
                });
                $('#slMPlatform').find('>option[value=' + button.data('server') + ']').attr('selected', true);
            } else {
                $('#slMPlatform').find('>option:selected').attr('selected', false);
            }
        };
        if (mth === 'POST') {
            toggle(this, false, 'btn-danger', 'btn-primary', "添加");
            // 重置表单
            setForm(button, true);
        } else {
            if (mth === 'PATCH') {
                toggle(this, false, 'btn-danger', 'btn-primary', "修改")
            } else if (mth === 'DELETE') {
                toggle(this, true, 'btn-primary', 'btn-danger', "删除")
            }
            // 填充表单
            setForm(button, false);
        }
    });

    $('#modalServer').on('show.bs.modal', function (event) {
        const button = $(event.relatedTarget)
        const server = button.parent().data('server')
        const method = button.data('method')
        $('#modalServer').find('form').attr('method', method)
        switch (method) {
            case "post":
                $("#modalServer").find("input,select,textarea").each(function () {
                    $(this).attr('disabled', false)
                });
                $("#btnSubmit").removeClass("btn-danger");
                $("#btnSubmit").removeClass("btn-primary");
                $("#btnSubmit").addClass("btn-secondary");
                $("#btnSubmit").text("添加");
                $("#inputID").val('');
                $("#inputName").val('');
                $("#inputAddress").val('');
                $("#inputPort").val('');
                $("#inputLogin").val('');
                $("#inputLoginType").val('');
                $("#inputPassword").text('');
                break;

            case "delete":
                $("#modalServer").find("input,select,textarea").each(function () {
                    $(this).attr('disabled', true)
                });
                $("#btnSubmit").removeClass("btn-secondary");
                $("#btnSubmit").removeClass("btn-primary");
                $("#btnSubmit").addClass("btn-danger");
                $("#btnSubmit").text("删除");
                $("#inputID").val(server.id);
                $("#inputName").val(server.name);
                $("#inputAddress").val(server.address);
                $("#inputPort").val(server.port);
                $("#inputLogin").val(server.login);
                $("#inputLoginType").val(server.login_type);
                $("#inputPassword").text(server.password);
                break;

            case "patch":
                $("#modalServer").find("input,select,textarea").each(function () {
                    $(this).attr('disabled', false)
                });
                $("#btnSubmit").removeClass("btn-secondary");
                $("#btnSubmit").removeClass("btn-danger");
                $("#btnSubmit").addClass("btn-primary");
                $("#btnSubmit").text("修改");
                $("#inputID").val(server.id);
                $("#inputName").val(server.name);
                $("#inputAddress").val(server.address);
                $("#inputPort").val(server.port);
                $("#inputLogin").val(server.login);
                $("#inputLoginType").val(server.login_type);
                $("#inputPassword").text(server.password);
                break

            default:
                break;
        }
    });

    $('#modalWebhook').on('show.bs.modal', function (event) {
        const modal = $('#modalWebhook')
        const button = $(event.relatedTarget)
        const webhook = button.parent().data('webhook')
        const method = button.data('method')
        $('#modalWebhook').find('form').attr('method', method)
        modal.find("#checkVerifySSL").removeAttr('checked');
        modal.find("#checkPushSuccess").removeAttr('checked');
        modal.find("#checkEnable").removeAttr('checked');
        switch (method) {
            case "post":
                modal.find("input,select,textarea").each(function () {
                    $(this).attr('disabled', false)
                });
                modal.find("#btnSubmit").removeClass("btn-danger");
                modal.find("#btnSubmit").removeClass("btn-primary");
                modal.find("#btnSubmit").addClass("btn-secondary");
                modal.find("#btnSubmit").text("添加");

                modal.find("#inputID").val('');
                modal.find("#inputPipelineID").val(button.data('pipeline'));
                modal.find("#inputURL").val('');
                modal.find("#selectMethod").val('');
                modal.find("#selectType").val('');
                modal.find("#textareaBody").text('');
                break;

            case "delete":
                modal.find("input,select,textarea").each(function () {
                    $(this).attr('disabled', true)
                });
                modal.find("#btnSubmit").removeClass("btn-secondary");
                modal.find("#btnSubmit").removeClass("btn-primary");
                modal.find("#btnSubmit").addClass("btn-danger");
                modal.find("#btnSubmit").text("删除");

                modal.find("#inputID").val(webhook.id);
                modal.find("#inputPipelineID").val(webhook.pipeline_id);
                modal.find("#inputURL").val(webhook.url);
                modal.find("#selectMethod").val(webhook.request_method);
                modal.find("#selectType").val(webhook.request_type);
                modal.find("#textareaBody").text(webhook.request_body);
                if (webhook.verify_ssl !== null) modal.find("#checkVerifySSL").attr('checked', 'checked');
                if (webhook.push_success !== null) modal.find("#checkPushSuccess").attr('checked', 'checked');
                if (webhook.enable !== null) modal.find("#checkEnable").attr('checked', 'checked');
                break;

            case "patch":
                $("#modalWebhook").find("input,select,textarea").each(function () {
                    $(this).attr('disabled', false)
                });
                $("#btnSubmit").removeClass("btn-secondary");
                $("#btnSubmit").removeClass("btn-danger");
                $("#btnSubmit").addClass("btn-primary");
                $("#btnSubmit").text("修改");

                modal.find("#inputID").val(webhook.id);
                modal.find("#inputPipelineID").val(webhook.pipeline_id);
                modal.find("#inputURL").val(webhook.url);
                modal.find("#selectMethod").val(webhook.request_method);
                modal.find("#selectType").val(webhook.request_type);
                modal.find("#textareaBody").text(webhook.request_body);
                if (webhook.verify_ssl) modal.find("#checkVerifySSL").attr('checked', 'checked');
                if (webhook.push_success) modal.find("#checkPushSuccess").attr('checked', 'checked');
                if (webhook.enable) modal.find("#checkEnable").attr('checked', 'checked');
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

function pipelineHandler(mth) {
    return ajaxUtil("/pipeline/", "#formAddPipeline", mth)
}

function serverHandler() {
    return ajaxUtil("/server/", "#formEditServer", $("#formEditServer").attr('method'))
}

function webhookHandler() {
    return ajaxUtil("/pipeline/webhook", "#formWebhook", $("#formWebhook").attr('method'))
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
    if (mth === "delete") {
        $(form).find(':input:disabled').removeAttr('disabled')
    }
    var ajaxUrl = mth === "delete" ? url + "?" + $(form).serialize() : url;
    var ajaxData = mth === "delete" ? {} : $(form).serialize();
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
