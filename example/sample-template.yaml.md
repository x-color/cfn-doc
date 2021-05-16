# sample-template.yaml

## Description

This template is sample CloudFormation template. It deploys EC2 Instance and Security Group of the instance.

## Parameters

| Name | Description | Type | Default |
| --- | --- | --- | --- |
| AmiId | - | AWS::EC2::Image::Id | - |
| InstanceType | EC2 instance type | String | t2.small |
| Subnet | Id of Subnet in your Virtual Private Cloud (VPC) | AWS::EC2::Subnet::Id | - |
| VPC | Id of your existing Virtual Private Cloud (VPC) | AWS::EC2::VPC::Id | - |

## Outputs

| Name | Description |
| --- | --- |
| InstanceId | Id of deployed instance |
| InstanceSecurityGroupId | - |

## Resources

| Resource | Service Type |
| --- | --- |
| Instance | AWS::EC2::Instance |
| InstanceSecurityGroup | AWS::EC2::SecurityGroup |
