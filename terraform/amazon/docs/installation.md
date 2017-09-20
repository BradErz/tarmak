# Installation Happy Path

## Prerequisites

- Linux / Mac OS machine
- Docker available
- AWS credentials are configured using these environment variables:
  ```
  $ env | grep AWS_
  AWS_ACCESS_KEY_ID=
  AWS_SECRET_ACCESS_KEY=
  AWS_SESSION_TOKEN=<optional>
  ```

## Setting up Hub & Tools

### Accept CentOS T's & C's

AWS CentOS images require the acceptance of T's & C's first, make sure you have done that for every account, where you want to use the toolkit

* Make sure your browser is logged into AWS console
* Go to https://aws.amazon.com/marketplace/pp/B00O7WM7QW
* Click on "Continue"
* Click on "Accept T's & C's"

### Build customised and updated CentOS AMI using packer

Now you can run packer:

```
$ make clean build packer_build
[...]
==> amazon-ebs: Prevalidating AMI Name...   (If the execution hangs here, make sure you have accepted the T's & C's)
[...]
==> Builds finished. The artifacts of successful builds are:
--> amazon-ebs: AMIs were created:

eu-central-1: ami-xxxxxxxx
```

### Upgrade ami-ids

Upgrade the ami id in `tfvars/global.tfvars`:

```
centos_ami = {
  eu-central-1 = "ami-xxxxxxxx"
}

```

### Ensure SSH key pairs exist

```
# Generate and upload an AWS key pair
# Generate jenkins
make clean build credentials_ensure
```

The Jenkins key pair needs to added to all Github repositories with read access

### Setup tfvar files for hub

#### `tfvars/${ENVIRONMENT}/hub/state.tfvars`

- `public_zone` has to be a zone, which is used and available publicly (including delegations)
- `bucket_prefix` needs to be string which prefixes all buckets

```
public_zone = "nonprod.p9s.jetstack.net"
bucket_prefix = "jetstack-p9s-"
```

#### `tfvars/${ENVIRONMENT}/hub/network.tfvars`

- `network should` be a private IPv4 CIDR which is not yet used by your organisation
- `private_zone` has be a zone, which can be arbitrary and is only available to services with a VPC

```
network = "10.99.128.0/20"
private_zone = "nonprod-private.p9s.jetstack.net"
```

#### `tfvars/${ENVIRONMENT}/hub/tools.tfvars`

- `foreman_admin_password` generate a new password used for the puppet dashboard of foreman
- `puppet_deploy_key` needs to cotain the public key from `credentials/jenkins_key_pair.pub`

```
foreman_admin_password = "secure123"
puppet_deploy_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDncOnSrQmQ+xZ0MLEEiCubalsBrZmaztkeC1CjVzJMbxlNCab9vkTGgzBdC9VgBk/DBagUbMqcBHZvz98ESOtrab/m3WPreTy4vMPqBt1LBORq4n9enIh/DUZqY4H6sY0y1e2wwgHthsXer5XgqkD6KkRCvgCggPZARYKKhjRQkai2p08e0U2SBcA6IC7lrZWZQTC6RToqXvRtjMxpd5t94SilnFfA42KJvnvaajH3NQqgFNilY+5uVjkQL88wb5uP/L7NrkZpZ2meDR0El4pHaZVjIf6dzyUcYn2+FMP5ux9Vfoab0RgWgq5L+25T2nZho3gwGtWNYGatYtfXq7Jrk/iaWAEWliOquIdiAXo5JAyc4CXZVvR3aK98/iHf0KWH0nqZcA/1PA071GnbKkDgCrHHauRNNtEsnF9076nz2m1jPSwoivsFtXo0j9siFITGy6IiAgts0EzGtLj5/pNlsy9Jw8UYpUmaeRny8kCwc79ZnPDVn6fNKOG/yODkNq2CjVyxrle3NYus3rNMT45+WGV930RnYlvzuzLIrAVRMjZxKFTp8+mNoNTyMbTBit9lBX8JNh2OT56OCeUWnoLh+DRTZ0B1+CY3TniAZlT6IbhB0ZprVUVGibAPPkXCkWTMJkese76Fm12Do7RSP0rghiQBlkL3SZFQG44tfOWm2w== puppetmaster"
```

#### `tfvars/${ENVIRONMENT}/hub/vault.tfvars`

- `consul_encypt` is a secret that is used by the Consul instances to secure
  communication within the cluster. Generate a unique secret for every
  environment using `openssl rand -base64 16`
- `consul_master_token` is the token to authorize Vault against the consul
  cluster. Generate for every environment using `uuidgen`

```
consul_encypt = "<16bytes-in-base64>"
consul_master_token = "<UUID>"
```

### Setup state hub

```
# Create network stack for hub (contains state buckets/dynamodb public zone)
export TERRAFORM_STACK=state TERRAFORM_NAME=hub TERRAFORM_ENVIRONMENT=nonprod

## plan
TERRAFORM_DISABLE_REMOTE_STATE=true make clean build terraform_sync terraform_plan

## apply
TERRAFORM_DISABLE_REMOTE_STATE=true make terraform_apply

## sync local state to remote state (response yes)
make terraform_plan
```

### Ensure AWS Certificates exist

The public domain zone (`public_zone = nonprod.p9s.jetstack.net`) which is
specified needs a AWS Wild Certificate created.

- Go to this AWS page: https://console.aws.amazon.com/acm/home
- Request a new certificate which is valid for these names:
  -  `*.public_zone`
  -  `*.devcluster.public_zone`
- Validate the validity of the request and wait till the certificate is issued

### Ensure Route 53 public zone delegation

The public domain zone needs to be delegated from the DNS root. You get the
nameservers for it from the output of the last `terraform apply`:

```
public_zone_name_servers = [
    ns-ww.awsdns-57.net,
    ns-xx.awsdns-19.org,
    ns-yy.awsdns-03.co.uk,
    ns-zz.awsdns-63.com
]
```

Test the delegation:

```
$ dig -t txt +short _tarmak.nonprod.p9s.jetstack.net
"delegation for nonprod-hub works"
```

### Setup network hub

```
# Create network stack for hub (contains state buckets/dynamodb)
export TERRAFORM_STACK=network TERRAFORM_NAME=hub TERRAFORM_ENVIRONMENT=nonprod

## plan
make clean build terraform_sync terraform_plan

## apply
make terraform_apply
```

### Setup tools hub

```
# Create network stack for hub (contains state buckets/dynamodb)
export TERRAFORM_STACK=tools TERRAFORM_NAME=hub TERRAFORM_ENVIRONMENT=nonprod

## plan
make clean build terraform_sync terraform_plan

## apply
make terraform_apply
```

### Ensure you can connect to the bastion instance

The only instance with a public IP address directly assigned is the bastion instance. It's used to connect all other instances.

```
# Connect to bastion instance
ssh -i credentials/aws_key_pair -o IdentitiesOnly=yes centos@bastion.nonprod.p9s.jetstack.net
```

### Initialise Jenkins

#### Setup jenkins itself

```
# SSH into jenkins instance
ssh -i credentials/aws_key_pair -o IdentitiesOnly=yes -o ProxyCommand="ssh -W %h:%p -i credentials/aws_key_pair -o IdentitiesOnly=yes centos@bastion.nonprod.p9s.jetstack.net" centos@jenkins.nonprod-private.p9s.jetstack.net

# Retrieve initial password
sudo cat /var/lib/jenkins/secrets/initialAdminPassword

# Go to Jenkins
> https://jenkins.nonprod.p9s.jetstack.net/

# Put initial password in

# Install suggested plugins

# Setup an admin user account
```

#### Bootstrap jobs/credentials

```
export TERRAFORM_ENVIRONMENT=nonprod
export JENKINS_URL=https://jenkins.nonprod.p9s.jetstack.net/
export JENKINS_USER=admin
export JENKINS_PASSWORD=admin
make clean build jenkins_initialize
```

### Setup Vault in the Hub using Jenkins job

* Run job `${ENVIRONMENT}/hub/terraform/vault`

### Build puppet_code

* Run job `${ENVIRONMENT}/puppet_code`


### Setup first Cluster `devcluster`

#### Setup parameters in `tfvars`

##### `tfvars/${ENVIRONMENT}/devcluster/network.tfvars`

* `network` to use for the cluster's VPC. Should be unique and non overlapping with other VPCs
* `vpc_peer_stack` this line just tells to peer with the `hub`

```
network = "10.61.16.0/20"
vpc_peer_stack = "hub"
```

##### `tfvars/${ENVIRONMENT}/devcluster/kubernetes.tfvars`

* Just create an empty file for now:

```

```

#### Setup cluster `devcluster` using Jenkins

* Run job `${ENVIRONMENT}/devcluster/terraform/network`
* Run job `${ENVIRONMENT}/devcluster/terraform/kubernetes`

## Known issues

### `Signature expired: TIMESTAMP is now earlier than TIMESTAMP`

Your machine (or the docker Virtual Machine's time is out of sync). This
happens mostly on Mac OS. Workaround:

```
docker run --rm --privileged alpine hwclock -s
```
