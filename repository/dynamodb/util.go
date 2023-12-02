package dynamodb

import (
	"fmt"

	awsdynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

func ddbString(input string) *awsdynamodbtypes.AttributeValueMemberS {
	return &awsdynamodbtypes.AttributeValueMemberS{
		Value: input,
	}
}

func ddbStrings(input []string) *awsdynamodbtypes.AttributeValueMemberSS {
	return &awsdynamodbtypes.AttributeValueMemberSS{
		Value: input,
	}
}

func stringFromDdbAttribute(attr awsdynamodbtypes.AttributeValue) (string, error) {
	val, ok := attr.(*awsdynamodbtypes.AttributeValueMemberS)
	if !ok {
		return "", errors.New(fmt.Sprintf("%s is not types.AttributeValueMemberS", attr))
	}

	return val.Value, nil
}

func stringsFromDdbAttribute(attr awsdynamodbtypes.AttributeValue) ([]string, error) {
	val, ok := attr.(*awsdynamodbtypes.AttributeValueMemberSS)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s is not types.AttributeValueMemberSS", attr))
	}

	return val.Value, nil
}
