// 获取radio button的值，入参为radio button的name
get_choose = function(radio) {
    var t = document.getElementsByName(radio);
    for(var i = 0; i < t.length; ++i) {
        if(t[i].checked) {
            return i;
        }
    }
    return -1;
}

// 获取一个量表内的所有回答（索引），入参为问题总数
get_form_choose = function(n) {
    var answers = new Array(n);
    answers[0] = get_choose("radio");
    for(var i = 1; i < n; ++i)
        answers[i] = get_choose("radio" + i.toString(10));
    return answers;
}

// 跳转到下一界面
jump_to_next = function (next_page) {
    var prefix = window.location.pathname.match("p([0-9]*)v([0-9]*)")[0];
    window.location.pathname = prefix + next_page;
}

// sas跳转到ess并记录用户所填内容
sas = function () {
    var n = 20;
    var name = document.getElementById("name").value;
    var serial_number = document.getElementById("serial_number").value;
    var gender = get_choose("radioG");
    var age = document.getElementById("age").value;
    var answers = get_form_choose(n);
    var sas_score = 0;
    for(var i = 0; i < n; ++i)
        sas_score += (answers[i] + 1);
    sas_score *= 1.25;
    
    var scale = {};
    scale["PhoneNumber"] = window.location.pathname.match("([0-9]{11})")[0];
    scale["Name"] = name;
    scale["SerialNumber"] = serial_number;
    scale["Gender"] = gender;
    scale["Age"] = parseInt(age);
    scale["SAS"] = sas_score;
    localStorage.setItem("scale", JSON.stringify(scale));  
    
    jump_to_next("/ess.html");
}

// ess跳转到isi并记录用户所填内容
ess = function () {
    var n = 8;
    var answers = get_form_choose(n);
    var ess_score = 0;
    for(var i = 0; i < n; ++i) {
        ess_score += answers[i];
    }

    var scale = JSON.parse(localStorage.getItem("scale"));
    if (scale["PhoneNumber"] == window.location.pathname.match("([0-9]{11})")[0]) {
        scale["ESS"] = ess_score;
        localStorage.setItem("scale", JSON.stringify(scale));  
        jump_to_next("/isi.html");
    } else {
        alert("用户校验错误，请重新登录");
    }
}

// isi跳转到sds并记录用户所填内容
isi = function () {
    var n = 8;
    var answers = get_form_choose(n);
    var isi_score = 0;
    for(var i = 0; i < 3; ++i) {
        if(answers[i] != -1)
            isi_score += answers[3];
    }
    for(var i = 4; i < n; ++i)
        isi_score += answers[i];

    var scale = JSON.parse(localStorage.getItem("scale"));
    if (scale["PhoneNumber"] != window.location.pathname.match("([0-9]{11})")[0]) {
        alert("用户校验错误，请重新登录");
        return;
    }
    scale["ISI"] = isi_score;
    localStorage.setItem("scale", JSON.stringify(scale));  
    jump_to_next("/sds.html");
}

// 页面跳转
sds = function () {
    var n = 20;
    var answers = get_form_choose(n);
    var postive = [0,2,3,6,7,8,9,12,14,18];
    var sds_score = 0;
    for(var i = 0; i < n; ++i) {
        var flag = 0;
        for(var j = 0; j < 10; ++j) {
            if(postive[j] == i) {
                flag = 1;
                break;
            }
        }
        if(flag === 1)
            sds_score += (answers[i] + 1);
        else
            sds_score += (4 - answers[i]);
    }
    sds_score *= 1.25;

    var scale = JSON.parse(localStorage.getItem("scale"));
    if (scale["PhoneNumber"] != window.location.pathname.match("([0-9]{11})")[0]) {
        alert("用户校验错误，请重新登录");
        return;
    }
    scale["SDS"] = sds_score;

    $.ajax({
        type:"POST",
        contentType:"application/json",
        dataType:"text",
        url:"/scale",
        data:JSON.stringify(scale),
        success:function (result) {
            jump_to_next("/upload.html");
        },
        error:function (xhr, textStatus, errorThrown) {
            jump_to_next("/failed_submit.html");
        }
    });

    localStorage.clear();  
}
