package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkCardanoNodeDevelopmentStackProps struct {
	awscdk.StackProps
}

func NewCdkCardanoNodeDevelopmentStack(scope constructs.Construct, id string, props *CdkCardanoNodeDevelopmentStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here
	ami := ec2.NewLookupMachineImage(&ec2.LookupMachineImageProps{
		Name:   jsii.String("whslabs-cardano-node-*"),
		Owners: jsii.Strings("102933037533"),
	})

	vpc := ec2.NewVpc(stack, jsii.String("Vpc"), &ec2.VpcProps{
		Cidr: jsii.String("10.0.0.0/16"),
	})

	key := ec2.NewCfnKeyPair(stack, jsii.String("Key"), &ec2.CfnKeyPairProps{
		KeyName: jsii.String("CdkCardanoNodeDevelopments"),
	})

	sg := ec2.NewSecurityGroup(stack, jsii.String("Sg"), &ec2.SecurityGroupProps{
		Vpc:              vpc,
		AllowAllOutbound: jsii.Bool(true),
	})

	sg.AddIngressRule(ec2.Peer_AnyIpv4(), ec2.Port_Tcp(jsii.Number(22)), jsii.String(""), jsii.Bool(false))

	userdata := ec2.MultipartUserData_ForLinux(&ec2.LinuxUserDataOptions{})

	userdata.AddCommands(jsii.String("systemctl disable cardano-node cardano-db-sync"))
	userdata.AddCommands(jsii.String("systemctl stop cardano-node cardano-db-sync"))

	userdata.AddCommands(jsii.String("systemctl enable cardano-node-testnet"))
	userdata.AddCommands(jsii.String("systemctl start cardano-node-testnet"))

	instance := ec2.NewInstance(stack, jsii.String("Instance"), &ec2.InstanceProps{
		InstanceType:  ec2.NewInstanceType(jsii.String("t3.medium")),
		KeyName:       key.KeyName(),
		MachineImage:  ami,
		SecurityGroup: sg,
		UserData:      userdata,
		Vpc:           vpc,
		BlockDevices: &[]*ec2.BlockDevice{
			&ec2.BlockDevice{
				DeviceName: jsii.String("/dev/sda1"),
				Volume:     ec2.BlockDeviceVolume_Ebs(jsii.Number(50), &ec2.EbsDeviceOptions{}),
			},
		},
		VpcSubnets: &ec2.SubnetSelection{
			SubnetType: ec2.SubnetType_PUBLIC,
		},
	})

	// example resource
	// queue := awssqs.NewQueue(stack, jsii.String("CdkCardanoNodeDevelopmentQueue"), &awssqs.QueueProps{
	// 	VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	// })

	awscdk.NewCfnOutput(stack, jsii.String("WillChangeIp"), &awscdk.CfnOutputProps{
		Value: instance.InstancePublicIp(),
	})

	awscdk.NewCfnOutput(stack, jsii.String("KeyId"), &awscdk.CfnOutputProps{
		Value: key.AttrKeyPairId(),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewCdkCardanoNodeDevelopmentStack(app, "CdkCardanoNodeDevelopmentStack", &CdkCardanoNodeDevelopmentStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
