# CloudFlare IP 列表更新工具

## 简介

这是一个用 Go 语言编写的自动化工具，用于定期更新 CloudFlare 的 IP 访问控制列表。该工具可以自动获取指定域名的 IP 地址，并将这些 IP 地址添加到 CloudFlare 的 IP 列表中，便于进行访问控制管理。

## 功能特点

- 自动查询指定域名的 IPv4 地址
- 支持添加固定 IP 地址
- 支持定期自动更新 IP 列表
- 可选择是否在每次更新前清空现有列表
- 详细的日志记录
- 完全基于配置文件，无需修改代码
- 支持多平台（Windows、Linux、macOS）和多架构（amd64、arm）

## 安装和使用

### 方法一：使用预编译的二进制文件（推荐）

1. 从 [发布页面](#) 下载适合您系统的预编译二进制文件
   - Windows: `cloudflare-custom-list-windows-amd64.exe` (Intel/AMD) 或 `cloudflare-custom-list-windows-arm.exe` (ARM)
   - Linux: `cloudflare-custom-list-linux-amd64` (Intel/AMD) 或 `cloudflare-custom-list-linux-arm` (ARM)
   - macOS: `cloudflare-custom-list-darwin-amd64` (Intel) 或 `cloudflare-custom-list-darwin-arm64` (Apple Silicon)

2. 创建配置文件 `config.yaml`（见下文）

3. 将二进制文件和配置文件放在同一目录下

4. 运行程序
   - Windows: 双击可执行文件或在命令行中运行 `cloudflare-custom-list-windows-amd64.exe`
   - Linux/macOS: 
     ```bash
     chmod +x cloudflare-custom-list-linux-amd64  # 添加执行权限
     ./cloudflare-custom-list-linux-amd64
     ```

### 方法二：从源码编译

#### 安装要求

- Go 1.13 或更高版本
- CloudFlare 账户及 API 凭证
- 以下 Go 依赖包：
  - github.com/cloudflare/cloudflare-go
  - gopkg.in/yaml.v2

## 快速开始

### 1. 克隆仓库或下载源代码

```bash
git clone [仓库地址]
cd [项目目录]
```


#### 编译步骤

1. 克隆仓库或下载源代码
   ```bash
   git clone [仓库地址]
   cd [项目目录]
   ```

2. 编译程序
   ```bash
   go build -o cloudflare-custom-list
   ```

   或使用提供的打包脚本编译多平台版本：
   ```bash
   chmod +x package.sh
   ./package.sh
   ```

3. 创建配置文件 `config.yaml`

4. 运行程序
   ```bash
   ./cloudflare-custom-list
   ```

## 配置文件

在程序同一目录下创建 `config.yaml` 文件，参考以下模板：

```yaml
api_key: xxxx # CloudFlare ApiKey
api_email: test@test.com # CloudFlare 账号
account_id: aaaa # CloudFlare 账户 ID
list_id: bbb # 管理账户，配置列表 ID
# CloudFlare 自定义列表只能添加 IP，此处域名会识别绑定的 ip 进行添加
domain_names:
  - example1.com
  - example2.com

# 不需要通过域名维护的，几乎固定不变的 IP
fixed_ips:
  - 1.1.1.1
# 定时更新的时间，单位分钟
interval_minutes: 5
# 是否替换列表，如果为 false，则会在列表后面追加 IP，如果为 true，则会先清空原有列表，再添加 IP
replace_list: true
```

### 3. 编译并运行程序

```bash
go build -o cf-ip-updater
./cf-ip-updater
```

## 配置参数说明

| 参数 | 类型 | 说明 |
|------|------|------|
| api_key | 字符串 | CloudFlare API 密钥 |
| api_email | 字符串 | CloudFlare 账户邮箱 |
| account_id | 字符串 | CloudFlare 账户 ID |
| list_id | 字符串 | 要更新的 IP 列表 ID，如果不存在则创建 |
| domain_names | 字符串数组 | 需要查询 IP 的域名列表，**注意：CloudFlare 自定义列表只能添加 IP，此处域名会识别绑定的 IP 进行添加** |
| fixed_ips | 字符串数组 | 固定添加的 IP 地址列表，适用于不需要通过域名维护的、几乎固定不变的 IP |
| interval_minutes | 整数 | 自动更新的时间间隔（分钟） |
| replace_list | 布尔值 | 是否在每次更新前清空现有列表。如果为 true，则会先清空原有列表再添加 IP；如果为 false，则会在列表后面追加 IP |

## 运行流程

1. 程序启动后，首先读取配置文件
2. 立即执行一次 IP 列表更新
3. 根据配置的时间间隔，定期自动更新 IP 列表
4. 所有操作日志保存在程序目录下的 `log.txt` 文件中

## 日志说明

程序运行时会在当前目录生成 `log.txt` 文件，记录所有操作过程和可能出现的错误。日志内容包括：

- 程序启动和配置读取信息
- IP 查询结果
- CloudFlare API 操作结果
- 可能出现的错误信息

## 使用场景

1. **动态 IP 环境**：当您的服务器使用动态 IP 地址时，可以通过域名自动追踪 IP 变化
2. **多服务器管理**：管理多个服务器或服务的 IP 访问控制
3. **安全访问控制**：只允许特定 IP 地址访问您的 CloudFlare 资源

## 注意事项

1. 确保您有足够的 CloudFlare API 权限
2. 域名解析可能返回多个 IP 地址，程序会添加所有 IPv4 地址
3. 默认情况下，更新间隔设置为 5 分钟，可根据需要在配置文件中调整
4. 建议开启 `replace_list` 选项以避免 IP 列表不断增长
5. 对于生产环境，建议将程序设置为系统服务自动运行

## 常见问题

### IP 列表不更新

- 检查 CloudFlare API 凭证是否正确
- 确认账户 ID 和列表 ID 是否正确
- 查看日志文件中是否有错误信息

### 程序无法启动

- 确保使用了正确版本的可执行文件（与您的操作系统和处理器架构匹配）
- 检查配置文件格式是否正确
- 确保配置文件与可执行文件在同一目录下
### 没有捕获到最新的 IP 变更

- 增加更新频率（减小 `interval_minutes` 值）
- 检查域名 DNS 记录是否已更新
- 检查日志文件中的 IP 查询结果
