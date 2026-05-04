# 中转设备 REST、身份验证与自动化集成指南

> **相关文档**：[05 单中转设备多组织同步配置指南](./05单中转设备多组织同步配置指南.md) | [02 REST API 参考](./02REST_API参考.md) | [04 中转设备方案详解](./04中转设备方案详解.md)

本文说明：设备二维码含义、系统设置里 GUI/API 认证与 REST 的关系、配置变更是否需重启、如何用 REST 替代手改 `config.xml` 做多组织扩展，以及「终端加入中转设备 C」与「小程序扫码」时应调用的接口与推荐架构。

---

## 1. 设备二维码是做什么的

Syncthing 在每台设备上生成的**设备二维码**，编码的是**该设备自身的设备 ID**（`XXXXXXX-…` 那一长串），与菜单「操作 → 显示 ID」中的字符串一致。

- **用途**：把本机设备 ID **方便地交给对方**，用于在另一台 Syncthing 上「添加远程设备」时填写，避免手抄。
- **不能单独实现**：扫一次码不会让多台机器自动组成全网；同步仍依赖各端**互相添加设备**、**共享并接受文件夹**等配置。
- **与组织/微信无关**：二维码里**没有**部门、openid、文件夹名；若业务上要绑定微信用户，需在**自有后端**完成映射。

---

## 2. 系统设置中的三项与 REST 请求的关系

在 **设置 → 图形用户界面** 中常见三项：

| 配置项 | 作用 | 与 `/rest/...` 的关系 |
|--------|------|------------------------|
| **API 密钥** | 供脚本、服务端、第三方客户端调用 REST | 请求头携带 `X-API-Key: <密钥>` 或 `Authorization: Bearer <密钥>` 即可访问受保护接口（与 GUI 用户名密码无关）。 |
| **GUI 身份验证用户 / 密码** | 保护浏览器打开的 Web 界面 | 浏览器登录或 HTTP Basic；**若请求已带有效 API 密钥，调 REST 时不需要再带 GUI 用户密码。** |

实现细节（本仓库与上游 Syncthing 一致）：认证中间件会**优先**校验 API 密钥；对 `/rest/` 的 CSRF 校验在携带有效 API 密钥时也会放行，便于服务端调用。

**实践建议**

- 自动化、中转服务、小程序后端：只用 **API 密钥** 访问中转 C 的 REST；密钥仅存服务端，勿写进小程序前端包。
- 为管理界面开启 **GUI 用户/密码**，避免把未鉴权的 Web UI 暴露给公网。

---

## 3. 修改配置后是否必须重启 Syncthing

- **通过 Web GUI 或 REST 修改**设备、文件夹、共享成员等：一般由运行中的进程**热应用**，**不需要**每加一台设备或一个组织就重启整个服务。
- **文档 05 中「先停止服务再改 `config.xml`」**：是为**整份手改文件**时的安全做法，不是「每次改配置都必须重启」的硬性规定。
- **注意**：若进程运行中**只改磁盘上的 `config.xml` 而不经 REST/GUI**，运行中的进程**不会自动重读文件**；应使用 REST，或改完后**重启一次**以加载文件。

与「是否重启」相关的多为**部分全局选项**（如监听地址等）变更；日常增加远程设备、调整文件夹共享列表，用 REST 即可。

---

## 4. 用 REST 替代手改 `config.xml`（对照文档 05）

文档 05 描述在中转设备 C 上为多个组织维护多个 `<folder>`，每个文件夹有独立的 `id`、`label`、`path` 与 `<device>` 列表。**等价操作**可通过 Syncthing 内置 REST 完成，无需停机整文件替换。

### 4.1 常用路径前缀

- Base：`http(s)://<C的地址>:8384/rest`（或 Unix socket，以实际部署为准）
- Header：`Content-Type: application/json`，以及 `X-API-Key: <C 的 API 密钥>`

### 4.2 与「文件夹 / 设备」直接相关的接口

| 操作 | 方法 | 路径 |
|------|------|------|
| 读取全部文件夹配置 | GET | `/rest/config/folders` |
| 读取单个文件夹 | GET | `/rest/config/folders/{folderId}` |
| 新建或整体替换单个文件夹配置 | PUT | `/rest/config/folders/{folderId}` |
| 删除文件夹 | DELETE | `/rest/config/folders/{folderId}` |
| 读取全部设备 | GET | `/rest/config/devices` |
| 读取单个设备 | GET | `/rest/config/devices/{deviceId}` |
| 新建或替换单个设备 | PUT | `/rest/config/devices/{deviceId}` |
| 删除设备 | DELETE | `/rest/config/devices/{deviceId}` |
| 待处理设备/文件夹（邀请等） | GET | `/rest/cluster/pending/devices`、`/rest/cluster/pending/folders` |

`{deviceId}`、`{folderId}` 为 URL 路径中的标识；设备 ID 含连字符时按常规 URL 编码即可。

### 4.3 新增一个组织（新文件夹）

1. 在 C 上创建好物理目录，并确定 `folderId`、`label`、`path`（与 05 一致）。
2. `PUT /rest/config/devices/{id}` 将该组织各终端登记为 C 的远程设备（若尚未存在）。
3. `PUT /rest/config/folders/{新folderId}` 提交**完整**文件夹 JSON：`id`、`label`、`path`、`type`（如 `sendreceive`）、`devices`（含本机 C 与各终端的 `deviceID`）等；字段含义与官方配置模型一致，详见 [02 REST API 参考](./02REST_API参考.md)。

**安全提示**：`PUT /rest/config/folders`（无 `{id}`）会替换**全部**文件夹列表，生产环境慎用；优先对单个 `PUT /rest/config/folders/{folderId}` 操作。

---

## 5. 加入中转设备 C：要调用哪些接口

「加入」是**双向**的：C 要知道终端，终端也要知道 C，并接受共享文件夹（若由 C 发起共享）。

### 5.1 在中转设备 C 上（登记某台电脑并纳入某文件夹）

适用于：例如把「研发部 A」加入 C 上 id 为 `development` 的文件夹（名称以你实际配置为准）。

1. **登记远程设备**（若该 `deviceId` 尚不存在于 C 的配置中）  
   `PUT /rest/config/devices/{研发部A的设备ID}`  
   Body：JSON，至少包含 `deviceID`、常用 `addresses`（如 `["dynamic"]`）、`name` 等，与 Web 界面「添加远程设备」等价。

2. **把该设备加入指定文件夹**  
   - `GET /rest/config/folders/development` 取出当前完整 JSON。  
   - 在 `devices` 数组中追加一项：`{ "deviceID": "<研发部A的设备ID>" }`（若已有加密等高级字段，保持原有结构）。  
   - `PUT /rest/config/folders/development` 写回**完整**文件夹对象，避免其它字段被误清空。

### 5.2 在终端本机（例如研发部 A）

1. **添加远程设备「中转 C」**  
   `PUT /rest/config/devices/{C的设备ID}`  
   Body：含 `deviceID`（C 的 ID）、`addresses`（如 `tcp://<C的IP>:22000` 或 `dynamic` 等，与网络环境一致）。

2. **接受来自 C 的文件夹**  
   当 C 已把文件夹共享给该终端后，本机会出现待接受项：需在 Web UI 确认，或通过本机 REST 增加对应文件夹配置；可先 `GET /rest/cluster/pending/folders` 了解待处理项，再按官方/02 文档完成与本机 `PUT /rest/config/folders/...` 等价的配置。

以上本机请求使用 **终端本机** 的 Base URL（如 `http://127.0.0.1:8384/rest`）及 **该终端** 的 API 密钥或已登录会话。

---

## 6. 微信小程序扫码场景：调谁、传什么

小程序**不应**把中转 C 的公网地址与 Syncthing **API 密钥**硬编码在客户端；应调用**你们自己的 HTTPS 业务 API**，由后端在受控网络内调用 C 的 REST。

### 6.1 扫码得到的内容

使用微信 `wx.scanCode` 等扫描 **终端 Syncthing 设备二维码** 后，结果一般为 **该终端的 Syncthing 设备 ID 字符串**（若你们自定义了二维码内容，需先解析再取出设备 ID）。

### 6.2 建议的小程序 → 自有后端

示例（路径与字段名可自定）：

- **方法**：`POST https://<你的域名>/api/v1/bind-device`（示例）
- **Header**：`Content-Type: application/json`；用户身份用小程序 `wx.login` + 后端 session / JWT 等，**不要**在 Body 里明文传 openid 若可从 token 解析。

**Body 示例**：

```json
{
  "syncthingDeviceId": "AAAAAAA-BBBBBBB-CCCCCCC-DDDDDDD-EEEEEEE-FFFFFFF-GGGGGGG-HHHHHHH",
  "org": "development"
}
```

| 字段 | 说明 |
|------|------|
| `syncthingDeviceId` | 扫码解析出的终端设备 ID |
| `org` | 业务侧部门/组织标识；后端映射到 C 上的 `folderId`（如 `finance` / `development`） |

后端校验权限后，对 C 执行 **第 5.1 节** 的 `PUT .../devices` 与 `GET` + `PUT .../folders/{folderId}`。

### 6.3 终端侧仍需「加入 C」

小程序与后端**无法**代替在内网终端上完成「本机添加设备 C」；仍需用户在终端 Web UI 操作一次，或由你们 **桌面客户端**（如 Electron）调用本机 `127.0.0.1:8384/rest` 自动添加 C 并接受共享（见第 5.2 节）。

---

## 7. 附录：curl 与最小 JSON 示例

说明：`PUT /rest/config/folders/...` 时，服务端会用 **默认文件夹模板** 与请求体 JSON 合并；因此请求体里**只需写与默认不同的关键字段**即可，不必贴整份 GUI 导出的长 JSON。

### 7.1 在 C 上登记设备 A（最小请求体）

将 `DEV-A-DEVICE-ID` 换成 A 的完整设备 ID；在 **C 本机**或能访问 C 的 `8384` 的环境执行。

```bash
curl -sS -X PUT "http://127.0.0.1:8384/rest/config/devices/DEV-A-DEVICE-ID" \
  -H "X-API-Key: YOUR_C_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"deviceID":"DEV-A-DEVICE-ID","name":"研发-A","addresses":["dynamic"],"compression":"metadata"}'
```

（其余布尔字段若省略，会沿用 Syncthing 默认设备模板；需要与 GUI 完全一致时再补全。）

### 7.2 在 C 上新建文件夹（与 [05](./05单中转设备多组织同步配置指南.md) 中 `development` / `finance` 对齐）

REST 里 `devices` 数组应对应 XML 里该 `<folder>` 下的各 `<device id="..."/>`：研发部文件夹为 **DEV-A、DEV-B、中转 C**；财务部为 **FIN-A、FIN-B、中转 C**。占位符需替换为真实设备 ID；`path` 改为 C 上实际路径。

- 研发部示例：[folder-development.min.json](./folder-development.min.json)（`rescanIntervalS` 1800、版本保留约 7 天，与文档 05 一致）
- 财务部示例：[folder-finance.min.json](./folder-finance.min.json)（`rescanIntervalS` 3600、版本保留约 1 年，与文档 05 一致）

```bash
curl -sS -X PUT "http://127.0.0.1:8384/rest/config/folders/development" \
  -H "X-API-Key: YOUR_C_API_KEY" \
  -H "Content-Type: application/json" \
  -d @folder-development.min.json
```

若某台成员（如 `DEV_B`）尚未就绪，可暂时从 `devices` 中删掉对应项，待设备注册后再 `GET` → 合并 `devices` → **整份** `PUT`。若文件夹 **已存在** 且只需追加成员，同样建议 `GET` 后合并再 `PUT`，避免覆盖其它字段。

---

## 8. 小结表

| 目标 | 主要 REST（在谁上调用） |
|------|-------------------------|
| C 登记某终端并把它加入某组织文件夹 | C：`PUT .../config/devices/{id}`，`GET` + `PUT .../config/folders/{folderId}` |
| 终端把 C 加为远程设备并接受同步 | 终端：`PUT .../config/devices/{C的id}`，并按需处理 pending / 本机 folder |
| 程序化改配置、日常扩容 | 优先 REST；避免运行中仅改磁盘 `config.xml` 却不重载 |
| 鉴权 | 程序化调用统一使用 **API 密钥**；GUI 用户名密码保护浏览器 |

更完整的端点列表、错误码与字段说明见 [02 REST API 参考](./02REST_API参考.md)；多组织目录与隔离原则见 [05 单中转设备多组织同步配置指南](./05单中转设备多组织同步配置指南.md)。
