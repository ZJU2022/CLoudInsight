# Kubernetes集群+Chaos Mesh部署核心精华总结

## 一、部署思路与关键数据
- **目标**：三台Ubuntu 22.04虚拟机，部署Kubernetes集群（1 Master+2 Worker），并安装Chaos Mesh实现混沌工程。
- **推荐配置**：每台4C8G，公网/内网IP任选，需互通。
- **节点示例**：
  - Master01: 106.75.163.110
  - Worker01: 106.75.163.117
  - Worker02: 106.75.163.102

## 二、落地命令精华（按顺序执行）
### 1. 基础环境准备（所有节点）
```bash
# 设置主机名
sudo hostnamectl set-hostname k8s-master   # Master01
sudo hostnamectl set-hostname k8s-worker1  # Worker01
sudo hostnamectl set-hostname k8s-worker2  # Worker02
# hosts解析（三台都执行）
sudo bash -c 'cat >> /etc/hosts <<EOF
106.75.163.110 k8s-master
106.75.163.117 k8s-worker1
106.75.163.102 k8s-worker2
EOF'
# 关闭swap和防火墙
sudo swapoff -a
sudo sed -i '/ swap / s/^[^#]/#&/' /etc/fstab
sudo ufw disable
# 安装Docker
sudo apt-get update && sudo apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update && sudo apt-get install -y docker-ce docker-ce-cli containerd.io
sudo systemctl enable docker && sudo systemctl start docker
# 安装kubeadm/kubelet/kubectl
sudo apt-get update && sudo apt-get install -y apt-transport-https ca-certificates curl
sudo curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update && sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
```

### 2. 初始化Kubernetes集群
```bash
# Master01执行
sudo kubeadm init --pod-network-cidr=10.244.0.0/16 --apiserver-advertise-address=106.75.163.110
# 配置kubectl
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
# 安装Flannel网络插件
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
```

### 3. Worker节点加入集群
```bash
# Worker01/Worker02执行（用Master01输出的join命令）
sudo kubeadm join 106.75.163.110:6443 --token <token> --discovery-token-ca-cert-hash sha256:<hash>
```

### 4. 检查集群状态
```bash
kubectl get nodes   # Master01执行，三节点Ready即成功
```

### 5. 安装Chaos Mesh（Master01）
```bash
# 安装Helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
# 添加Chaos Mesh仓库
helm repo add chaos-mesh https://charts.chaos-mesh.org
helm repo update
# 安装Chaos Mesh
kubectl create namespace chaos-mesh
helm install chaos-mesh chaos-mesh/chaos-mesh -n chaos-mesh --set dashboard.create=true
# 验证
kubectl get pods -n chaos-mesh
# 访问Dashboard
kubectl port-forward -n chaos-mesh svc/chaos-mesh-dashboard 2333:2333
# 浏览器访问 http://localhost:2333
```

---

## 三、常见问题与SOP
### 1. 节点NotReady/网络不通
- 检查防火墙/安全组，确保6443、2379-2380、10250等端口开放。
- Flannel未正常运行：`kubectl get pods -n kube-system`，重装网络插件。
- 检查`/etc/hosts`和主机名一致。

### 2. kubeadm join失败
- Token过期：Master01上用`kubeadm token create --print-join-command`重新获取。
- CA哈希不符：用`openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | sha256sum | awk '{print $1}'`获取。

### 3. Chaos Mesh Pod不Running
- 资源不足：建议每节点4G+内存。
- 镜像拉取失败：检查网络或用国内镜像源。
- `kubectl describe pod ...` 查看详细报错。

### 4. kubectl命令无权限
- 检查`$HOME/.kube/config`权限和内容。

### 5. 端口转发/访问Dashboard失败
- 检查本地与Master01网络连通。
- 确认端口未被占用。

---

## 四、最佳实践Tips
- 所有命令建议用root或sudo执行，避免权限问题。
- 重要操作前快照虚拟机，便于回滚。
- 生产环境建议多Master高可用，测试/学习单Master足够。
- 关注[官方文档](https://kubernetes.io/zh/docs/setup/)和[Chaos Mesh文档](https://chaos-mesh.org/zh/docs/)

---

如需更多数据库/存储/高可用等集群部署方案，可随时补充！
