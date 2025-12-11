# K8s集群与Flannel网络核心问题与解决总结

## 1. 典型问题

### 1.1 Pod网络不通/Service ClusterIP不可达
- 现象：
  - curl/nc 10.96.0.1:443 卡住
  - Pod间网络不通，Service无法访问
  - flannel.1 网卡缺失
  - Pod Pending/CrashLoopBackOff
- 原因：
  - CNI插件（如flannel）未安装或异常
  - 节点内核参数/网络模块未加载（如br_netfilter）
  - 节点iptables残留、路由异常
  - 节点未加载CNI插件镜像

### 1.2 flannel Pod 启动失败
- 现象：
  - flannel Pod Pending/CrashLoopBackOff
  - 日志报错 `Failed to create SubnetManager: ... dial tcp 10.96.0.1:443: i/o timeout`
  - 没有 flannel.1 网卡
- 原因：
  - 节点无法访问apiserver Service IP（10.96.0.1:443）
  - CNI目录/配置缺失或权限错误
  - 镜像未加载到所有节点
  - 内核参数未设置（如br_netfilter）

### 1.3 CNI插件未安装/未生效
- 现象：
  - Pod网络完全不通
  - flannel.1、cni0 网卡缺失
- 原因：
  - 未安装flannel（kubectl apply -f kube-flannel.yml）
  - 节点CNI目录未创建

### 1.4 节点重装/重置后网络异常
- 现象：
  - 节点reset后Pod网络不通
  - flannel Pod无法启动
- 原因：
  - CNI残留未清理
  - 节点未重新join集群

---

## 2. 排查与解决思路

### 2.1 检查CNI与flannel状态
- 检查flannel Pod、网卡、路由、iptables、br_netfilter、sysctl
- 检查 /etc/cni/net.d/10-flannel.conflist 配置
- 检查所有节点是否已加载flannel相关镜像

### 2.2 彻底清理CNI残留
- 停止kubelet
- 删除 /etc/cni/net.d/*、/var/lib/cni/*、/run/flannel/*
- 重启docker/containerd、kubelet
- 删除并重建flannel Pod

### 2.3 检查内核参数与网络模块
- 加载br_netfilter模块
- 设置sysctl net.bridge.bridge-nf-call-iptables=1

### 2.4 镜像分发
- 所有节点都需本地加载flannel及相关镜像（docker load）

### 2.5 节点reset与重新加入
- 如极端死锁，reset节点并重新join集群

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
- 检查CNI配置：
  ```bash
  cat /etc/cni/net.d/10-flannel.conflist
  ls -l /opt/cni/bin
  ```
- 镜像分发：
  ```bash
  docker save -o flannel-v0.27.2.tar ghcr.io/flannel-io/flannel:v0.27.2
  scp ...
  docker load -i ...
  ```
- 清理CNI残留：
  ```bash
  sudo systemctl stop kubelet
  sudo rm -rf /etc/cni/net.d/*
  sudo rm -rf /var/lib/cni/
  sudo rm -rf /run/flannel/*
  sudo systemctl restart docker || sudo systemctl restart containerd
  sudo systemctl start kubelet
  kubectl delete pod -n kube-flannel -l app=flannel
  ```
- 检查内核参数：
  ```bash
  sudo modprobe br_netfilter
  sudo sysctl -w net.bridge.bridge-nf-call-iptables=1
  ```
- 节点reset与join：
  ```bash
  sudo kubeadm reset -f
  kubeadm token create --print-join-command
  sudo kubeadm join ...
  ```

---

## 4. 最终落地SOP

1. 检查所有节点flannel Pod、网卡、路由、iptables、CNI配置
2. 如无flannel.1，彻底清理CNI残留并重启kubelet/docker
3. 检查/加载br_netfilter模块，设置sysctl参数
4. 所有节点本地加载flannel相关镜像
5. 如极端死锁，reset节点并join
6. 检查Pod网络连通性（kubectl run nettest --image=busybox ...）

---

## 5. 总结
- K8s集群网络问题多为CNI插件未安装/异常、节点内核参数/网络模块未设置、镜像未分发、CNI残留未清理
- 排查需结合Pod状态、节点网络、CNI配置、内核参数、镜像分布等多维度
- 解决方案需结合实际环境灵活选用，测试环境可用emptyDir，生产需持久化
- 关键是理解K8s调度、CNI、网络、镜像的本质机制