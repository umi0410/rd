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
		userTableName:   rdDynamodbConfig.UserTableName,
		userTablePkName: rdDynamodbConfig.UserTablePkName,
	}

	return repository
}

type DynamodbAuthRepository struct {
	cli             *awsdynamodb.Client
	userTableName   string
	userTablePkName string
}

func (r *DynamodbAuthRepository) GetUser(ctx context.Context, username string) (*entity.User, error) {
	user := new(entity.User)
	output, err := r.cli.GetItem(ctx, &awsdynamodb.GetItemInput{
		TableName: aws.String(r.userTableName),
		Key: map[string]types.AttributeValue{
			r.userTablePkName: ddbString(username),
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := awsdynamodbattribute.UnmarshalMap(output.Item, user); err != nil {
		return nil, errors.WithStack(err)
	}

	groupNames, err := stringsFromDdbAttribute(output.Item["groups"])
	if err != nil {
		return nil, err
	}

	for _, groupName := range groupNames {
		user.Groups = append(user.Groups, &entity.Group{Name: groupName})
	}

	return user, nil
}
