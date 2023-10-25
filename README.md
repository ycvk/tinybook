# Geek_homework

Golang class homework in Geek Space.

## Table of Contents

- [Week01: 实现切片的删除操作](#Week01)
- [Week02: 实现用户信息编辑功能](#Week02)
- [Week03: 部署方案修改](#Week03)
- [Week04：引入本地缓存](#Week04)

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

<h2 id="Week04">Week04：引入本地缓存</h2>

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

#### 热门本地缓存库一览

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
- [cache 层](https://github.com/ycvk/geek_homework/blob/d6f037683884b41ccb56c20bc8f6425fe8fe8ad6/tinybook/internal/repository/cache/code.go#L29-L94) (
  使用 `theine-go` 作为本地缓存, 逻辑详见代码注释)
- [wire DI 层](https://github.com/ycvk/geek_homework/blob/d6f037683884b41ccb56c20bc8f6425fe8fe8ad6/tinybook/wire.go#L25) (
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