// 获取随机验证码
var canvas;

$("#VerificationCode").click(() => {
    // 清空验证码
    canvas = "";
    let ctx = $("#canvas")[0].getContext("2d");
    // 清空画布
    ctx.clearRect(0, 0, 400, 400);
    // 设置字体
    ctx.font = "128px bold 黑体";
    // 设置垂直对齐方式
    ctx.textBaseline = "top";
    // 设置颜色
    ctx.fillStyle = randomColor();
    // 绘制文字（参数：要写的字，x坐标，y坐标）
    ctx.fillText(n(getRandomNum(10)), 0, getRandomNum(50));
    ctx.fillStyle = randomColor();
    ctx.fillText(n(getRandomNum(10)), 50, getRandomNum(50));
    ctx.fillStyle = randomColor();
    ctx.fillText(n(getRandomNum(10)), 100, getRandomNum(50));
    ctx.fillStyle = randomColor();
    ctx.fillText(n(getRandomNum(10)), 150, getRandomNum(50));

    $("#VerificationCode").val("看不清楚,换一张")
    // 
    $("#verify").focus();
})
function randomColor() {
    let colorValue = "0,1,2,3,4,5,6,7,8,9,a,b,c,d,e,f";
    let colorArray = colorValue.split(",");
    let color = "#";
    for (let i = 0; i < 6; i++) {
        color += colorArray[Math.floor(Math.random() * 16)];
    }
    return color;
}
function getRandomNum(n) {
    return function () {
        return parseInt(Math.random() * n);
    }()
}
// 获取js生成的验证码
let n = function (num) {
    canvas += num;
    return num;
}