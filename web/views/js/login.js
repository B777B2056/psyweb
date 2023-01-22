// 登录按钮响应事件，测试是否满足登录条件
test_condition = function (phone) {
    var ret_val = 0
    var verification_code = document.getElementById("password").value;
    if (verification_code === "") {
        alert("请填写验证码！");
        return 0;
    }

    $.ajax({
        type:"POST",
        contentType:"application/json",
        dataType:"text",
        url:"login/verification",
        data:JSON.stringify({ 'PhoneNumber': phone, 'Code': verification_code }),
        success:function (result) {
            var j = JSON.parse(result)
            if (j['UserStatus'] === 0) {
                alert("用户不存在，请检查是否发送了验证码");
            } else {
                if (j['VerificationResult']) {
                    ret_val = 1;
                    // 页面跳转
                    id = "p" + phone + "v" + verification_code;
                    if (j['UserStatus'] === 1)
                        window.location.href = id + "/sas.html";     // 新用户，未提交过任何数据
                    else if (j['UserStatus'] === 2)
                        window.location.href = id + "/wait.html";    // 诊断未完成，需等待
                    else
                        window.location.href = id + "/index.html";   // 诊断已完成
                } else {
                    ret_val = 0;
                    alert("验证码错误！");
                }
            }
        },
        error:function (xhr, textStatus, errorThrown) {
           window.location.href = "/failed_submit.html";
        }
    });

    return ret_val;
}

// 发送验证码
send_verification_code = function(phone) {
    if (phone === "") {
        alert("请输入手机号！");
        return 0;
    }

    count_down = function() {
        var button = document.getElementById("verification_code_btn");
        button.disabled = true;
        var time = 60;
        var timer = setInterval(function() {
            button.innerHTML = time + "秒后重新获取";
            time--;
            if (time == 0) {
                clearInterval(timer);
                button.disabled = false;
                button.innerHTML = "获取验证码";
            }
        }, 1000);
    }
    
    $.ajax({
        type:"POST",
        contentType:"application/json",
        dataType:"text",
        url:"login/get_verification_code",
        data:JSON.stringify({ 'PhoneNumber': phone }),
        success:function (result) {
            count_down();
        },
        error:function (xhr, textStatus, errorThrown) {
           window.location.href = "/failed_submit.html";
        }
    });
    return 1;
}
