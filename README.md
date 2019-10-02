# Arbitrary Fault-Tolerant and Locality-Aware MapReduce

This repository contains a *[Go](https://golang.org)* implementation of an **Arbitrary Fault-Tolerant and Locality-Aware MapReduce (AFT-LA) system** whose architecture is been fully described in this *[Report](https://github.com/AndreaG93/SDCC-Project-Report)*.

Project for: *Sistemi Distribuiti e Cloud Computing (SDCC) - A.A. 2018/19 Università degli Studi di Roma “Tor Vergata”*

## How to build and run client application

In order to build our application, make sure to have *[Go](https://golang.org)* and *[git](https://git-scm.com/)* properly installed.
After that, add required dependencies using following commands:

    go get -u github.com/samuel/go-zookeeper/zk
    go get -u github.com/aws/aws-sdk-go/service/s3/...
    go get -u github.com/aws/aws-sdk-go/aws/...
    go get -u github.com/Sirupsen/logrus

At this point, you can easly build client application as explained below:

    git -C ./go/src clone https://github.com/AndreaG93/SDCC-Project
    go build -o $HOME/wcclient $HOME/go/src/SDCC-Project/aftmapreduce/process/client/main/wcclient.go

To start client use:

    $HOME/wcclient [input text file] [IPv4 address] [IPv4 address] [IPv4 address] ...

> [IPv4 address] arguments are the public IPv4 addresses associated to the *[Amazon EC2](https://aws.amazon.com/it/ec2/)* instances where primary processes are running.

## How to build and run server side processes

We prefer to omit a detailed description about server side processes building because, having adopted some Amazon AWS services, *OUR* Amazon AWS token is required to run processes; then, unfortunately, it's not possible to test our application out of the box (we are sorry).

However following Amazon AWS services are required:

1. *[Amazon EC2](https://aws.amazon.com/it/ec2/)*
1. *[Amazon Elastic IP](https://docs.aws.amazon.com/en_us/AWSEC2/latest/UserGuide/elastic-ip-addresses-eip.html)* 
1. *[Amazon S3](https://aws.amazon.com/it/s3/)*

In general it's enough to create a certain amount of EC2 instances, at least two of which must be associated with different *Amazon Elastic IPs* in order to be used by clients worldwide. Each instance must be tagged as `PP-x` (Primary Process) or `WP-x` (Worker Process) where `x` is a number representing its ID (all process must to have different IDs, PP and WP can share, instead, same IDs). A S3 bucket must be provided. All processes can be configured using a `json` configuration file and IAM roles of EC2 instances where PPs are running must be propelry configured in order to access to Amazon S3.

The script which we have used to automatically configure OUR Amazon EC2 instances is available *[here](https://github.com/AndreaG93/SDCC-Project/blob/master/script/AWS-SystemConfiguration.sh)*. BE CAREFULL: this script is been design for a system made up of 3 PPs and 6 WPs, some changes are required in order to be used on different configurations. 

Techically is possible to test our application locally, but also in this case our Amazon AWS key is required. The script which we have used to build and run locally PPs and WPs is avalaible *[here](https://github.com/AndreaG93/SDCC-Project/blob/master/script/LocalSystemConfiguration.sh)* and *[here](https://github.com/AndreaG93/SDCC-Project/blob/master/script/StartSystemLocally.sh)* 
