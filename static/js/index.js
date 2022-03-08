// 获取警示框
let div_Warning = $("#panel");
// type参数是Success或Error两种情况,也对应两种图标
let Warning = function (type, txt) {
    div_Warning.css("display", "block"); // 显示警示框
    div_Warning.html(`
        <img src="../static/image/${type}.png" width="30px" alt="图片未加载成功" class="img_error"> 
        ${txt}
    `);
    setTimeout(() => {
        div_Warning.fadeOut(2000);
    }, 1000)
    return false;
}

// 注册button
{

    function getSubmit() {
        let user = $("#user").val();
        let phonenum = $("#ipnum").val();//手机号码
        let pwd = $("#pwd").val();//密码
        let repwd = $("#repwd").val();//确认密码
        let email = $("#email").val(); //邮箱
        let retPhonenum = $("#ipnum1").val();//推荐人手机号码,可选
        let verify = $("#verify").val();//验证码

        // 都是true才能往下走
        if (!(user != "" && func_user(user))) {
            if (user === "") {
                $("#user").focus();
                return Warning("Error", "用户名不能为空");
            }
            return;
        }

        if (!(phonenum != "" && func_phone(phonenum))) {
            if (phonenum === "") {
                $("#ipnum").focus();
                return Warning("Error", "手机号码不能为空");
            }
            return;
        }
        if (!(pwd != "" && repwd != "" && func_passwd(pwd, repwd))) {
            if (pwd === "") {
                $("#pwd").focus();
                return Warning("Error", "密码不能为空");
            } else if (repwd === "") {
                $("#repwd").focus();
                return Warning("Error", "确认密码不能为空");
            }
            return;
        }
        if (!(email != "" && func_email(email))) {
            if (email === "") {
                $("#email").focus();
                return Warning("Error", "邮箱不能为空");
            }
            return;
        }
        if (!(retPhonenum != "" && phonenum != retPhonenum && func_phone(retPhonenum))) {
            if (retPhonenum === "") {
                $("#ipnum1").focus();
                return Warning("Error", "推荐人手机号码不能为空");
            } else if (phonenum === retPhonenum) {
                $("#ipnum1").focus();
                $("#ipnum1").val("");
                return Warning("Error", "推荐人不能是自己");
            }
            return;
        }
        if (!(verify != "" && func_verify(verify, canvas))) {
            if (verify === "") {
                $("#verify").focus();
                return Warning("Error", "验证码不能为空");
            }
            return;
        }

        console.log("success");
        // 密码加密 sha1加密(安全哈希算法)
        let sha1_pwd = sha1(pwd)
        req(user, phonenum, sha1_pwd, email, retPhonenum);
    }
    // ajax接口
    let req = function (user, phonenum, pwd, email, retPhonenum) {
        $.ajax({
            type: 'post',
            url: '/',
            data: {
                "username": user,
                "phonenum": phonenum,
                "password": pwd,
                "email": email,
                "retPhonenum": retPhonenum
            },
            success: function (data) {
                if (data.error != 0) {   // 0为成功
                    Warning("Error", `添加失败!`);
                } else {
                    Warning("Success", "添加成功!");
                }
            },
            error: function (data) {
                console.log(data);
                Warning("Error", `添加失败!`);
            }
        })
    }

    // 判断用户名
    let func_user = function (username) {
        if (username.length >= 6 && username.length <= 12) {
            let reg = /^[\u4E00-\u9FA5A-Za-z0-9_]+$/;
            if (!reg.test($("#user").val())) {
                return Warning("Error", "用户名可以是中文、英文、数字包括下划线");
            }
            return true;
        } else {
            return Warning("Error", "用户名长度为6~12位");
        }
        // $.get(URL,data,function(data,status,xhr),dataType);  
    };

    // 密码验证
    let func_passwd = function (pwd, repwd) {
        if (pwd === repwd && pwd != "" && pwd != undefined) {
            // 密码(以字母开头，长度在6~18之间，只能包含字母、数字和下划线)
            let reg = /^[a-zA-Z]\w{5,17}$/;
            if (!reg.test(pwd)) {
                return Warning("Error", "密码以字母开头,长度在6~18之间,只能包含字母、数字和下划线");
            }
            return true;
        } else {
            return Warning("Error", "密码需要一致");
        }
    }
    // 验证码验证
    let func_verify = function (verify, canvas) {
        if (verify === canvas) {
            return true;
        } else {
            return Warning("Error", "验证码输入错误"); // 报错
        }
    }
    // 验证Email地址
    let func_email = function (email) {
        let reg = /^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/;
        // return reg.test(email);
        if (!reg.test(email)) {
            return Warning("Error", "邮箱地址错误,请重新输入");
        }
        return true;
    }
    //手机号码
    let func_phone = function (phonenum) {
        let reg = /^(13[0-9]|14[5|7]|15[0|1|2|3|4|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9]|19[0|1|2|3|5|6|7|8|9])\d{8}$/;
        if (!reg.test(phonenum)) {
            return Warning("Error", "电话号码格式错误,请重新输入");
        }
        return true;
    }

}
// 眼睛开闭
function eyes(i) {
    // 使用JavaScript获取路径或URL的最后一段
    let thePath = i.src;
    let lastItem = thePath.substring(thePath.lastIndexOf('/') + 1);
    // eye图片的input
    if (lastItem != "eye-open.png") {
        i.parentNode.children[0].type = "text";
        i.src = "../static/image/eye-open.png";
    } else {
        i.parentNode.children[0].type = "password";
        i.src = "../static/image/eye-close.png";
    }
}