# Silicom STS Special Resource Usage
![alt text](spec/sts-sro.png "Special Resource Operator")

# Silicom STS Operator
![alt text](spec/sts-operator.png "STS Overview")

# Silicom STS Operator deployments
![alt text](spec/sts-deployments.png "STS Deployments")

## Table of Contents
- [STS Operator](#sts-operator)
- [STS Discovery](#sts-discovery)
- [StsConfig](#stsconfig)
- [Quick Start](#quick-start)

## STS Operator
Sts Operator, runs in `sts-silicom` namespace, manages cluster wide STS configurations. It offers `StsConfig` and `StsOperatorConfig` CRDs and creates `tsyncd` to apply node specific STS config.

## STS Discovery daemonset
Once NFD operator has labelled the nodes, this daemonset queries the network interfaces and STS specific information and accordingly labels the nodes. (feature.node.kubernetes.io/custom-silicom.sts.devices: "true")

## STS daemonset
This consists of the following daemons on labelled nodes, all of these running in the same pod. Nodes labelled with the sts.silicom.com/config: "gm-1"
* tsyncd
* GPSd (if in T-GM.8275.1 mode)
* ts2phcs
* phc2sys
* grpc-tsyncd

## StsOperatorConfig
Example (default)
```yaml
apiVersion: sts.silicom.com/v1alpha1
kind: StsOperatorConfig
metadata:
  name: stsoperatorconfig
  namespace: sts-silicom
spec:
  imageRegistry: quay.io/silicom
  stsVersion: 2.0.1.0
  iceVersion: 1.6.4
  grpcSvcPort: 50051
  gpsSvcPort: 2947

```

## StsConfig
Example (Grand Master mode)
```yaml
apiVersion: sts.silicom.com/v1alpha1
kind: StsConfig
metadata:
  name: gm-1
  namespace: sts-silicom
spec:
  name: gm-1
  nodeSelector:
    sts.silicom.com/config: "gm-1"
  mode: T-GM.8275.1
  namespace: sts-silicom
  interfaces:
    - ethName: enp2s0f0
      synce: true
      holdoff: 500
      mode: Slave
    - ethName: enp2s0f1
      synce: true
      holdoff: 500
      mode: Master

```

```
fb@g9:~$ oc get stsconfig.sts.silicom.com  -n sts-silicom -o yaml
apiVersion: v1
items:
- apiVersion: sts.silicom.com/v1alpha1
  kind: StsConfig
  metadata:
    annotations:
      kubectl.kubernetes.io/last-applied-configuration: |
        {"apiVersion":"sts.silicom.com/v1alpha1","kind":"StsConfig","metadata":{"annotations":{},"name":"gm-1","namespace":"sts-silicom"},"spec":{"interfaces":[{"ethName":"enp2s0f0","holdoff":500,"synce":true},{"ethName":"enp2s0f1","holdoff":500,"synce":true}],"mode":"gm","name":"gm-1","namespace":"sts-silicom","nodeSelector":{"sts.silicom.com/config:"gm-1"}}}
    creationTimestamp: "2021-10-25T13:40:30Z"
    generation: 1
  spec:
    interfaces:
    - ethName: enp2s0f0
      holdoff: 500
      synce: true
      mode: Slave
    - ethName: enp2s0f1
      holdoff: 500
      synce: true
      mode: Master
    mode: T-GM.8275.1
    name: gm-1
    namespace: sts-silicom
    nodeSelector:
      sts.silicom.com/config: "gm-1"
kind: List
metadata:
  resourceVersion: ""
  selfLink: ""

```
