# kubesplit

Breaks down single-page YAML kubernetes specs into individual YAML files.

## Motivation

Why? https://github.com/operator-framework/operator-sdk/issues/4930

While developing controllers with operator-sdk, there seems to not be an obvious way to go from the generated files directly to templating for a standard helm chart.

This simple tool tries to make it as much as programmatically as possible so it can be chained and plugged in the operator-sdk make targets easily.

## How to use

e.g. just pipe it to your kustomize, to break down files by type:

```
$ cat some-big-file.yaml | kubesplit chart
Adding 'chart/configmap.yaml' (ConfigMap)
Adding 'chart/service.yaml' (Service)
Adding 'chart/service.yaml' (Service)
Adding 'chart/deployment.yaml' (Deployment)
Adding 'chart/certificates.yaml' (Certificate)
Adding 'chart/certificates.yaml' (Issuer)
Adding 'chart/mutatingwebhook.yaml' (MutatingWebhookConfiguration)

$ tree chart                                         
chart
├── certificates.yaml
├── configmap.yaml
├── deployment.yaml
├── mutatingwebhook.yaml
└── service.yaml

```

## License

Copyright (c) 2022 Spectro Cloud

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
