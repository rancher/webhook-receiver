# webhook-receiver

Redirect webhook requests from prometheus alertmanager to targets

## Building

`make` 

## Running

`./bin/webhook-receiver -config /myconfig.yaml` 

``` 
providers:
  dingtalk:
    type: DINGTALK
    webhook_url: <webhook_url>
    secret: <optional_secret>
    proxy_url: <optional_proxy_url>
  msteams:
    type: MICROSOFT_TEAMS
    webhook_url: <webhook_url>
    proxy_url: <optional_proxy_url>
  aliyunsms:
    type: ALIYUN_SMS
    access_key_id: <access_key_id>
    access_key_secret: <access_key_secret>
    sign_name: <sign_name>
    template_code: <template_code>
    proxy_url: <optional_proxy_url>

receivers:
  test1:
    provider: dingtalk
  test2:
    provider: msteams
  test3:
    provider: aliyunsms
    to:
      - <phone_number_1>
      - <phone_number_2>

logLevel: Info
```

## License

Copyright (c) 2014-2020 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License"); 
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, 
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

