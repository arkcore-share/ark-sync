# 无网页（API-only）打包指南

> **适用目标**：不提供 Web GUI 页面，只保留 REST API 能力。  
> **相关文档**：[07 安装与使用文档](./07安装与使用文档.md) | [02 REST API 参考](./02REST_API参考.md)

---

## 1. 目标说明

本文档用于构建“无网页前端资源”的 Syncthing 可执行文件：

- 不打包 `gui` 静态页面资源
- 不访问 Web 页面（`/`、`/index.html` 返回 404）
- 保留 `/rest/*` API，供自动化系统调用

适合中转设备、后端服务节点、仅 API 管理场景。

---

## 2. 代码基础（已完成）

仓库已具备以下支持：

1. 使用 `noassets` 构建标签切换到无前端资源模式。  
2. 在 `noassets` 模式下，静态资源处理器返回 404。  
3. API 路由与处理逻辑不受影响。

---

## 3. 打包命令

### 3.1 使用仓库构建脚本（推荐）

```bash
go run build.go -tags "noassets" -build-out ./bin/arksync build syncthing
```

Windows 建议显式输出 `.exe` 文件名：

```powershell
go run build.go -tags "noassets" -build-out ./bin/arksync.exe build syncthing
```

产物为 `./bin/arksync`。

### 3.2 直接 Go 构建

```bash
go build -tags noassets -o ./bin/arksync ./cmd/syncthing
```

---

## 4. 运行方式

```bash
./bin/arksync --no-browser
```

或 systemd / 守护进程方式运行均可。`--no-browser` 仅避免自动打开浏览器，不影响 API。

---

## 5. 行为验证

假设服务监听在 `127.0.0.1:8384`。

### 5.1 页面不可访问（预期）

```bash
curl -i http://127.0.0.1:8384/
curl -i http://127.0.0.1:8384/index.html
```

预期：`404 Not Found`。

### 5.2 API 可访问（预期）

```bash
curl -s http://127.0.0.1:8384/rest/noauth/health
curl -s -H "X-API-Key: <API_KEY>" http://127.0.0.1:8384/rest/system/ping
```

预期：返回健康状态和 `{"ping":"pong"}`（或等价 JSON）。

---

## 6. 注意事项

1. `noassets` 下，部分“依赖 GUI 页面 200 返回”的测试会失败，这是预期行为，不影响 API-only 目标。  
2. 若你使用外部监控探活，请改为探活 `/rest/noauth/health`，不要探活 `/`。  
3. 若通过反向代理暴露服务，建议仅开放 `/rest/` 路径并结合 IP 白名单与 API Key。  
4. 当前仓库已关闭自动升级和手动升级入口（GUI/CLI/REST），如需恢复请单独回退相关补丁。

---

## 7. 建议的生产参数

- 启动参数：`--no-browser`
- 绑定地址：GUI/API 只监听内网或 localhost
- 鉴权：启用 API Key，密钥仅保存在服务端
- 运维：统一通过 REST 进行配置与状态查询

---

## 8. GitHub Actions 自动发布（Linux / Windows / macOS）

仓库已提供 [`.github/workflows/release-ark.yaml`](../.github/workflows/release-ark.yaml)，在 GitHub 上编译 **API-only**（`noassets`）的 `arksync` 并上传到 Release。

### 8.1 发布步骤

1. 提交并推送你的修改到 `arkcore-share/ark-sync`。
2. 打版本标签并推送（版本号来自 git tag，与 `go run build.go version` 一致）：

```bash
git tag v1.0.0
git push origin v1.0.0
```

3. 打开 GitHub → **Actions** → **Release ark-sync**，确认 workflow 成功。
4. 在 **Releases** 页面下载对应平台包，例如：
   - `arksync-linux-amd64-v1.0.0.tar.gz`
   - `arksync-windows-amd64-v1.0.0.zip`
   - `arksync-macos-universal-v1.0.0.zip`（Intel + Apple Silicon）
5. 使用 `SHA256SUMS.txt` 校验下载文件。

也可在 Actions 里 **Run workflow**，填写已存在的 tag（如 `v1.0.0`）手动重跑发布。

### 8.2 产物说明

- 压缩包前缀为 `arksync-`，包内可执行文件为 `arksync` / `arksync.exe`。
- 默认带 `noassets` 标签，与本文档第 3 节本地打包一致。
- 若 Release 需要带 Web GUI，编辑 workflow 中 `TAGS` / `TAGS_LINUX`，去掉 `noassets` 后重新打 tag 发布。

### 8.3 与官方 Syncthing Release 的区别

- 不能使用 [syncthing/syncthing releases](https://github.com/syncthing/syncthing/releases)（未包含本仓库修改）。
- 本 workflow 不依赖 Syncthing 组织的签名密钥；macOS/Windows 包为 **未公证/未签名** 版本，首次运行可能需系统安全提示放行。

---

如需“API-only Docker 镜像构建模板（含健康检查、最小权限、只读根文件系统）”，可在本仓库再补一份 `Dockerfile` 示例。
