[![Go Report Card](https://goreportcard.com/badge/github.com/LuciferInLove/kubernetes-nginx-operator)](https://goreportcard.com/report/github.com/LuciferInLove/kubernetes-nginx-operator)
![Build status](https://github.com/LuciferInLove/kubernetes-nginx-operator/workflows/Build/badge.svg)

# kubernetes-nginx-operator

An example of Kubernetes operator based on kubebuilder that creates nginx deployment.

Operator creates such resources:
* deployment
* service
* ingress
* certificate

## Installation

The repository contains a [helm](https://helm.sh/) chart to deploy **kubernetes-nginx-operator**.
You should have **ingress-controller** and **cert-manager** installed in the target Kubernetes cluster.
The chart contains dependencies, but they're skipped by default. As an **ingress-controller** the chart contains **nginx-ingress-controller** dependency.
See **charts/kubernetes-nginx-operator** folder for details.

## Usage

To deploy new Nginx to your cluster using **kubernetes-nginx-operator** you should create an Nginx resource. Example:

```yaml
apiVersion: nginx.custom-nginx.org/v1beta1
kind: Nginx
metadata:
  name: nginx-sample
spec:
  replicas: 1
  host: example.org
  image: "nginx:latest"
  serviceAccount: default
  certManagerIssuer: my-ca-issuer
```

Fields **host** and **certManagerIssuer** are mandatory, other fields are optional. You can see default values in CRD.

## Quick start for Linux users

* Install [kind](https://kind.sigs.k8s.io/docs/user/quick-start/) and [podman](https://podman.io/)
* [Enable cgroupsv2](https://kind.sigs.k8s.io/docs/user/rootless/) (enabled by default on Fedora)
* Create kind cluster config (optional). Example:
    ```yaml
    # three node (two workers) cluster config
    kind: Cluster
    apiVersion: kind.x-k8s.io/v1alpha4
    nodes:
    - role: control-plane
    - role: worker
    - role: worker
    ```
* Create kind cluster:
    ```shell
    KIND_EXPERIMENTAL_PROVIDER=podman kind create cluster --config cluster.yaml
    ```
* Deploy **kubernetes-nginx-operator** with dependencies using helm chart from the repository. **nginx-ingress-controller** and **cert-manager** will be also installed.
    ```shell
    $ helm dependency update
    $ helm install kubernetes-nginx-operator ./charts/kubernetes-nginx-operator \
        --set installDependencies=true
    ```
* Create Own-CA Issuer for **cert-manager**.
    ```yaml
    apiVersion: cert-manager.io/v1
    kind: ClusterIssuer
    metadata:
      name: selfsigned-issuer
    spec:
      selfSigned: {}
    ---
    apiVersion: cert-manager.io/v1
    kind: Certificate
    metadata:
      name: my-selfsigned-ca
    spec:
      isCA: true
      commonName: my-selfsigned-ca
      secretName: root-secret
      privateKey:
        algorithm: ECDSA
        size: 256
      issuerRef:
        name: selfsigned-issuer
        kind: ClusterIssuer
        group: cert-manager.io
    ---
    apiVersion: cert-manager.io/v1
    kind: Issuer
    metadata:
      name: my-ca-issuer
    spec:
      ca:
        secretName: root-secret
    ```
* Create Nginx resource. You can use the example above as is.
* Add the next string to /etc/hosts:
    ```shell
    127.0.0.1 example.org
    ```
* Expose **ingress-controller** port externally. You can use port-forward:
    ```shell
    kubectl port-forward service/nginx-ingress-controller 8443:443
    ```
* Try to open https://example.org:8443 in your browser. Skip certificate warning. You should see "Welcome to Nginx" page.

## Local development with kind/podman/buildah

* Build **kubernetes-nginx-operator** after changing the code.
    ```shell
    $ buildah bud -f local-development/Dockerfile .
    ```
* Save the built image to the tar file.
    ```shell
    $ podman tag built_image_id kubernetes-nginx-operator:v0.0.1
    $ podman save kubernetes-nginx-operator:v0.0.1 -o controller.tar
    ```
* Load the image from the tar file to the kind cluster.
    ```shell
    $ kind load image-archive controller.tar
    ```
* Update version in **local-development/values.yaml** file if needed.
* Install/Upgrade **kubernetes-nginx-operator** using chart.
    ```shell
    $ helm dependency update
    $ helm install kubernetes-nginx-operator ./charts/kubernetes-nginx-operator \
        --set installDependencies=true \
        -f local-development/values.yaml
    ```
