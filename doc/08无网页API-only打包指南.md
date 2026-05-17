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

如需“API-only Docker 镜像构建模板（含健康检查、最小权限、只读根文件系统）”，可在本仓库再补一份 `Dockerfile` 示例。
