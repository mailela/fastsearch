# fastsearch

fastsearch 一个golang实现的全文检索引擎，支持持久化和单机亿级数据毫秒级查找。
- 支持持久化                 
- 基于磁盘+内存缓存             
- 支持表达式            
- 原生二进制，无外部依赖           
- 自带中文分词和词库             
- 自带可视化管理界面             
- 基于Golang原生可执行文件，内存非常小 
- 默认可以不加任何参数启动，并且提供少量配置 
- 快速检索

- 主动防御监测非法关键词
- 禁用搜索非法关键词
- 负面词管理
- 负面消息推送
- 接口可以通过http调用。

详见 [API文档](./docs/api.md)





## 文档

+ [示例](./docs/example.md)
+ [API文档](./docs/api.md)
+ [索引原理](./docs/index.md)
+ [配置文档](./docs/config.md)
+ [持久化](./docs/storage.md)
+ [编译部署](./docs/compile.md)


## fastsearch在线管理后台Demo
[http://127.0.0.1:5679/admin](http://127.0.0.1:5679/admin)



## 二进制文件下载

> 支持Windows、Linux、macOS、（amd64和arm64）和苹果M1 处理器

## 技术栈

+ 二分法查找
+ 快速排序法
+ 倒排索引
+ 正排索引
+ 文件分片
+ golang-jieba分词
+ leveldb
+ 支持表达式的条件过滤
+ 支持表达式权限排序

### 为何要用golang实现一个全文检索引擎？

+ 正如其名，`fastsearch`去探索全文检索的世界，一个小巧精悍的全文检索引擎，支持持久化和单机亿级数据毫秒级查找。

+ 传统的项目大多数会采用`ElasticSearch`来做全文检索，因为`ElasticSearch`够成熟，社区活跃、资料完善。缺点就是配置繁琐、基于JVM对内存消耗比较大。

+ 所以我们需要一个更高效的搜索引擎，而又不会消耗太多的内存。 以最低的内存达到全文检索的目的，相比`ElasticSearch`，`fastsearch`是原生编译，会减少系统资源的消耗。而且对外无任何依赖。

## 安装和启动

> 下载好源码之后，进入到源码目录，执行下列两个命令
>

+ 编译

> 直接下载 [可执行文件](https://gitee.com/rachel_os/fastsearch/fastsearch/releases) 可以不用编译，省去这一步。

```shell
go get && go build
```

+ 启动

```shell
./fastsearch --addr=:8080 --data=./data
```

+ docker部署

```shell
docker build -t fastsearch .
docker run -d --name fastsearch -p 5678:5679 -v /mnt/data/fastsearch:/usr/local/go_found/data fastsearch:latest
```

+ 其他命令
  参考 [配置文档](./docs/config.md)

## 多语言SDK

> 使用fastsearch的多语言SDK，可以在不同语言中使用fastsearch。但是请注意，版本号与fastsearch需要一致。主版本和子版本号，修订版不一致不影响。

[API文档](./docs/api.md)用HTTP请求实现。

## 和ES比较

| ES          | fastsearch               |
|-------------|-----------------------|
| 支持持久化       | 支持持久化                 |
| 基于内存索引      | 基于磁盘+内存缓存             |
| 支持表达式      | 支持表达式            |
| 需要安装JDK     | 原生二进制，无外部依赖           |
| 需要安装第三方分词插件 | 自带中文分词和词库             |
| 默认没有可视化管理界面 | 自带可视化管理界面             |
| 内存占用大       | 基于Golang原生可执行文件，内存非常小 |
| 配置复杂        | 默认可以不加任何参数启动，并且提供少量配置 |

## 功能特性

- 支持持久化                 
- 基于磁盘+内存缓存             
- 支持表达式            
- 原生二进制，无外部依赖           
- 自带中文分词和词库             
- 自带可视化管理界面             
- 基于Golang原生可执行文件，内存非常小 
- 默认可以不加任何参数启动，并且提供少量配置 
- 快速检索

## 新增功能
- 主动防御监测非法关键词
- 禁用搜索非法关键词
- 负面词管理
- 负面消息推送

## 待办

[TODO](docs/TODO.md)

