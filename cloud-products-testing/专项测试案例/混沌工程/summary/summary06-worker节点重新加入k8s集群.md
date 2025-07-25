# worker1 节点 reset 并重新加入 K8s 集群 SOP

## 1. 在 worker1 节点 reset
```bash
sudo kubeadm reset -f
sudo systemctl restart docker || sudo systemctl restart containerd
sudo systemctl restart kubelet
```

## 2. 在 master 节点获取 join 命令
```bash
# 获取 join 命令（含 token 和 ca hash）
kubeadm token create --print-join-command
```

## 3. 在 worker1 节点执行 join 命令
```bash
# 复制上一步输出的 join 命令，在 worker1 执行
sudo kubeadm join <master-ip>:6443 --token <token> --discovery-token-ca-cert-hash sha256:<hash>
```

## 4. 检查节点状态
```bash
kubectl get nodes
```

---

**注意：**
- join 命令只能在 master 节点上生成。
- join 时如遇证书等报错，可加 `--v=5` 查看详细日志。
