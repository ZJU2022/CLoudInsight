# 大型现代应用软件架构的基本层次

![Alt text](img/%E5%A4%A7%E5%9E%8B%E5%BA%94%E7%94%A8%E8%BD%AF%E4%BB%B6%E6%9E%B6%E6%9E%84%E5%9B%BE.png)

![Alt text](img/%E7%8E%B0%E4%BB%A3%E8%BD%AF%E4%BB%B6%E6%9E%B6%E6%9E%84.png)

在讲云计算之前，我们可以先了解一下大型现代应用软件架构的基本层次，通常可以分为四个层次：

## 1. 应用层
这是直接与用户交互的部分，所有的用户体验和业务逻辑都在这一层实现。无论是我们开发的电商平台、客户管理系统，还是移动应用，所有的前端功能和服务都属于应用层。这个层次的目标是为用户提供直接的业务价值和交互体验。

## 2. 中台层
阿里巴巴提出的概念，这一层并非所有系统都必须具备，但它的存在极大地增强了我们的业务灵活性和复用性。中台层通过将常见的业务逻辑和数据服务抽象出来，使我们能够快速响应市场需求，减少重复开发。

- **业务中台**：支撑多个应用层软件的通用模块，这些模块可以被不同的应用复用，减少重复开发。
  - 例如：《咨询预约》、《账单流水》、《在线支付》、《客户评价》、《在线签约》等功能模块。这些模块往往是多个不同应用的共同需求，如《客户评价》可以用于多个服务型平台如餐饮、旅游、医疗等。

- **数据中台**：提供统一的数据管理和分析能力，解决数据孤岛问题、数据统一性和一致性。
  - 例如：美团点评旗下有多种业务，如外卖、酒店、旅游、电影票务等，每个业务都产生大量的数据。
    1. **数据整合和治理**：将不同业务线的数据整合在一个平台上，进行数据清洗和治理，确保数据的一致性和质量。
    2. **用户画像**：通过数据中台，美团能够生成更为精准的用户画像，支持各业务线的个性化服务。
    3. **运营分析**：数据中台为美团的各个业务线提供实时运营数据分析，帮助优化业务运营和决策。
  
  **作用**：数据中台使美团能够更好地了解用户需求，优化各业务线的运营，并通过数据驱动的方式提升整体业务效率。

## 3. 技术底座
这是为之上的业务提供支撑的部分，包含了各种底层技术组件，比如：
- 容器平台
- 监控平台
- 数据库
- 消息中间件
- 日志平台（如 Uxiao）
- 智能告警平台（如 UNOC）
- 统一认证鉴权服务
- Gitlab 等

## 4. 基础设施
IaaS 层，包括主机、存储、网络等基础资源。

---

# 云厂商提供的服务

![Alt text](img/%E4%BA%91%E6%9C%8D%E5%8A%A1.png)

既然大型软件架构是这样划分的，那么云厂商提供的服务主要分为三种：

## 1. 基础设施即服务（IaaS - Infrastructure as a Service）
- **描述**：云计算厂商提供基础的计算、存储和网络资源，用户负责在这些资源上搭建和管理操作系统、数据库、中间件以及其上的应用软件。
- **用户责任**：用户需要管理虚拟机、操作系统、数据库、应用软件的安装、配置、更新和安全。
- **示例**：Amazon EC2、Microsoft Azure VM、Google Compute Engine。

## 2. 平台即服务（PaaS - Platform as a Service）
- **描述**：云计算厂商不仅提供计算、存储和网络资源，还包括管理的数据库、中间件、开发框架等，用户可以在这些基础设施上开发、运行和管理应用程序。
- **用户责任**：用户主要负责应用程序的开发、部署和数据管理，而不需要管理底层的操作系统、数据库或中间件。
- **示例**：Google App Engine、Microsoft Azure App Services、Heroku。

## 3. 软件即服务（SaaS - Software as a Service）
- **描述**：云计算厂商提供完整的软件解决方案，用户直接使用这些软件服务，无需管理底层的硬件、操作系统或应用程序。
- **用户责任**：用户只需要管理和配置应用软件中的用户数据和业务规则，其他一切由云厂商负责。
- **示例**：Google Workspace（Gmail、Docs、Drive）、Salesforce CRM、Microsoft 365。

_注：架构图可参考《悟空聊架构-教你画架构图》中的示例。_

## 4. 多云部署（公有云/私有云/混合云/专有云/专属云）
- **描述**：多云部署是指企业根据自身业务需求，选择不同类型的云服务进行组合使用。这些云类型包括公有云、私有云、混合云、专有云和专属云。尽管它们在商业模式和部署方式上有所区别，但在底层技术和服务类型上本质上是相同的，都是基于云计算的资源和服务。

  - **公有云**：由第三方云服务提供商管理，资源共享给多个租户使用。公有云的优势在于成本效益高、弹性好，并且用户无需管理底层硬件。
  
  - **私有云**：专门为某个组织或企业内部使用，通常由企业自己管理或由第三方提供专属资源。私有云提供更高的安全性、控制力和定制化支持，适合对数据隐私和合规性要求较高的组织。
  
  - **混合云**：结合了公有云和私有云的优点，通过统一的管理平台进行协调，允许企业将敏感数据或关键任务放在私有云中，而将其他工作负载放在公有云中，优化成本和资源利用率。
  
  - **专有云**：通常是为特定客户提供的公有云实例，独享硬件资源，提供与私有云相似的安全性和控制力，但仍由云服务提供商管理。专有云适用于对性能和安全有较高要求的企业，但不愿承担私有云的管理负担。
  
  - **专属云**：专属云是公有云的一种高级形式，提供给特定客户的隔离环境，拥有更高的定制化、性能和安全性。它通常适用于大型企业或政府机构，提供与传统数据中心类似的独立性和控制力，但享受云计算的灵活性和扩展性。

- **用户责任**：
  - 公有云：用户主要负责管理操作系统、应用程序、数据和服务的配置与管理，而不需担心硬件基础设施。
  - 私有云：用户需负责私有云环境的全部管理，包括硬件、软件、网络以及安全措施的实施。
  - 混合云：用户需协调和管理公有云和私有云之间的数据和工作负载的迁移和同步。
  - 专有云：用户负责应用层面的管理和配置，云提供商负责底层基础设施的维护。
  - 专属云：用户享有与专有云相似的责任，但拥有更高的定制化选项和对资源的控制力。

- **示例**：
  - **公有云**：Amazon Web Services (AWS)、Microsoft Azure、Google Cloud Platform (GCP)
  - **私有云**：VMware vSphere、OpenStack、IBM Cloud Private
  - **混合云**：Microsoft Azure Stack、AWS Outposts、Google Anthos
  - **专有云**：AWS Dedicated Hosts、Azure Dedicated Hosts
  - **专属云**：阿里云专有云、腾讯云TCE（腾讯云企业版）、UCloud专属云




---
[参考链接](https://www.51cto.com/article/717315.html)

---

# 知识图谱
我们将对以下知识图谱中涉及到的云计算与云原生相关技术产品和开源项目的基本概念、应用场景、使用示例进行讲解，此知识图谱主要根据Ucloud公有云涉及的产品与生态软件进行归类，目前还是一份初稿，后续将会横向与纵向地延展，欢迎大家对内容提供建议与反馈。后期的使用教程欢迎大家补充

![Alt text](img/%E4%BA%91%E5%8E%9F%E7%94%9F%E7%9F%A5%E8%AF%86%E5%9B%BE%E8%B0%B1.jpg)

[参考链接]([Title](cloud))

# 快速链接
- 计算：[Title](cloud/%E8%AE%A1%E7%AE%97.md)
- 网络：[Title](cloud/%E7%BD%91%E7%BB%9C.md)
- 云盘：[Title](cloud/%E4%BA%91%E7%9B%98.md)
- 存储：[Title](cloud/%E5%AF%B9%E8%B1%A1%E5%AD%98%E5%82%A8.md)
- 容器：[Title](cloud/%E5%AE%B9%E5%99%A8.md)
- 运维：[Title](cloud/%E8%BF%90%E7%BB%B4%E4%B8%8E%E7%9B%91%E6%8E%A7.md)
- 弹性伸缩：[Title](cloud/%E5%BC%B9%E6%80%A7%E4%BC%B8%E7%BC%A9%E4%B8%8E%E8%B4%9F%E8%BD%BD%E5%9D%87%E8%A1%A1.md)

# 代码逻辑
- demo 文件夹包含基于 Go 语言及相关框架的 Web 开发示例，包括接口设计和插件使用等内容。

# 参考文献
- 《深入浅出云计算》
- 《激荡十年：云计算的过去、现在和未来》
- 《图解云计算架构，基础设施和API》
- [GitHub: cloud-computing-101](https://github.com/ZJU2022/cloud-computing-101)
- [GitHub: PassJava-Learning](https://github.com/Jackson0714/PassJava-Learning)

# 贡献者欢迎
-  本项目为我们旅程的上半部分，主要聚焦于云计算相关知识。下半部分将深入探讨UCloud公有云产品线的实际业务经验。
-  请注意，项目中的内容是基于我个人对云计算的理解，因此可能存在错误或遗漏之处。欢迎大家积极反馈，提交 issue 或者 pull request（PR），共同完善和扩展此项目。
-  非常欢迎您的加入！




