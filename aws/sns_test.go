package aws_test

import (
	"testing"

	"github.com/sabariramc/goserverbase/aws"
	"github.com/sabariramc/goserverbase/utils"
)

func TestSNSClient(t *testing.T) {
	arn := AWSTestConfig.AWS.SNS_ARN
	ctx := GetCorrelationContext()
	snsClient := aws.GetDefaultSNSClient(AWSTestLogger)
	message := utils.NewMessage("event", "sns.test")
	message.AddPayload("payment", &utils.Payload{
		"entity": map[string]interface{}{
			"id":     "pay_14341234",
			"amount": 123,
		},
	})
	message.AddPayload("bank", &utils.Payload{
		"entity": map[string]interface{}{
			"id":                "bank_fadsfas",
			"bankAccountNumber": "0000021312",
		},
	})
	message.AddPayload("customer", &utils.Payload{
		"entity": map[string]interface{}{
			"id": "cust_fasdfsa",
		},
	})
	err := snsClient.PublishWithContext(ctx, &arn, nil, message, nil)
	if err != nil {
		t.Fatal(err)
	}

}
