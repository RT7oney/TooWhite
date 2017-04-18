# TooWhite
一个websocket消息分组和个人的广播系统，使用mongodb作为后端数据库服务支持

## url:
ws://xx.xxx.xx.xx:1234

##  支持格式
    json

##  请求参数
| 字段                     |   必选            |   类型及范围    | 说明                               |
|:-------------------------|:----------------- |:----------------|:-----------------------------------|
|msgtype|true|int|消息类型|
|from|true|string|消息发起人的唯一标识ID|
|target|true|string|消息接受者（个人/群组）的唯一标识ID|
|data|true|string|字符串化的数据（可根据之后的需求自由定制）|


##  请求示例
```javascript
so = new WebSocket(url);
so.send('{"MsgType":0,"From":"XXXXX","Target":"XXXX","Data":"XXXX"}');
```

##  注意事项
### 用户加入
* msgtype 值为 0
* from 传入用户在系统中的唯一标识ID
* target null
* data 传入用户的昵称

### 增加一个新的分组
* msgtype 值为 1
* from 传入用户在系统中的唯一标识ID
* target null 
* data 分组的名称

### 用户加入一个分组
* msgtype 值为2
* from 传入用户在系统中的唯一标识ID
* target null
* data 需要加入的分组的名称

### 点对点消息
* msgtype 值为3
* from 传入用户在系统中的唯一标识ID
* target 传入用户在系统中的唯一标识ID
* data 消息内容

### 群组消息
* msgtype 值为4
* from 传入用户在系统中的唯一标识ID
* target 传入分组在系统中的唯一标识ID
* data 消息内容

### 删除分组
* msgtype 值为5
* from 传入用户在系统中的唯一标识ID
* target null
* data 需要删除的分组的名称（需要该用户为分组的创建人）

### 用户离开分组
* msgtype 值为6
* from 传入用户在系统中的唯一标识ID
* target null
* data 需要离开的分组的名称

