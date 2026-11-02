
const loginForm = document.querySelector("#loginForm");
const result = document.querySelector("#result");

loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const username = document.querySelector("#username").value.trim();
    const password = document.querySelector("#password").value;

    result.textContent = "正在登录……";

    try {
        const response = await request("/api/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        });

        result.textContent = "登录成功！";
        window.location.href = "/profile";
    } catch (error) {
        console.error("登录错误：", error);
        result.textContent = error.message;
    }
});
