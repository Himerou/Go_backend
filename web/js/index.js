// 首页逻辑
document.addEventListener("DOMContentLoaded", async () => {
    // 1. 处理从受限路由重定向过来的未登录提示
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.get("error") === "unauthorized") {
        toast("请先登录账号后再访问该功能", "warning");
        window.history.replaceState({}, "", "/");
    }

    // 2. 检查当前用户登录状态
    const currentUser = await getCurrentUser();
    initNav(currentUser);

    if (currentUser) {
        initUserWorkspace(currentUser);
    } else {
        initGuestView();
    }
});

/**
 * 初始化导航栏
 */
function initNav(user) {
    const navLinks = document.getElementById("navAuthLinks");
    if (!navLinks) return;

    if (user) {
        navLinks.innerHTML = `
            <span class="nav-user">你好，<strong>${escapeHtml(user.username)}</strong></span>
            <a href="/profile" class="btn btn-sm btn-outline">个人中心</a>
            <button type="button" class="btn btn-sm btn-secondary" id="navLogoutBtn">退出</button>
        `;
        document.getElementById("navLogoutBtn").onclick = doLogout;
    } else {
        navLinks.innerHTML = `
            <a href="/login" class="btn btn-sm btn-outline">登录</a>
            <a href="/register" class="btn btn-sm btn-primary">注册</a>
        `;
    }
}

/**
 * 初始化未登录访客界面
 */
function initGuestView() {
    const guestSection = document.getElementById("guestSection");
    if (guestSection) guestSection.hidden = false;

    // 公开资料查询
    const publicSearchForm = document.getElementById("publicSearchForm");
    const searchUsername = document.getElementById("searchUsername");
    const publicSearchResult = document.getElementById("publicSearchResult");

    if (publicSearchForm) {
        publicSearchForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            const uname = searchUsername.value.trim();
            if (!uname) return;

            publicSearchResult.hidden = false;
            publicSearchResult.innerHTML = `<p style="color:var(--text-muted);font-size:0.9rem;">正在查询...</p>`;

            try {
                const res = await request(`/api/users/${encodeURIComponent(uname)}/profile`);
                const data = await res.json();

                publicSearchResult.innerHTML = `
                    <div style="background:var(--bg-page); border:1px solid var(--border); border-radius:var(--radius-md); padding:1rem;">
                        <div style="font-weight:600; font-size:1.05rem; margin-bottom:0.5rem; color:var(--text-main);">
                            ${escapeHtml(data.username)}
                            <span class="badge badge-secondary" style="margin-left:0.5rem;">UID: ${data.id}</span>
                        </div>
                        <div style="font-size:0.9rem; color:var(--text-muted); margin-bottom:0.4rem;">
                            注册时间：${formatDate(data.created_at)}
                        </div>
                        <div style="font-size:0.92rem; color:var(--text-main); margin-top:0.5rem; white-space:pre-wrap;">
                            ${escapeHtml(data.profile || "（该用户暂未填写个人简介）")}
                        </div>
                    </div>
                `;
            } catch (err) {
                publicSearchResult.innerHTML = `
                    <div style="color:var(--danger); font-size:0.9rem; padding:0.5rem 0;">
                        未找到用户 "${escapeHtml(uname)}" 的公开资料
                    </div>
                `;
            }
        });
    }
}

/**
 * 初始化已登录工作台
 */
let allTodos = [];
let currentFilter = "all";

function initUserWorkspace(user) {
    const userSection = document.getElementById("userSection");
    if (userSection) userSection.hidden = false;

    // 渲染快捷资料
    const currentUserDisplay = document.getElementById("currentUserDisplay");
    const currentUserCreatedAt = document.getElementById("currentUserCreatedAt");
    if (currentUserDisplay) currentUserDisplay.textContent = user.username;
    if (currentUserCreatedAt) currentUserCreatedAt.textContent = formatDate(user.created_at);

    const quickLogoutBtn = document.getElementById("quickLogoutBtn");
    if (quickLogoutBtn) quickLogoutBtn.onclick = doLogout;

    // 筛选 Tab 切换
    const tabBtns = document.querySelectorAll(".filter-tabs .tab-btn");
    tabBtns.forEach(btn => {
        btn.addEventListener("click", () => {
            tabBtns.forEach(b => b.classList.remove("active"));
            btn.classList.add("active");
            currentFilter = btn.dataset.filter;
            renderTodoList();
        });
    });

    // 刷新按钮
    const refreshBtn = document.getElementById("refreshTodosBtn");
    if (refreshBtn) refreshBtn.onclick = fetchTodos;

    // 创建 Todo 表单
    const createForm = document.getElementById("createTodoForm");
    const titleInput = document.getElementById("newTodoTitle");
    const contentInput = document.getElementById("newTodoContent");
    const createBtn = document.getElementById("createTodoBtn");

    createForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const title = titleInput.value.trim();
        const content = contentInput.value.trim();

        if (!title) {
            toast("任务标题不能为空", "warning");
            return;
        }

        createBtn.disabled = true;
        createBtn.textContent = "创建中...";

        try {
            await request("/api/todos", {
                method: "POST",
                body: JSON.stringify({
                    title: title,
                    content: content,
                    completed: false
                })
            });

            toast("任务创建成功", "success");
            titleInput.value = "";
            contentInput.value = "";
            fetchTodos();
        } catch (err) {
            toast(err.message || "创建任务失败", "error");
        } finally {
            createBtn.disabled = false;
            createBtn.textContent = "创建待办";
        }
    });

    // 初次加载 Todo 列表
    fetchTodos();
}

/**
 * 获取 Todo 列表
 */
async function fetchTodos() {
    const listEl = document.getElementById("todoList");
    const emptyEl = document.getElementById("todoEmpty");

    try {
        const res = await request("/api/todos");
        allTodos = await res.json() || [];
        updateStats();
        renderTodoList();
    } catch (err) {
        toast("获取任务失败: " + err.message, "error");
    }
}

/**
 * 更新统计指标
 */
function updateStats() {
    const total = allTodos.length;
    const completed = allTodos.filter(t => t.completed).length;
    const active = total - completed;

    const elTotal = document.getElementById("statTotal");
    const elActive = document.getElementById("statActive");
    const elCompleted = document.getElementById("statCompleted");

    if (elTotal) elTotal.textContent = total;
    if (elActive) elActive.textContent = active;
    if (elCompleted) elCompleted.textContent = completed;
}

/**
 * 渲染 Todo 列表
 */
function renderTodoList() {
    const listEl = document.getElementById("todoList");
    const emptyEl = document.getElementById("todoEmpty");
    if (!listEl) return;

    listEl.innerHTML = "";

    let filtered = allTodos;
    if (currentFilter === "active") {
        filtered = allTodos.filter(t => !t.completed);
    } else if (currentFilter === "completed") {
        filtered = allTodos.filter(t => t.completed);
    }

    if (filtered.length === 0) {
        if (emptyEl) emptyEl.hidden = false;
        return;
    }
    if (emptyEl) emptyEl.hidden = true;

    filtered.forEach(todo => {
        const li = document.createElement("li");
        li.className = `todo-item ${todo.completed ? 'completed' : ''}`;
        li.dataset.id = todo.id;

        li.innerHTML = `
            <input type="checkbox" class="todo-checkbox" ${todo.completed ? 'checked' : ''} title="切换状态">
            <div class="todo-body">
                <div class="todo-title">${escapeHtml(todo.title)}</div>
                ${todo.content ? `<div class="todo-desc">${escapeHtml(todo.content)}</div>` : ''}
                <div class="todo-meta">
                    <span>${todo.completed ? '<span class="badge badge-success">已完成</span>' : '<span class="badge badge-warning">未完成</span>'}</span>
                    <span>更新于: ${formatDate(todo.updated_at || todo.created_at)}</span>
                </div>
            </div>
            <div class="todo-actions">
                <button type="button" class="btn btn-sm btn-outline edit-btn">编辑</button>
                <button type="button" class="btn btn-sm btn-danger-outline delete-btn">删除</button>
            </div>
        `;

        // 复选框事件：直接切换完成状态
        const checkbox = li.querySelector(".todo-checkbox");
        checkbox.addEventListener("change", async () => {
            await toggleTodoStatus(todo, checkbox.checked);
        });

        // 编辑按钮事件
        li.querySelector(".edit-btn").addEventListener("click", () => {
            openEditTodoModal(todo);
        });

        // 删除按钮事件
        li.querySelector(".delete-btn").addEventListener("click", () => {
            confirmDeleteTodo(todo);
        });

        listEl.appendChild(li);
    });
}

/**
 * 快速切换任务完成状态
 */
async function toggleTodoStatus(todo, completed) {
    try {
        await request(`/api/todos/${todo.id}`, {
            method: "PUT",
            body: JSON.stringify({
                title: todo.title,
                content: todo.content,
                completed: completed
            })
        });

        todo.completed = completed;
        updateStats();
        renderTodoList();
        toast(completed ? "任务标记为已完成" : "任务标记为未完成", "info", 1500);
    } catch (err) {
        toast("更新状态失败: " + err.message, "error");
        fetchTodos(); // 还原状态
    }
}

/**
 * 模态框编辑任务
 */
function openEditTodoModal(todo) {
    const bodyHtml = `
        <div class="form-group">
            <label class="form-label" for="editTodoTitle">任务标题</label>
            <input class="form-input" type="text" id="editTodoTitle" value="${escapeHtml(todo.title)}" required>
        </div>
        <div class="form-group">
            <label class="form-label" for="editTodoContent">任务描述</label>
            <textarea class="form-textarea" id="editTodoContent" rows="4">${escapeHtml(todo.content || "")}</textarea>
        </div>
        <div class="form-group" style="display:flex; align-items:center; gap:0.5rem; margin-top:0.75rem;">
            <input type="checkbox" id="editTodoCompleted" ${todo.completed ? 'checked' : ''} style="width:1.1rem; height:1.1rem;">
            <label for="editTodoCompleted" style="cursor:pointer; font-size:0.9rem;">标记为已完成</label>
        </div>
    `;

    showModal({
        title: "编辑待办任务",
        bodyHtml: bodyHtml,
        confirmText: "保存更改",
        onConfirm: async (modalEl) => {
            const newTitle = modalEl.querySelector("#editTodoTitle").value.trim();
            const newContent = modalEl.querySelector("#editTodoContent").value.trim();
            const newCompleted = modalEl.querySelector("#editTodoCompleted").checked;

            if (!newTitle) {
                toast("任务标题不能为空", "warning");
                return false;
            }

            try {
                await request(`/api/todos/${todo.id}`, {
                    method: "PUT",
                    body: JSON.stringify({
                        title: newTitle,
                        content: newContent,
                        completed: newCompleted
                    })
                });
                toast("任务已更新", "success");
                fetchTodos();
                return true;
            } catch (err) {
                toast("保存失败: " + err.message, "error");
                return false;
            }
        }
    });
}

/**
 * 确认删除任务
 */
function confirmDeleteTodo(todo) {
    showModal({
        title: "删除待办任务",
        bodyHtml: `<p>确定要删除任务「<strong>${escapeHtml(todo.title)}</strong>」吗？此操作无法撤销。</p>`,
        confirmText: "确认删除",
        isDanger: true,
        onConfirm: async () => {
            try {
                await request(`/api/todos/${todo.id}`, {
                    method: "DELETE"
                });
                toast("任务已删除", "success");
                fetchTodos();
                return true;
            } catch (err) {
                toast("删除失败: " + err.message, "error");
                return false;
            }
        }
    });
}

/**
 * HTML 转义
 */
function escapeHtml(str) {
    if (!str) return "";
    return String(str)
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}
