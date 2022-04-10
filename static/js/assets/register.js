$(document).ready(() => {
    // $("#username").val() 
    // $("#phonenumber").val()
    // $("#password").val()
    // $("#repassword").val()
    $("#but_register").click(() => {
        // 发送数据到后台
        $.ajax({
            url: "/user/register",
            type: "POST",
            data: {
                username: $("#username").val(),
                phonenum: $("#phonenumber").val(),
                password: $("#password").val()
            },
            success: (data) => {
                console.log(data)
                if (data.status == "success") {
                    window.location.href = "/login.html"
                } else {
                    alert("Wrong")
                }
            }
        })
    })
})