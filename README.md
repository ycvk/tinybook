# geek_homework
Golang class homework in geek space

## [小微书](https://github.com/ycvk/geek_homework/tree/main/tinybook)
### week02作业：
#### 实现编辑功能

你需要完善 /users/edit 对应的接口。

要求：

允许用户补充基本个人信息，包括：

昵称：字符串，你需要考虑允许的长度。

生日：前端输入为 1992-01-01 这种字符串。

个人简介：一段文本，你需要考虑允许的长度。

尝试校验这些输入，并且返回准确的信息。

修改 /users/profile 接口，确保这些信息也能输出到前端。

不要求你开发前端页面。提交作业的时候，顺便提交 postman 响应截图。

加一个 README 文件，里面贴个图。

就是补充 record 分支上的 Edit 和 Profile 接口。

**post结果**:
![post_01](https://i.mji.rip/2023/10/02/73405f3b359c19579beaaa5fb4fb588e.png
)
![post_02](https://i.mji.rip/2023/10/02/2c01cc2c383c90dfea1d2ff39612d0c0.png
)

**前端请求预览**：
![web_01](https://i.mji.rip/2023/10/02/23b5761e808f0d6b12a3582d8fa39dbf.png
)
![web_02](https://i.mji.rip/2023/10/02/f8b1662852a50f852884534bbb4b1876.png
)


## [week01](https://github.com/ycvk/geek_homework/tree/main/week01)
作业：实现切片的删除操作
实现删除切片特定下标元素的方法。

要求一：能够实现删除操作就可以。

要求二：考虑使用比较高性能的实现。

要求三：改造为泛型方法。

要求四：支持缩容，并设计缩容机制。
