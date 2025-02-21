package aws_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/sabariramc/goserverbase/aws"
)

func TestS3(t *testing.T) {
	ctx := GetCorrelationContext()
	s3Client := aws.GetDefaultS3Client(AWSTestLogger)
	path := fmt.Sprintf("dev/temp/goserverbasetest/%v.pdf", uuid.NewString())
	s3Bucket := AWSTestConfig.AWS.S3_BUCKET
	err := s3Client.PutFile(ctx, s3Bucket, path, "./testdata/sample_aadhaar.pdf")
	if err != nil {
		t.Fatal(err)
	}
	err = s3Client.GetFile(ctx, s3Bucket, path, "./testdata/test.pdf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s3Client.CreatePresignedUrlGET(ctx, s3Bucket, path, 10*60)
	if err != nil {
		t.Fatal(err)
	}
	path = fmt.Sprintf("dev/temp/goserverbasetest/%v.pdf", uuid.NewString())
	_, err = s3Client.CreatePresignedUrlPUT(ctx, s3Bucket, path, 10)
	if err != nil {
		t.Fatal(err)
	}
}

func TestS3PII(t *testing.T) {
	ctx := GetCorrelationContext()
	keyArn := AWSTestConfig.AWS.KMS_ARN
	s3Client, err := aws.GetDefaultS3PIIClient(AWSTestLogger, keyArn)
	if err != nil {
		t.Fatal(err)
	}
	path := fmt.Sprintf("dev/temp/goserverbasetest/%v.pdf", uuid.NewString())
	s3Bucker := AWSTestConfig.AWS.S3_BUCKET
	err = s3Client.PutFile(ctx, s3Bucker, path, "./testdata/sample_aadhaar.pdf")
	if err != nil {
		t.Fatal(err)
	}
	err = s3Client.GetFile(ctx, s3Bucker, path, "./testdata/testpii.pdf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s3Client.GetFileCache(ctx, s3Bucker, path, "dev", "testCache")
	if err != nil {
		t.Fatal(err)
	}
}
