package webook

/*

1. Web 层
职责: Web 层主要负责处理来自客户端的 HTTP 请求。它通常包含控制器（Controllers），负责解析请求、调用相应的业务逻辑、处理用户输入和返回响应结果。
角色:
输入解析: 从 HTTP 请求中解析出需要的数据（如表单数据、URL 参数等）。
响应处理: 将业务逻辑处理后的结果封装成 HTTP 响应（如 JSON、HTML 等）返回给客户端。
依赖注入: 通过依赖注入调用 Service 层的接口，以实现业务逻辑的解耦。

2. Service 层
职责: Service 层主要负责处理应用程序的业务逻辑。它封装了业务规则，可能会与数据访问层（例如存储库模式，Repository）进行交互。
角色:
业务逻辑: 实现核心业务逻辑，比如用户注册、登录、数据处理等。
跨服务交互: 可能调用其他微服务、外部 API 或者数据访问层以完成特定业务功能。
事务管理: 如果业务操作涉及多个数据库操作，Service 层通常负责事务的开始、提交或回滚。

3. Web 层与 Service 层的关系
调用关系: Web 层通过调用 Service 层的接口来完成具体的业务操作。Service 层专注于业务逻辑的实现，而 Web 层则负责处理请求和响应，并与 Service 层解耦。
松耦合: 由于 Web 层依赖于接口而不是具体的实现，这使得两者之间的耦合度很低，便于维护和扩展。

4. 设计思路
Web 层 (web):
UserHandler 是 Web 层的一部分，主要职责是处理用户相关的请求（如注册、登录、编辑用户信息等）。
UserHandler 中定义了两个服务接口的依赖：UserService 和 CodeService。这些接口的实现将被注入到 UserHandler 中，以便在处理请求时调用相应的业务逻辑。

Service 层 (service):
service/code.go 定义了 CodeService 接口，抽象了与验证码相关的业务逻辑。
Service 层专注于业务逻辑的实现，而不是具体的输入或输出处理。CodeService 只定义了接口，而其具体实现则可以在不同的模块中进行，这种设计提供了灵活性。
依赖注入与解耦:

UserHandler 并不直接依赖于 CodeService 的具体实现，而是依赖于 CodeService 接口。这种设计方式使得 UserHandler 与业务逻辑实现解耦，便于测试和维护。




 +-------------------------+
 |      Web Layer           |
 |--------------------------|       +--------------------------+
 | +-----------------------+|       |  Service Layer            |
 | |  UserHandler          |<------->|  UserService Interface   |
 | |                       ||       |  +-----------------------+|
 | |  - Handle Requests    ||       |  |  UserServiceImpl       ||
 | |  - Validate Input     ||       |  |                       ||
 | |  - Send Response      ||       |  |  - Core Business Logic ||
 | +-----------------------+|       |  |  - Call Repositories   ||
 +--------------------------+       |  |                       ||
                                     |  +-----------------------+|
                                     +--------------------------+

*/

/*
1. 前后端联调请求流程

第一步：前端输入：img/前端输入.png
第二步：控制台HTTP请求：img/HTTP请求.png
第三步：后端Web层，接受请求并校验 webook/web/user.go/SignUp
第四步：调用业务逻辑处理请求 webook/service/user.go
第五步：根据业务逻辑处理结果返回响应


*/
