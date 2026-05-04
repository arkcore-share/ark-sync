# Syncthing 完整指南

> 📚 **本指南包含**: 功能介绍 | 使用教程 | 架构设计 | 快速参考  
> 🔗 **API 文档**: [REST API 参考](02REST_API参考.md)

---

# 目录

## 第一部分：入门指南
- [什么是 Syncthing](#什么是-syncthing)
- [快速开始](#快速开始)
- [安装部署](#安装部署)

## 第二部分：核心功能
- [文件同步](#文件同步)
- [设备管理](#设备管理)
- [网络与连接](#网络与连接)
- [安全与加密](#安全与加密)
- [版本控制](#版本控制)

## 第三部分：使用教程
- [基本使用](#基本使用)
- [高级配置](#高级配置)
- [故障排除](#故障排除)
- [最佳实践](#最佳实践)

## 第四部分：技术架构
- [整体架构](#整体架构)
- [核心模块](#核心模块)
- [数据库设计](#数据库设计)
- [性能优化](#性能优化)

## 第五部分：快速参考
- [常用命令](#常用命令)
- [默认端口](#默认端口)
- [环境变量](#环境变量)
- [安全建议](#安全建议)

---

# 第一部分：入门指南

## 什么是 Syncthing

Syncthing 是一个**开源的连续文件同步程序**，用于在两台或多台计算机之间同步文件。

### 设计目标（按重要性排序）

1. **数据安全** - 保护用户数据免受损失
2. **安全防护** - 防止未授权访问
3. **易于使用** - 简单直观
4. **自动化** - 最小化用户交互
5. **通用可用** - 跨平台支持

### 核心特点

- ✅ **免费开源**: MPLv2 许可证
- ✅ **去中心化**: P2P 架构，无需服务器
- ✅ **安全**: 端到端 TLS 加密
- ✅ **隐私**: 数据只在你的设备间传输
- ✅ **跨平台**: Windows, macOS, Linux, BSD, Solaris
- ✅ **自动化**: 配置后无需干预

### 技术特性

- **语言**: Go (Golang)
- **架构**: 去中心化 P2P
- **协议**: BEP (Block Exchange Protocol)
- **加密**: TLS 1.2+

---

## 快速开始

### 1. 安装

#### Ubuntu/Debian
```bash
curl -fsSL https://syncthing.net/release-key.txt | sudo gpg --dearmor -o /usr/share/keyrings/syncthing-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/syncthing-archive-keyring.gpg] https://apt.syncthing.net/ syncthing stable" | sudo tee /etc/apt/sources.list.d/syncthing.list
sudo apt update && sudo apt install syncthing
```

#### macOS
```bash
brew install syncthing
```

#### Windows
从 https://syncthing.net/downloads/ 下载安装包

#### Docker
```bash
docker run -d \
  --name syncthing \
  -p 8384:8384 -p 22000:22000/tcp -p 22000:22000/udp -p 21027:21027/udp \
  -v ./config:/var/syncthing \
  -v ./sync:/var/syncthing/Sync \
  syncthing/syncthing:latest
```

### 2. 启动

```bash
syncthing                          # 启动
syncthing --no-browser             # 不打开浏览器
syncthing --version                # 查看版本
syncthing deviceid                 # 查看设备ID
```

### 3. 访问 Web 界面

打开浏览器访问: http://127.0.0.1:8384

---

## 安装部署

### 配置路径

| 平台 | 配置目录 | 数据目录 |
|------|----------|----------|
| **Linux** | `~/.config/syncthing/` | `~/.local/state/syncthing/` |
| **macOS** | `~/Library/Application Support/Syncthing/` | 同配置目录 |
| **Windows** | `%LOCALAPPDATA%\Syncthing\` | 同配置目录 |

### 主要文件

```
config.xml          - 配置文件
cert.pem            - TLS 证书
key.pem             - TLS 私钥
index-v0.14.0.db/   - SQLite 数据库
```

### systemd 服务

```bash
# 用户服务
systemctl --user start syncthing
systemctl --user enable syncthing
systemctl --user status syncthing

# 查看日志
journalctl --user -u syncthing -f
```

---

# 第二部分：核心功能

## 文件同步

### 同步模式

| 模式 | 说明 | 适用场景 |
|------|------|----------|
| **发送接收** | 双向同步（默认） | 大多数场景 |
| **仅发送** | 只发送不接收 | 备份服务器 |
| **仅接收** | 只接收不发送 | 只读客户端 |
| **接收加密** | 接收加密数据 | 不可信设备 |

### 同步机制

1. **文件扫描**: 检测文件变更
2. **索引交换**: 设备间交换文件列表
3. **块传输**: 分块传输文件
4. **冲突解决**: 自动处理冲突文件

### 忽略模式

创建 `.stignore` 文件：

```bash
# 忽略所有 .log 文件
*.log

# 忽略特定目录
node_modules/
.tmp/

# 忽略操作系统文件
.DS_Store
Thumbs.db

# 例外：保留重要日志
!important.log
```

---

## 设备管理

### 设备 ID

- 基于 TLS 证书生成
- 52 字符的唯一标识
- 格式: `AAAAA-BBBBB-CCCCC-DDDDD-EEEEE-FFFFF-GGGGG-HHHHH`

### 添加设备

1. 在 Web 界面点击"添加远程设备"
2. 输入设备 ID 和名称
3. 选择要共享的文件夹
4. 在另一台设备上接受

### 设备配置

- **名称**: 便于识别
- **地址**: 自动发现或手动指定
- **压缩**: 启用/禁用数据压缩
- **速率限制**: 限制同步速度

---

## 网络与连接

### 连接方式

Syncthing 按以下优先级尝试连接：

1. **局域网直连** (LAN) - 最快
2. **广域网直连** (WAN) - 需要 NAT 穿透
3. **中继连接** (Relay) - 最后选择

### 传输协议

- **TCP**: 端口 22000
- **QUIC**: 端口 22000 (UDP)
- **WebSocket**: 可选

### 服务发现

#### 本地发现
- 协议: UDP 广播
- 端口: 21027
- 范围: 同一局域网

#### 全局发现
- 协议: HTTPS
- 服务器: discovery.syncthing.net
- 功能: 跨网络发现设备

### 中继系统

**作用**: 当设备无法直连时，作为数据中转站

**特点**:
- ❌ 不存储文件（只转发）
- ✅ 端到端加密
- ✅ 社区提供（免费）
- ⚠️ 速度较慢

> 💡 **深入了解**: 中继池服务器见下方 [基础设施组件](#基础设施组件)

---

## 安全与加密

### TLS 加密

- 所有通信使用 TLS 1.2+
- 自动生成证书
- 设备 ID 基于证书

### 端到端加密

- 数据在中继服务器上也是加密的
- 中继服务器无法解密

### 访问控制

- **API 密钥**: REST API 认证
- **GUI 密码**: Web 界面保护
- **设备验证**: 手动确认设备 ID

### 不可信设备

支持端到端加密文件夹：
- 数据在发送前加密
- 不可信设备无法读取内容
- 只能存储加密数据

---

## 版本控制

当文件被修改或删除时，自动保留历史版本。

### 版本控制类型

| 类型 | 适用场景 | 保留策略 |
|------|----------|----------|
| **无** | 不需要历史 | 不保留 |
| **简单** | 有限历史 | 保留最近 N 个版本 |
| **回收站** | 防止误删 | 保留 N 天后清理 |
| **交错** | 长期追踪（推荐） | 智能时间间隔保留 |
| **外部** | Git 集成 | 自定义脚本 |

### 配置示例

**简单版本控制**:
- 保留数量: 5
- 清理天数: 30

**交错版本控制**:
- 1 小时内: 每 30 秒
- 1 天内: 每小时
- 30 天内: 每天
- 超过 30 天: 每周

### 版本存储

- **位置**: `.stversions` 文件夹
- **命名**: `file.txt~20240101-120000.txt`
- **结构**: 保留原目录结构

### 集成 Git

使用外部版本控制 + Git 脚本，可以自动提交到 GitHub/GitLab：

```bash
#!/bin/bash
FOLDER_PATH="$1"
FILE_PATH="$2"

cd "$FOLDER_PATH"
git add -A
git commit -m "Auto-sync: $FILE_PATH at $(date)"
git push origin main
```

---

# 第三部分：使用教程

## 基本使用

### 添加文件夹

1. 点击"添加文件夹"
2. 设置文件夹标签和 ID
3. 选择本地路径
4. 选择要共享的设备
5. 配置版本控制（可选）

### 查看同步状态

- **仪表盘**: 显示所有设备和文件夹
- **传输速度**: 实时显示上传/下载速度
- **进度条**: 显示同步进度
- **最近更改**: 查看文件变更历史

### 暂停/恢复

```bash
# CLI 方式
syncthing cli operations pause myfolder
syncthing cli operations resume myfolder

# API 方式
curl -X POST -H "X-Api-Key: KEY" "http://127.0.0.1:8384/rest/system/pause?folder=myfolder"
```

---

## 高级配置

### 速率限制

```xml
<options>
  <maxSendKbps>1000</maxSendKbps>    <!-- 上传限制 1MB/s -->
  <maxRecvKbps>5000</maxRecvKbps>    <!-- 下载限制 5MB/s -->
</options>
```

### 自定义监听地址

```xml
<gui>
  <address>0.0.0.0:8384</address>    <!-- 允许外部访问 -->
  <tls>true</tls>                     <!-- 启用 HTTPS -->
</gui>
```

### 代理设置

```bash
# 环境变量
export http_proxy=http://proxy:8080
export https_proxy=http://proxy:8080
```

---

## 故障排除

### 设备无法连接

**检查清单**:
1. ✅ 防火墙是否开放端口 22000
2. ✅ 路由器是否配置端口转发（UPnP）
3. ✅ 设备 ID 是否正确
4. ✅ 是否在同一个网络或可访问

**解决方案**:
```bash
# 手动指定地址
# 在设备配置中添加:
tcp://192.168.1.100:22000
```

### 同步速度慢

**可能原因**:
- 使用中继而非直连
- 速率限制过低
- 网络带宽不足

**优化建议**:
1. 检查连接类型（直连 > 中继）
2. 增加并发传输数
3. 调整速率限制

### 文件冲突

**原因**: 两个设备同时修改同一文件

**解决**:
- 查看 `.sync-conflict` 文件
- 手动选择保留版本
- 启用版本控制保留历史

### 重置数据库

```bash
# 方法 1: CLI
syncthing debug reset-database

# 方法 2: 手动删除
rm -rf ~/.local/state/syncthing/index-v0.14.0.db/

# 重启后会重新扫描
```

### Web GUI 无法访问

```bash
# 检查服务状态
systemctl --user status syncthing

# 查看监听端口
netstat -tlnp | grep 8384

# 重置 GUI 密码
# 编辑 config.xml，删除 <user> 和 <password> 字段
```

---

## 最佳实践

### 1. 安全配置

- ✅ 设置 GUI 密码
- ✅ 使用防火墙限制访问
- ✅ 定期更新 Syncthing
- ✅ 保护证书文件
- ✅ 验证设备 ID

### 2. 性能优化

- 优先使用直连而非中继
- 合理设置速率限制
- 定期清理旧版本
- 使用 SSD 存储数据库

### 3. 备份策略

- 定期备份 config.xml
- 备份 TLS 证书
- 使用版本控制
- 考虑外部备份（Git/云存储）

---

# 第四部分：技术架构

## 整体架构

### 分层架构

```
┌─────────────────────────────────────────────────┐
│                   用户界面层                      │
│  Web GUI  │  REST API  │  CLI Tools            │
└─────────────────────────────────────────────────┘
                        │
┌─────────────────────────────────────────────────┐
│                   应用服务层                      │
│   Model   │    API    │   Config               │
└─────────────────────────────────────────────────┘
                        │
┌─────────────────────────────────────────────────┐
│                   协议与连接层                    │
│  Protocol │ Connections │  Discovery            │
└─────────────────────────────────────────────────┘
                        │
┌─────────────────────────────────────────────────┐
│                   基础设施层                      │
│  Database │  Filesystem │    Events             │
└─────────────────────────────────────────────────┘
```

### 组件关系

```
用户操作 → Web GUI/CLI
              ↓
         lib/api (HTTP 服务)
              ↓
        lib/model (核心逻辑)
              ↓
    ┌─────────┼─────────┐
    ↓         ↓         ↓
lib/config lib/protocol lib/connections
    ↓         ↓         ↓
lib/events lib/scanner lib/discover
    ↓         ↓         ↓
lib/db    lib/fs      lib/relay
```

---

## 核心模块

### lib/model - 核心业务逻辑

**职责**:
- 文件夹管理
- 设备协调
- 同步调度
- 冲突解决

### lib/protocol - BEP 协议

**功能**:
- 索引交换
- 文件请求
- 块传输
- 加密解密

### lib/connections - 连接管理

**特性**:
- TCP/QUIC 支持
- 中继连接
- 连接池管理
- 优先级调度

### lib/config - 配置管理

**特点**:
- XML 格式
- 自动迁移
- 验证修复
- 热加载

### lib/scanner - 文件扫描

**机制**:
- 文件系统监控
- 哈希计算
- 增量扫描
- 批量处理

### lib/versioner - 版本控制

**实现**:
- Simple: 保留 N 个版本
- Trashcan: 回收站模式
- Staggered: 智能时间间隔
- External: 外部脚本

### lib/events - 事件系统

**类型**:
- 设备事件（连接/断开）
- 文件夹事件（同步完成/错误）
- 文件事件（开始/完成）
- 系统事件（启动/关闭）

### lib/fs - 文件系统抽象

**支持**:
- 本地文件系统
- 虚拟文件系统
- 权限处理
- 符号链接

### lib/ignore - 忽略模式

**功能**:
- .stignore 解析
- 通配符匹配
- 排除规则
- 性能优化

---

## 数据库设计

### SQLite 数据库

**位置**: `index-v0.14.0.db/`

**特点**:
- WAL 模式（Write-Ahead Logging）
- 增量清理
- 自动优化

**存储内容**:
- 文件索引
- 同步元数据
- 设备信息
- 全局版本向量

### 数据库结构

```sql
-- 文件信息表
CREATE TABLE file_info (
    folder_id TEXT,
    file_name TEXT,
    version BIGINT,
    size BIGINT,
    modified_at TIMESTAMP,
    ...
);

-- 全局版本向量
CREATE TABLE global_version (
    device_id TEXT,
    folder_id TEXT,
    version BIGINT,
    ...
);
```

---

## 性能优化

### 并发模型

- Goroutine 池
- 异步 I/O
- 批量处理
- 优先级队列

### 内存管理

- LRU 缓存
- 对象池
- 延迟加载
- 增量扫描

### 网络优化

- 连接复用
- 数据压缩
- 块大小调整
- 并行传输

---

# 第五部分：快速参考

## 常用命令

### 启动与管理

```bash
syncthing                          # 启动
syncthing --no-browser             # 不打开浏览器
syncthing --version                # 查看版本
syncthing deviceid                 # 查看设备ID
syncthing paths                    # 查看路径
syncthing generate                 # 生成配置
```

### CLI 操作

```bash
syncthing cli config show          # 显示配置
syncthing cli show devices         # 显示设备
syncthing cli show folders         # 显示文件夹
syncthing cli operations pause myfolder    # 暂停
syncthing cli operations resume myfolder   # 恢复
```

---

## 默认端口

| 端口 | 协议 | 用途 |
|------|------|------|
| **8384** | TCP | Web GUI |
| **22000** | TCP | 同步 (TCP) |
| **22000** | UDP | 同步 (QUIC) |
| **21027** | UDP | 本地发现 |

---

## 环境变量

```bash
STHOMEDIR=/path/to/home          # 主目录
STCONFDIR=/path/to/config        # 配置目录
STDATADIR=/path/to/data          # 数据目录
STGUIADDRESS=0.0.0.0:8384        # GUI地址
STNOUPGRADE=1                    # 禁用升级
STTRACE=model,connections        # 调试日志
STLOGLEVEL=DEBUG                 # 日志级别
```

---

## 防火墙配置

### UFW (Ubuntu)

```bash
sudo ufw allow 22000/tcp
sudo ufw allow 22000/udp
sudo ufw allow 21027/udp
sudo ufw allow 8384/tcp  # 如需远程访问 GUI
```

### firewalld (CentOS/RHEL)

```bash
sudo firewall-cmd --permanent --add-port=22000/tcp
sudo firewall-cmd --permanent --add-port=22000/udp
sudo firewall-cmd --permanent --add-port=21027/udp
sudo firewall-cmd --reload
```

---

## 安全建议

### ✅ 必须做

1. **设置 GUI 密码**
2. **使用防火墙**
3. **定期更新**
4. **保护证书文件**
5. **验证设备 ID**

### ⚠️ 建议做

6. 启用端到端加密（不可信设备）
7. 限制 GUI 访问 IP
8. 使用 HTTPS
9. 定期备份配置
10. 监控日志

### ❌ 不要做

- 不要暴露 GUI 到公网（除非必要）
- 不要忽略设备 ID 验证
- 不要禁用 TLS
- 不要删除证书文件
- 不要共享 config.xml

---

## 外部资源

- **官方网站**: https://syncthing.net/
- **官方文档**: https://docs.syncthing.net/
- **社区论坛**: https://forum.syncthing.net/
- **GitHub**: https://github.com/syncthing/syncthing
- **API 文档**: [REST API 参考](02REST_API参考.md)

---

**文档版本**: v2.0  
**最后更新**: 2024-01  
**包含内容**: 功能介绍 + 使用教程 + 架构设计 + 快速参考
