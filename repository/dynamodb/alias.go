package dynamodb

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsdynamodbattribute "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsdynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"rd/config"
	"rd/entity"
)

type DynamodbAliasRepository struct {
	cli              *awsdynamodb.Client
	aliasTableName   string
	aliasTablePkName string
	//aliasHitEventTableName   string
	//aliasHitEventTablePkName string
}

const (
	AliasTablePk = "group_alias"
)

func NewDynamodbAliasRepository(rdDynamodbConfig config.DynamodbConfig, awsConfig aws.Config) *DynamodbAliasRepository {
	cli := awsdynamodb.NewFromConfig(awsConfig)

	return &DynamodbAliasRepository{
		cli:              cli,
		aliasTableName:   rdDynamodbConfig.AliasTableName,
		aliasTablePkName: rdDynamodbConfig.AliasTablePkName,
	}
}

func (r *DynamodbAliasRepository) Create(alias *entity.Alias) (*entity.Alias, error) {
	ctx := context.TODO()
	deletedAt := alias.DeletedAt.Time.String()
	if !alias.DeletedAt.Valid {
		deletedAt = ""
	}
	_, err := r.cli.PutItem(ctx, &awsdynamodb.PutItemInput{
		Item: map[string]awsdynamodbtypes.AttributeValue{
			"group":       ddbString(alias.AliasGroup),
			"name":        ddbString(alias.Name),
			"destination": ddbString(alias.Destination),
			"created_at":  ddbString(alias.CreatedAt.String()),
			"updated_at":  ddbString(alias.UpdatedAt.String()),
			"deleted_at":  ddbString(deletedAt),
		},
		TableName:                aws.String(r.aliasTableName),
		ExpressionAttributeNames: map[string]string{"#group": "group", "#name": "name"},
		ConditionExpression:      aws.String(fmt.Sprintf("attribute_not_exists(#group) AND attribute_not_exists(#name)")),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return alias, nil
}

func (r *DynamodbAliasRepository) List() []*entity.Alias {
	panic("implement me")
	//ctx := context.TODO()
	//aliases := make([]*entity.Alias, 0, 32)
	//output, err := r.cli.Query(ctx, &awsdynamodb.QueryInput{
	//	TableName: aws.String(r.aliasTableName),
	//})
	//if err != nil {
	//	log.Error(err)
	//	return nil
	//}
	//
	//for _, item := range output.Items {
	//	alias := new(entity.Alias)
	//	group, name, err := entity.GetGroupAndNameFromDynamodbPk(item[r.aliasTablePkName])
	//}
	//return
}

// TODO(umi0410): This method using dynamodb cannot support sorting by recentHitCount
// Another solution such as Redis might help.
func (r *DynamodbAliasRepository) ListByGroup(group string, recentHitCountSince time.Time) []*entity.Alias {
	ctx := context.TODO()
	aliases := make([]*entity.Alias, 0, 32)
	output, err := r.cli.Query(ctx, &awsdynamodb.QueryInput{
		TableName:                aws.String(r.aliasTableName),
		KeyConditionExpression:   aws.String("#group = :group"),
		ExpressionAttributeNames: map[string]string{"#group": "group"},
		ExpressionAttributeValues: map[string]awsdynamodbtypes.AttributeValue{
			":group": ddbString(group),
		},
	})
	if err != nil {
		log.Error(err)
		return nil
	}

	for _, item := range output.Items {
		aliases = append(aliases, r.mapDynamoAttributeValueMapToAliasEntity(item))
	}

	return aliases
}

func (r *DynamodbAliasRepository) ListByGroupAndAlias(group, alias string) []*entity.Alias {
	ctx := context.TODO()
	aliases := make([]*entity.Alias, 0, 32)
	output, err := r.cli.Query(ctx, &awsdynamodb.QueryInput{
		TableName:                aws.String(r.aliasTableName),
		KeyConditionExpression:   aws.String("#group = :group AND #name = :name"),
		ExpressionAttributeNames: map[string]string{"#group": "group", "#name": "name"},
		ExpressionAttributeValues: map[string]awsdynamodbtypes.AttributeValue{
			":group": ddbString(group),
			":name":  ddbString(alias),
		},
	})
	if err != nil {
		log.Error(err)
		return nil
	}

	for _, item := range output.Items {
		aliases = append(aliases, r.mapDynamoAttributeValueMapToAliasEntity(item))
	}

	return aliases
}

func (r *DynamodbAliasRepository) mapDynamoAttributeValueMapToAliasEntity(m map[string]awsdynamodbtypes.AttributeValue) *entity.Alias {
	alias := new(entity.Alias)

	if err := awsdynamodbattribute.UnmarshalMap(m, alias); err != nil {
		log.Error(errors.WithStack(err))
	}

	return alias
}

func (r *DynamodbAliasRepository) Get(id int) (*entity.Alias, error) {
	panic("implement me")
}

func (r *DynamodbAliasRepository) Delete(id int) (*entity.Alias, error) {
	panic("implement me")
}

func (r *DynamodbAliasRepository) Close() error {
	panic("implement me")
}
