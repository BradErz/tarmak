package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (a *AWS) RemoteStateName() string {
	return fmt.Sprintf(
		"%s%s-terraform-state",
		a.conf.AWS.BucketPrefix,
		a.Region(),
	)
}

const DynamoDBKey = "LockID"

// TODO: remove me, deprecated
func (a *AWS) RemoteStateBucketName() string {
	return a.RemoteStateName()
}

func (a *AWS) RemoteState(namespace string, clusterName string, stackName string) string {
	return fmt.Sprintf(`terraform {
  backend "s3" {
    bucket = "%s"
    key = "%s"
    region = "%s"
    lock_table ="%s"
  }
}`,
		a.RemoteStateName(),
		fmt.Sprintf("%s/%s/%s.tfstate", namespace, clusterName, stackName),
		a.Region(),
		a.RemoteStateName(),
	)
}

func (a *AWS) RemoteStateBucketAvailable() (bool, error) {
	svc, err := a.S3()
	if err != nil {
		return false, err
	}

	_, err = svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(a.RemoteStateName()),
	})
	if err == nil {
		return true, nil
	} else if strings.HasPrefix(err.Error(), "NotFound:") {
		return false, nil
	}

	return false, fmt.Errorf("error while checking if remote state is available: %s", err)
}

func (a *AWS) RemoteStateAvailable(bucketName string) (bool, error) {
	sess, err := a.Session()
	if err != nil {
		return false, fmt.Errorf("error getting session: %s", err)
	}

	svc := s3.New(sess)
	_, err = svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: &bucketName,
	})
	if err == nil {
		return true, nil
	} else if strings.HasPrefix(err.Error(), "NotFound:") {
		return false, nil
	} else {
		return false, fmt.Errorf("error while checking if remote state is available: %s", err)
	}
}
func (a *AWS) initRemoteStateBucket() error {
	svc, err := a.S3()
	if err != nil {
		return err
	}

	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(a.RemoteStateName()),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(a.Region()),
		},
	})
	if err != nil {
		return err
	}

	_, err = svc.PutBucketVersioning(&s3.PutBucketVersioningInput{
		Bucket: aws.String(a.RemoteStateName()),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: aws.String("Enabled"),
		},
	})
	return err
}

func (a *AWS) validateRemoteStateBucket() error {
	svc, err := a.S3()
	if err != nil {
		return err
	}

	_, err = svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(a.RemoteStateName()),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "NotFound" {
				return a.initRemoteStateBucket()
			}
		}
		return fmt.Errorf("error looking for terraform state bucket: %s", err)
	}

	location, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: aws.String(a.RemoteStateName()),
	})
	if err != nil {
		return err
	}
	if bucketRegion, myRegion := *location.LocationConstraint, a.Region(); bucketRegion != myRegion {
		return fmt.Errorf("bucket region is wrong, actual: %s expected: %s", bucketRegion, myRegion)
	}

	versioning, err := svc.GetBucketVersioning(&s3.GetBucketVersioningInput{
		Bucket: aws.String(a.RemoteStateName()),
	})
	if err != nil {
		return err
	}
	if *versioning.Status != "Enabled" {
		a.log.Warnf("state bucket %s has versioning disabled", a.RemoteStateName())
	}

	return nil

}

func (a *AWS) initRemoteStateDynamoDB() error {
	svc, err := a.DynamoDB()
	if err != nil {
		return err
	}

	_, err = svc.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(a.RemoteStateName()),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String(DynamoDBKey),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			&dynamodb.KeySchemaElement{
				AttributeName: aws.String(DynamoDBKey),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	})
	return err
}

func (a *AWS) validateRemoteStateDynamoDB() error {
	svc, err := a.DynamoDB()
	if err != nil {
		return err
	}

	describeOut, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(a.RemoteStateName()),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ResourceNotFoundException" {
				return a.initRemoteStateDynamoDB()
			}
		}
		return fmt.Errorf("error looking for terraform state dynamodb: %s", err)
	}

	attributeFound := false
	for _, params := range describeOut.Table.AttributeDefinitions {
		if *params.AttributeName == DynamoDBKey {
			attributeFound = true
		}
	}
	if !attributeFound {
		return fmt.Errorf("the DynamoDB table '%s' doesn't contain a parameter named '%s'", a.RemoteStateName(), DynamoDBKey)
	}

	return nil
}
