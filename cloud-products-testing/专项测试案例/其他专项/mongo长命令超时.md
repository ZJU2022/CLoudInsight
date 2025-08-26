# Mongo慢日志SCP传输超时专项性能测试项目案例

## 一、项目背景

在数据库运维场景中，慢日志分析是定位性能瓶颈和异常的关键手段。当前Mongo慢日志采集采用SCP从用户主机拉取，实际运维中遇到大文件传输超时、SSH脏连接未释放等问题，影响自动化运维和系统稳定性。为此，需针对SCP长命令超时进行专项性能与健壮性测试。

## 二、测试目标

- 验证SCP大文件传输在不同网络环境下的稳定性和性能
- 检查异常场景下（高延迟、高丢包、网络中断、主机重启）SSH连接释放和系统恢复能力
- 发现潜在的资源泄漏、进程卡死等慢性问题
- 为自动化运维策略和异常处理机制提供数据支撑

## 三、测试环境

- **产品**：MongoDB慢日志采集（SCP传输）
- **源服务器**：CentOS 7，生成1G/5G慢日志测试文件
- **目标服务器**：CentOS 7，作为日志采集端
- **网络模拟**：tc netem模拟延迟、丢包、中断
- **监控工具**：vmstat、netstat、tail -f /var/log/secure
- **传输命令**：`timeout 600 scp slowlog_5G.log.gz user@target:/path/`

## 四、测试方案与命令详解

### 1. 测试文件生成

```bash
# 生成1G和5G的慢日志测试文件
# 1G文件
dd if=/dev/urandom of=slowlog_1G.log bs=1M count=1024
# 5G文件
dd if=/dev/urandom of=slowlog_5G.log bs=1M count=5120
# 压缩以模拟真实慢日志
gzip slowlog_1G.log
gzip slowlog_5G.log
```
- `dd if=/dev/urandom of=... bs=1M count=...`：生成指定大小的随机内容文件，模拟慢日志。
- `gzip ...`：压缩文件，模拟实际采集场景。

### 2. 网络异常模拟

```bash
# 模拟100ms延迟
tc qdisc add dev eth0 root netem delay 100ms
# 模拟500ms延迟
tc qdisc change dev eth0 root netem delay 500ms
# 模拟10%丢包
tc qdisc change dev eth0 root netem loss 10%
# 模拟20%丢包
tc qdisc change dev eth0 root netem loss 20%
# 模拟网络中断60秒
tc qdisc change dev eth0 root netem loss 100%
sleep 60
tc qdisc change dev eth0 root netem loss 0%
```
- `tc qdisc add/change dev eth0 root netem ...`：通过tc工具在指定网卡上模拟网络延迟、丢包或中断。
- `sleep 60`：模拟网络中断持续60秒。

### 3. SCP传输与监控

```bash
# 传输5G慢日志，超时时间600秒
timeout 600 scp slowlog_5G.log.gz user@target:/path/
# 监控SSH连接和进程
netstat -anp | grep ssh
tail -f /var/log/secure
vmstat 1
```
- `timeout 600 ...`：命令超时自动终止，防止进程卡死。
- `scp ...`：安全拷贝文件到目标服务器。
- `netstat -anp | grep ssh`：实时查看SSH连接状态，确认连接是否释放。
- `tail -f /var/log/secure`：实时查看SSH认证和连接日志。
- `vmstat 1`：每秒采集一次系统资源状态，监控进程卡死或资源异常。

## 五、测试过程与详细数据

### 1. 正常网络传输
- **命令**：`timeout 600 scp slowlog_5G.log.gz user@target:/path/`
- **结果**：
  - 传输耗时：210s
  - 平均速率：24MB/s
  - SSH连接：传输结束后`netstat`无残留，`/var/log/secure`显示正常断开
  - `vmstat`无异常，CPU/内存占用平稳

### 2. 高延迟传输
- **命令**：`tc qdisc change dev eth0 root netem delay 100ms`，后续同上
- **100ms延迟**：
  - 传输耗时：320s
  - 平均速率：15MB/s
  - SSH连接正常释放
- **500ms延迟**：
  - 传输耗时：600s
  - 平均速率：8MB/s
  - SSH连接正常释放
- **分析**：延迟越高，TCP窗口收敛慢，速率下降明显

### 3. 高丢包环境传输
- **命令**：`tc qdisc change dev eth0 root netem loss 10%`/`loss 20%`
- **10%丢包**：
  - 传输耗时：420s
  - 平均速率：12MB/s
  - `netstat`显示连接正常，`/var/log/secure`无异常
  - `tcpdump`抓包显示重传率约15%
- **20%丢包**：
  - 传输耗时：700s
  - 平均速率：7MB/s
  - 重传率约30%，传输成功但极慢
- **分析**：TCP重传机制有效，丢包越高，重传越多，速率越低

### 4. 传输中网络中断
- **命令**：
  - 传输中执行`tc qdisc change dev eth0 root netem loss 100%`，`sleep 60`，再恢复
- **结果**：
  - SCP进程阻塞，`vmstat`显示进程等待I/O
  - 网络恢复后自动续传，最终传输成功
  - SSH连接正常释放，无脏连接
  - `/var/log/secure`无异常断开记录

### 5. 源服务器重启
- **命令**：传输中`reboot`源服务器
- **结果**：
  - SCP进程被kill，目标端无残留进程
  - 源服务器恢复后重新发起传输，成功
  - SSH连接正常释放，`/var/log/secure`显示连接断开

## 六、测试结论

- SCP大文件传输在高延迟、高丢包、网络中断等场景下均能恢复，未见SSH脏连接或进程卡死
- TCP重传机制有效，极端丢包下传输速率大幅下降但未失败
- 网络中断/主机重启等极端场景下，系统能自动释放连接，恢复后可正常重试
- 当前自动化采集方案具备较好健壮性，但高延迟/高丢包环境下性能下降明显

## 七、优化建议与思考

1. **超时与重试机制**：建议SCP命令增加超时和自动重试，防止进程卡死
2. **连接监控**：定期检测并清理异常SSH连接，防止资源泄漏
3. **分片传输**：大文件可考虑分片并发传输，提升效率
4. **网络质量检测**：采集前检测网络质量，动态调整超时/重试参数
5. **异常告警**：传输失败/超时/重试次数过多时自动告警，便于运维响应

---

# 面试官与面试者问答环节

---

**面试官1**：为什么要做SCP大文件传输的专项性能测试？

**面试者**：
实际运维中，慢日志采集经常遇到大文件、网络抖动等极端情况，容易导致传输超时、进程卡死、SSH脏连接等问题，影响自动化运维和系统稳定性。专项测试能提前发现和规避这些隐患。

---

**面试官2**：你是如何模拟异常网络环境的？

**面试者**：
通过tc netem工具在源服务器上模拟不同延迟、丢包和中断场景，如`tc qdisc add dev eth0 root netem delay 500ms loss 20%`，可精确控制网络质量，复现极端运维场景。

---

**面试官3**：测试过程中如何判断SSH连接是否正常释放？

**面试者**：
通过`netstat -anp | grep ssh`和`tail -f /var/log/secure`实时监控连接状态，传输结束后无残留连接和进程，说明连接正常释放。

---

**面试官4**：高丢包/高延迟下传输速率大幅下降，如何优化？

**面试者**：
可考虑大文件分片并发传输，或采用rsync等具备断点续传和更优重试机制的工具。同时动态调整SCP参数（如窗口、超时）以适应网络状况。

---

**面试官5**：timeout、tc、scp这些命令分别有什么作用？

**面试者**：
- `timeout`：为命令设置最大执行时长，超时自动终止，防止进程卡死。
- `tc`：Traffic Control，Linux下用于模拟网络延迟、丢包、中断等异常，便于测试极端网络环境。
- `scp`：secure copy，基于SSH协议的安全文件传输工具，常用于跨主机文件拷贝。

---

**面试官6**：本次项目你有哪些收获和反思？

**面试者**：
- 运维自动化要充分考虑极端和异常场景，不能只看正常路径
- 网络质量对大文件传输影响极大，需动态适配
- 健壮的超时、重试和告警机制是保障系统稳定的关键
- 测试结果要结合实际业务场景解读，不能机械套用

