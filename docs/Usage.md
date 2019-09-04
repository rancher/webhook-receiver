# webhook-receiver使用文档

## rancher平台下，使用应用商店进行部署，[chart地址](https://github.com/gangchang/pandaria-catalog)

## NOTE

- 部署的时候，将配置内容base64编码后，对应表单的key为config，value为base64编码后的配置
- rancher必须开启监控（数据源来自监控的alert manager组件）
- 默认开启9094端口
- 添加通知者的时候，url的最后一个路径对应配置文件中的receiver的名称


## 使用示例

webhook-receiver部署完成，其namespace为kube-system，config对应的值为下方配置示例内容base64编码后的值。添加通知者，类型为webhook，url为http://webhook-receiver:9094/test1,则该通知者收到的告警消息都会发送至dingtalk的机器人，对应示例配置中的receivers.test1的配置。


## 配置示例
```yaml
providers:
  alibaba:
    access_key_id: real_access_key_id
    access_key_secret: real_access_key_secret
    sign_name: real_sign_name
    template_code: real_template_code
  dingTalk:
    webhook_url: real_webhook_url
receivers:
  test1:
    provider: dingTalk
  test2:
    provider: alibaba
    to:
    - 测试的手机号
```
