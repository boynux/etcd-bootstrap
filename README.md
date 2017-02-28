[![CircleCI](https://circleci.com/gh/stylight/etcd-bootstrap/tree/master.svg?style=svg&circle-token=9bbd5dd7529b97a08f07eaacdafd122caa09990f)](https://circleci.com/gh/stylight/etcd-bootstrap/tree/master)

# Sample Cloudformation stack
[![Launch Stack](https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png)](https://console.aws.amazon.com/cloudformation/home?region=eu-west-1#/stacks/new?stackName=etcd-bootstrap-test&templateURL=https://s3-eu-west-1.amazonaws.com/packages.stylight.net/cloudformation/etcd-bootstrap/example-1-cloudformation.yaml)


# What is it?
This is a single Go binary that bootstraps ETCD cluster in AWS Autoscaling Group

# Why we made this:
Maintaining and managing ETCD cluster has a significant operational cost. We wanted to make it as simple as possible. So we decided to write this app to make it easier to bootstrap an ETCD cluster within AWS Autoscaling groups.

We also wanted to be able to point our etcd clients to one single endpoint. Which is not possible.

# What it does:
When starting up, an ec2 instance in an autoscaling group calls etcd-bootstrap.
etcd-bootstrap checks -
- if other ec2 instances of the austoscalling group are already forming a cluster. It creates an etcd  config to join the cluster.
- if the cluster formed by other instances in the autoscalling group mention members does do not exist anymore (their ec2 instances was removed maybe). It then removes those members from the cluster.

# How to use it:
The binary if runs inside an EC2 instance that belongs to autoscaling group will output either ENV variables or command line arguments for Etcd2 application. So in the simplest way it can be used like this:

    $ etcd2 $(etcd-boostrap)

If you want to use it with CloudInit and coreOS perhaps the best option is to use ENV variables.

    $ etcd-bootstrap -output env > /etc/etcd.env

Sample CloudFormation is provided in examples directory or you can click on the button above to launch a cluster.

Command line args:

    Usage of ./etcd-bootstrap:
      -add-member
        	Add this etcd member explicitly to the cluster (default true)
      -client-port int
        	ETCD Cient port (default 2379)
      -output string
        	Output format. Available options: args, env (default "args")
      -public
        	Use EC2 Public IP for client URLs if available
      -quiet
        	Disable log output
      -region string
        	Region to initialize the script. (default "eu-west-1")
      -schema string
        	Schema to communicate to the cluster, currently only 'http' works (default "http")


# State:
This is a work in progress and not ready for production.

# Contribution:
Pull requests and issues are welcome.

### Brought to you proudly by Cloud Team @Stylight
