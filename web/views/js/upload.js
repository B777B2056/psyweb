check_file_name = function(filePath) {
    var fileName = filePath.substring(filePath.lastIndexOf('\\')+1);
    var name = fileName.substring(0, fileName.lastIndexOf('.')); 
    return (/1\d{10}/.test(name));
}

check_file_ext = function(filePath, targetExt) {
    var fileName = filePath.substring(filePath.lastIndexOf('\\')+1);
    var extName = fileName.substring(fileName.lastIndexOf('.')+1); 
    return extName === targetExt;
}

function upload() {
    set_file = document.getElementById("set_file");
    fdt_file = document.getElementById("fdt_file");
    // 校验文件数量
    if (set_file.files.length < 1) {
        alert("未选择.set文件！")
        return; 
    }
    if (fdt_file.files.length < 1) {
        alert("未选择.fdt文件！")
        return; 
    }
    if (set_file.files.length > 1) {
        alert(".set文件最多只能上传一个！")
        return; 
    }
    if (fdt_file.files.length > 1) {
        alert(".fdt文件最多只能上传一个！")
        return; 
    }
    // 校验文件后缀
    if (!check_file_ext(set_file.files[0].name, "set")) {
        alert(".set文件类型错误！")
        return;
    }
    if (!check_file_ext(fdt_file.files[0].name, "fdt")) {
        alert(".fdt文件类型错误！")
        return; 
    }
    document.getElementById("ups").submit();
}