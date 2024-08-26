# FastSerach

`FastSearch`启动之后，会监听一个TCP端口，接收来自客户端的搜索请求。处理http请求部分使用`gin`框架。


## 多数据库支持

从1.1版本开始，我们支持了多数据库，API接口中通过get参数来指定数据库。

如果不指定，默认数据库为`default`。

如：`api/index?database=db1` 其他post参数不变

如果指定的数据库名没有存在，将会自动创建一个新的数据库。如果需要删除，直接删除改数据库目录，然后重启FastSearch即可。
## 表达式支持

| 操作符 | 含义       |示例|
|------|-------------|----| 
| Trim |  清除空格   
| IN |  包含 |
| NotIN | 不包含   |
| Like | 模模糊符匹配 |
| ==  |  相等 |   
| < | 小于|
| <= | 小且等于|
| >= | 大且等于
| >|  大于 |
| != | 不等于 
| + | 加|
| -   | 减|
| * | 乘|
| / | 除|
| % | 模|
| &&  | 逻辑与|
| \|\|  | 逻辑或|
| ! | 逻辑非|
| ( |  左括号|
| ) |  右括号|
| [ |  左中括号|
| ] |  右中括号|
| , |  逗号|  


```
示例:
说明：索引添加的document对象，均可以作为变量使用
{
  "id": "id",
  "flag": "标识",
  "tags": "标签",
  "document": {
    title:"标题",
    hot:10,
    url: "http://www.baidu.com"
  }
}
过滤条件示例1： 标题为"中国" 且 热度大于10
(document.title Like  '中国') && (document.hot > 10) 
过滤条件示例1： 标题为"中国" 且 热度大于10

```


## 增加/修改索引

需要在query参数中指定数据库名`database=default`

| 接口地址 | /api/index       |
|------|------------------|
| 请求方式 | POST             |
| 请求类型 | application/json |

### 请求

| 字段       | 类型     | 必选  | 描述                                |
|----------|--------|-----|-----------------------------------|
| id       | string | 是   | 文档的主键id，需要保持唯一性，如果id重复，将会覆盖直接的文档。 |
| title    | string | 是   | 标题|
| flag     | string | 否   | 标识
| tags     | string | 否   | 标签                   |
| text     | string | 是   | 需要索引的文本块                          |
| document | object | 是   | 附带的文档数据，json格式，搜索的时候原样返回          |

query参数(params-data)

| 字段       | 类型     | 必选  | 描述     |
|----------|--------|-----|--------|
| database | string | 是   | 指定数据库名 |

+ POST /api/index

```json
{
  "id": "no-1",
  "text": "我爱中国",
  "title":"我爱中国",
  "flag":"标识",
  "tags":" 标签,
  "document": {
    "title": "我爱中国",
    "number": 223
  }
}
```

+ 命令行

```bash
curl -H "Content-Type:application/json" -X POST --data '{"id":88888,"text":"我爱中国","document":{"title":"我爱中国","number":223}}' http://127.0.0.1:5678/api/index?database=default
```

### 响应

```json
{
  "state": true,
  "message": "success"
}
```

## 批量增加/修改索引
与添加单个索引一样，也需要在query参数中指定数据库名`database=default`

| 接口地址 | /api/index/batch |
|------|------------------|
| 请求方式 | POST             |
| 请求类型 | application/json |

参数与单个一致，只是需要用数组包裹多个json对象，例如：

```json
[
  {
    "id": "s-1",
    "text": "我爱中国",
    "title":"我爱中国",
    "flag":"标识",
    "tags":" 标签,
    "document": {
      "title": "我爱中国",
      "number": 223
    }
  },
  {
    "id": "s-2",
    "text": "我爱长沙",
    "title":"我爱中国",
    "flag":"标识",
    "tags":" 标签,
    "document": {
      "title": "我爱长沙",
      "number": 123
    }
  }
]
```

## 删除索引
与以上接口一样，也需要在query参数中指定数据库名`database=default`

| 接口地址 | /api/index/remove |
|------|-------------------|
| 请求方式 | POST              |
| 请求类型 | application/json  |

### 请求

| 字段  | 类型     | 必选  | 描述      |
|-----|--------|-----|---------|
| id  | string | 是   | 文档的主键id |

+ POST /api/remove

```json
{
  "id": 88888
}
```

+ 命令行

```bash
curl -H "Content-Type:application/json" -X POST --data '{"id":88888}' http://127.0.0.1:5678/api/remove?database=default
```

### 响应

```json
{
  "state": true,
  "message": "success"
}
```

## 查询索引

`FastSearch`提供了一种查询方式，按照文本查询。与其他Nosql数据库不同，`FastSearch`不支持按照文档的其他查询。

| 接口地址 | /api/query       |
|------|------------------|
| 请求方式 | POST             |
| 请求类型 | application/json |

### 请求

| 字段        | 类型     | 必选  | 描述                                                                                           |
|-----------|--------|-----|----------------------------------------------------------------------------------------------|
| query     | string | 是   | 查询的关键词，都是or匹配                                                                                |
| page      | int    | 否   | 页码，默认为1                                                                                      |
| limit     | int    | 否   | 返回的文档数量，默认为100，没有最大限制，最好不要超过1000，超过之后速度会比较慢，内存占用会比较多                                         |
| order     | string | 否   | 排序方式，取值`asc`和`desc`，默认为`desc`，按id排序，然后根据结果得分排序                                               |
| highlight | object | 否   | 关键字高亮，相对text字段中的文本                                                                           |
| negative  | object | 否   | 过滤负向词开关                                                                           |
| scoreExp  | string | 否   | 根据文档的字段计算分数，然后再进行排序，例如：score+document.hot*10，表达式中score为关键字的分数,document.hot为document中的hot字段 |
| filterExp  | string | 否   | 根据表达式过滤结果，然后再进行排序，例如：document.flag Like '中国'
| maxlimit     | int  | 否   | 当使用表达式时，限制最大输出结果，默认为1000，此值设置太大影响性能                               |



query参数(params-data)

| 字段     | 类型   | 必选  |  描述                  |
|----------|--------|-----|---------------------|
| database | string | 否   | 指定数据库名，不填默认为default |


### negative

> 配置以后，符合条件的负向关键词将会被过滤

| 字段      | 描述    | 值  |
|---------|-------|-------|
| query  | 关键词过滤  请求时满足条件直接拦截 |  true/false |
| content| 结果过滤   搜索结果满足条件时过滤 |   true/false |
### highlight

> 配置以后，符合条件的关键词将会被preTag和postTag包裹

| 字段      | 描述    | 示例  |
|---------|-------|-------|
| preTag  | 关键词前缀 | <span style='color:red'> |
| postTag | 关键词后缀 | </span> |

+ 示例

```json
{
  "query": "我爱中国",
  "page": 1,
  "limit": 10,
  "maxlimit": 1000,
  "filterExp": "flag Like '中国'",
  "scoreExp": "score+document.hot*10",
  "order": "desc",
  "negative": {
    "query" : true,
    "content" : false
  },
  "highlight": {
    "preTag": "<span style='color:red'>",
    "postTag": "</span>"
  }
}
```

+ POST /api/query

```json
{
  "query": "我爱中国",
  "page": 1,
  "limit": 10,
  "order": "desc"
}
```

+ 命令行

```bash
curl -H "Content-Type:application/json" -X POST --data '{"query":"我爱中国","page":1,"limit":10,"order":"desc"}' http://127.0.0.1:5678/api/query
```

### 响应

| 字段        | 类型      | 描述                      |
|-----------|---------|-------------------------|
| time      | float32 | 搜索文档用时                  |
| total     | int     | 符合条件的数量                 |
| pageCount | int     | 页总数                     |
| page      | int     | 当前页码                    |
| limit     | int     | 每页数量                    |
| documents | array   | 文档列表，[参考索引文档](#增加/修改索引) |

```json
{
  "state": true,
  "message": "success",
  "data": {
    "time": 2.75375,
    "total": 13487,
    "pageCount": 1340,
    "page": 1,
    "limit": 10,
    "documents": [
      {
        "id": 525810194,
        "text": "我爱中国",
        "document": {
          "id": "e489dd19dce0de2c9f4e59c969ec9ec0",
          "title": "我爱中国"
        },
        "score": 1
      }
    ],
    "words": [
      "我",
      "爱"
      "中国",
    ]
  }
}
```

## 查询状态

| 接口地址 | /api/status      |
|------|------------------|
| 请求方式 | GET              |

### 请求

```bash
curl http://127.0.0.1:5678/api/status
```

### 响应

```json
{
  "state": true,
  "message": "success",
  "data": {
    "index": {
      "queue": 0,
      "shard": 10,
      "size": 531971
    },
    "memory": {
      "alloc": 1824664656,
      "heap": 1824664656,
      "heap_idle": 10008625152,
      "heap_inuse": 2100068352,
      "heap_objects": 3188213,
      "heap_released": 9252003840,
      "heap_sys": 12108693504,
      "sys": 12700504512,
      "total": 11225144273040
    },
    "status": "ok",
    "system": {
      "arch": "arm64",
      "cores": 10,
      "os": "darwin",
      "version": "go1.18"
    }
  }
}
```

## 删除数据库

| 接口地址 | /api/db/drop |
|------|--------------|
| 请求方式 | GET          |

### 请求

```bash
curl http://127.0.0.1:5678/api/db/drop?database=db_name
```

### 响应

```json
{
  "state": true,
  "message": "success"
}
```

## 在线分词

| 接口地址 | /api/word/cut   |
|------|-----------------|
| 请求方式 | GET             |

### 请求参数

| 字段  | 类型     | 必选  | 描述  |
|-----|--------|-----|-----|
| q   | string | 关键词 |

### 请求

```bash
curl http://127.0.0.1:5678/api/word/cut?q=我爱中国
```

### 响应

```json
{
  "state": true,
  "message": "success",
  "data": [
    "我",
    "爱",
    "中国",
  ]
}
```