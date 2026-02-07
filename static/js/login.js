const usernameInput = document.getElementById("username");
const passwordInput = document.getElementById("password");
const loginForm = document.getElementsByClassName("login-btn");

async function login() {
    const username = usernameInput.value.trim();
    const password = passwordInput.value.trim();
    if (!username || !password) {
        alert("请输入账号和密码");
        return;
    }

    try {
        const response = await fetch("/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ username, password })
        });

        const result = await response.json();
        alert(result.message);
        if (result.token) {
            localStorage.setItem("token", result.token);
            window.location.href = "/home";
        }

    } catch (error) {
        if (error.response) {
            const errMsg = error.response.data.message || "账号或密码错误";
            alert(`登录失败：${errMsg}`);
        } else {
            alert("网络异常，请检查 Gin 服务是否正常运行");
        }
    }
}

loginForm[0].addEventListener("click", function(event) {
    event.preventDefault();
    login();
});