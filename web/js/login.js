// 登录逻辑
const loginForm = document.getElementById("loginForm");
const submitBtn = document.getElementById("submitBtn");

loginForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("password").value;

    if (!username || !password) {
        toast("用户名和密码不能为空", "warning");
        return;
    }

    submitBtn.disabled = true;
    submitBtn.textContent = "登录中...";

    try {
        await request("/api/login", {
            method: "POST",
            body: JSON.stringify({ username, password })
        });

        toast("登录成功，正在进入...", "success");
        setTimeout(() => {
            window.location.href = "/profile";
        }, 350);
    } catch (err) {
        toast(err.message || "登录失败，请检查用户名或密码", "error");
        submitBtn.disabled = false;
        submitBtn.textContent = "登录";
    }
});
