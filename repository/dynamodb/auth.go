package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsdynamodbattribute "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"rd/config"
	"rd/entity"
)

func NewDynamodbAuthRepository(rdDynamodbConfig config.DynamodbConfig, awsConfig aws.Config) *DynamodbAuthRepository {
	cli := awsdynamodb.NewFromConfig(awsConfig)
	repository := &DynamodbAuthRepository{
		cli:             cli,
		tableName:       rdDynamodbConfig.UserTableName,
		userInfoSortKey: "userInfo",
	}

	return repository
}

type DynamodbAuthRepository struct {
	cli             *awsdynamodb.Client
	tableName       string
	userInfoSortKey string
}

func (r *DynamodbAuthRepository) GetUser(ctx context.Context, username string) (*entity.User, error) {
	user := &entity.User{Username: username}
	output, err := r.cli.GetItem(ctx, &awsdynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": ddbString(user.GetDynamodbPartitionKey()),
			"SK": ddbString(r.userInfoSortKey),
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := awsdynamodbattribute.UnmarshalMap(output.Item, user); err != nil {
		return nil, errors.WithStack(err)
	}
	pk, err := stringFromDdbAttribute(output.Item["PK"])
	if err != nil {
		return nil, err
	}
	user.Username = entity.GetUserNameFrom(pk)

	groupNames, err := stringsFromDdbAttribute(output.Item["groups"])
	if err != nil {
		return nil, err
	}

	for _, groupName := range groupNames {
		user.Groups = append(user.Groups, &entity.Group{Name: groupName})
	}

	return user, nil
}
