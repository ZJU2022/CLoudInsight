# K8s 网络/认证故障排查操作SOP

## 1. 问题现象
- kubectl 在 worker 节点执行 `kubectl get node` 报错：
  - `tls: failed to verify certificate: x509: certificate signed by unknown authority`
  - `the server has asked for the client to provide credentials`

## 2. 排查思路
- 检查 kubeconfig 文件的证书配置（certificate-authority、client-certificate、client-key 或 base64 data）。
- 判断是证书链问题还是认证凭证问题。
- 检查 API Server 地址和端口连通性。
- 检查 kubectl 版本一致性。

## 3. 关键数据
- kubeconfig 文件内容（certificate-authority-data、client-certificate-data、client-key-data 是否完整且未损坏）。
- `/etc/kubernetes/pki/ca.crt` 内容一致。
- server 字段指向的 API Server 地址可达。

## 4. 解决办法
### 4.1 证书链问题
- 如果 kubeconfig 使用 `certificate-authority` 路径，需将 master 节点的 `/etc/kubernetes/pki/ca.crt` 拷贝到 worker 节点相同路径，或修改为实际存在的路径。
- 推荐使用 base64 型 `certificate-authority-data`，避免路径依赖。

### 4.2 认证凭证问题
- kubeconfig 必须包含有效的 `client-certificate-data` 和 `client-key-data`（或对应文件路径和文件）。
- **最保险做法**：直接将 master 节点 `/etc/kubernetes/admin.conf` 拷贝到 worker 节点 `~/.kube/config`，并 `chmod 600 ~/.kube/config`。
- 确认 `users` 字段为 `kubernetes-admin`，且为 admin 权限。

### 4.3 其他建议
- 确认 kubectl 版本一致。
- 确认 config 文件权限和 owner 正确。
- 遇到 `the server has asked for the client to provide credentials`，优先怀疑 kubeconfig 不是 admin 权限或内容损坏。

## 5. 落地操作步骤
1. 在 master 节点执行：
   ```bash
   scp /etc/kubernetes/admin.conf root@k8s-worker1:/root/.kube/config
   ```
2. 在 worker 节点执行：
   ```bash
   chmod 600 ~/.kube/config
   kubectl get nodes
   ```
3. 如仍有问题，检查 API Server 地址连通性、kubectl 版本、kubeconfig 内容完整性。

---

如遇特殊报错或疑难问题，建议贴出 `kubectl config view` 和 `kubectl version` 结果，便于进一步分析。