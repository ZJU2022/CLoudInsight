# Zookeeper 在 Kubernetes 集群部署与排障核心总结

## 1. 问题现象
- Zookeeper StatefulSet 部署后，Pod 状态为 `ErrImageNeverPull` 或 `CrashLoopBackOff`，或长时间 `Pending` 无法调度。
- 日志报错：`serverid zookeeper-test-0 is not a number`。
- Pod 事件提示：`pod has unbound immediate PersistentVolumeClaims`。

## 2. 关键排查与数据
- 镜像已在所有节点用 `docker load` 导入，但 Pod 仍报找不到镜像。
- 使用 containerd 作为容器运行时，需用 `ctr` 导入镜像：
  ```bash
  sudo ctr -n k8s.io images import /path/to/zookeeper_latest.tar
  sudo ctr -n k8s.io images ls | grep zookeeper
  ```
- Pod 日志：
  ```
  Invalid config, exiting abnormally
  org.apache.zookeeper.server.quorum.QuorumPeerConfig$ConfigException: Error processing /conf/zoo.cfg
  Caused by: java.lang.IllegalArgumentException: serverid zookeeper-test-0 is not a number
  ```
- Pod describe 事件：
  ```
  Warning  FailedScheduling  pod has unbound immediate PersistentVolumeClaims
  ```

## 3. 问题分析
- **镜像问题**：K8s 使用 containerd，需用 `ctr` 导入镜像，`docker load` 无效。
- **环境变量问题**：Zookeeper 期望 `ZOO_MY_ID` 为数字（0,1,2），原配置用 Pod 名称导致报错。
- **存储问题**：StatefulSet 默认用 PVC，若无可用 PV，Pod 会 Pending。

## 4. 解决思路与落地操作

### 4.1 镜像问题解决
- 在所有节点用 `ctr` 导入镜像：
  ```bash
  sudo ctr -n k8s.io images import /path/to/zookeeper_latest.tar
  sudo ctr -n k8s.io images ls | grep zookeeper
  ```

### 4.2 环境变量修正
- `ZOO_MY_ID` 应为数字。推荐用：
  ```yaml
  env:
    - name: ZOO_MY_ID
      valueFrom:
        fieldRef:
          fieldPath: metadata.annotations['apps.kubernetes.io/pod-index']
  ```

### 4.3 存储问题快速绕过（开发/测试环境）
- 使用 `emptyDir` 临时卷，避免 PVC/PV 绑定问题：
  ```yaml
  volumes:
    - name: zookeeper-data
      emptyDir: {}
    - name: zookeeper-logs
      emptyDir: {}
  ```
- 这样 Pod 可直接调度运行，适合功能验证和临时集群。

### 4.4 持久化存储（生产建议）
- 需先创建 StorageClass 和本地 PV，确保 PVC 能绑定。
- 每个节点需有本地目录并授权：
  ```bash
  sudo mkdir -p /data/zookeeper-0 /data/zookeeper-1 /data/zookeeper-2
  sudo chmod 777 /data/zookeeper-*
  ```
- 创建 PV 时需指定 nodeAffinity 绑定到对应节点。

## 5. 验证与排查
- Pod 启动后，检查状态：
  ```bash
  kubectl get pods -n chaos-test -l app=zookeeper-test
  kubectl logs -n chaos-test zookeeper-test-0
  ```
- 检查 PVC/PV 绑定：
  ```bash
  kubectl get pvc -n chaos-test
  kubectl get pv
  ```
- 检查环境变量：
  ```bash
  kubectl exec -n chaos-test zookeeper-test-0 -- env | grep ZOO
  ```

## 6. 总结
- **镜像导入要用 containerd 工具**。
- **Zookeeper 集群 ID 必须为数字**，用 pod-index 注解自动赋值。
- **开发测试可用 emptyDir 跳过存储问题，生产需配置本地 PV。**
- **遇到 Pending 多半是 PVC 没绑定 PV，需先排查存储。**

> 本文档为混沌工程环境下 Zookeeper 在 K8s 部署与排障的实战精华总结。

---

## 7. 实战对话精华补充

### 7.1 典型问题现象
- Pod 一直 CrashLoopBackOff 或 Pending，PVC 处于 Pending，PV 状态 Released/Available。
- 日志报错：`serverid zookeeper-test-0 is not a number`、`UnknownHostException`、`unbound immediate PersistentVolumeClaims`。
- describe pod 发现 init container、环境变量、存储等问题。

### 7.2 关键排查与数据
- `kubectl get pv/pvc/pods` 显示 PVC 未绑定 PV，或 PV 状态 Released/Available。
- `kubectl describe pod` 发现 ZOO_MY_ID 配置为 Pod 名称或未为数字。
- `kubectl logs` 发现 ZooKeeper 配置异常、主机名无法解析（UnknownHostException）。
- StatefulSet 配置中 volumeClaimTemplates 数量与 PV 数量不匹配。
- PV 只创建了 3 个，但每个 Pod 需要 2 个 PVC，实际应创建 6 个 PV。
- 使用 containerd 运行时，镜像需用 ctr 导入。

### 7.3 问题分析
- **ZOO_MY_ID 配置错误**：不能用 Pod 名称，需为数字（0/1/2），否则 ZooKeeper 启动报错。
- **PV/PVC 数量不匹配**：每个 Pod 2 个 PVC，需 6 个 PV，否则 Pod Pending。
- **主机名解析失败**：ZOO_SERVERS 需用 FQDN（如 zookeeper-test-0.zookeeper-test.chaos-test.svc.cluster.local），否则集群通信失败。
- **镜像拉取失败**：containerd 环境下 docker load 无效，需用 ctr。

### 7.4 解决思路与落地操作
- **ZOO_MY_ID 自动赋值**：推荐用 initContainer 提取 $HOSTNAME 数字写入 /data/myid，主容器用 ZOO_MY_ID_FILE=/data/myid。
  ```yaml
  initContainers:
    - name: init-zookeeper
      image: zookeeper:latest
      imagePullPolicy: Never
      command: ['sh', '-c', 'echo $HOSTNAME | sed "s/zookeeper-test-//" > /data/myid']
      volumeMounts:
        - name: zookeeper-data
          mountPath: /data
  containers:
    - name: zookeeper-test
      env:
        - name: ZOO_MY_ID_FILE
          value: "/data/myid"
  ```
- **PV/PVC 配置**：每个 Pod 2 个 PVC（data/logs），需 6 个 PV，且 storageClass、hostPath 匹配。
  ```bash
  # 创建 6 个本地 PV
  sudo mkdir -p /data/zookeeper-data-0 /data/zookeeper-logs-0 /data/zookeeper-data-1 /data/zookeeper-logs-1 /data/zookeeper-data-2 /data/zookeeper-logs-2
  sudo chmod 777 /data/zookeeper-*
  # PV yaml 需 storageClass/hostPath 匹配 PVC
  ```
- **主机名解析**：ZOO_SERVERS 用 FQDN：
  ```yaml
  - name: ZOO_SERVERS
    value: "server.0=zookeeper-test-0.zookeeper-test.chaos-test.svc.cluster.local:2888:3888;2181 server.1=zookeeper-test-1.zookeeper-test.chaos-test.svc.cluster.local:2888:3888;2181 server.2=zookeeper-test-2.zookeeper-test.chaos-test.svc.cluster.local:2888:3888;2181"
  ```
- **镜像导入**：containerd 环境用 ctr 导入镜像。
  ```bash
  sudo ctr -n k8s.io images import /path/to/zookeeper_latest.tar
  ```
- **开发测试可用 emptyDir 跳过存储问题**。

### 7.5 验证与排查建议
- 检查 Pod、PVC、PV 状态，确保全部 Bound/Running。
- 检查 ZOO_MY_ID 是否为数字，主机名解析是否正常。
- 检查 ZooKeeper 日志，确认集群正常启动。

> 本节为本次多轮实战对话的核心精华补充，涵盖典型问题、排查、分析、解决与落地操作。