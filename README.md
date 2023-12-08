# Geek_homework

Golang class homework in Geek Space.

## Table of Contents

- [Week01: 实现切片的删除操作](#week01-实现切片的删除操作)
- [Week02: 实现用户信息编辑功能](#week02-实现用户信息编辑功能)
- [Week03: 部署方案修改](#week03-部署方案修改)
- [Week04: 引入本地缓存](#week04-引入本地缓存)
- [Week05: 同步转异步的容错机制](#week05-同步转异步的容错机制)
- [Week06: 优化打印日志的部分](#week06-优化打印日志的部分)
- [Week07: 找出点赞数量前N的数据](#week07-找出点赞数量前N的数据)

---

<h2 id="Week01">Week01: 实现切片的删除操作</h2>

[GitHub Link](https://github.com/ycvk/geek_homework/tree/main/week01)

### 作业要求

- **实现删除切片特定下标元素的方法。**
- **性能要求：** 实现相对高性能的删除操作。
- **泛型化：** 改造为支持泛型的方法。
- **缩容机制：** 添加缩容支持，并设计缩容机制。

---

<h2 id="Week02">Week02: 实现用户信息编辑功能</h2>

[GitHub Link](https://github.com/ycvk/geek_homework/tree/main/tinybook)

### 作业要求

1. **完善 `/users/edit` 接口**
2. **允许的用户基本信息**

    - 昵称：字符串（限制长度）
    - 生日：日期字符串（如 "1992-01-01"）
    - 个人简介：文本（限制长度）
3. **输入校验**：实现输入内容的校验，并返回准确信息。
4. **用户资料接口：** 修改 `/users/profile` 接口以展示新添加的信息。
5. **响应截图：** 提交 Postman 响应截图。
6. **附加文档：** 添加 README 文件，并附带示意图。
7. **代码更新：** 在 `record` 分支上补充 `Edit` 和 `Profile` 接口。

### Postman 响应截图

<details>
  <summary>点击展开图片</summary>

![post_01](https://i.mji.rip/2023/10/02/73405f3b359c19579beaaa5fb4fb588e.png)
![post_02](https://i.mji.rip/2023/10/02/2c01cc2c383c90dfea1d2ff39612d0c0.png)

</details>

### 前端请求预览

<details>
  <summary>点击展开图片</summary>

<img src="https://i.mji.rip/2023/10/02/23b5761e808f0d6b12a3582d8fa39dbf.png" width="50%" height="50%" alt="web_01" />

<img src="https://i.mji.rip/2023/10/02/f8b1662852a50f852884534bbb4b1876.png" width="50%" height="50%" alt="web_02" />

<img src="https://i.mji.rip/2023/10/12/c298dd3635e8b41562a377f98be29cb1.png" width="50%" height="50%" alt="web_03" />

<img src="https://i.mji.rip/2023/10/12/3fb6326899142d4b0903ec785dd646c2.png" width="50%" height="50%" alt="web_04" />

<img src="https://i.mji.rip/2023/10/12/523ea485027beee9f6e381a53f8db630.png" width="50%" height="50%" alt="web_05" />

<img src="https://i.mji.rip/2023/10/12/b8377b4813aee91d9997b6d07291c744.png" width="50%" height="50%" alt="web_06" />

</details>

---

<h2 id="Week03">Week03: 部署方案修改</h2>

[GitHub Link](https://github.com/ycvk/geek_homework)

### 作业要求

1. **项目端口：** 修改项目启动端口为 8081。
2. **Pod 数量：** 将项目部署为 2 个 Pod。
3. **Redis 端口：** 修改 Redis 访问端口为 6380。
4. **MySQL 端口：** 修改 MySQL 访问端口为 3308。

### 需要提交的内容

- `kubectl get services` 的截图。
- `kubectl get pods` 的截图。
- 通过浏览器访问 Tinybook 项目并获得正确响应的截图。

### kubectl 截图

![kubectl_01](https://i.mji.rip/2023/10/15/95d24d6ba5ecba314592afe22bacb45b.png)

### 浏览器访问截图

<details>
  <summary>点击展开图片</summary>

#### Profile 页面

![web_01](https://i.mji.rip/2023/10/15/fbed29fab3a8267054635fbbb893e6e9.png)

#### Edit 页面

![web_02](https://i.mji.rip/2023/10/15/fe1d30f80d88b5f10f284df3e6a2149f.png)

#### Login 页面

![web_03](https://i.mji.rip/2023/10/15/deddbe2bf427c55e2837d45497e6329b.png)

</details>

---

<h2 id="Week04">Week04: 引入本地缓存</h2>

[GitHub Link](https://github.com/ycvk/geek_homework/tree/week04)

### 作业要求

1. **重构现有的CodeCache**：

    - 将当前的 `CodeCache` 改名为 `CodeRedisCache`。
2. **实现本地缓存的CodeCache**：

    - 自由选择本地缓存技术，注意体会技术选型的关键因素。
3. **并发安全**：

    - 确保在单机或开发环境下并发安全。

<details>
  <summary>👉 点击展开结果</summary>

### 技术选型

#### 目前热门本地缓存库

[freecache](https://github.com/coocood/freecache)

[bigcache](https://github.com/allegro/bigcache)

[fastcache](https://github.com/VictoriaMetrics/fastcache)

[ristretto](https://github.com/dgraph-io/ristretto)

[go-cache](https://github.com/patrickmn/go-cache)

[theine-go](https://github.com/Yiling-J/theine-go)

| 缓存库       | 优点                              | 缺点                                                                   | 是否支持TTL | 内存效率 | 适用场景                | 并发安全 | 社区活跃度 |
|-----------|---------------------------------|----------------------------------------------------------------------|---------|------|---------------------|------|-------|
| freecache | 近似LRU淘汰，支持Key设置TTL              | 需要提前知道缓存大小，可能导致内存浪费                                                  | 是       | 中等   | 高并发、内存敏感环境          | 是    | 中等    |
| bigcache  | 不需要提前知道缓存大小，能动态扩展               | 有序列化开销，缓存淘汰效率差，无法为每个key设置TTL，会在内存中分配大数组用以达到 0 GC 的目的，一定程度上会影响到 GC 频率 | 是       | 高    | 动态数据量，需要快速扩展的场景     | 是    | 高     |
| fastcache | 性能高，分片降低锁粒度，索引存储优化              | 不支持TTL                                                               | 否       | 高    | 高性能需求，不需要TTL管理      | 是    | 高     |
| ristretto | 高性能，有准入政策和SampledLFU驱逐政策        | 对GC无优化，内部使用 sync.map                                                 | 是       | 高    | 高性能需求，需要精细控制淘汰策略的场景 | 是    | 高     |
| go-cache  | 易于使用，长时间维护                      | 长久未更新，可能存在潜在的安全和性能问题                                                 | 是       | 低    | 简单缓存需求，不关心长期维护和扩展性  | 是    | 低     |
| theine-go | 支持TTL与持久化，自适应W-TinyLFU淘汰策略，高命中率 | 相对较新，社区支持可能较少                                                        | 是       | 高    | 需要TTL管理和持久化，高命中率要求  | 是    | 不确定   |

综上所述，本次作业可以选用 `ristretto` 或 `theine-go` 作为本地缓存。

##### 参考链接

[性能敏感场景下，Go 三方库的选型思路和案例分析](https://blog.csdn.net/kevin_tech/article/details/125437607)

[golang本地缓存(bigcache/freecache/fastcache等)选型对比及原理总结 - 知乎](https://zhuanlan.zhihu.com/p/487455942)

### 实现与测试

#### 代码实现

- [service 层](https://github.com/ycvk/geek_homework/blob/week04/tinybook/internal/service/code.go)
- [repository 层](https://github.com/ycvk/geek_homework/blob/week04/tinybook/internal/repository/code.go)
- [cache 层](https://github.com/ycvk/geek_homework/blob/80690ff380c90b9bf1b01f7f7e3e39f176561f32/tinybook/internal/repository/cache/code.go#L31-L102) (
  使用 `theine-go` 作为本地缓存, 逻辑详见代码注释)
- [wire DI 层](https://github.com/ycvk/geek_homework/blob/80690ff380c90b9bf1b01f7f7e3e39f176561f32/tinybook/wire.go#L25) (
  依赖注入时, 使用 `LocalCodeCache` 替换 `CodeRedisCache`)

#### 测试结果

##### 1. 发送验证码与登录

![test_01](https://github.com/ycvk/PicDemo/blob/main/8325afc6715b05b8290ef82597ddd98a.png?raw=true)

##### 2. 再次使用此验证码登录

![test_02](https://github.com/ycvk/PicDemo/blob/main/26874c3dafaa801849828b3b057d3391.png?raw=true)

##### 3. 点击登录超过 3 次

![test_03](https://github.com/ycvk/PicDemo/blob/main/d46be533a394741ec42730c58eb4e536.png?raw=true)

##### 4. 短时间内发送验证码超过 3 次

![test_04](https://github.com/ycvk/PicDemo/blob/main/WeChat4dbea418d336ac0b3bb35dc63de2296c.jpg?raw=true)

</details>

---

<h2 id="Week05">Week05: 同步转异步的容错机制</h2>

[GitHub Link](https://github.com/ycvk/geek_homework/tree/main/week05)

### 作业要求

- **容错机制设计：** 当触发限流或服务商崩溃时，将请求转储到数据库，并异步处理。
- **服务商崩溃判定：** 设计新的判定机制，不使用课程给定方法，并阐述决策理由。
- **异步重试控制：** 允许重试 N 次，重试间隔自由决策，参数需可控。
- **编码风格：** 保持面向接口和依赖注入的编程风格。
- **适用场景分析：** 明确容错机制适用场景及其优缺点。
- **改进方案：** 针对缺点，提出改进措施。
- **测试编写：** 提供单元测试以验证功能。

---

### 容错机制设计详解

#### 服务商崩溃判定逻辑

- **判定依据：** 在一定时间窗口内（比如最近5分钟），监控成功与失败的比例，而不是连续失败的次数。如果错误率超过了设定的阈值（例如30%），则触发容错机制。此阈值还可以根据实际情况进行动态调整。
- **决策理由：** 错误率的增加往往是服务崩溃的前兆，此方法不依赖于连续错误，而是整体服务表现，更能反映非连续性的问题。

#### 异步重试机制

- **用户控制：** 用户可以设置基本时间单元、最大间隔时间以及最大重试次数。
- **机制说明：** 采用**指数退避策略**，初始重试间隔设置为一个基本时间单元，比如2秒，之后每次重试间隔时间翻倍，直到达到最大间隔时间，或者重试次数达到上限。

    - **指数退避策略：**
        - 在指数退避策略中，重试时间通常按照指定的乘数增加。
        - 例如，如果基本时间单元为0.5秒，乘数为2，最大重试次数为5，则重试间隔为：0.5秒、1秒、2秒、4秒、8秒。

### 适用场景及优缺点分析

#### 适用场景

- 高可用性系统，消息传递系统，以及需要高度可靠性的服务调用。

#### 优点

- 提高了系统的弹性，即使在服务商出现问题时也能保证消息最终被处理。

#### 缺点

- 可能导致数据库负担加重，如果消息量大而服务商长时间不可用，可能会造成大量的消息堆积。

### 改进方案

- **消息队列：** 使用消息队列服务来代替数据库存储消息，这样可以更有效地处理大量消息，并且可以很容易地进行水平扩展。
- **服务降级：** 当检测到服务商出现问题时，可以暂时降级服务，比如使用备用的短信服务商，或者减少发送频率/停止发送等。

### 重试逻辑

在我的 [重试模块](https://github.com/ycvk/geek_homework/blob/week05/tinybook/internal/service/sms/failover/retry/retry_task.go)
中，`BaseInterval` 是重试的起始间隔，`Multiplier`
是每次重试时增加的倍数。由于还加入了一个随机化因素（`RandomizationFactor`），每次重试的实际间隔会在计算出的指数退避间隔基础上有所波动。

具体来说，每次重试的间隔计算公式为：

$$\text{Interval} = \text{BaseInterval} \times \text{Multiplier}^\text{retry number} \times (1 -
\text{RandomizationFactor} + \text{random factor})$$

其中，`random factor` 是一个 [0, 2 * `RandomizationFactor`] 范围内的随机数。

假设没有随机化因素，即 `RandomizationFactor` 为 0，那么在 `BaseInterval` 为 1 秒，`Multiplier` 为 2 的情况下，每次重试的理想间隔时间将是：

- 第1次重试: \( $1 \times 2^0 = 1 $\) 秒
- 第2次重试: \( $1 \times 2^1 = 2 $\) 秒
- 第3次重试: \( $1 \times 2^2 = 4 $\) 秒
- 第4次重试: \( $1 \times 2^3 = 8 $\) 秒
- 第5次重试: \( $1 \times 2^4 = 16 $\) 秒

但实际上，每次重试的间隔还会包含随机化因素的影响。这意味着每次的实际间隔会在上述理论值的基础上上下浮动。例如，如果 `RandomizationFactor`
是 0.5，那么实际间隔将随机地增加或减少最多50%。

以下是基于模拟并考虑随机化因素的每次重试的大约时间（单位为秒）：

<details>
  <summary>👉 点击展开模拟代码</summary>

~~~python
import random

# 定义重试配置参数
base_interval = 1  # 基础间隔时间，单位秒
multiplier = 2     # 间隔增加的倍数
max_retries = 5    # 最大重试次数
randomization_factor = 0.5  # 随机化因素

# 模拟重试间隔时间计算
def simulate_retries(base_interval, multiplier, max_retries, randomization_factor):
    intervals = []
    for retry in range(max_retries):
        # 计算理论上的间隔时间
        interval = base_interval * (multiplier ** retry)
        # 加入随机化因素
        random_factor = random.uniform(-randomization_factor, randomization_factor)
        interval_with_randomization = interval * (1 + random_factor)
        intervals.append(interval_with_randomization)
    return intervals

# 获取每次重试的大约时间
retry_intervals = simulate_retries(base_interval, multiplier, max_retries, randomization_factor)
retry_intervals

~~~

</details>

- 第1次重试: 约 \(0.63\) 秒
- 第2次重试: 约 \(2.14\) 秒
- 第3次重试: 约 \(3.06\) 秒
- 第4次重试: 约 \(7.12\) 秒
- 第5次重试: 约 \(8.90\) 秒

这些值包含了随机化因素，实际的间隔会在每次运行时有所不同。这样设计是为了避免在服务出现问题时多个客户端同时进行重试，从而可能导致的“群集效应”。

### 代码实现

以下是实现的异步重试容错机制的组件及其逻辑：

- **错误率监控器**:
  [error_rate_monitor.go](https://github.com/ycvk/geek_homework/blob/week05/tinybook/internal/service/sms/failover/retry/error_rate_monitor.go)
- **重试模块**:
  [retry_task.go](https://github.com/ycvk/geek_homework/blob/week05/tinybook/internal/service/sms/failover/retry/retry_task.go)
- **异步重试逻辑层**:
  [failover_async.go](https://github.com/ycvk/geek_homework/blob/week05/tinybook/internal/service/sms/failover/retry/failover_async.go)
- **单元测试**:
  [failover_async_test.go](https://github.com/ycvk/geek_homework/blob/week05/tinybook/internal/service/sms/failover/retry/failover_async_test.go)

这些组件协同工作，实现了以下特点：

- 实现了`sms.Service`的接口，类似于装饰器模式包含了实际调用的`send`接口。
- 包含了一个`sms.Service`实例，用于实际发送短信。
- 包含了一个`AsyncRetry`实例，用于异步重试。
- 包含了一个`ErrorRateMonitor`实例，用于监控错误率。
- 包含了一个`limiter.Limiter`实例，用于限流。

具体流程如下：

#### 初始化

1. **错误率监控器** (`NewErrorRateMonitor`)

    - 设置初始阈值和窗口开始时间。
    - 启动定时任务调整错误率和阈值 (`adjustErrorRateAndThreshold`).
2. **异步重试机制** (`NewAsyncRetry`)

    - 设置重试间隔和重试次数。
3. **限流器** (`NewLimiter`)

    - 设置限流阈值和窗口大小。

#### 发送短信

4. **执行`Send`方法发送短信**
    - 使用`Limiter`方法检查是否超过限流阈值。
        - 如果**超过限流阈值**:
            - 触发容错机制。
            - 存储消息到数据库。
            - 启动`AsyncRetry`进行异步重试。
            - 记录错误信息。
        - 如果**未超过限流阈值**:
            - 调用`Send`方法，发送短信。
    - 记录发送结果。
        - 使用`RecordResult`方法传入发送结果。
        - 如果错误率超过阈值，触发容错机制。
            - 将消息存储到数据库。
            - 启动`AsyncRetry`进行异步重试。
            - 记录错误信息。
    - 如果发送成功，返回结果。

#### 容错机制

5. **定时调整错误率和阈值** (`adjustErrorRateAndThreshold`)

    - 每分钟执行。
    - 清理过时结果 (`cleanUpOldResults`)。
    - 计算当前错误率 (`calculateErrorRate`)。
    - 自适应调整阈值 (`adjustThreshold`)。
6. **定期检查错误率** (`CheckErrorRate`)

    - 比较当前错误率与阈值。
    - 如果错误率超过阈值，触发容错机制。
        - 存储消息到数据库。
        - 启动重试机制（使用goroutine或服务）。

### k8s 项目实测结果

<details>

  <summary>👉 点击展开测试结果</summary>

#### 1. 测试限流下的异步重试机制

k8s 项目启动初始化时，在 [InitSMSService](https://github.com/ycvk/geek_homework/blob/week05/tinybook/ioc/sms.go)
中将限流器设置为 1 分钟只能发送 1 条短信，并选用本地测试的短信服务`localsms.NewService()`，来测试限流后，异步重试机制的效果。

结果如下：

![test_01](https://github.com/ycvk/PicDemo/blob/main/WeChat55ce105ae905f7098d035f861459c98b.jpg?raw=true)

可以看到确实是触发了限流后，操作数据库insert了一条记录，并且启动了异步重试机制。

异步重试成功后，日志打印了验证码信息和重试成功的信息，随后删除了数据库中的记录。

#### 2. 测试重复手机号的异步重试机制

我在发送验证码前，往数据库中insert一条跟要发送的手机号一样的测试记录，来测试重复手机号的情况。

结果如下：

![test_02](https://github.com/ycvk/PicDemo/blob/main/WeChatb9cf24f1d446d2ff85703f4dc7e580d7.jpg?raw=true)

可以在日志看到，在insert插入重复数据时，匹配到mysql的1062唯一索引错误，返回重复手机号的错误信息，并且没有触发异步重试机制。

#### 3. 测试错误率超过阈值下的连续失败异步重试机制

首先在 [InitSMSService](https://github.com/ycvk/geek_homework/blob/week05/tinybook/ioc/sms.go)
中将错误率监控器改为 `monitor := retry.NewErrorRateMonitor(0.01, 0.5, 10*time.Second)`
表示错误率阈值为 1%，窗口大小为 10 秒，来测试错误率超过阈值后，异步重试机制的效果。

随后在localSMSservice的本地测试的 [send](https://github.com/ycvk/geek_homework/blob/week05/tinybook/internal/service/sms/localsms/service.go)
方法中，让其返回错误，来模拟错误率超过阈值的情况。

结果如下：

![test03](https://github.com/ycvk/PicDemo/blob/main/WeChat018eabf5090e7c9943e0a68d4ac4d55f.jpg?raw=true)

可以看到再连续发送了几条短信后，因为已经将`send()`方法改为永远返回错误，所以在错误率监控器中，错误率为100%，超过了阈值1%，触发了异步重试机制。

insert进数据库后，开始重试，重试超过了最大次数，重试彻底失败，日志打印了重试失败的堆栈信息。

</details>

---

<h2 id="Week06">Week06: 优化打印日志的部分</h2>

[GitHub Link](https://github.com/ycvk/geek_homework/blob/week06/tinybook/internal/web/middleware/error.go)

### 作业要求

- 在 Web 的 Handler 部分，有很多 if-else 分支，基本上都是在判定 err !=nil。每一个 if 里面都要打印日志。
- 现在要求优化这些打印日志的逻辑，避免每一处 err !=nil 的时候，都得手动打一个日志。

### 优化思路

- 使用中间件进行错误捕获和日志记录。
- 在 Gin 中，中间件可以用来拦截请求并进行预处理，也可以在处理流程的最后统一处理错误。可以定义一个中间件来捕获处理过程中的任何错误，并在这里统一进行日志记录。

### 代码实现

1. 定义错误处理中间件

    - 首先，定义一个中间件函数。这个函数将在请求的处理链中被调用。使用 ctx.Next() 来调用链中的其他处理函数。之后，中间件会检查是否有任何错误被加入到
      Gin 的 Context 中。
        - ```go
             func ErrorHandler() gin.HandlerFunc {
                return func(ctx *gin.Context) {
                    ctx.Next()
                    // 检查是否有错误被加入到 Gin 的 Context 中
                    err := ctx.Errors.Last()
                    if err != nil {
                        // 如果有错误，记录日志
                        log.Printf("Error: %s", err.Error())
                    }
                }
             }
            ```

2. 在路由中应用中间件

- `router.Use(ErrorLogger()) // 应用中间件`
- 这样，每个经过这个路由的请求都会经过错误处理中间件。

3. 在处理函数中添加错误到上下文

    - 在处理函数中，如果有错误，可以使用 `ctx.Error(err)` 将错误添加到 Gin 的 Context 中。
        - ```go
           func MyHandler(c *gin.Context) {
                err := someFunction() // 假设这是一个可能产生错误的函数
            if err != nil {
                c.Error(err) // 将错误添加到上下文
                return
            }
           // 正常的处理逻辑...
           }
            ```
    - 这样，错误就会被传递到错误处理中间件中。

### 优点

使用这种方法的主要优点是：

- 集中管理：你可以在一个地方集中处理所有的错误日志记录，而不需要在每个处理函数中重复相同的逻辑。
- 灵活性：你可以轻松地调整错误处理和日志记录的策略，而无需修改大量的处理函数。
- 代码清洁：这使得处理函数更加专注于它们的主要职责，从而使代码更加清晰和易于理解。

---

<h2 id="Week07">Week07: 找出点赞数量前N的数据</h2>

[GitHub Link](https://github.com/ycvk/geek_homework/tree/week07)

### 作业要求

本次的核心是设计一个能够找出按照点赞数量前 N 个数据的高性能方案。方案应考虑以下要点：

1. 综合缓存利用：整合 Redis 和本地缓存的使用，确保方案在面试中具有吸引力和竞争力。
2. 业务折中：清晰阐述任何业务上的折中和权衡。
3. 自主设计：独立设计解决方案，遵循自主搜索和讨论的原则。

### 设计思路

- Redis ZSet：使用 Redis 的 ZSet 数据结构来存储和更新文章点赞数。例如，通过命令 `zincrby article:like_count 1 id`
  更新点赞数，并使用 `ZRevRange` 获取指定排名数据。


- 本地缓存与 Kafka 结合：采用 **theine-go** 高性能本地缓存，结合 Kafka 消息队列。
  在点赞或取消点赞时，发送消息到 Kafka。消费者接收消息后，每固定时间，向redis缓存中设置或更新一个布尔型键值，表示排行榜数据的变化。


- 定时器同步机制：设置可调节的定时器，周期性检查redis缓存中的布尔键。如键存在，则表明排行榜数据有变化，触发本地缓存从 Redis
  拉取最新数据，拉取后重置该键。


- 控制获取排行榜的TopN数量：控制只获取最多排行榜前N个数据，避免数据量过大。也避免了非法请求 N 过大，造成获取数据量过大的情况。


- 关于更新本地缓存的key：这里使用redis来存储并更新此key，而不是使用本地缓存。
  因为我k8s项目中使用的是多pod实例部署，如果使用本地缓存来存储此key，那么每个pod都会有一份此key，无法保证key唯一，会导致数据不一致。


- 关于使用kafka：我使用的是`kafka-go`库，可以通过设置`CommitInterval`来控制消费者的消费频率。
  这样配合定时器，可以控制消费者每隔一段时间才消费这时间段的所有消息，配合锁可以保证此时间段内只向redis请求一次，避免消费者消费过快，造成redis压力过大。

### 业务折中

- **Redis ZSet 与本地缓存的折中：** 由于 Redis ZSet 本身就是一个有序集合，可以直接使用 ZSet 来存储排行榜数据。但是，如果每次请求都直接从
  Redis
  拉取数据，会导致 Redis 压力过大。因此，可以使用本地缓存来缓存排行榜数据，减轻 Redis 压力。但是，本地缓存的数据可能会与 Redis
  中的数据不一致，因此需要定时器同步机制来保证数据的一致性。


- **实时缓存 与 定时缓存 的折中：**
    - 由于本地缓存是一个内存缓存，如果每次点赞或取消点赞都直接更新本地缓存，很可能会导致缓存数据与 Redis
      中的数据不一致。
    - 但如果每次点赞或取消点赞都直接获取 Redis 中的数据，再更新本地缓存，那跟直接从 Redis 拉取数据有什么区别呢？
    - 因此，可以使用 Kafka 消息队列来解耦点赞和取消点赞的操作。当点赞或取消点赞时，发送消息到
      Kafka。消费者接收消息后，在一个可设置的时间段内，只会向redis缓存中设置或更新一次布尔型键值，表示排行榜数据的变化。
    - 如果每次点赞或取消点赞都直接更新本地缓存，那么 Kafka 消息队列的作用就只是解耦了点赞和取消点赞的操作，而没有减轻本地缓存的压力。
    - 因此，可以设置一个定时器，周期性检查redis缓存中的布尔键。如键存在，则表明排行榜数据有变化，触发本地缓存从 Redis
      拉取最新数据，拉取后重置该键。
    - 以上这样，就可以减轻缓存的压力，又能保证数据的最终一致性。
    - 但是，这样做也有一个缺点，就是可能会导致数据的实时性降低。因为如果定时器的周期设置得太长，那么可能会有一段时间内，本地缓存中的数据是旧的。
      因此，需要在实际使用中，根据实际情况来调节定时器的周期。

### 接口调用流程

1. 用户点赞或取消点赞时，发送消息到 Kafka。
2. 消费者接收消息后，在设置的时间段内，只会向redis缓存中设置或更新一次布尔型键值，表示排行榜数据的变化。
3. 定时器周期性检查redis缓存中的布尔键。如键存在，则表明排行榜数据有变化，触发本地缓存从 Redis 拉取最新数据，拉取后重置该键。
4. 用户请求排行榜数据时，从本地缓存中获取数据。
5. 如果本地缓存中没有数据，则从 Redis 拉取数据，并发送消息到 Kafka。
6. 消费者接收消息后，在设置的时间段内，只会向redis缓存中设置或更新一次布尔型键值，表示排行榜数据的变化。
7. 如果redis有数据，返回redis排行榜数据。
8. 如果redis没有数据，去数据库拉取数据，返回数据库排行榜数据。
9. 在以上返回数据库排行榜数据的同时，异步更新redis排行榜数据。
10. redis更新完成排行榜数据后，发送消息到 Kafka。

### 代码实现

- handler层: 接口调用流程的实现
    - [article_handler](https://github.com/ycvk/geek_homework/blob/week07/tinybook/internal/web/article_handler.go#L326-L350)
- service层: 业务逻辑的实现
    - [interactive_service](https://github.com/ycvk/geek_homework/blob/week07/tinybook/internal/service/interactive.go#L30-L65)
- repository层: 数据库与缓存操作的实现
    - [interactive_repository](https://github.com/ycvk/geek_homework/blob/b89b00f471642aac670c2f8d2082955fead93e4b/tinybook/internal/repository/interactive.go#L35-L60)
- dao层: 数据库操作的实现
    - [interactive_dao](https://github.com/ycvk/geek_homework/blob/b89b00f471642aac670c2f8d2082955fead93e4b/tinybook/internal/repository/dao/interactive.go#L63-L71)
- cache层: 缓存操作的实现
    - [interactive_cache](https://github.com/ycvk/geek_homework/blob/week07/tinybook/internal/repository/cache/interactive.go#L42-L108)
- kafka层: kafka消息队列的实现
    - [producer](https://github.com/ycvk/geek_homework/blob/week07/tinybook/internal/events/interactive/producer.go)
    - [consumer](https://github.com/ycvk/geek_homework/blob/week07/tinybook/internal/events/interactive/consumer.go)
    - 定时器 ticker
      的实现也在其中 [ticker](https://github.com/ycvk/geek_homework/blob/week07/tinybook/internal/events/interactive/consumer.go#L118-L159)
    - 配合定时器的固定时间尺度，比如 **(1分钟/8小时/1天)** ，可以做到 **每固定时间** 去检查一次redis缓存中的布尔键。

    - 如果键存在，则表明排行榜数据有变化，触发本地缓存从 Redis 拉取最新数据，拉取后重置该键。如果键不存在，则表明排行榜数据没有变化，不更新本地缓存。
    - 而`consumer`中的`Call()`方法，则保证了在给定时间内，无论消费者多少次消费，也只会最多执行一次去更新redis中的`bool key`
      的操作。

### UML时序图

![UML时序图](https://github.com/ycvk/PicDemo/blob/main/1750319627.png?raw=true)

### 测试结果

#### 测试环境

- MacBook Pro 2021 M1Max 32G
- macOS Sonoma 14.1.2
- Go 1.21.4
- Redis 7
- Kafka 3.6.0
- MySQL 8.0.27

#### 1. 没有缓存时，从数据库拉取数据

![test_01](https://github.com/ycvk/PicDemo/blob/main/2108972058.png?raw=true)

#### 2. 有redis缓存时，从redis拉取数据

![test_02](https://github.com/ycvk/PicDemo/blob/main/764976255.png?raw=true)

#### 3. 有本地缓存时，从本地缓存拉取数据

![test_03](https://github.com/ycvk/PicDemo/blob/main/44177685.png?raw=true)

#### 结论

可看到，当有本地缓存时，平均响应时间直接降低了50%以上。
QPS为8151.66，比没有缓存时提高了5倍多。且没有任何错误。

直接请求数据库时，可以看到后面的请求响应时间都是timeout，错误率超过50%，比走redis慢了4倍，比走本地缓存慢了5倍多。

---