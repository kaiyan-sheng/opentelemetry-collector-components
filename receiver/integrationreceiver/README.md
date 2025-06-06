# Integration receiver

<!-- status autogenerated section -->
| Status        |           |
| ------------- |-----------|
| Stability     | [development]: logs, metrics, traces   |
| Distributions | [] |
| Issues        | [![Open issues](https://img.shields.io/github/issues-search/open-telemetry/opentelemetry-collector-contrib?query=is%3Aissue%20is%3Aopen%20label%3Areceiver%2Fintegration%20&label=open&color=orange&logo=opentelemetry)](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues?q=is%3Aopen+is%3Aissue+label%3Areceiver%2Fintegration) [![Closed issues](https://img.shields.io/github/issues-search/open-telemetry/opentelemetry-collector-contrib?query=is%3Aissue%20is%3Aclosed%20label%3Areceiver%2Fintegration%20&label=closed&color=blue&logo=opentelemetry)](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues?q=is%3Aclosed+is%3Aissue+label%3Areceiver%2Fintegration) |
| Code coverage | [![codecov](https://codecov.io/github/open-telemetry/opentelemetry-collector-contrib/graph/main/badge.svg?component=receiver_integration)](https://app.codecov.io/gh/open-telemetry/opentelemetry-collector-contrib/tree/main/?components%5B0%5D=receiver_integration&displayType=list) |
| [Code Owners](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/CONTRIBUTING.md#becoming-a-code-owner)    | [@jsoriano](https://www.github.com/jsoriano) |

[development]: https://github.com/open-telemetry/opentelemetry-collector/blob/main/docs/component-stability.md#development
<!-- end autogenerated section -->

This receiver can be used to create other receivers based on pipelines
built on the combination of other receiver and a chain of processors. These
pipelines are instantiated internally when the receiver is started. This
approach allows the creation of new receivers based on existing components,
without writing Go code.

To use it, you need integrations. Integrations are files that look very similar to
the collector configuration. They contain receivers, processors and
pipelines, and they can contain placeholders for parametrization. There are
different extensions that can be used to obtain integrations during runtime.

## Configuration

**name**

Name of the integration to instantiate. Finder extensions are used to resolve
it.

**pipelines**

Set of pipelines to instantiate as part of this receiver. These pipelines need
to exist in the integration referenced by its `name`. If not set, all pipelines
in the integration are used.

**parameters**

Parameters used to resolve placeholders in the integrations.

## Integrations

Integrations are defined in YAML format, with a structure similar to the
collector configuration. They can define receivers, processors, and pipelines.

For example the following integration would define a logs receiver based on the
`filelog` receiver and the `transform` processor:
```yaml
receivers:
  filelog:
    include: ${var:paths}
    start_at: beginning

processors:
  transform/resource:
    log_statements:
      - set(log.attributes["resource"], "${var:resource}")

pipelines:
  logs:
    receiver: filelog
    processors:
      - transform/resource
```

If this template would be installed for example in the `/etc/opentelemetry/integrations/somelog.yaml`
file, it could be used with the following configuration:
```yaml
extensions:
  file_integrations:
    path: "/etc/opentelemetry/integrations"

receivers:
  integration:
    name: "somelog"
    pipelines: ["logs"]
    parameters:
      paths: "/var/log/somelog-*.log"
      resource: "example"

exporters: ...

service:
  extensions: [file_integrations]
  pipelines:
    logs:
      receivers: [integration]
      exporters: ...
```

The configuration above would create a receiver that would use internally the
`filelog` processor to collect logs, and the `transform` processor to attach an
attribute.
