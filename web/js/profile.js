// ==================== 个人资料页 ====================
const profileLoading = document.querySelector("#profileLoading");
const profileContent = document.querySelector("#profileContent");
const profileError = document.querySelector("#profileError");
const profileMessage = document.querySelector("#profileMessage");

const profileForm = document.querySelector("#profileForm");
const profileText = document.querySelector("#profileText");
const profileDisplay = document.querySelector("#profileDisplay");
const editProfileBtn = document.querySelector("#editProfileBtn");
const cancelProfileBtn = document.querySelector("#cancelProfileBtn");

const profileId = document.querySelector("#profileId");
const profileUsername = document.querySelector("#profileUsername");
const profileCreatedAt = document.querySelector("#profileCreatedAt");

const editUsernameBtn = document.querySelector("#editUsernameBtn");
const editPasswordBtn = document.querySelector("#editPasswordBtn");
const logoutBtn = document.querySelector("#logoutBtn");

let currentProfile = "";

// 优先使用公共 Toast；如果 Toast 尚未加载，则退回页面内提示区域。
function notify(type, title, message) {
    if (typeof window.showToast === "function") {
        window.showToast({ title, message, type });
        return;
    }

    console.warn("Toast 组件不可用，使用页面内提示：", title, message);

    if (type === "error") {
        profileError.textContent = `${title}：${message}`;
        profileError.hidden = false;
        profileMessage.hidden = true;
    } else {
        profileMessage.textContent = `${title}：${message}`;
        profileMessage.hidden = false;
        profileError.hidden = true;
    }
}

// 清除页面内提示信息
function clearMessage() {
    profileError.hidden = true;
    profileMessage.hidden = true;
}

// 展示简介：为空时显示浅色提示，不展示空白文本框。
function renderProfile(profile) {
    currentProfile = profile || "";

    if (currentProfile.trim() === "") {
        profileDisplay.className = "profile-empty";
        profileDisplay.textContent = "这个人很懒，什么都没写";
    } else {
        profileDisplay.className = "profile-text-display";
        profileDisplay.textContent = currentProfile;
    }
}

// 加载当前登录用户的个人资料
async function loadProfile() {
    profileLoading.hidden = false;
    profileContent.hidden = true;
    clearMessage();

    try {
        const response = await request("/api/profile");
        const data = await response.json();

        profileId.textContent = data.id;
        profileUsername.textContent = data.username;
        profileText.value = data.profile || "";

        // 根据简介内容更新展示区域
        renderProfile(data.profile);

        profileCreatedAt.textContent = data.created_at
            ? new Date(data.created_at).toLocaleString()
            : "暂无";

        // 加载完成后回到只读展示状态
        profileForm.hidden = true;
        profileDisplay.hidden = false;
        editProfileBtn.hidden = false;
        profileContent.hidden = false;
    } catch (error) {
        console.error("获取个人资料失败：", error);
        notify("error", "获取个人资料失败", error.message || "请稍后重试");
    } finally {
        profileLoading.hidden = true;
    }
}

// 点击“编辑简介”后才显示文本框
editProfileBtn.addEventListener("click", () => {
    profileText.value = currentProfile;

    profileDisplay.hidden = true;
    editProfileBtn.hidden = true;
    profileForm.hidden = false;

    profileText.focus();
});

// 取消编辑，不向后端发送请求
cancelProfileBtn.addEventListener("click", () => {
    profileText.value = currentProfile;

    profileForm.hidden = true;
    profileDisplay.hidden = false;
    editProfileBtn.hidden = false;
});

// 修改个人简介：只绑定一次 submit 事件，避免重复请求
profileForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const profile = profileText.value.trim();
    const saveBtn = profileForm.querySelector('button[type="submit"]');

    clearMessage();
    saveBtn.disabled = true;

    try {
        await request("/api/profile", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                msg: profile
            })
        });

        // 保存成功后更新页面展示
        renderProfile(profile);

        profileForm.hidden = true;
        profileDisplay.hidden = false;
        editProfileBtn.hidden = false;

        notify("success", "保存成功", "个人简介已更新");
    } catch (error) {
        console.error("修改个人简介失败：", error);
        notify("error", "保存失败", error.message || "请稍后重试");
    } finally {
        saveBtn.disabled = false;
    }
});

// 修改用户名
editUsernameBtn.addEventListener("click", async () => {
    const newUsername = prompt(
        "请输入新的用户名：",
        profileUsername.textContent
    );

    if (newUsername === null) {
        return;
    }

    const username = newUsername.trim();

    if (!username) {
        notify("error", "修改失败", "用户名不能为空");
        return;
    }

    if (username === profileUsername.textContent) {
        notify("info", "无需修改", "用户名没有变化");
        return;
    }

    clearMessage();
    editUsernameBtn.disabled = true;

    try {
        await request("/api/username", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                msg: username
            })
        });

        // 重新获取资料，避免页面显示与后端不一致
        await loadProfile();

        notify("success", "修改成功", "用户名已更新");
    } catch (error) {
        console.error("修改用户名失败：", error);
        notify("error", "修改用户名失败", error.message || "请稍后重试");
    } finally {
        editUsernameBtn.disabled = false;
    }
});

// 修改密码：操作状态使用 Toast，不使用 alert
editPasswordBtn.addEventListener("click", async () => {
    const newPassword = prompt("请输入新密码：");

    if (newPassword === null) {
        return;
    }

    if (!newPassword.trim()) {
        notify("error", "修改失败", "密码不能为空");
        return;
    }

    const confirmPassword = prompt("请再次输入新密码：");

    if (confirmPassword === null) {
        return;
    }

    if (newPassword !== confirmPassword) {
        notify("error", "修改失败", "两次输入的密码不一致");
        return;
    }

    editPasswordBtn.disabled = true;

    try {
        await request("/api/password", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                msg: newPassword
            })
        });

        notify("success", "修改成功", "密码已更新");
    } catch (error) {
        console.error("修改密码失败：", error);
        notify("error", "修改密码失败", error.message || "请稍后重试");
    } finally {
        editPasswordBtn.disabled = false;
    }
});

// 退出登录：保留操作前确认，结果提示使用 Toast
logoutBtn.addEventListener("click", async () => {
    if (!confirm("确定要退出登录吗？")) {
        return;
    }

    logoutBtn.disabled = true;

    try {
        await request("/api/logout", {
            method: "POST"
        });

        notify("success", "退出成功", "即将返回首页");
        window.location.href = "/";
    } catch (error) {
        console.error("退出登录失败：", error);
        notify("error", "退出登录失败", error.message || "请稍后重试");
        logoutBtn.disabled = false;
    }
});

// 页面加载后获取资料
loadProfile();

// ==================== Todo 列表 ====================

const deleteDialog = document.querySelector("#deleteDialog");
const deleteMessage = document.querySelector("#deleteMessage");
const cancelDeleteBtn = document.querySelector("#cancelDeleteBtn");
const confirmDeleteBtn = document.querySelector("#confirmDeleteBtn");

// 暂存当前准备删除的任务 ID
let pendingDeleteId = null;

async function loadTodos() {
    const todoList = document.querySelector("#todoList");

    try {
        const response = await request("/api/todos");
        const todos = await response.json();

        todoList.replaceChildren();

        if (todos === null || todos.length === 0) {
            todoList.textContent = "目前没有任务";
            return;
        }

        for (const todo of todos) {
            const li = document.createElement("li");
            li.className = "todo-item";
            li.dataset.id = todo.id;

            const title = document.createElement("h3");
            title.textContent = todo.title;

            const content = document.createElement("p");
            content.textContent = todo.content || "暂无详细内容";

            const status = document.createElement("p");
            status.textContent = todo.completed
                ? "状态：已完成"
                : "状态：未完成";

            const createdAt = document.createElement("p");
            createdAt.textContent =
                `创建时间：${new Date(todo.created_at).toLocaleString()}`;

            const updatedAt = document.createElement("p");
            updatedAt.textContent =
                `更新时间：${new Date(todo.updated_at).toLocaleString()}`;

            // 创建操作按钮容器
            const actions = document.createElement("div");
            actions.className = "todo-actions";

            const editBtn = document.createElement("button");
            editBtn.type = "button";
            editBtn.textContent = "编辑";
            editBtn.className = "edit-btn";

            editBtn.addEventListener("click", () => {
                // 创建标题输入框
                const titleInput = document.createElement("input");
                titleInput.type = "text";
                titleInput.value = todo.title;
                titleInput.name = "title";
                titleInput.required = true;

                // 创建内容输入框
                const contentInput = document.createElement("textarea");
                contentInput.value = todo.content || "";
                contentInput.name = "content";

                // 创建完成状态复选框
                const completedInput = document.createElement("input");
                completedInput.type = "checkbox";
                completedInput.checked = todo.completed;
                completedInput.name = "completed";

                // 创建保存按钮
                const saveBtn = document.createElement("button");
                saveBtn.type = "button";
                saveBtn.textContent = "保存";

                // 创建取消按钮
                const cancelBtn = document.createElement("button");
                cancelBtn.type = "button";
                cancelBtn.textContent = "取消";

                // 替换原来的显示内容
                title.replaceWith(titleInput);
                content.replaceWith(contentInput);
                status.replaceWith(completedInput);

                // 隐藏编辑按钮，并添加保存、取消按钮
                editBtn.hidden = true;
                actions.append(saveBtn, cancelBtn);

                // 点击取消：恢复原来的显示
                cancelBtn.addEventListener("click", () => {
                    titleInput.replaceWith(title);
                    contentInput.replaceWith(content);
                    completedInput.replaceWith(status);

                    saveBtn.remove();
                    cancelBtn.remove();
                    editBtn.hidden = false;
                });

                // 保存修改
                saveBtn.addEventListener("click", async () => {
                    const newTitle = titleInput.value.trim();
                    const newContent = contentInput.value.trim();
                    const newCompleted = completedInput.checked;

                    if (!newTitle) {
                        alert("任务标题不能为空");
                        return;
                    }

                    saveBtn.disabled = true;

                    try {
                        await request(`/api/todos/${todo.id}`, {
                            method: "PUT",
                            headers: {
                                "Content-Type": "application/json"
                            },
                            body: JSON.stringify({
                                title: newTitle,
                                content: newContent,
                                completed: newCompleted
                            })
                        });

                        // 保存成功后重新加载任务列表
                        await loadTodos();
                    } catch (error) {
                        console.error("更新 Todo 失败：", error);
                        alert(`保存失败：${error.message}`);
                        saveBtn.disabled = false;
                    }
                });
            });

            const deleteBtn = document.createElement("button");
            deleteBtn.type = "button";
            deleteBtn.textContent = "删除";
            deleteBtn.className = "delete-btn";

            deleteBtn.addEventListener("click", () => {
                // 记录准备删除的任务
                pendingDeleteId = todo.id;

                // 显示具体任务名称
                deleteMessage.textContent =
                    `确定要删除任务「${todo.title}」吗？此操作无法撤销。`;

                // 显示弹窗
                deleteDialog.hidden = false;
            });

            actions.append(editBtn, deleteBtn);
            li.append(title, content, status, createdAt, updatedAt, actions);
            todoList.appendChild(li);
        }
    } catch (error) {
        console.error("获取 Todo 失败：", error);
        todoList.textContent = error.message;
    }
}

// 点击按钮重新加载任务
document.querySelector("#loadBtn").addEventListener("click", loadTodos);

// 创建 Todo
const todoForm = document.querySelector("#todoForm");
const todoResult = document.querySelector("#todoResult");

todoForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const title = document.querySelector("#todoTitle").value.trim();
    const content = document.querySelector("#todoContent").value.trim();

    todoResult.textContent = "正在创建任务……";

    try {
        await request("/api/todos", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                title: title,
                content: content
            })
        });

        todoResult.textContent = "创建成功！";
        todoForm.reset();

        await loadTodos();
    } catch (error) {
        console.error("创建任务失败：", error);
        todoResult.textContent = `创建失败：${error.message}`;
    }
});

// 页面加载后获取任务列表
loadTodos();

// 取消删除
cancelDeleteBtn.addEventListener("click", () => {
    // 关闭弹窗，不删除任务
    deleteDialog.hidden = true;
    pendingDeleteId = null;
});

// 确认删除
confirmDeleteBtn.addEventListener("click", async () => {
    // 没有选中任务，就不执行删除
    if (pendingDeleteId === null) {
        return;
    }

    confirmDeleteBtn.disabled = true;

    try {
        await request(`/api/todos/${pendingDeleteId}`, {
            method: "DELETE"
        });

        // 删除成功，关闭弹窗并刷新列表
        deleteDialog.hidden = true;
        pendingDeleteId = null;

        await loadTodos();
    } catch (error) {
        console.error("删除 Todo 失败：", error);
        alert(`删除失败：${error.message}`);
    } finally {
        confirmDeleteBtn.disabled = false;
    }
});