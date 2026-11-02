const params = new URLSearchParams(window.location.search);

if (params.get("error") === "unauthorized") {
    alert("你尚未登录，请尝试登录或注册使用服务");

    // 清除 URL 中的提示参数，避免刷新时重复弹窗
    window.history.replaceState({}, "", "/");
}