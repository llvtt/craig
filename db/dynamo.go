package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDBAccess implements DataAccess for DynamoDB.
type DynamoDBAccess struct {
	TableName string
	Client    *dynamodb.DynamoDB
}

type DynamoDBAccessManager struct {
	client *dynamodb.DynamoDB
}

func NewDynamoDBAccessManager(client *dynamodb.DynamoDB) *DynamoDBAccessManager {
	return &DynamoDBAccessManager{client}
}

func (mgr *DynamoDBAccessManager) Table(tableName string) *DynamoDBAccess {
	return &DynamoDBAccess{tableName, mgr.client}
}

func (acc *DynamoDBAccess) List(ctx context.Context) (it Iterator, err error) {
	input := &dynamodb.ScanInput{TableName: aws.String(acc.TableName)}

	var docs []map[string]*dynamodb.AttributeValue
	err = acc.Client.ScanPagesWithContext(ctx, input, func(output *dynamodb.ScanOutput, lastPage bool) bool {
		docs = append(docs, output.Items...)

		return !lastPage
	})
	it = &DynamoAccessIterator{scannedItems: docs}

	return
}

func (acc *DynamoDBAccess) Upsert(ctx context.Context, record interface{}, previousRecord ...interface{}) error {
	if len(previousRecord) > 1 {
		return fmt.Errorf("up to one previousRecord may be provided, got %d", len(previousRecord))
	}

	item, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		return err
	}

	var returnValues *string
	if len(previousRecord) > 0 {
		returnValues = aws.String(dynamodb.ReturnValueAllOld)
	}

	output, err := acc.Client.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName:    aws.String(acc.TableName),
		ReturnValues: returnValues,
		Item:         item,
	})

	if err != nil || output.Attributes == nil {
		return err
	}

	if previousRecord != nil {
		if err = dynamodbattribute.UnmarshalMap(output.Attributes, previousRecord[0]); err != nil {
			return err
		}
	}

	return err
}

// DynamoAccessIterator implements Iterator for DynamoAccess.
type DynamoAccessIterator struct {
	scannedItems []map[string]*dynamodb.AttributeValue
	position     int
}

func (it *DynamoAccessIterator) Next(out interface{}) (err error) {
	if it.position >= len(it.scannedItems) && err == nil {
		err = IteratorExhausted
	}

	if it.position < len(it.scannedItems) {
		err = dynamodbattribute.UnmarshalMap(it.scannedItems[it.position], out)
		it.position++
	}

	return
}
