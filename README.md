# Geek_homework

Golang class homework in Geek Space.

---

## [Week01：实现切片的删除操作](https://github.com/ycvk/geek_homework/tree/main/week01)

### 作业要求

- **实现删除切片特定下标元素的方法。**

- **性能要求：** 实现相对高性能的删除操作。

- **泛型化：** 改造为支持泛型的方法。

- **缩容机制：** 添加缩容支持，并设计缩容机制。

---

## [Week02：实现用户信息编辑功能](https://github.com/ycvk/geek_homework/tree/main/tinybook)

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

## [Week03：部署方案修改](https://github.com/ycvk/geek_homework/tree/main/tinybook)

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