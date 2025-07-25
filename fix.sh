## 🚀 一键修复脚本（推荐使用）

创建一个名为 `fix_k8s_network.sh` 的脚本文件：

```bash
#!/bin/bash

echo "=== Kubernetes 网络插件修复脚本 ==="
echo "正在强制删除 Calico 并安装 Weave Net..."

# 强制删除 Calico 资源
echo "1. 强制删除 Calico 资源..."
kubectl delete daemonset calico-node -n kube-system --force --grace-period=0 2>/dev/null || true
kubectl delete deployment calico-kube-controllers -n kube-system --force --grace-period=0 2>/dev/null || true
kubectl delete configmap calico-config -n kube-system --force --grace-period=0 2>/dev/null || true
kubectl delete serviceaccount calico-node -n kube-system --force --grace-period=0 2>/dev/null || true
kubectl delete serviceaccount calico-kube-controllers -n kube-system --force --grace-period=0 2>/dev/null || true
kubectl delete clusterrole calico-node --force --grace-period=0 2>/dev/null || true
kubectl delete clusterrole calico-kube-controllers --force --grace-period=0 2>/dev/null || true
kubectl delete clusterrolebinding calico-node --force --grace-period=0 2>/dev/null || true
kubectl delete clusterrolebinding calico-kube-controllers --force --grace-period=0 2>/dev/null || true

# 删除 CRD
echo "2. 删除 Calico CRD..."
kubectl delete crd felixconfigurations.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd ipamblocks.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd blockaffinities.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd ipamhandles.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd ipamconfigs.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd bgppeers.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd bgpconfigurations.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd ippools.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd hostendpoints.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd clusterinformations.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd globalnetworkpolicies.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd globalnetworksets.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd networkpolicies.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd networksets.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true

# 清理网络配置
echo "3. 清理网络配置..."
sudo rm -rf /etc/cni/net.d/*
sudo rm -rf /opt/cni/bin/*
sudo rm -rf /var/lib/calico
sudo rm -rf /var/run/calico

# 重启 containerd
echo "4. 重启 containerd..."
sudo systemctl restart containerd

# 等待一下
echo "5. 等待服务重启..."
sleep 5

# 安装 Weave Net
echo "6. 安装 Weave Net..."
kubectl apply -f "https://github.com/weaveworks/weave/releases/download/latest_release/weave-daemonset-k8s.yaml"

# 等待 Weave 启动
echo "7. 等待 Weave Net 启动..."
sleep 10

# 检查状态
echo "8. 检查安装状态..."
echo "=== Weave Net Pods ==="
kubectl get pods -n kube-system | grep weave

echo "=== 节点状态 ==="
kubectl get nodes

echo "=== 所有 Pods 状态 ==="
kubectl get pods --all-namespaces

echo "=== 修复完成！==="
```

### 使用方法：

```bash
# 创建脚本文件
cat > fix_k8s_network.sh << 'EOF'
#!/bin/bash

echo "=== Kubernetes 网络插件修复脚本 ==="
echo "正在强制删除 Calico 并安装 Weave Net..."

# 强制删除 Calico 资源
echo "1. 强制删除 Calico 资源..."
kubectl delete daemonset calico-node -n kube-system --force --grace-period=0 2>/dev/null || true
kubectl delete deployment calico-kube-controllers -n kube-system --force --grace-period=0 2>/dev/null || true
kubectl delete configmap calico-config -n kube-system --force --grace-period=0 2>/dev/null || true
kubectl delete serviceaccount calico-node -n kube-system --force --grace-period=0 2>/dev/null || true
kubectl delete serviceaccount calico-kube-controllers -n kube-system --force --grace-period=0 2>/dev/null || true
kubectl delete clusterrole calico-node --force --grace-period=0 2>/dev/null || true
kubectl delete clusterrole calico-kube-controllers --force --grace-period=0 2>/dev/null || true
kubectl delete clusterrolebinding calico-node --force --grace-period=0 2>/dev/null || true
kubectl delete clusterrolebinding calico-kube-controllers --force --grace-period=0 2>/dev/null || true

# 删除 CRD
echo "2. 删除 Calico CRD..."
kubectl delete crd felixconfigurations.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd ipamblocks.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd blockaffinities.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd ipamhandles.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd ipamconfigs.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd bgppeers.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd bgpconfigurations.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd ippools.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd hostendpoints.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd clusterinformations.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd globalnetworkpolicies.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd globalnetworksets.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd networkpolicies.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true
kubectl delete crd networksets.crd.projectcalico.org --force --grace-period=0 2>/dev/null || true

# 清理网络配置
echo "3. 清理网络配置..."
sudo rm -rf /etc/cni/net.d/*
sudo rm -rf /opt/cni/bin/*
sudo rm -rf /var/lib/calico
sudo rm -rf /var/run/calico

# 重启 containerd
echo "4. 重启 containerd..."
sudo systemctl restart containerd

# 等待一下
echo "5. 等待服务重启..."
sleep 5

# 安装 Weave Net
echo "6. 安装 Weave Net..."
kubectl apply -f "https://github.com/weaveworks/weave/releases/download/latest_release/weave-daemonset-k8s.yaml"

# 等待 Weave 启动
echo "7. 等待 Weave Net 启动..."
sleep 10

# 检查状态
echo "8. 检查安装状态..."
echo "=== Weave Net Pods ==="
kubectl get pods -n kube-system | grep weave

echo "=== 节点状态 ==="
kubectl get nodes

echo "=== 所有 Pods 状态 ==="
kubectl get pods --all-namespaces

echo "=== 修复完成！==="
EOF

# 给脚本执行权限
chmod +x fix_k8s_network.sh

# 运行脚本
./fix_k8s_network.sh
```

**直接复制上面的命令到 Master 节点执行即可！**