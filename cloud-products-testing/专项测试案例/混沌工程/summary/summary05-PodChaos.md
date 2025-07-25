# PodChaos混沌工程核心精华总结

## 一、问题与目标
- **目标**：验证MySQL及其业务链路在外部依赖（如Zookeeper、Redis）Pod级别故障、资源瓶颈等极端场景下的高可用、容灾与降级能力。
- **业务链路**：api → core → action（备份）/resource，所有流量真实走业务链路。
- **环境**：三台Ubuntu自建K8s集群，研发已将公有云MySQL外部依赖指向本集群zk/redis。

---

## 二、核心混沌测试场景与落地细节

### 1. Zookeeper 故障
- **意义**：验证MySQL及业务在注册中心不可用时的容错与降级能力
- **注入方式**：PodChaos删除zk节点（单点/多点），或NetworkChaos阻断zk节点间通信
- **YAML示例（全量）**：
  ```yaml
  apiVersion: chaos-mesh.org/v1alpha1
  kind: PodChaos
  metadata:
    name: zk-all-failure
    namespace: chaos-test
  spec:
    action: pod-failure
    mode: all
    duration: 2m
    selector:
      labelSelectors:
        "app": "zookeeper-test"
  ```
- **验证点**：业务接口是否可用、降级逻辑是否生效、告警是否准确
- **排查要点**：kubectl get pod -n chaos-test | grep zookeeper

### 2. Redis 故障
- **意义**：验证缓存失效时业务降级能力
- **注入方式**：PodChaos删除redis主节点，或NetworkChaos模拟redis端口丢包
- **YAML示例（主节点故障）**：
  ```yaml
  apiVersion: chaos-mesh.org/v1alpha1
  kind: PodChaos
  metadata:
    name: redis-master-failure
    namespace: chaos-test
  spec:
    action: pod-failure
    mode: one
    duration: 2m
    selector:
      labelSelectors:
        "app": "redis"
  ```
- **验证点**：业务接口降级、本地缓存/直连、告警准确性
- **排查要点**：kubectl get pod -n chaos-test | grep redis

### 3. 服务雪崩/级联故障
- **意义**：验证多组件同时故障时系统极限降级能力
- **注入方式**：PodChaos同时删除zk和redis
- **YAML示例**：分别apply zk和redis的PodChaos
- **验证点**：业务是否还能以最小可用方式运行，降级链路是否健全

### 4. 流量洪峰/资源瓶颈
#### 4.1 管控API/业务服务压力测试（需Chaos Mesh）
- **意义**：验证限流、熔断、资源隔离
- **注入方式**：StressChaos对api/core/action/resource等服务施加CPU/内存压力
- **YAML示例（以api为例）**：
  ```yaml
  apiVersion: chaos-mesh.org/v1alpha1
  kind: StressChaos
  metadata:
    name: api-cpu-stress
    namespace: chaos-test
  spec:
    mode: one
    duration: 2m
    selector:
      labelSelectors:
        "app": "api"
    stressors:
      cpu:
        workers: 8
        load: 80
  ```
- **验证点**：接口响应、限流熔断、核心业务优先级
- **排查要点**：kubectl top pod -n chaos-test | grep api

#### 4.2 网络带宽/延迟（需Chaos Mesh）
- **意义**：验证组件间网络瓶颈、慢链路、超时降级
- **注入方式**：NetworkChaos对zk/redis/mysql等Pod限速或加延迟
- **YAML示例（zk限速）**：
  ```yaml
  apiVersion: chaos-mesh.org/v1alpha1
  kind: NetworkChaos
  metadata:
    name: zk-bandwidth-limit
    namespace: chaos-test
  spec:
    mode: one
    duration: 2m
    selector:
      labelSelectors:
        "app": "zookeeper-test"
    action: bandwidth
    bandwidth:
      rate: 1mbps
      limit: 1000
      buffer: 1000
  ```
- **验证点**：接口超时、降级、监控
- **排查要点**：kubectl describe networkchaos -n chaos-test

### 5. MySQL主节点故障
#### 5.1 直接kill主库进程/Pod（可用shell或接口）
- **意义**：验证主备切换、数据一致性
- **操作**：接口触发主库备份/切换，或shell kill主库进程/Pod
- **Shell示例**：
  ```sh
  kubectl delete pod <mysql主库pod名> -n <ns>
  # 或
  kubectl exec -it <mysql主库pod名> -n <ns> -- kill -9 <mysqld进程号>
  ```
- **验证点**：备库提升、业务无感知、数据一致

#### 5.2 复杂主库假死/网络分区（需底层网络权限，K8s下较难，建议跳过或专项推进）
- **意义**：模拟主库进程存活但网络不可达，验证脑裂、切换策略
- **注入方式**：需iptables等底层操作，K8s下高风险

### 6. 其他场景
- **存储故障/磁盘写满/元数据损坏**：涉及底层PVC/宿主机，K8s下高风险，建议跳过
- **复制异常/查询层/优雅重启/备库延迟/人为误操作**：部分可用接口+shell+SQL完成，部分需DBA配合，建议专项推进

---

## 三、哪些场景必须用Chaos Mesh？
- Pod级别故障注入（如zk/redis/mysql Pod故障）
- 压力/资源瓶颈注入（CPU/内存/网络带宽/延迟/丢包等）
- 多组件级联故障/雪崩场景
- 这些场景用Chaos Mesh最方便、可控、可回滚。

## 四、你可以用接口+shell完成的场景
- 主库进程/Pod kill
- 业务接口级别的流量/功能测试
- 部分SQL注入/参数变更

---

## 五、如需帮助
- 需要具体YAML模板、业务链路梳理、接口自动化脚本、混沌实验批量执行脚本等，随时告知你的具体需求和服务标签/Pod命名规范，我可以帮你定制。

---

## 六、常见问题与排查
- PodChaos selector未命中Pod：需确保Pod标签与selector一致。
- 资源apply报错（Cannot update chaos spec）：需先delete再apply。
- Pod重建慢/ErrImageNeverPull：pause等基础镜像未本地准备，需提前拉取。
- 故障注入无效：确认chaos-mesh组件、CRD、权限、目标Pod状态。

---

## 七、验证点与落地效果
- 业务接口可用性、降级链路、主备切换、数据一致性、告警准确性。
- 结合接口自动化/监控/日志，验证故障期间业务表现与恢复过程。
- 典型现象与数据可用于面试答辩、复盘、优化。

---

## 八、面试答辩要点
- 强调“真实业务链路+PodChaos/StressChaos/NetworkChaos”验证容灾的系统性和实战性。
- 展示对K8s集群、混沌工程工具、问题排查的深刻理解。
- 结合具体YAML、操作、现象，突出可落地、可复盘、可推广。

---

# 附录：典型混沌实验YAML与实战建议（摘自chaos.md）

## 1. 典型YAML片段

### Zookeeper单点Pod故障
```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: zk-leader-failure
  namespace: chaos-test
spec:
  action: pod-failure
  mode: one
  duration: 1m
  selector:
    labelSelectors:
      "statefulset.kubernetes.io/pod-name": "zookeeper-test-0"
```

### Redis主节点Pod故障
```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: redis-master-failure
  namespace: chaos-test
spec:
  action: pod-failure
  mode: one
  duration: 1m
  selector:
    labelSelectors:
      "app": "redis"
```

### 注册中心级联故障（服务雪崩）
```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: service-registry-cascade-failure
  namespace: chaos-test
spec:
  action: pod-failure
  mode: all
  duration: 2m
  selector:
    labelSelectors:
      "app": "zookeeper-test"
```

### 管控API压力测试（StressChaos）
```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: StressChaos
metadata:
  name: api-pressure-test
  namespace: chaos-test
spec:
  mode: one
  duration: 2m
  selector:
    labelSelectors:
      "app": "api-gateway"
  stressors:
    cpu:
      workers: 8
      load: 80
    memory:
      workers: 4
      size: 256MB
```

### 网络带宽瓶颈模拟（NetworkChaos）
```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: NetworkChaos
metadata:
  name: bandwidth-bottleneck
  namespace: chaos-test
spec:
  mode: one
  duration: 1m
  selector:
    labelSelectors:
      "app": "api-gateway"
  action: bandwidth
  bandwidth:
    rate: 1mbps
    limit: 1000
    buffer: 1000
```

## 2. 场景对比与实战建议

- **简单中断 vs 服务雪崩 vs 流量洪峰**：
  - 简单中断：单个Pod故障，验证基础容错（如主从切换、故障转移）。
  - 服务雪崩：多组件级联故障，验证系统极限降级、本地缓存、配置热更新。
  - 流量洪峰：资源瓶颈/压力，验证限流、熔断、资源隔离、优先级保护。

- **测试建议**：
  1. 先做简单中断，验证基础容错能力。
  2. 再做服务雪崩，验证降级策略。
  3. 最后做流量洪峰，验证流控机制。
  4. 可组合多种场景，验证系统极限。

- **注意事项**：
  - 所有实验前确保相关镜像已本地准备，避免Pod恢复慢。
  - 监控业务接口、降级链路、主备切换、数据一致性、告警准确性。
  - 不建议在生产环境直接演练高风险场景。

- **使用说明**：
  1. 根据实际环境修改labelSelectors。
  2. 调整duration和压力参数以适应测试需求。
  3. 可组合使用多个场景进行混合故障测试。
  4. 监控系统指标验证容灾和降级机制是否生效。