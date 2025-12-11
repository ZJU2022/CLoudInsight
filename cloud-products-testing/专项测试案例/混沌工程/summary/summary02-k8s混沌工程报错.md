# K8s混沌工程实战报错与解决核心总结

## 1. 典型问题

### 1.1 Chaos Mesh Webhook 报错
- 报错信息：
  - `failed calling webhook "mpodchaos.kb.io": ... connect: no route to host`
  - `connect: connection refused`
- 原因：
  - chaos-mesh-controller-manager Pod未正常监听443端口或Pod未Running
  - Service targetPort与Pod实际监听端口不一致
  - K8s网络（CNI/flannel）异常，Pod间网络不通

### 1.2 Pod 网络不通/Service ClusterIP不可达
- 现象：
  - curl/nc 10.96.0.1:443 卡住
  - flannel.1 网卡缺失
  - Pod Pending/CrashLoopBackOff
- 原因：
  - CNI插件未安装或异常
  - 节点内核参数/网络模块未加载（如br_netfilter）
  - 节点iptables残留、路由异常

### 1.3 StatefulSet/PVC绑定失败
- 现象：
  - Pod Pending，提示 `pod has unbound immediate PersistentVolumeClaims`
  - StorageClass为no-provisioner，需手动创建PV
- 原因：
  - 没有动态存储，PVC无法自动绑定

### 1.4 ErrImageNeverPull
- 现象：
  - Pod调度到某节点，报 `ErrImageNeverPull`
- 原因：
  - 该节点本地无镜像，且imagePullPolicy: Never

### 1.5 Zookeeper CrashLoopBackOff
- 现象：
  - Pod反复重启，Exit Code: 2
- 原因：
  - ZOO_MY_ID 环境变量配置错误，不能用metadata.name或metadata.ordinal

---

## 2. 排查与解决思路

### 2.1 网络与CNI问题
- 检查flannel Pod、网卡、路由、iptables、br_netfilter、sysctl
- 如无flannel.1网卡，彻底清理CNI残留并重启kubelet/docker
- 如极端死锁，reset节点并重新join

### 2.2 Chaos Mesh Webhook问题
- 检查controller-manager Pod状态、Service/Endpoints、Pod监听端口
- 确认Service targetPort与Pod实际监听端口一致
- 检查Pod日志定位端口/证书/网络问题

### 2.3 PVC与存储问题
- StorageClass为no-provisioner时，需手动创建PV
- 测试环境可用emptyDir临时存储规避绑定问题

### 2.4 镜像拉取问题
- 所有节点都需本地加载镜像（docker load）
- 或用nodeSelector强制调度到有镜像的节点

### 2.5 Zookeeper环境变量问题
- ZOO_MY_ID不能用fieldRef: metadata.name/ordinal
- 推荐用ConfigMap或initContainer动态生成

---

## 3. 关键数据与命令

- 检查flannel状态：
  ```bash
  kubectl get pods -n kube-flannel -o wide
  ip a | grep flannel
  ip route
  sudo systemctl restart kubelet
  sudo systemctl restart docker
  ```
- 检查Service/Endpoints/Pod端口：
  ```bash
  kubectl -n chaos-mesh get svc chaos-mesh-controller-manager
  kubectl -n chaos-mesh get endpoints chaos-mesh-controller-manager
  kubectl -n chaos-mesh exec -it <pod> -- netstat -tnlp
  ```
- 镜像分发：
  ```bash
  docker save -o zookeeper-latest.tar zookeeper:latest
  scp ...
  docker load -i ...
  ```
- PVC/存储：
  ```bash
  kubectl get pvc -n <ns>
  kubectl get pv
  ```

---

## 4. 最终落地SOP

### 4.1 网络与CNI恢复
1. 检查flannel Pod、网卡、路由
2. 如无flannel.1，清理CNI并重启kubelet/docker
3. 如极端死锁，reset节点并join

### 4.2 Chaos Mesh Webhook修复
1. 检查controller-manager Pod状态
2. 检查Service targetPort与Pod监听端口一致
3. Pod日志定位端口/证书/网络问题

### 4.3 PVC与存储
1. 测试环境用emptyDir
2. 生产环境需手动创建PV

### 4.4 镜像拉取
1. 所有节点都需本地docker load镜像
2. 或用nodeSelector强制调度

### 4.5 Zookeeper多副本配置
1. ZOO_MY_ID用ConfigMap或initContainer动态生成
2. 避免用fieldRef: metadata.name/ordinal

---

## 5. 总结
- K8s混沌工程常见问题多为网络、存储、镜像、环境变量配置
- 排查需结合Pod状态、事件、日志、节点网络、存储、镜像分布等多维度
- 解决方案需结合实际环境灵活选用，测试环境可用emptyDir，生产需持久化
- 关键是理解K8s调度、CNI、存储、镜像、环境变量的本质机制