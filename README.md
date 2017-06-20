# tarmak

#### Table of Contents

1. [Description](#description)
2. [Classes](#classes)
## Description
Tarmak

This is the top-level class for the tarmak project. It's not including
any component. It's just setting global variables for the cluster

## Classes

### `tarmak`

Puppernetes

This is the top-level class for the tarmak project. It's not including
any component. It's just setting global variables for the cluster

#### Parameters

##### `dest_dir`

* path to installation directory for components
* Type: `String`
* Default: `'/opt'`

##### `bin_dir`

* path to the binary directory for components
* Type: `String`
* Default: `'/opt/bin'`

##### `cluster_name`

* a DNS compatible name for the cluster
* Type: `String`
* Default: `$tarmak::params::cluster_name`

##### `systemctl_path`

* absoulute path to systemctl binary
* Type: `String`
* Default: `$tarmak::params::systemctl_path`

##### `role`

* which role to build
* Type: `Enum['etcd','master','worker', nil]`
* Default: `nil`

##### `kubernetes_version`

* Kubernetes version to install
* Type: `String`
* Default: `$tarmak::params::kubernetes_version`

##### `kubernetes_user`

* Type: `String`
* Default: `'kubernetes'`

##### `kubernetes_group`

* Type: `String`
* Default: `'kubernetes'`

##### `kubernetes_uid`

* Type: `Integer`
* Default: `837`

##### `kubernetes_gid`

* Type: `Integer`
* Default: `837`

##### `kubernetes_ca_name`

* Type: `String`
* Default: `'k8s'`

##### `kubernetes_ssl_dir`

* Type: `String`
* Default: `'/etc/kubernetes/ssl'`

##### `kubernetes_config_dir`

* Type: `String`
* Default: `'/etc/kubernetes'`

##### `kubernetes_api_insecure_port`

* Type: `Integer[1,65535]`
* Default: `6443`

##### `kubernetes_api_secure_port`

* Type: `Integer[1,65535]`
* Default: `8080`

##### `kubernetes_pod_network`

* Type: `String`
* Default: `'10.234.0.0/16'`

##### `kubernetes_api_url`

* Type: `String`
* Default: `nil`

##### `kubernetes_api_prefix`

* Type: `String`
* Default: `'api'`

##### `kubernetes_authorization_mode`

* Type: `Array[Enum['AlwaysAllow', 'ABAC', 'RBAC']]`
* Default: `[]`

##### `dns_root`

* Type: `String`
* Default: `$tarmak::params::dns_root`

##### `hostname`

* Type: `String`
* Default: `$tarmak::params::hostname`

##### `etcd_cluster`

* Type: `Array[String]`
* Default: `[]`

##### `etcd_start_index`

* Type: `Integer[0,1]`
* Default: `1`

##### `etcd_user`

* Type: `String`
* Default: `'etcd'`

##### `etcd_group`

* Type: `String`
* Default: `'etcd'`

##### `etcd_uid`

* Type: `Integer`
* Default: `873`

##### `etcd_gid`

* Type: `Integer`
* Default: `873`

##### `etcd_home`

* Type: `String`
* Default: `'/etc/etcd'`

##### `etcd_ssl_dir`

* Type: `String`
* Default: `'/etc/etcd/ssl'`

##### `etcd_instances`

* Type: `Integer`
* Default: `3`

##### `etcd_advertise_client_network`

* Type: `String`
* Default: `$tarmak::params::etcd_advertise_client_network`

##### `etcd_overlay_client_port`

* Type: `Integer[1,65535]`
* Default: `2359`

##### `etcd_overlay_peer_port`

* Type: `Integer[1,65535]`
* Default: `2360`

##### `etcd_overlay_ca_name`

* Type: `String`
* Default: `'etcd-overlay'`

##### `etcd_overlay_version`

* Type: `String`
* Default: `'2.3.7'`

##### `etcd_k8s_main_client_port`

* Type: `Integer[1,65535]`
* Default: `2379`

##### `etcd_k8s_main_peer_port`

* Type: `Integer[1,65535]`
* Default: `2380`

##### `etcd_k8s_main_ca_name`

* Type: `String`
* Default: `'etcd-k8s'`

##### `etcd_k8s_main_version`

* Type: `String`
* Default: `'3.0.15'`

##### `etcd_k8s_events_client_port`

* Type: `Integer[1,65535]`
* Default: `2369`

##### `etcd_k8s_events_peer_port`

* Type: `Integer[1,65535]`
* Default: `2370`

##### `etcd_k8s_events_ca_name`

* Type: `String`
* Default: `'etcd-k8s'`

##### `etcd_k8s_events_version`

* Type: `String`
* Default: `'3.0.15'`

##### `cloud_provider`

* Type: `Enum['aws', '']`
* Default: `''`

##### `helper_path`

* Type: `String`
* Default: `$tarmak::params::helper_path`

##### `systemd_dir`

* Type: `String`
* Default: `'/etc/systemd/system'`

#### Examples

##### Declaring the class

```
include ::tarmak
```
##### Overriding the kubernetes version

```
class{'tarmak':
  kubernetes_version => '1.5.4'
}
```

### `tarmak::etcd`




### `tarmak::master`



#### Parameters

##### `disable_kubelet`

* Type: `Any`
* Default: `true`

##### `disable_proxy`

* Type: `Any`
* Default: `true`


### `tarmak::overlay_calico`




### `tarmak::params`

Defines parameters for other classes to reuse


### `tarmak::single_node`



#### Parameters

##### `dns_root`

* Type: `String`
* Default: `$tarmak::params::dns_root`

##### `cluster_name`

* Type: `String`
* Default: `$tarmak::params::cluster_name`

##### `etcd_advertise_client_network`

* Type: `String`
* Default: `$tarmak::params::etcd_advertise_client_network`

##### `kubernetes_api_url`

* Type: `String`
* Default: `nil`

##### `kubernetes_version`

* Type: `String`
* Default: `$tarmak::params::kubernetes_version`

##### `kubernetes_authorization_mode`

* Type: `Array[Enum['AlwaysAllow', 'ABAC', 'RBAC']]`
* Default: `[]`


### `tarmak::worker`


