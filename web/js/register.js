
const registerForm = document.querySelector("#registerForm");
const result = document.querySelector("#result");

registerForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const username = document.querySelector("#username").value.trim();
    const password = document.querySelector("#password").value;

    result.textContent = "正在注册……";

    try {
        await request("/api/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        });

        result.textContent = "注册成功！正在跳转登录页……";
        window.location.href = "/profile";
    } catch (error) {
        console.error("注册错误：", error);
        result.textContent = error.message;
    }
});
