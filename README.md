# webhook-receiver

webhook-receiver可以对接prometheus的alertmanager，通过配置可以将alertmanager的告警发送至对应的后端

## Support

- 阿里云短信服务
- 钉钉机器人

## Installation

通过rancher的应用商店进行部署，[应用商店地址](https://github.com/cnrancher/pandaria-catalog), 其下的webhook-receiver及为该服务，详细的配置参考其对应的readme

## [Sample config](https://github.com/cnrancher/webhook-receiver/blob/master/examples/config.yaml)


# License

Copyright (c) 2014-2019 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.