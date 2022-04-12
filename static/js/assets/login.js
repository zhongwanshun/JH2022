$(document).ready(()=> {
    $("#but_login").click(()=> {
         $.ajax({
                url: "/user/login",
                type: "POST",
                data: {
                    username: $("#username").val(),
                    password: $("#password").val()
                },
                success: (data)=> {
                        console.log(data)
                        if(data.status == "success"){
                            console.log(data.info)
                            window.location.href = "/"
                        }
                        else{
                            alert("Wrong username or password")
                        }
                }
            })
    })
})