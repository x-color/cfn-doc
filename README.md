# cfn-doc

This is a tool for generating a document of CloudFormation template.

## Usage

```
cfn-doc is a tool for generating a document of CloudFormation template.

Usage:
  cfn-doc [OPTION] <TEMPLATE FILE>

Options:
  -o string
    	File name. Write to file instead of stdout
  -v	Show version
```

## Install

```sh
$ go install github.com/x-color/cfn-doc@latest
```

\*Required Go >= 1.16

## Sample

Generated document

```markdown
# sample-template.yaml

## Description

This template is sample CloudFormation template. It deploys EC2 Instance and Security Group of the instance.

## Parameters

| Name         | Description                                      | Type                 | Default  |
| ------------ | ------------------------------------------------ | -------------------- | -------- |
| AmiId        | -                                                | AWS::EC2::Image::Id  | -        |
| InstanceType | EC2 instance type                                | String               | t2.small |
| Subnet       | Id of Subnet in your Virtual Private Cloud (VPC) | AWS::EC2::Subnet::Id | -        |
| VPC          | Id of your existing Virtual Private Cloud (VPC)  | AWS::EC2::VPC::Id    | -        |

## Outputs

| Name                    | Description             |
| ----------------------- | ----------------------- |
| InstanceId              | Id of deployed instance |
| InstanceSecurityGroupId | -                       |

## Resources

| Resource              | Service Type            |
| --------------------- | ----------------------- |
| Instance              | AWS::EC2::Instance      |
| InstanceSecurityGroup | AWS::EC2::SecurityGroup |
```

\*cfn-doc does not format tables in a generated document. You use markdown formatter if you want to format it.

Source template

```yaml
AWSTemplateFormatVersion: 2010-09-09

Description: This template is sample CloudFormation template.
  It deploys EC2 Instance and Security Group of the instance.

Parameters:
  AmiId:
    Type: AWS::EC2::Image::Id

  Subnet:
    Description: Id of Subnet in your Virtual Private Cloud (VPC)
    Type: AWS::EC2::Subnet::Id

  VPC:
    Description: Id of your existing Virtual Private Cloud (VPC)
    Type: AWS::EC2::VPC::Id

  InstanceType:
    Description: EC2 instance type
    Type: String
    Default: t2.small

Resources:
  Instance:
    Type: AWS::EC2::Instance
    Properties:
      ImageId: !Ref AmiId
      InstanceType: !Ref InstanceType
      SubnetId: !Ref Subnet

  InstanceSecurityGroup:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: Allow http to client host
      VpcId: !Ref VPC
      SecurityGroupIngress:
        - IpProtocol: tcp
          FromPort: 80
          ToPort: 80
          CidrIp: 0.0.0.0/0
      SecurityGroupEgress:
        - IpProtocol: tcp
          FromPort: 80
          ToPort: 80
          CidrIp: 0.0.0.0/0

Outputs:
  InstanceId:
    Description: Id of deployed instance
    Value: !Ref Instance
    Export:
      Name: !Sub ${AWS::StackName}-InstanceId

  InstanceSecurityGroupId:
    Value: !Ref InstanceSecurityGroup
    Export:
      Name: !Sub ${AWS::StackName}-InstanceSecurityGroupId
```
