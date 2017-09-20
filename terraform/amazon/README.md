# Kubernetes Platform Terraform

## Terraform

![terraform architecture](docs/terraform-diagram.png "Terraform Architecture")

### Stacks

#### network

Is used for every cluster and for the tools and vault in the hub

Contains:

* VPC, IGW, NATGW, Route Tables Subnet
* Optional: private public DNS zones
* Optional: terraform state bucket

#### tools

* Jenkins
  * [Jenkins (non-prod)]:(https://jenkins.nonprod.jetsky.jetstack.net/)
* Puppet
  * [Foreman Dashboard (non-prod)]:(https://foreman.nonprod.jetsky.jetstack.net/)
* Bastion

#### kubernetes

* Kubernetes Master ASGs
* Kubernetes Worker ASGs
* ETCD nodes
* ELBs

#### vault

## Rake tricks

### AWS

* To login with MFA you can use a temporary token like that:

```
eval $(bundle exec rake aws:login_mfa )
```

It will read the MFA serial from `.aws/config` and generate temporary tokens that are exported to then environment

### Terraform


#### Plan Hub Network and Tools

```
bundle exec rake terraform:plan TERRAFORM_NAME=hub TERRAFORM_ENVIRONMENT=nonprod TERRAFORM_STACK=network
bundle exec rake terraform:plan TERRAFORM_NAME=hub TERRAFORM_ENVIRONMENT=nonprod TERRAFORM_STACK=tools
```
