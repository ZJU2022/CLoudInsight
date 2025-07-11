# 数据库自动化测试设计思路

## 概述

本文档详细描述了数据库产品线的自动化测试设计思路，包括功能点梳理、自动化交付产物设计、测试策略制定以及具体实施方案。

**具体功能模块参考img/功能模块细化案例.png**

## 第一步：功能点与实例类型梳理

### 实例类型分类

#### 1. 存储介质
- **物理机**: 传统物理服务器部署
- **NVME（虚拟机）**: 基于NVME存储的虚拟机部署
- **云盘**: 基于云存储的部署方式

#### 2. 计费项
- **历史计费模式**: 以前未开放CPU计费
- **新计费模式**: 现在开放CPU计费

#### 3. 数据库版本
- **MySQL 5.5**: 经典版本
- **MySQL 5.6**: 稳定版本
- **MySQL 5.7**: 主流版本
- **MySQL 8.0**: 最新版本
- **Percona**: 企业级分支版本

#### 4. 部署模式
- **单点部署**: 单实例部署
- **高可用部署**: 主从架构部署
  - 跨可用区高可用
  - 同可用区高可用

### 功能模块分类

#### 1. 核心功能模块
- **升降级**: 实例规格调整
- **单点升级HA**: 单点实例升级为高可用
- **只读从库**: 读写分离功能
- **备份回档**: 数据备份与恢复
- **批量升降级**: 批量实例规格调整
- **大小版本升级**: 数据库版本升级
- **日志管理**: 日志查看和管理
- **磁盘热升级**: 在线磁盘扩容
- **Binlog转储**: 二进制日志导出

#### 2. 运维功能模块
- **计费管理**: 费用计算和账单
- **容灾功能**: 灾难恢复能力
- **基础功能**: 
  - 重启、关闭、启动
  - 删除实例
  - 监控告警
  - SSL证书管理

## 第二步：自动化交付产物设计

### 自动化交付策略

#### 1. 基础功能回归 + 关联集合模板
- **限时要求**: 30分钟内完成
- **并行测试**: 支持多个工作流同时运行
- **智能回归**: 
  - 当前改动影响的功能：运行细分场景
  - 未影响的功能：正常回归测试

#### 2. 关联集合模板设计
- **基础回归验证**: 覆盖MySQL基本功能
- **微服务覆盖**: 确保所有相关微服务正常启动
- **快速验证**: 发布时快速验证核心功能

#### 3. 回归测试工作流
设计六条核心工作流，根据发布可用区自动选择：
- 存储类型：NVME/云盘
- 部署模式：单点/高可用
- 实例规格：根据测试需求选择
- 数据库版本：根据兼容性要求选择

## 第三步：测试策略制定

### 测试分层设计

#### 1. 黑盒层测试
- **接口连续调用**: 模拟真实用户操作流程
- **SQL场景构造**: 通过SQL语句构造测试场景
- **异常场景设计**: 覆盖各种异常情况

#### 2. 灰盒层测试
- **微服务集成**: 验证服务间交互
- **数据一致性**: 确保数据在不同服务间的一致性
- **性能监控**: 监控关键性能指标

#### 3. 白盒层测试（较少）
- **代码分支覆盖**: 覆盖关键代码路径
- **只读从库逻辑**: 重点覆盖只读从库相关代码

#### 4. 平行层测试（例）
- **容灾自动化**: 
  - 宿主机关闭测试
  - 外部组件故障测试
- **网络测试**: iperf打流测试
- **内核性能测试**: 
  - sysbench基准测试
  - unixbench系统性能测试
  - iperf网络性能测试
  - fio磁盘性能测试
  - 圆周率计算测试
- **主机连通性测试**: telnet连接测试

### 测试场景示例

#### 完整测试流程示例
```
CreateHost → DescribeHost → telnet → CreateDatabase → Upgrade → 
GetMetric → CRUD → CreateReadonly → GetLog → BackupLogic → 
CreateFromBackup → StopDatabase → DeleteDatabase
```

## 第四步：实施计划

### 1. 周会对齐机制
- **频率**: 每周定期拉会
- **参与方**: 测试团队与业务方
- **内容**: 对齐每个场景的具体细节
- **目标**: 确保测试场景与实际业务需求一致

### 2. 场景构造方法
- **接口连续调用**: 设计完整的业务流程
- **SQL辅助构造**: 通过SQL语句构造特定场景
- **异常场景覆盖**: 设计各种异常情况的测试用例

### 3. 测试数据管理
- **模板化配置**: 提供标准化的测试配置模板
- **参数化设计**: 支持灵活的参数调整
- **数据隔离**: 确保测试数据不影响生产环境



## 总结

数据库自动化测试设计采用分层测试策略，通过功能点梳理、自动化交付产物设计、测试策略制定和实施计划，构建了一个完整的自动化测试体系。该体系能够有效保障数据库产品的质量，提高测试效率，支持快速迭代和发布。
