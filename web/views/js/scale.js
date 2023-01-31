// 获取radio button的值，入参为radio button的name
function get_choose(radio) {
    var t = document.getElementsByName(radio);
    for(var i = 0; i < t.length; ++i) {
        if(t[i].checked) {
            return i;
        }
    }
    return -1;
}

// 获取一个量表内的所有回答（索引），入参为问题总数
function get_form_choose(n) {
    var answers = new Array(n);
    answers[0] = get_choose("radio");
    for(var i = 1; i < n; ++i)
        answers[i] = get_choose("radio" + i.toString(10));
    return answers;
}

// sas跳转到ess并记录用户所填内容
function sas () {
    var n = 20;
    var name = document.getElementById("name").value;
    var serial_number = document.getElementById("serial_number").value;
    var gender = get_choose("radio20");
    var age = document.getElementById("age").value;
    var answers = get_form_choose(n);
    var sas_score = 0;
    for(var i = 0; i < n; ++i)
        sas_score += (answers[i] + 1);
    sas_score *= 1.25;
    
    var scale = {};
    scale["Name"] = name;
    scale["SerialNumber"] = serial_number;
    scale["Gender"] = gender;
    scale["Age"] = parseInt(age);
    scale["SAS"] = sas_score;
    localStorage.setItem("scale", JSON.stringify(scale));  
}

// ess跳转到isi并记录用户所填内容
function ess () {
    var n = 8;
    var answers = get_form_choose(n);
    var ess_score = 0;
    for(var i = 0; i < n; ++i) {
        ess_score += answers[i];
    }

    var scale = JSON.parse(localStorage.getItem("scale"));
    scale["ESS"] = ess_score;
    localStorage.setItem("scale", JSON.stringify(scale));  
}

// isi跳转到sds并记录用户所填内容
function isi () {
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
    scale["ISI"] = isi_score;
    localStorage.setItem("scale", JSON.stringify(scale));  
}

// 页面跳转
function sds () {
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
    scale["SDS"] = sds_score;

    $.ajax({
        type:"POST",
        contentType:"application/json",
        dataType:"text",
        url:"/scale",
        data:JSON.stringify(scale),
        success:function (result) {
            localStorage.clear();  
        },
        error:function (xhr, textStatus, errorThrown) {
            var err = eval("(" + xhr.responseText + ")");
            alert(err.Message);
        }
    });
}
