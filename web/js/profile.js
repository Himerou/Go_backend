// 个人中心逻辑
let profileData = null;
let profileTodos = [];
let profileCurrentFilter = "all";

document.addEventListener("DOMContentLoaded", async () => {
    // 绑定退出登录
    const logoutBtn = document.getElementById("logoutBtn");
    if (logoutBtn) logoutBtn.onclick = doLogout;

    // 加载个人资料
    await loadUserProfile();

    // 绑定修改用户名
    initEditUsername();

    // 绑定修改简介
    initEditBio();

    // 绑定修改密码
    initEditPassword();

    // 绑定注销账户
    initDeleteAccount();

    // 绑定查询他人公开资料
    initPublicSearch();

    // 绑定个人中心 Todo 任务管理
    initProfileTodos();
});

/**
 * 加载个人资料
 */
async function loadUserProfile() {
    const loadingEl = document.getElementById("profileLoading");
    const contentEl = document.getElementById("profileContent");

    try {
        const res = await request("/api/profile");
        profileData = await res.json();

        document.getElementById("profileId").textContent = profileData.id;
        document.getElementById("profileUsername").textContent = profileData.username;
        document.getElementById("profileCreatedAt").textContent = formatDate(profileData.created_at);

        renderBioDisplay(profileData.profile);

        loadingEl.hidden = true;
        contentEl.hidden = false;
    } catch (err) {
        toast("获取个人资料失败: " + err.message, "error");
        if (err.status === 401) {
            setTimeout(() => {
                window.location.href = "/?error=unauthorized";
            }, 800);
        }
    }
}

/**
 * 渲染个人简介展示
 */
function renderBioDisplay(bio) {
    const displayEl = document.getElementById("profileDisplay");
    if (!bio || !bio.trim()) {
        displayEl.innerHTML = '<span class="profile-empty-text">这个人很懒，什么都没留下...</span>';
    } else {
        displayEl.textContent = bio;
    }
}

/**
 * 修改个人简介
 */
function initEditBio() {
    const openBtn = document.getElementById("openEditBioBtn");
    const bioForm = document.getElementById("bioForm");
    const bioInput = document.getElementById("bioInput");
    const cancelBtn = document.getElementById("cancelBioBtn");
    const displayEl = document.getElementById("profileDisplay");
    const saveBtn = document.getElementById("saveBioBtn");

    openBtn.onclick = () => {
        bioInput.value = profileData.profile || "";
        displayEl.hidden = true;
        openBtn.hidden = true;
        bioForm.hidden = false;
        bioInput.focus();
    };

    cancelBtn.onclick = () => {
        bioForm.hidden = true;
        displayEl.hidden = false;
        openBtn.hidden = false;
    };

    bioForm.onsubmit = async (e) => {
        e.preventDefault();
        const newBio = bioInput.value.trim();

        saveBtn.disabled = true;
        saveBtn.textContent = "保存中...";

        try {
            const res = await request("/api/profile", {
                method: "PUT",
                body: JSON.stringify({ msg: newBio })
            });
            const updated = await res.json();
            profileData.profile = updated.profile;
            renderBioDisplay(profileData.profile);

            bioForm.hidden = true;
            displayEl.hidden = false;
            openBtn.hidden = false;
            toast("个人简介已更新", "success");
        } catch (err) {
            toast("更新简介失败: " + err.message, "error");
        } finally {
            saveBtn.disabled = false;
            saveBtn.textContent = "保存简介";
        }
    };
}

/**
 * 修改用户名
 */
function initEditUsername() {
    const btn = document.getElementById("editUsernameBtn");
    btn.onclick = () => {
        showModal({
            title: "修改用户名",
            bodyHtml: `
                <div class="form-group">
                    <label class="form-label" for="newUsernameInput">新用户名</label>
                    <input class="form-input" type="text" id="newUsernameInput" value="${escapeHtml(profileData.username)}" required>
                    <div style="font-size:0.8rem; color:var(--text-muted); margin-top:0.4rem;">用户名用于登录，不可与其他用户重复。</div>
                </div>
            `,
            confirmText: "确认修改",
            onConfirm: async (modalEl) => {
                const newUsername = modalEl.querySelector("#newUsernameInput").value.trim();
                if (!newUsername) {
                    toast("用户名不能为空", "warning");
                    return false;
                }
                if (newUsername === profileData.username) {
                    toast("新用户名与当前用户名相同", "info");
                    return false;
                }

                try {
                    const res = await request("/api/username", {
                        method: "PUT",
                        body: JSON.stringify({ msg: newUsername })
                    });
                    const updated = await res.json();
                    profileData.username = updated.username;
                    document.getElementById("profileUsername").textContent = updated.username;
                    toast("用户名修改成功", "success");
                    return true;
                } catch (err) {
                    toast("修改失败: " + err.message, "error");
                    return false;
                }
            }
        });
    };
}

/**
 * 修改密码
 */
function initEditPassword() {
    const btn = document.getElementById("editPasswordBtn");
    btn.onclick = () => {
        showModal({
            title: "修改密码",
            bodyHtml: `
                <div class="form-group">
                    <label class="form-label" for="newPwdInput">新密码</label>
                    <input class="form-input" type="password" id="newPwdInput" placeholder="请输入新密码" required autocomplete="new-password">
                </div>
                <div class="form-group">
                    <label class="form-label" for="confirmPwdInput">确认新密码</label>
                    <input class="form-input" type="password" id="confirmPwdInput" placeholder="请再次输入新密码" required autocomplete="new-password">
                </div>
            `,
            confirmText: "确认修改",
            onConfirm: async (modalEl) => {
                const newPwd = modalEl.querySelector("#newPwdInput").value;
                const confirmPwd = modalEl.querySelector("#confirmPwdInput").value;

                if (!newPwd) {
                    toast("新密码不能为空", "warning");
                    return false;
                }
                if (newPwd !== confirmPwd) {
                    toast("两次输入的密码不一致", "warning");
                    return false;
                }

                try {
                    await request("/api/password", {
                        method: "PUT",
                        body: JSON.stringify({ msg: newPwd })
                    });
                    toast("密码修改成功，请牢记您的新密码", "success");
                    return true;
                } catch (err) {
                    toast("修改密码失败: " + err.message, "error");
                    return false;
                }
            }
        });
    };
}

/**
 * 注销账户
 */
function initDeleteAccount() {
    const btn = document.getElementById("deleteAccountBtn");
    btn.onclick = () => {
        showModal({
            title: "注销账户确认",
            bodyHtml: `
                <p style="color:var(--danger); font-weight:600; margin-bottom:0.5rem;">此操作不可逆！</p>
                <p>注销后，您的账号以及该账号下的所有待办任务将全部被永久删除，无法找回。</p>
            `,
            confirmText: "确认注销",
            cancelText: "考虑一下",
            isDanger: true,
            onConfirm: async () => {
                try {
                    await request("/api/user", { method: "DELETE" });
                    toast("账户已成功注销", "info");
                    setTimeout(() => {
                        window.location.href = "/";
                    }, 500);
                    return true;
                } catch (err) {
                    toast("注销失败: " + err.message, "error");
                    return false;
                }
            }
        });
    };
}

/**
 * 查询他人公开资料
 */
function initPublicSearch() {
    const form = document.getElementById("publicSearchForm");
    const input = document.getElementById("searchUsername");
    const resultBox = document.getElementById("publicSearchResult");

    form.onsubmit = async (e) => {
        e.preventDefault();
        const name = input.value.trim();
        if (!name) return;

        resultBox.hidden = false;
        resultBox.innerHTML = `<p style="color:var(--text-muted);font-size:0.9rem;">正在查询...</p>`;

        try {
            const res = await request(`/api/users/${encodeURIComponent(name)}/profile`);
            const data = await res.json();
            resultBox.innerHTML = `
                <div style="background:var(--bg-page); border:1px solid var(--border); border-radius:var(--radius-md); padding:0.85rem;">
                    <div style="font-weight:600; font-size:0.95rem; margin-bottom:0.25rem;">
                        ${escapeHtml(data.username)}
                        <span class="badge badge-secondary" style="margin-left:0.4rem;">UID: ${data.id}</span>
                    </div>
                    <div style="font-size:0.82rem; color:var(--text-muted); margin-bottom:0.4rem;">
                        注册时间：${formatDate(data.created_at)}
                    </div>
                    <div style="font-size:0.88rem; color:var(--text-main); white-space:pre-wrap;">
                        ${escapeHtml(data.profile || "（该用户暂未填写个人简介）")}
                    </div>
                </div>
            `;
        } catch (_) {
            resultBox.innerHTML = `
                <div style="color:var(--danger); font-size:0.88rem;">
                    未找到用户 "${escapeHtml(name)}"
                </div>
            `;
        }
    };
}

/**
 * 个人中心 Todo 任务管理
 */
function initProfileTodos() {
    // 筛选 Tab
    const tabBtns = document.querySelectorAll(".filter-tabs [data-pfilter]");
    tabBtns.forEach(btn => {
        btn.addEventListener("click", () => {
            tabBtns.forEach(b => b.classList.remove("active"));
            btn.classList.add("active");
            profileCurrentFilter = btn.dataset.pfilter;
            renderProfileTodoList();
        });
    });

    // 刷新
    const refreshBtn = document.getElementById("profileRefreshTodosBtn");
    if (refreshBtn) refreshBtn.onclick = fetchProfileTodos;

    // 创建任务
    const form = document.getElementById("profileTodoForm");
    const titleIn = document.getElementById("profileTodoTitle");
    const contentIn = document.getElementById("profileTodoContent");
    const btn = document.getElementById("profileCreateTodoBtn");

    form.onsubmit = async (e) => {
        e.preventDefault();
        const title = titleIn.value.trim();
        const content = contentIn.value.trim();

        if (!title) {
            toast("任务标题不能为空", "warning");
            return;
        }

        btn.disabled = true;
        btn.textContent = "创建中...";

        try {
            await request("/api/todos", {
                method: "POST",
                body: JSON.stringify({ title, content, completed: false })
            });
            toast("任务创建成功", "success");
            titleIn.value = "";
            contentIn.value = "";
            fetchProfileTodos();
        } catch (err) {
            toast(err.message || "创建任务失败", "error");
        } finally {
            btn.disabled = false;
            btn.textContent = "立即创建";
        }
    };

    fetchProfileTodos();
}

/**
 * 获取 Todo 列表
 */
async function fetchProfileTodos() {
    try {
        const res = await request("/api/todos");
        profileTodos = await res.json() || [];
        updateProfileStats();
        renderProfileTodoList();
    } catch (err) {
        toast("获取任务失败: " + err.message, "error");
    }
}

/**
 * 更新统计指标
 */
function updateProfileStats() {
    const total = profileTodos.length;
    const completed = profileTodos.filter(t => t.completed).length;
    const active = total - completed;

    const elTotal = document.getElementById("pStatTotal");
    const elActive = document.getElementById("pStatActive");
    const elCompleted = document.getElementById("pStatCompleted");

    if (elTotal) elTotal.textContent = total;
    if (elActive) elActive.textContent = active;
    if (elCompleted) elCompleted.textContent = completed;
}

/**
 * 渲染 Todo 列表
 */
function renderProfileTodoList() {
    const listEl = document.getElementById("profileTodoList");
    const emptyEl = document.getElementById("profileTodoEmpty");
    if (!listEl) return;

    listEl.innerHTML = "";

    let filtered = profileTodos;
    if (profileCurrentFilter === "active") {
        filtered = profileTodos.filter(t => !t.completed);
    } else if (profileCurrentFilter === "completed") {
        filtered = profileTodos.filter(t => t.completed);
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

        // 切换完成状态
        const cb = li.querySelector(".todo-checkbox");
        cb.onchange = async () => {
            try {
                await request(`/api/todos/${todo.id}`, {
                    method: "PUT",
                    body: JSON.stringify({
                        title: todo.title,
                        content: todo.content,
                        completed: cb.checked
                    })
                });
                todo.completed = cb.checked;
                updateProfileStats();
                renderProfileTodoList();
                toast(cb.checked ? "任务已完成" : "任务恢复为未完成", "info", 1500);
            } catch (err) {
                toast("更新失败: " + err.message, "error");
                fetchProfileTodos();
            }
        };

        // 编辑
        li.querySelector(".edit-btn").onclick = () => {
            showModal({
                title: "编辑任务",
                bodyHtml: `
                    <div class="form-group">
                        <label class="form-label" for="modalTodoTitle">任务标题</label>
                        <input class="form-input" type="text" id="modalTodoTitle" value="${escapeHtml(todo.title)}" required>
                    </div>
                    <div class="form-group">
                        <label class="form-label" for="modalTodoContent">任务描述</label>
                        <textarea class="form-textarea" id="modalTodoContent" rows="3">${escapeHtml(todo.content || "")}</textarea>
                    </div>
                    <div class="form-group" style="display:flex; align-items:center; gap:0.5rem; margin-top:0.75rem;">
                        <input type="checkbox" id="modalTodoCompleted" ${todo.completed ? 'checked' : ''} style="width:1.1rem; height:1.1rem;">
                        <label for="modalTodoCompleted" style="cursor:pointer; font-size:0.9rem;">标记为已完成</label>
                    </div>
                `,
                confirmText: "保存更改",
                onConfirm: async (modalEl) => {
                    const newTitle = modalEl.querySelector("#modalTodoTitle").value.trim();
                    const newContent = modalEl.querySelector("#modalTodoContent").value.trim();
                    const newCompleted = modalEl.querySelector("#modalTodoCompleted").checked;

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
                        toast("任务更新成功", "success");
                        fetchProfileTodos();
                        return true;
                    } catch (err) {
                        toast("更新失败: " + err.message, "error");
                        return false;
                    }
                }
            });
        };

        // 删除
        li.querySelector(".delete-btn").onclick = () => {
            showModal({
                title: "删除任务",
                bodyHtml: `<p>确定要删除任务「<strong>${escapeHtml(todo.title)}</strong>」吗？</p>`,
                confirmText: "确认删除",
                isDanger: true,
                onConfirm: async () => {
                    try {
                        await request(`/api/todos/${todo.id}`, { method: "DELETE" });
                        toast("任务已删除", "success");
                        fetchProfileTodos();
                        return true;
                    } catch (err) {
                        toast("删除失败: " + err.message, "error");
                        return false;
                    }
                }
            });
        };

        listEl.appendChild(li);
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
