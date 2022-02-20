# Kubernetes-Nginx-Operator

![Version: 0.0.1](https://img.shields.io/badge/Version-0.0.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

An example of Kubernetes operator that deploys nginx

## Installing the Chart ignoring dependencies

```shell
$ helm dependency update
$ helm install kubernetes-nginx-operator ./charts/kubernetes-nginx-operator
```

## Installing the Chart with dependencies

```shell
$ helm dependency update
$ helm install kubernetes-nginx-operator ./charts/kubernetes-nginx-operator \
    --set installDependencies=true
```

## Uninstalling the Chart

```shell
$ helm uninstall kubernetes-nginx-operator
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| containerSecurityContext.allowPrivilegeEscalation | bool | `false` |  |
| deploymentStrategy.rollingUpdate.maxSurge | int | `0` |  |
| deploymentStrategy.rollingUpdate.maxUnavailable | int | `1` |  |
| deploymentStrategy.type | string | `"RollingUpdate"` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"quay.io/luciferinlove/kubernetes-nginx-operator"` |  |
| image.tag | string | `"v0.0.1"` |  |
| installDependencies | bool | `false` |  |
| kubeRbacProxy.image.pullPolicy | string | `"IfNotPresent"` |  |
| kubeRbacProxy.image.repository | string | `"gcr.io/kubebuilder/kube-rbac-proxy"` |  |
| kubeRbacProxy.image.tag | string | `"v0.8.0"` |  |
| podSecurityContext.runAsNonRoot | bool | `true` |  |
| replicas | int | `1` |  |
| resources.limits.cpu | string | `"500m"` |  |
| resources.limits.memory | string | `"128Mi"` |  |
| resources.requests.cpu | string | `"10m"` |  |
| resources.requests.memory | string | `"64Mi"` |  |
| serviceAccount.name | string | `""` |  |
| terminationGracePeriodSeconds | int | `10` |  |
| tolerations | list | `[]` |  |
