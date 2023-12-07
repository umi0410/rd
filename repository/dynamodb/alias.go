package dynamodb

import (
	"context"
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
			"PK":          ddbString(alias.GetDynamodbPartitionKey()),
			"SK":          ddbString(alias.GetDynamodbSortKey()),
			"destination": ddbString(alias.Destination),
			"created_at":  ddbString(alias.CreatedAt.String()),
			"updated_at":  ddbString(alias.UpdatedAt.String()),
			"deleted_at":  ddbString(deletedAt),
		},
		TableName: aws.String(r.aliasTableName),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return alias, nil
}

func (r *DynamodbAliasRepository) List() []*entity.Alias {
	log.Errorf("%+v", errors.New("Not supported method"))
	return nil
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
	a := &entity.Alias{AliasGroup: group}
	output, err := r.cli.Query(ctx, &awsdynamodb.QueryInput{
		TableName:                aws.String(r.aliasTableName),
		KeyConditionExpression:   aws.String("#PK = :PK"),
		ExpressionAttributeNames: map[string]string{"#PK": "PK"},
		ExpressionAttributeValues: map[string]awsdynamodbtypes.AttributeValue{
			":PK": ddbString(a.GetDynamodbPartitionKey()),
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
	a := &entity.Alias{AliasGroup: group, Name: alias}
	output, err := r.cli.Query(ctx, &awsdynamodb.QueryInput{
		TableName:                aws.String(r.aliasTableName),
		KeyConditionExpression:   aws.String("#PK = :PK AND #SK = :SK"),
		ExpressionAttributeNames: map[string]string{"#PK": "PK", "#SK": "SK"},
		ExpressionAttributeValues: map[string]awsdynamodbtypes.AttributeValue{
			":PK": ddbString(a.GetDynamodbPartitionKey()),
			":SK": ddbString(a.GetDynamodbSortKey()),
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
	alias.AliasGroup = entity.GetAliasGroupFrom(m["PK"].(*awsdynamodbtypes.AttributeValueMemberS).Value)
	alias.Name = entity.GetAliasName(m["SK"].(*awsdynamodbtypes.AttributeValueMemberS).Value)
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
