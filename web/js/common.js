// 通用工具模块 (原生 JS，简洁无过度依赖)

/**
 * 封装 Fetch 请求
 */
async function request(url, options = {}) {
    const defaultOptions = {
        headers: {
            "Content-Type": "application/json",
            ...(options.headers || {})
        }
    };

    const finalOptions = { ...defaultOptions, ...options };

    // 如果没有 body，且不是 POST/PUT/PATCH，删除 Content-Type 避免不必要的预检
    if (!finalOptions.body && (finalOptions.method === "GET" || finalOptions.method === "DELETE" || !finalOptions.method)) {
        delete finalOptions.headers["Content-Type"];
    }

    const response = await fetch(url, finalOptions);

    if (!response.ok) {
        let errMsg = `请求失败 (${response.status})`;
        try {
            const data = await response.json();
            if (data && data.error) {
                errMsg = data.error;
            }
        } catch (_) {
            // 响应非 JSON 格式
        }
        const error = new Error(errMsg);
        error.status = response.status;
        throw error;
    }

    return response;
}

/**
 * 极简轻量级 Toast 提示
 */
function toast(message, type = "info", duration = 3000) {
    let container = document.getElementById("toastContainer");
    if (!container) {
        container = document.createElement("div");
        container.id = "toastContainer";
        container.className = "toast-container";
        document.body.appendChild(container);
    }

    const item = document.createElement("div");
    item.className = `toast toast-${type}`;
    
    const msgEl = document.createElement("span");
    msgEl.className = "toast-msg";
    msgEl.textContent = message;

    const closeBtn = document.createElement("button");
    closeBtn.className = "toast-close-btn";
    closeBtn.innerHTML = "&times;";
    closeBtn.onclick = () => item.remove();

    item.appendChild(msgEl);
    item.appendChild(closeBtn);
    container.appendChild(item);

    if (duration > 0) {
        setTimeout(() => {
            if (item.parentNode) {
                item.remove();
            }
        }, duration);
    }
}

// 兼容旧代码引用
window.showToast = function({ title, message, type }) {
    const text = title ? `${title}: ${message}` : message;
    toast(text, type);
};

/**
 * 格式化 ISO 时间字符串
 */
function formatDate(dateStr) {
    if (!dateStr) return "-";
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return dateStr;
    const pad = (n) => String(n).padStart(2, "0");
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

/**
 * 通用模态对话框 Promise
 */
function showModal({ title, bodyHtml, confirmText = "确定", cancelText = "取消", isDanger = false, onConfirm }) {
    return new Promise((resolve) => {
        const modalBackdrop = document.createElement("div");
        modalBackdrop.className = "modal-backdrop";

        modalBackdrop.innerHTML = `
            <div class="modal-card">
                <div class="modal-header">
                    <h3 class="modal-title">${title}</h3>
                    <button class="modal-close" type="button">&times;</button>
                </div>
                <div class="modal-body">${bodyHtml}</div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary modal-cancel-btn">${cancelText}</button>
                    <button type="button" class="btn ${isDanger ? 'btn-danger' : 'btn-primary'} modal-confirm-btn">${confirmText}</button>
                </div>
            </div>
        `;

        document.body.appendChild(modalBackdrop);

        const close = (result = null) => {
            modalBackdrop.remove();
            resolve(result);
        };

        modalBackdrop.querySelector(".modal-close").onclick = () => close(null);
        modalBackdrop.querySelector(".modal-cancel-btn").onclick = () => close(null);
        
        const confirmBtn = modalBackdrop.querySelector(".modal-confirm-btn");
        confirmBtn.onclick = async () => {
            if (onConfirm) {
                const res = await onConfirm(modalBackdrop);
                if (res !== false) {
                    close(res !== undefined ? res : true);
                }
            } else {
                close(true);
            }
        };
    });
}

/**
 * 获取当前登录用户信息（未登录则返回 null）
 */
async function getCurrentUser() {
    try {
        const res = await request("/api/profile");
        return await res.json();
    } catch (_) {
        return null;
    }
}

/**
 * 统一退出登录
 */
async function doLogout() {
    try {
        await request("/api/logout", { method: "POST" });
        toast("已退出登录", "info");
        setTimeout(() => {
            window.location.href = "/";
        }, 300);
    } catch (e) {
        toast("退出失败: " + e.message, "error");
    }
}
