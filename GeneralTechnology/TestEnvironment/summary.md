# 测试环境方法论

## 核心理念

**目标**：快速搭建稳定、可维护的测试环境
**原则**：标准化 + 自动化 + 可观测

## 方法论框架

```
┌─────────────────────────────────────────────────────────────┐
│                    测试环境方法论                            │
├─────────────────────────────────────────────────────────────┤
│  规划 → 搭建 → 部署 → 监控 → 维护                           │
└─────────────────────────────────────────────────────────────┘
```

## 第一阶段：环境规划

### 1.1 资源规划表
| 角色 | 配置 | 数量 | 用途 |
|------|------|------|------|
| 应用服务器 | 16核64G+960G SSD | 3台 | 应用服务、负载均衡 |
| 数据库服务器 | 16核128G+4T SSD | 2台 | MySQL、Redis、MongoDB |
| 监控服务器 | 8核32G+1T HDD | 1台 | Prometheus、Grafana、ELK |

### 1.2 网络架构
```
Internet
    │
    ▼
┌─────────────┐
│  负载均衡   │ (app-server-01:80)
└─────────────┘
    │
    ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  应用集群   │───▶│  数据库集群 │───▶│  监控系统   │
│ (8080-8090) │    │ (3306,6379) │    │ (9090,3000) │
└─────────────┘    └─────────────┘    └─────────────┘
```

## 第二阶段：环境搭建

### 2.1 快速搭建脚本
```bash
#!/bin/bash
# 一键搭建脚本

# 1. 安装Docker和K8s
curl -fsSL https://get.docker.com | bash
apt-get update && apt-get install -y kubelet kubeadm kubectl

# 2. 初始化K8s集群
kubeadm init --pod-network-cidr=10.244.0.0/16
mkdir -p $HOME/.kube && cp -i /etc/kubernetes/admin.conf $HOME/.kube/config

# 3. 安装网络插件
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml

# 4. 安装监控
kubectl apply -f https://github.com/prometheus-operator/kube-prometheus/raw/main/manifests/setup/
kubectl apply -f https://github.com/prometheus-operator/kube-prometheus/raw/main/manifests/
```

### 2.2 服务部署模板
```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: app
        image: myapp:latest
        ports:
        - containerPort: 8080
```

## 第三阶段：自动化部署

### 3.1 CI/CD流水线
```
代码提交 → 自动构建 → 测试 → 部署 → 验证
    │         │        │      │      │
    ▼         ▼        ▼      ▼      ▼
   Git    Jenkins   Jest   K8s   Health
```

### 3.2 Jenkins Pipeline
```groovy
pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                sh 'docker build -t myapp:$BUILD_NUMBER .'
                sh 'docker push myapp:$BUILD_NUMBER'
            }
        }
        stage('Deploy') {
            steps {
                sh 'kubectl set image deployment/app-deployment app=myapp:$BUILD_NUMBER'
                sh 'kubectl rollout status deployment/app-deployment'
            }
        }
    }
}
```

## 第四阶段：监控告警

### 4.1 监控架构
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  应用日志   │───▶│  ELK Stack   │───▶│   告警     │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       ▼                   ▼                   ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  系统指标   │───▶│ Prometheus  │───▶│  钉钉/邮件  │
└─────────────┘    └─────────────┘    └─────────────┘
```

### 4.2 关键监控指标
| 指标类型 | 监控项 | 告警阈值 | 检查频率 |
|----------|--------|----------|----------|
| 系统指标 | CPU使用率 | >80% | 1分钟 |
| 系统指标 | 内存使用率 | >85% | 1分钟 |
| 系统指标 | 磁盘使用率 | >90% | 5分钟 |
| 应用指标 | 响应时间 | >500ms | 30秒 |
| 应用指标 | 错误率 | >5% | 1分钟 |

### 4.3 监控脚本
```bash
#!/bin/bash
# 健康检查脚本

check_health() {
    # 检查CPU
    cpu=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)
    if (( $(echo "$cpu > 80" | bc -l) )); then
        send_alert "CPU使用率过高: ${cpu}%"
    fi
    
    # 检查服务状态
    for service in nginx mysql redis; do
        if ! systemctl is-active --quiet $service; then
            send_alert "服务 $service 已停止"
            systemctl restart $service
        fi
    done
}

send_alert() {
    curl -X POST "https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN" \
         -H "Content-Type: application/json" \
         -d "{\"text\":\"$1\"}"
}

check_health
```

## 第五阶段：日常维护

### 5.1 维护日历
```
周一：系统优化 + 安全更新
周二：备份验证 + 数据清理
周三：安全检查 + 权限审计
周四：性能分析 + 瓶颈优化
周五：文档更新 + 经验总结
```

### 5.2 自动化维护脚本
```bash
#!/bin/bash
# 自动维护脚本

# 清理日志
find /var/log -name "*.log" -mtime +7 -delete

# 清理Docker
docker system prune -f

# 数据库备份
mysqldump --all-databases > /backup/mysql_$(date +%Y%m%d).sql

# 清理旧备份
find /backup -name "*.sql" -mtime +7 -delete
```

### 5.3 定时任务配置
```bash
# /etc/crontab
0 2 * * * root /scripts/db_backup.sh      # 每天2点备份
0 * * * * root /scripts/health_check.sh   # 每小时检查
0 3 * * * root /scripts/auto_cleanup.sh   # 每天3点清理
```

## 常见问题速查

### 问题1：服务无法访问
```bash
# 排查步骤
1. ping 目标IP
2. telnet 目标端口
3. systemctl status 服务名
4. kubectl get pods -n 命名空间
```

### 问题2：磁盘空间不足
```bash
# 快速清理
df -h                    # 查看磁盘使用
du -sh /* | sort -hr     # 查看大目录
find /var/log -name "*.log" -mtime +7 -delete  # 清理日志
docker system prune -a   # 清理Docker
```

### 问题3：内存不足
```bash
# 内存优化
free -h                  # 查看内存使用
ps aux --sort=-%mem | head -10  # 查看内存占用进程
sync && echo 3 > /proc/sys/vm/drop_caches  # 清理缓存
```

## 团队分工

| 角色 | 人数 | 主要职责 | 技能要求 |
|------|------|----------|----------|
| 运维工程师 | 1-2 | 系统维护、监控、故障处理 | Linux、Docker、K8s |
| 开发工程师 | 3-5 | 应用部署、配置管理 | CI/CD、脚本编写 |
| 测试工程师 | 2-3 | 环境验证、自动化测试 | 测试框架、数据管理 |

## 评估指标

| 指标 | 目标值 | 检查频率 |
|------|--------|----------|
| 环境可用性 | ≥99.5% | 实时 |
| 部署成功率 | ≥95% | 每次部署 |
| 故障恢复时间 | ≤30分钟 | 故障时 |
| 测试覆盖率 | ≥80% | 每次发布 |

## 工具清单

### 必需工具
- **容器化**：Docker + Kubernetes
- **CI/CD**：Jenkins + Git
- **监控**：Prometheus + Grafana
- **日志**：ELK Stack

### 可选工具
- **配置管理**：Ansible
- **服务网格**：Istio
- **安全扫描**：SonarQube

## 快速开始

### 1. 环境准备（30分钟）
```bash
# 克隆配置模板
git clone https://github.com/your-org/test-env-templates
cd test-env-templates

# 一键部署
./deploy.sh
```

### 2. 验证环境（10分钟）
```bash
# 检查服务状态
kubectl get pods --all-namespaces
kubectl get services

# 访问监控面板
open http://your-ip:3000  # Grafana
open http://your-ip:9090  # Prometheus
```

### 3. 配置告警（5分钟）
```bash
# 配置钉钉告警
cp config/alertmanager.yml.example config/alertmanager.yml
# 编辑配置文件，填入钉钉token
kubectl apply -f config/alertmanager.yml
```

## 总结

**核心要点**：
1. **标准化**：统一配置模板和部署流程
2. **自动化**：脚本化操作，减少人工干预
3. **可观测**：全面监控，及时发现问题
4. **可维护**：定期维护，持续优化

**成功标准**：
- 环境搭建时间 ≤ 1小时
- 故障发现时间 ≤ 5分钟
- 故障恢复时间 ≤ 30分钟
- 团队满意度 ≥ 90%

这套方法论注重实用性和可操作性，通过标准化的流程和自动化的工具，确保测试环境的高效管理和稳定运行。
