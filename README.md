# TempLog - 零配置临时日志库

TempLog 是一个简单易用的 Go 语言日志库，提供零配置的临时日志功能。只需导入包即可开始记录日志，无需任何初始化配置。
---

**注意**: TempLog 设计用于开发和临时环境。对于生产环境，建议根据具体需求进行更详细的配置。

## 特性

- 🚀 **零配置** - 导入即用，无需初始化
- 📁 **自动文件管理** - 自动创建日志目录和文件
- 📊 **分级日志** - 支持 Debug、Info、Warn、Error 级别
- 🔄 **日志轮转** - 自动进行日志轮转和压缩
- 🖥️ **控制台输出** - 同时输出到控制台和文件
- 📍 **调用者信息** - 自动记录调用文件和行号
- 🧹 **自动清理** - 自动关闭文件句柄

## 安装

```bash
go get github.com/leijux/templog
```

## 快速开始

```go
package main

import (
    _ "github.com/leijux/templog" // 只需导入即可启用日志
    "log/slog"
)

func main() {
    // 直接使用标准库的 slog 记录日志
    slog.Debug("这是一条调试日志", "key", "value")
    slog.Info("程序启动", "version", "1.0.0")
    slog.Warn("警告信息", "threshold", 80)
    slog.Error("错误发生", "err", "file not found")
}
```

## 日志文件结构

导入 TempLog 后，会自动在项目根目录创建 `logs/` 目录，结构如下：

```
logs/
├── debug/
│   └── debug.log      # 所有调试级别日志
├── info/
│   └── info.log       # 所有信息级别日志
├── warn/
│   └── warn.log       # 所有警告级别日志
├── error/
│   └── error.log      # 所有错误级别日志
└── all.log            # 所有级别的汇总日志
```

## 配置参数

虽然 TempLog 是零配置的，但内部使用以下默认配置：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| 日志目录 | `logs/` | 日志文件存储目录 |
| 文件最大大小 | 100 MB | 单个日志文件最大大小 |
| 最大备份数 | 3 | 保留的备份文件数量 |
| 最大保留天数 | 30 | 日志文件保留天数 |
| 压缩 | true | 是否压缩备份文件 |

## 使用标准库 slog

TempLog 与 Go 1.21+ 的标准库 `log/slog` 完全兼容。你可以使用所有标准的 slog 功能：

```go
import (
    "log/slog"
    "context"

    _ "github.com/leijux/templog"
)

func example() {
    // 基本日志记录
    slog.Info("用户登录", "user_id", 123, "ip", "192.168.1.1")
    
    // 带上下文的日志
    ctx := context.Background()
    slog.LogAttrs(ctx, slog.LevelInfo, "操作完成", 
        slog.String("operation", "create"),
        slog.Int("duration_ms", 150))
    
    // 结构化日志
    logger := slog.Default()
    logger.With("service", "api").Info("请求处理")
}
```

## 高级用法

### 自定义日志级别

```go
import (
    _ "github.com/leijux/templog"
    "log/slog"
)

func main() {
    // 设置日志级别（默认是 Debug）
    // 注意：TempLog 默认启用所有级别，但你可以通过 slog.SetLogLoggerLevel 控制输出
    slog.SetLogLoggerLevel(slog.LevelInfo)
}
```


## 开发指南

### 项目结构

```
templog/
├── templog.go     # 主要实现
├── go.mod         # 模块定义
└── README.md      # 本文档
```

### 依赖

- `github.com/rs/zerolog` - 高性能日志库
- `github.com/samber/slog-zerolog/v2` - slog 到 zerolog 的适配器
- `gopkg.in/natefinch/lumberjack.v2` - 日志轮转库
