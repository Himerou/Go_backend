// 注册逻辑
const registerForm = document.getElementById("registerForm");
const submitBtn = document.getElementById("submitBtn");

registerForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("password").value;
    const confirmPassword = document.getElementById("confirmPassword").value;

    if (!username) {
        toast("请输入用户名", "warning");
        return;
    }

    if (!password) {
        toast("请输入密码", "warning");
        return;
    }

    if (password !== confirmPassword) {
        toast("两次输入的密码不一致", "warning");
        return;
    }

    submitBtn.disabled = true;
    submitBtn.textContent = "注册中...";

    try {
        await request("/api/register", {
            method: "POST",
            body: JSON.stringify({ username, password })
        });

        toast("注册成功，正在为您登录...", "success");
        setTimeout(() => {
            window.location.href = "/profile";
        }, 350);
    } catch (err) {
        toast(err.message || "注册失败，请稍后重试", "error");
        submitBtn.disabled = false;
        submitBtn.textContent = "注册账户";
    }
});
