## 背景

这几年，一直尝试使用 golang ，写出一个易用可靠的服务器框架

在 github 上开了以下几个坑：

深坑                                                   |  制作思路                                                                            | 主要问题
:------------------------------------------------------|:-------------------------------------------------------------------------------------|:-----------------------------------
[go-x](https://github.com/fananchong/go-x)             | 使用 k8s 做服务发现</br>服务节点能知道其他服务节点加入或离开                         | 需要部署 k8s 才能用，限制太大</br>需要自己处理服务节点间的互连、通信</br>没有区分框架层与应用层
[go-xserver](https://github.com/fananchong/go-xserver) | 内置连接管理器，做服务发现</br>应用层使用插件加载</br>内置处理服务节点间的互连、通信 | 框架层代码结构混乱，维护困难</br>无法轻易的替换功能实现。比如服务发现功能</br>框架的使用面板不简洁、不够优雅

虽然在挖坑，但也在不断的总结经验教训，看到 go-xserver 先天不足，很难深入，于是中断了 3 个月左右

这期间，一直在 github 上寻找阅读可能的更优解决方案

最终遇到了 [micro/go-micro](https://github.com/micro/go-micro) ，该框架突出的特点在于：
- 框架层代码结构清晰，可插拔
- 框架的使用面板简洁、优美

于是用 [micro/go-micro](https://github.com/micro/go-micro) 尝试制作了聊天服，可以快速开发

但是性能测试下来，也有致命缺陷：
- 同步调用效率太低（必然结果，因为一定会阻塞网络层）
  - 性能测试请参考： [test_go-micro_qps](https://github.com/fananchong/test_go-micro_qps)
  - 莫名的 `call timeout: context deadline exceeded` 与 `error selecting xxxx node: not found`
- 异步调用基于消息队列 pub/sub 模式
  - 多机房部署时，会遇到困惑。消息队列需暴露到外网
  - 服务节点间所有消息处理都经消息队列的话，没有过这种经验（是否可行，线上没听过有这种做法）
- 无法简单的实现自己的异步 call 插件
  - [micro/go-micro](https://github.com/micro/go-micro) 已经有众多插件，接口层面已经无法修改
  - [micro/go-micro](https://github.com/micro/go-micro) 新增异步接口，同样会动摇众多插件的底层实现，基本上不可能

因此有了本深坑3


## 计划

- 计划1 ： 通过复刻其主要接口，并有自己的实现
- 计划2 ： 摈弃同步调用接口，改为异步调用接口


## 目标

- 目标1 ： 本坑将持续深挖，不轻易弃坑
- 目标2 ： [micro/go-micro](https://github.com/micro/go-micro) 主要针对微服务，对性能要求不是很高。因此重点改进的就是这里
- 目标3 ： 趟一遍 [micro/go-micro](https://github.com/micro/go-micro) 主要接口，力求完全体会掌握 [micro/go-micro](https://github.com/micro/go-micro) 精髓
- 目标4 ： 真正实现一个易用可靠的服务器框架

## 参考版本

- 分支：   master
- 版本号： 5b327ce72374615404e2fac8c73aaaa7066ca22c
