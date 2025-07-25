# K8s节点无法直接拉取镜像时，flannel镜像本地导入SOP

## 1. 在本地Mac拉取flannel相关镜像

```bash
# 拉取flannel主镜像
docker pull ghcr.io/flannel-io/flannel:v0.27.2
# 拉取flannel-cni-plugin镜像
docker pull ghcr.io/flannel-io/flannel-cni-plugin:v1.7.1-flannel1
```

## 2. 保存镜像为tar文件

```bash
docker save -o flannel-v0.27.2.tar ghcr.io/flannel-io/flannel:v0.27.2
docker save -o flannel-cni-plugin-v1.7.1-flannel1.tar ghcr.io/flannel-io/flannel-cni-plugin:v1.7.1-flannel1
```

## 3. 上传镜像到K8s所有节点（master和worker）

```bash
# 以scp为例，假设节点用户名为ubuntu，IP为10.0.0.2
scp flannel-v0.27.2.tar ubuntu@106.75.163.110:/home/ubuntu/chaos
scp flannel-cni-plugin-v1.7.1-flannel1.tar ubuntu@106.75.163.110:/home/ubuntu/chaos
# 其他节点同理
scp flannel-v0.27.2.tar ubuntu@106.75.163.102:/home/ubuntu/chaos
scp flannel-cni-plugin-v1.7.1-flannel1.tar ubuntu@106.75.163.102:/home/ubuntu/chaos
```

## 4. 在每台K8s节点上加载镜像

```bash
# 登录到节点
ssh ubuntu@10.0.0.2
# 加载镜像
sudo docker load -i flannel-v0.27.2.tar
sudo docker load -i flannel-cni-plugin-v1.7.1-flannel1.tar
```

## 5. 配置K8s使用本地镜像（无需修改yaml，kubelet会优先用本地镜像）

- 默认情况下，`imagePullPolicy` 为 `IfNotPresent`，本地有镜像就不会拉取远端。
- 如需强制使用本地镜像，可编辑flannel DaemonSet：

```bash
kubectl -n kube-flannel edit ds kube-flannel-ds
```

- 找到 `imagePullPolicy` 字段，改为：

```yaml
imagePullPolicy: IfNotPresent
```

- 保存退出即可。

## 6. 检查Pod启动

```bash
kubectl get pods -n kube-flannel -w
```

---

**注意：所有K8s节点都要加载镜像，否则对应节点的flannel Pod会拉取失败。**
