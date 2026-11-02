
(() => {
    // 创建并显示 Toast
    function showToast({
                           title = "提示",
                           message = "",
                           type = "warning",
                           duration = 4500
                       } = {}) {
        const toast = document.createElement("div");

        const allowedTypes = ["warning", "success", "error", "info"];
        const toastType = allowedTypes.includes(type) ? type : "info";

        toast.className = `toast toast-${toastType}`;
        toast.setAttribute("role", "status");
        toast.setAttribute("aria-live", "polite");

        // 使用 textContent 展示文本，避免把提示内容当成 HTML
        const icon = document.createElement("div");
        icon.className = "toast-icon";
        icon.textContent = toastType === "success"
            ? "✓"
            : toastType === "error"
                ? "×"
                : toastType === "warning"
                    ? "!"
                    : "i";

        const content = document.createElement("div");
        content.className = "toast-content";

        const titleElement = document.createElement("div");
        titleElement.className = "toast-title";
        titleElement.textContent = title;

        const messageElement = document.createElement("div");
        messageElement.className = "toast-message";
        messageElement.textContent = message;

        content.append(titleElement, messageElement);

        const closeButton = document.createElement("button");
        closeButton.className = "toast-close";
        closeButton.type = "button";
        closeButton.setAttribute("aria-label", "关闭提示");
        closeButton.textContent = "×";

        const progress = document.createElement("div");
        progress.className = "toast-progress";

        toast.append(icon, content, closeButton, progress);
        document.body.appendChild(toast);

        // 下一帧添加显示状态，触发进入动画
        requestAnimationFrame(() => {
            toast.classList.add("toast-visible");
        });

        let dismissed = false;
        let timer = null;
        let remaining = duration;
        let startTime = 0;

        function startTimer() {
            if (dismissed || timer !== null) return;

            startTime = Date.now();
            timer = setTimeout(dismiss, remaining);
            progress.style.animationPlayState = "running";
        }

        function pauseTimer() {
            if (dismissed || timer === null) return;

            clearTimeout(timer);
            timer = null;
            remaining -= Date.now() - startTime;
            progress.style.animationPlayState = "paused";
        }

        function dismiss() {
            if (dismissed) return;

            dismissed = true;
            clearTimeout(timer);

            toast.classList.remove("toast-visible");
            toast.classList.add("toast-hiding");

            toast.addEventListener("transitionend", () => {
                toast.remove();
            }, { once: true });

            // 动画未触发时也确保清理
            setTimeout(() => toast.remove(), 500);
        }

        closeButton.addEventListener("click", dismiss);
        toast.addEventListener("mouseenter", pauseTimer);
        toast.addEventListener("mouseleave", startTimer);

        startTimer();
    }

    // 对外暴露接口，供其他 JS 文件调用
    window.showToast = showToast;

    // 保留原来的未登录提示功能
    const url = new URL(window.location.href);

    if (url.searchParams.get("error") === "unauthorized") {
        // 清除提示参数，避免刷新时重复弹出
        url.searchParams.delete("error");

        window.history.replaceState(
            {},
            "",
            url.pathname + url.search + url.hash
        );

        function showUnauthorizedToast() {
            showToast({
                title: "需要登录",
                message: "你尚未登录，请尝试登录或注册使用服务",
                type: "warning"
            });
        }

        if (document.readyState === "loading") {
            document.addEventListener(
                "DOMContentLoaded",
                showUnauthorizedToast,
                { once: true }
            );
        } else {
            showUnauthorizedToast();
        }
    }
})();
