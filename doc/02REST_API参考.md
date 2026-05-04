# Syncthing REST API 参考文档

> 📘 **相关文档**: [项目功能详解](项目功能详解.md) | [使用指南](使用指南.md) | [架构与模块详解](架构与模块详解.md) | [快速参考](快速参考.md)

## 目录
- [概述](#概述)
- [认证](#认证)
- [系统 API](#系统-api)
- [配置 API](#配置-api)
- [数据库 API](#数据库-api)
- [统计 API](#统计-api)
- [事件 API](#事件-api)
- [服务 API](#服务-api)
- [错误处理](#错误处理)
- [使用示例](#使用示例)

---

## 概述

Syncthing 提供完整的 RESTful API，允许程序化控制和监控 Syncthing 实例。

### 基本信息

- **Base URL**: `http://127.0.0.1:8384/rest`
- **数据格式**: JSON
- **字符编码**: UTF-8
- **HTTP 方法**: GET, POST, PUT, DELETE

### 内容类型

**请求**:
```
Content-Type: application/json
```

**响应**:
```
Content-Type: application/json
```

---

## 认证

### API Key 认证

所有 API 请求（除 `/rest/noauth/` 外）都需要 API Key。

#### 方式 1: HTTP Header (推荐)
```http
X-Api-Key: your-api-key-here
```

#### 方式 2: Header (备用)
```http
Authorization: Basic base64(user:password)
```

#### 方式 3: Cookie
通过 Web 界面登录后的会话 Cookie

### 获取 API Key

**方法 1**: Web 界面
```
操作 -> 高级 -> 选项 -> GUI -> API Key
```

**方法 2**: 配置文件
```xml
<config>
  <gui>
    <apikey>your-api-key-here</apikey>
  </gui>
</config>
```

**方法 3**: API 端点
```bash
curl http://127.0.0.1:8384/rest/system/config | jq '.gui.apikey'
```

### 生成新 API Key

```bash
curl -X POST http://127.0.0.1:8384/rest/system/reset
```

---

## 系统 API

### GET /rest/system/status

获取系统状态信息。

**响应**:
```json
{
  "alloc": 25639808,
  "connectionServiceStatus": {
    "dynamic+https://relays.syncthing.net/endpoint": {
      "error": null
    }
  },
  "cpuPercent": 1.2,
  "discoveryEnabled": true,
  "discoveryErrors": {},
  "discoveryMethods": 3,
  "goroutines": 48,
  "lastDialStatus": {},
  "myID": "ABCDEF1-ABCDEF1-ABCDEF1-ABCDEF1-ABCDEF1-ABCDEF1-A",
  "pathSeparator": "/",
  "startTime": "2024-01-01T00:00:00Z",
  "sys": 52428800,
  "themes": ["default", "dark", "light", "black"],
  "tilde": "/home/user",
  "upgradesAllowed": true,
  "upgradesMetadata": {
    "latestRelease": "v1.27.0"
  },
  "version": "v1.27.0"
}
```

### GET /rest/system/ping

简单的连通性测试。

**响应**:
```json
{
  "ping": "pong"
}
```

### POST /rest/system/pause

暂停同步。

**参数**:
- `folder` (可选): 文件夹 ID
- `device` (可选): 设备 ID

**示例**:
```bash
# 暂停所有
curl -X POST http://127.0.0.1:8384/rest/system/pause \
  -H "X-Api-Key: KEY"

# 暂停特定文件夹
curl -X POST "http://127.0.0.1:8384/rest/system/pause?folder=myfolder" \
  -H "X-Api-Key: KEY"

# 暂停特定设备
curl -X POST "http://127.0.0.1:8384/rest/system/pause?device=DEVICEID" \
  -H "X-Api-Key: KEY"
```

### POST /rest/system/resume

恢复同步。

**参数**: 同 pause

### POST /rest/system/restart

重启 Syncthing。

```bash
curl -X POST http://127.0.0.1:8384/rest/system/restart \
  -H "X-Api-Key: KEY"
```

### POST /rest/system/shutdown

关闭 Syncthing。

```bash
curl -X POST http://127.0.0.1:8384/rest/system/shutdown \
  -H "X-Api-Key: KEY"
```

### POST /rest/system/reset

重置配置或数据库。

**参数**:
- `folder` (可选): 重置特定文件夹
- 无参数: 重置所有

```bash
# 重置特定文件夹
curl -X POST "http://127.0.0.1:8384/rest/system/reset?folder=myfolder" \
  -H "X-Api-Key: KEY"

# 重置所有
curl -X POST http://127.0.0.1:8384/rest/system/reset \
  -H "X-Api-Key: KEY"
```

### POST /rest/system/upgrade

执行升级。

```bash
curl -X POST http://127.0.0.1:8384/rest/system/upgrade \
  -H "X-Api-Key: KEY"
```

### GET /rest/system/version

获取版本信息。

**响应**:
```json
{
  "arch": "amd64",
  "longVersion": "syncthing v1.27.0 ...",
  "os": "linux",
  "version": "v1.27.0"
}
```

### GET /rest/system/log

获取日志。

**响应**:
```json
{
  "messages": [
    {
      "when": "2024-01-01T00:00:00Z",
      "message": "Starting syncthing",
      "level": "INFO"
    }
  ]
}
```

### GET /rest/system/debug

获取/设置调试设施。

**GET 响应**:
```json
{
  "enabled": ["model", "connections"],
  "disabled": ["scanner", "versioner"]
}
```

**POST 设置**:
```bash
curl -X POST http://127.0.0.1:8384/rest/system/debug \
  -H "X-Api-Key: KEY" \
  -d '{"enable": "model,connections", "disable": "scanner"}'
```

### GET /rest/system/discovered

获取发现的设备。

**响应**:
```json
{
  "DEVICEID": {
    "addresses": [
      {
        "address": "192.168.1.100:22000",
        "type": "tcp"
      }
    ]
  }
}
```

### GET /rest/system/connections

获取连接状态。

**响应**:
```json
{
  "connections": {
    "DEVICEID": {
      "address": "192.168.1.100:22000",
      "at": "2024-01-01T00:00:00Z",
      "clientVersion": "v1.27.0",
      "connected": true,
      "crypto": "TLS1.3-TLS_AES_128_GCM_SHA256",
      "inBytesTotal": 1234567,
      "outBytesTotal": 7654321,
      "paused": false,
      "type": "tcp-server"
    }
  },
  "total": {
    "at": "2024-01-01T00:00:00Z",
    "inBytesTotal": 1234567,
    "outBytesTotal": 7654321
  }
}
```

### GET /rest/system/config

获取完整配置。

**响应**: 完整 Configuration JSON

### PUT /rest/system/config

更新完整配置。

```bash
curl -X PUT http://127.0.0.1:8384/rest/system/config \
  -H "X-Api-Key: KEY" \
  -H "Content-Type: application/json" \
  -d @config.json
```

### POST /rest/system/config

部分更新配置。

```bash
curl -X POST http://127.0.0.1:8384/rest/system/config \
  -H "X-Api-Key: KEY" \
  -H "Content-Type: application/json" \
  -d '{"options": {"maxSendKbps": 1000}}'
```

### GET /rest/system/error

获取错误列表。

**响应**:
```json
{
  "errors": [
    {
      "when": "2024-01-01T00:00:00Z",
      "message": "Failed to connect to device"
    }
  ]
}
```

### POST /rest/system/error

清除错误。

```bash
curl -X POST http://127.0.0.1:8384/rest/system/error/clear \
  -H "X-Api-Key: KEY"
```

### GET /rest/system/browse

浏览文件系统路径。

**参数**: `current` (可选): 当前路径

**响应**:
```json
["/home", "/home/user", "/home/user/Sync"]
```

### GET /rest/noauth/health

健康检查（无需认证）。

**响应**:
```json
{
  "status": "OK"
}
```

---

## 配置 API

### GET /rest/config

获取配置。

**同**: `/rest/system/config`

### PUT /rest/config

更新配置。

### GET /rest/config/folders

获取所有文件夹配置。

**响应**:
```json
[
  {
    "id": "myfolder",
    "label": "My Folder",
    "path": "/home/user/Sync",
    "type": "sendreceive",
    "devices": [
      {
        "deviceID": "DEVICEID"
      }
    ],
    "rescanIntervalS": 3600,
    "ignorePerms": false,
    "versioning": {
      "type": "staggered",
      "params": {
        "maxAge": "31536000"
      }
    },
    "order": "random",
    "paused": false
  }
]
```

### GET /rest/config/folders/{folderID}

获取特定文件夹配置。

### PUT /rest/config/folders/{folderID}

创建或更新文件夹。

```bash
curl -X PUT http://127.0.0.1:8384/rest/config/folders/myfolder \
  -H "X-Api-Key: KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "myfolder",
    "label": "My Folder",
    "path": "/home/user/Sync",
    "type": "sendreceive",
    "devices": [{"deviceID": "DEVICEID"}],
    "rescanIntervalS": 3600
  }'
```

### DELETE /rest/config/folders/{folderID}

删除文件夹。

```bash
curl -X DELETE http://127.0.0.1:8384/rest/config/folders/myfolder \
  -H "X-Api-Key: KEY"
```

### GET /rest/config/devices

获取所有设备配置。

**响应**:
```json
[
  {
    "deviceID": "DEVICEID",
    "name": "My Laptop",
    "addresses": ["dynamic"],
    "compression": "metadata",
    "introducer": false,
    "paused": false
  }
]
```

### GET /rest/config/devices/{deviceID}

获取特定设备配置。

### PUT /rest/config/devices/{deviceID}

创建或更新设备。

```bash
curl -X PUT http://127.0.0.1:8384/rest/config/devices/DEVICEID \
  -H "X-Api-Key: KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceID": "DEVICEID",
    "name": "My Laptop",
    "addresses": ["dynamic"],
    "compression": "metadata"
  }'
```

### DELETE /rest/config/devices/{deviceID}

删除设备。

### GET /rest/config/defaults/folder

获取默认文件夹配置。

### PUT /rest/config/defaults/folder

设置默认文件夹配置。

### GET /rest/config/defaults/device

获取默认设备配置。

### PUT /rest/config/defaults/device

设置默认设备配置。

### GET /rest/config/defaults/ignores

获取默认忽略模式。

### PUT /rest/config/defaults/ignores

设置默认忽略模式。

```bash
curl -X PUT http://127.0.0.1:8384/rest/config/defaults/ignores \
  -H "X-Api-Key: KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "lines": ["*.log", ".DS_Store", "Thumbs.db"]
  }'
```

---

## 数据库 API

### GET /rest/db/status

获取文件夹数据库状态。

**参数**: `folder` (必需)

**响应**:
```json
{
  "globalBytes": 1234567890,
  "globalDeleted": 100,
  "globalDirectories": 50,
  "globalFiles": 1000,
  "globalSymlinks": 10,
  "globalTotalItems": 1160,
  "ignorePatterns": false,
  "inSyncBytes": 1234567890,
  "inSyncFiles": 1000,
  "localBytes": 1234567890,
  "localDeleted": 100,
  "localDirectories": 50,
  "localFiles": 1000,
  "localSymlinks": 10,
  "localTotalItems": 1160,
  "needBytes": 0,
  "needDeletes": 0,
  "needDirectories": 0,
  "needFiles": 0,
  "needSymlinks": 0,
  "needTotalItems": 0,
  "pullErrors": 0,
  "receiveOnlyChangedBytes": 0,
  "receiveOnlyChangedDeletes": 0,
  "receiveOnlyChangedDirectories": 0,
  "receiveOnlyChangedFiles": 0,
  "receiveOnlyChangedSymlinks": 0,
  "sequence": 12345,
  "state": "idle",
  "stateChanged": "2024-01-01T00:00:00Z",
  "version": 12345
}
```

### GET /rest/db/completion

获取完成状态。

**参数**:
- `folder` (可选)
- `device` (可选)

**响应**:
```json
{
  "completion": 100,
  "globalBytes": 1234567890,
  "needBytes": 0,
  "needItems": 0
}
```

### GET /rest/db/file

获取文件信息。

**参数**:
- `folder` (必需)
- `file` (必需): 文件路径

**响应**:
```json
{
  "name": "path/to/file.txt",
  "type": "FILE",
  "permissions": "0644",
  "modified": "2024-01-01T00:00:00Z",
  "modifiedBy": "DEVICEID",
  "size": 12345,
  "version": [...],
  "sequence": 123,
  "blockSize": 131072,
  "blocks": [
    {
      "offset": 0,
      "size": 131072,
      "hash": "base64hash..."
    }
  ]
}
```

### GET /rest/db/localchanged

获取本地更改的文件。

**参数**: `folder` (必需)

**响应**:
```json
{
  "files": [
    {
      "name": "file.txt",
      "type": "FILE",
      "size": 12345,
      "modified": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "perpage": 1000,
  "total": 1
}
```

### GET /rest/db/need

获取需要同步的文件。

**参数**: `folder` (必需)

**响应**: 类似 localchanged

### GET /rest/db/remotechanged

获取远程更改的文件。

### POST /rest/db/scan

触发文件夹扫描。

**参数**:
- `folder` (必需)
- `sub` (可选): 子路径
- `next` (可选): 下次扫描延迟（秒）

```bash
# 扫描整个文件夹
curl -X POST "http://127.0.0.1:8384/rest/db/scan?folder=myfolder" \
  -H "X-Api-Key: KEY"

# 扫描特定路径
curl -X POST "http://127.0.0.1:8384/rest/db/scan?folder=myfolder&sub=path/to/dir" \
  -H "X-Api-Key: KEY"
```

### POST /rest/db/override

覆盖更改（仅发送模式）。

**参数**: `folder` (必需)

```bash
curl -X POST "http://127.0.0.1:8384/rest/db/override?folder=myfolder" \
  -H "X-Api-Key: KEY"
```

### POST /rest/db/revert

还原更改（仅接收模式）。

**参数**: `folder` (必需)

### GET /rest/db/ignores

获取忽略模式。

**参数**: `folder` (必需)

**响应**:
```json
{
  "ignore": ["*.log", ".DS_Store"],
  "expanded": ["*.log", ".DS_Store"]
}
```

### PUT /rest/db/ignores

设置忽略模式。

```bash
curl -X PUT http://127.0.0.1:8384/rest/db/ignores \
  -H "X-Api-Key: KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "folder": "myfolder",
    "lines": ["*.log", ".DS_Store", "Thumbs.db"]
  }'
```

---

## 统计 API

### GET /rest/stats/folder

获取文件夹统计。

**响应**:
```json
{
  "myfolder": {
    "lastScan": "2024-01-01T00:00:00Z",
    "lastFile": "path/to/file.txt",
    "lastFileReceivedAt": "2024-01-01T00:00:00Z"
  }
}
```

### GET /rest/stats/device

获取设备统计。

**响应**:
```json
{
  "DEVICEID": {
    "lastSeen": "2024-01-01T00:00:00Z",
    "lastSeenDays": 0.5
  }
}
```

---

## 事件 API

### GET /rest/events

获取事件流（Server-Sent Events）。

**参数**:
- `since` (必需): 上一个事件 ID
- `limit` (可选): 返回事件数量
- `timeout` (可选): 超时时间（秒）
- `types` (可选): 事件类型过滤

#### 获取历史事件

```bash
curl "http://127.0.0.1:8384/rest/events?since=0&limit=10" \
  -H "X-Api-Key: KEY"
```

**响应**:
```json
[
  {
    "id": 1,
    "type": "Startup",
    "time": "2024-01-01T00:00:00Z",
    "data": {}
  },
  {
    "id": 2,
    "type": "DeviceConnected",
    "time": "2024-01-01T00:01:00Z",
    "data": {
      "id": "DEVICEID",
      "addr": "192.168.1.100:22000",
      "name": "My Laptop"
    }
  }
]
```

#### 实时事件流

```bash
curl -N "http://127.0.0.1:8384/rest/events?since=last" \
  -H "X-Api-Key: KEY"
```

#### 过滤事件类型

```bash
curl "http://127.0.0.1:8384/rest/events?since=0&types=DeviceConnected,DeviceDisconnected" \
  -H "X-Api-Key: KEY"
```

### 主要事件类型

#### 设备事件

**DeviceConnected**:
```json
{
  "type": "DeviceConnected",
  "data": {
    "id": "DEVICEID",
    "addr": "192.168.1.100:22000",
    "name": "My Laptop",
    "clientVersion": "v1.27.0"
  }
}
```

**DeviceDisconnected**:
```json
{
  "type": "DeviceDisconnected",
  "data": {
    "id": "DEVICEID",
    "error": "unexpected EOF"
  }
}
```

**DeviceDiscovered**:
```json
{
  "type": "DeviceDiscovered",
  "data": {
    "device": "DEVICEID",
    "addrs": ["192.168.1.100:22000"]
  }
}
```

#### 文件夹事件

**FolderCompletion**:
```json
{
  "type": "FolderCompletion",
  "data": {
    "device": "DEVICEID",
    "folder": "myfolder",
    "completion": 100,
    "globalBytes": 1234567890,
    "needBytes": 0
  }
}
```

**FolderErrors**:
```json
{
  "type": "FolderErrors",
  "data": {
    "folder": "myfolder",
    "errors": [
      {
        "path": "file.txt",
        "error": "permission denied"
      }
    ]
  }
}
```

**FolderSummary**:
```json
{
  "type": "FolderSummary",
  "data": {
    "folder": "myfolder",
    "state": "idle",
    "globalFiles": 1000,
    "localFiles": 1000,
    "needFiles": 0
  }
}
```

#### 文件事件

**ItemStarted**:
```json
{
  "type": "ItemStarted",
  "data": {
    "folder": "myfolder",
    "item": "path/to/file.txt",
    "type": "file",
    "action": "update"
  }
}
```

**ItemFinished**:
```json
{
  "type": "ItemFinished",
  "data": {
    "folder": "myfolder",
    "item": "path/to/file.txt",
    "type": "file",
    "action": "update",
    "error": null,
    "duration": 123456789
  }
}
```

#### 系统事件

**Startup**:
```json
{
  "type": "Startup",
  "data": {}
}
```

**StateChanged**:
```json
{
  "type": "StateChanged",
  "data": {
    "folder": "myfolder",
    "from": "scanning",
    "to": "idle",
    "when": "2024-01-01T00:00:00Z"
  }
}
```

**ConfigSaved**:
```json
{
  "type": "ConfigSaved",
  "data": {
    "version": 52,
    "folders": 3,
    "devices": 2
  }
}
```

**DownloadProgress**:
```json
{
  "type": "DownloadProgress",
  "data": {
    "folder": "myfolder",
    "updates": [
      {
        "name": "file.txt",
        "total": 10,
        "pulling": 3,
        "copiedFromOrigin": 0,
        "reused": 2,
        "pulled": 5,
        "copiedFromElsewhere": 0,
        "pullingFrom": 0,
        "bytesDone": 655360,
        "bytesTotal": 1310720
      }
    ]
  }
}
```

---

## 服务 API

### GET /rest/svc/report

获取使用报告。

### GET /rest/svc/deviceid

将设备名称转换为设备 ID。

**参数**: `name` (必需)

```bash
curl "http://127.0.0.1:8384/rest/svc/deviceid?name=DEVICEID" \
  -H "X-Api-Key: KEY"
```

**响应**:
```json
{
  "id": "ABCDEF1-ABCDEF1-ABCDEF1-ABCDEF1-ABCDEF1-ABCDEF1-A"
}
```

### GET /rest/svc/lang

获取支持的语言列表。

**响应**:
```json
["en", "zh-CN", "zh-TW", "ja", "ko", "fr", "de", ...]
```

### GET /rest/svc/langs

获取语言详细信息。

### GET /rest/svc/ignores

解析忽略模式。

**POST 请求**:
```bash
curl -X POST http://127.0.0.1:8384/rest/svc/ignores \
  -H "X-Api-Key: KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "folder": "myfolder",
    "lines": ["*.log", "!important.log"]
  }'
```

### GET /rest/svc/random/string

生成随机字符串。

**参数**: `length` (可选，默认 16)

```bash
curl "http://127.0.0.1:8384/rest/svc/random/string?length=32" \
  -H "X-Api-Key: KEY"
```

**响应**:
```json
{
  "random": "aBcDeFgHiJkLmNoPqRsTuVwXyZ123456"
}
```

### GET /rest/svc/lsan

获取 LAN 设备地址。

### POST /rest/svc/report

发送使用报告。

---

## 错误处理

### HTTP 状态码

- **200**: 成功
- **400**: 请求错误
- **403**: 禁止访问（认证失败）
- **404**: 未找到
- **500**: 服务器错误

### 错误响应格式

```json
{
  "error": "Error message here"
}
```

### 常见错误

#### 403 Forbidden
```json
{
  "error": "incorrect API key"
}
```

**原因**: API Key 不正确

**解决**: 检查 X-Api-Key header

#### 404 Not Found
```json
{
  "error": "folder not found"
}
```

**原因**: 文件夹或设备不存在

**解决**: 检查 ID 是否正确

#### 400 Bad Request
```json
{
  "error": "invalid configuration"
}
```

**原因**: 请求参数无效

**解决**: 检查 JSON 格式和必填字段

### 重试策略

```python
import time
import requests

def api_request_with_retry(url, api_key, max_retries=3):
    headers = {"X-Api-Key": api_key}
    
    for attempt in range(max_retries):
        try:
            response = requests.get(url, headers=headers, timeout=10)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            if attempt == max_retries - 1:
                raise
            time.sleep(2 ** attempt)  # 指数退避
```

---

## 使用示例

### Python 示例

#### 基础客户端

```python
import requests

class SyncthingAPI:
    def __init__(self, base_url, api_key):
        self.base_url = base_url.rstrip('/')
        self.headers = {
            'X-Api-Key': api_key,
            'Content-Type': 'application/json'
        }
    
    def get(self, endpoint, params=None):
        url = f"{self.base_url}{endpoint}"
        response = requests.get(url, headers=self.headers, params=params)
        response.raise_for_status()
        return response.json()
    
    def post(self, endpoint, data=None):
        url = f"{self.base_url}{endpoint}"
        response = requests.post(url, headers=self.headers, json=data)
        response.raise_for_status()
        return response.json()
    
    def get_status(self):
        return self.get('/rest/system/status')
    
    def get_folders(self):
        return self.get('/rest/config/folders')
    
    def get_devices(self):
        return self.get('/rest/config/devices')
    
    def pause_folder(self, folder_id):
        return self.post(f'/rest/system/pause?folder={folder_id}')
    
    def resume_folder(self, folder_id):
        return self.post(f'/rest/system/resume?folder={folder_id}')
    
    def scan_folder(self, folder_id, sub=None):
        params = {'folder': folder_id}
        if sub:
            params['sub'] = sub
        return self.post('/rest/db/scan', params)

# 使用示例
api = SyncthingAPI('http://127.0.0.1:8384', 'your-api-key')

# 获取状态
status = api.get_status()
print(f"Version: {status['version']}")

# 获取文件夹
folders = api.get_folders()
for folder in folders:
    print(f"Folder: {folder['label']} ({folder['id']})")

# 暂停文件夹
api.pause_folder('myfolder')
```

#### 监控事件

```python
import json
import requests

def monitor_events(api_key, since=0):
    url = 'http://127.0.0.1:8384/rest/events'
    headers = {'X-Api-Key': api_key}
    params = {'since': since, 'timeout': 60}
    
    while True:
        response = requests.get(url, headers=headers, params=params, stream=True)
        
        for line in response.iter_lines():
            if line:
                event = json.loads(line)
                print(f"Event: {event['type']}")
                print(f"Data: {event['data']}")
                print("---")
                
                since = event['id']
                params['since'] = since

monitor_events('your-api-key')
```

### Bash 脚本示例

#### 同步状态监控

```bash
#!/bin/bash

API_KEY="your-api-key"
BASE_URL="http://127.0.0.1:8384"

echo "=== Syncthing Status ==="
echo

# 获取系统信息
version=$(curl -s -H "X-Api-Key: $API_KEY" \
  $BASE_URL/rest/system/version | jq -r '.version')
echo "Version: $version"
echo

# 获取文件夹状态
echo "Folders:"
curl -s -H "X-Api-Key: $API_KEY" \
  $BASE_URL/rest/folders | jq -r '.[] | "  - \(.label) (\(.id)): \(.type)"'
echo

# 获取设备状态
echo "Devices:"
curl -s -H "X-Api-Key: $API_KEY" \
  $BASE_URL/rest/system/connections | jq -r '.connections | to_entries[] | 
  "  - \(.value.name): \(.value.connected)"'
echo

# 获取需要同步的文件数
echo "Pending Items:"
curl -s -H "X-Api-Key: $API_KEY" \
  "$BASE_URL/rest/db/status?folder=myfolder" | jq '.needTotalItems'
```

#### 自动备份配置

```bash
#!/bin/bash

API_KEY="your-api-key"
BASE_URL="http://127.0.0.1:8384"
BACKUP_DIR="/backup/syncthing"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# 备份配置
curl -s -H "X-Api-Key: $API_KEY" \
  $BASE_URL/rest/system/config > \
  $BACKUP_DIR/config_$DATE.json

# 压缩
cd $BACKUP_DIR
tar czf config_$DATE.tar.gz config_$DATE.json
rm config_$DATE.json

echo "Backup created: config_$DATE.tar.gz"

# 保留最近7天的备份
find $BACKUP_DIR -name "config_*.tar.gz" -mtime +7 -delete
```

### JavaScript/Node.js 示例

```javascript
const fetch = require('node-fetch');

class SyncthingClient {
  constructor(baseUrl, apiKey) {
    this.baseUrl = baseUrl.replace(/\/$/, '');
    this.apiKey = apiKey;
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseUrl}${endpoint}`;
    const response = await fetch(url, {
      ...options,
      headers: {
        'X-Api-Key': this.apiKey,
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    return response.json();
  }

  async getStatus() {
    return this.request('/rest/system/status');
  }

  async getFolders() {
    return this.request('/rest/config/folders');
  }

  async pauseFolder(folderId) {
    return this.request(`/rest/system/pause?folder=${folderId}`, {
      method: 'POST',
    });
  }

  async scanFolder(folderId) {
    return this.request(`/rest/db/scan?folder=${folderId}`, {
      method: 'POST',
    });
  }
}

// 使用示例
const client = new SyncthingClient(
  'http://127.0.0.1:8384',
  'your-api-key'
);

client.getStatus()
  .then(status => console.log('Version:', status.version))
  .catch(err => console.error(err));
```

---

## 最佳实践

### 1. 错误处理

```python
try:
    response = requests.get(url, headers=headers)
    response.raise_for_status()
    data = response.json()
except requests.exceptions.HTTPError as e:
    if response.status_code == 403:
        print("API Key 错误")
    elif response.status_code == 404:
        print("资源不存在")
    else:
        print(f"HTTP 错误: {e}")
except requests.exceptions.ConnectionError:
    print("连接失败，检查 Syncthing 是否运行")
except requests.exceptions.Timeout:
    print("请求超时")
```

### 2. 批量操作

```python
# 批量暂停文件夹
folders = api.get_folders()
for folder in folders:
    try:
        api.pause_folder(folder['id'])
        print(f"Paused: {folder['label']}")
    except Exception as e:
        print(f"Failed to pause {folder['label']}: {e}")
```

### 3. 事件监控

```python
# 持续监控事件
last_id = 0
while True:
    try:
        events = get_events(since=last_id)
        for event in events:
            process_event(event)
            last_id = event['id']
    except Exception as e:
        print(f"Error: {e}")
        time.sleep(5)  # 错误后等待
```

### 4. 速率限制

API 本身没有限速，但建议：
- 不要过于频繁地轮询
- 使用事件订阅代替轮询
- 批量操作减少请求数

### 5. 安全性

- ✅ 使用 HTTPS (生产环境)
- ✅ 保护 API Key
- ✅ 限制 GUI 访问地址
- ✅ 定期更换 API Key

---

## 总结

Syncthing REST API 提供了完整的控制能力：

- **系统管理**: 启动、停止、重启、升级
- **配置管理**: 文件夹、设备的 CRUD 操作
- **状态监控**: 实时状态和统计信息
- **事件订阅**: 实时事件流
- **数据库查询**: 文件和索引信息
- **服务工具**: 辅助功能

通过 API，可以实现：
- 自动化管理脚本
- 自定义监控面板
- 集成到其他系统
- 批量操作管理

开始使用 API 吧！🚀
