package recorder_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"

	"github.com/go-spectest/aws/mocks"
	"github.com/go-spectest/aws/recorder"
	"github.com/go-spectest/spectest"
)

type user struct {
	Name       string `json:"name"`
	Registered bool   `json:"registered"`
}

func TestDynamoDBRecorderPutItem(t *testing.T) {
	attr, _ := dynamodbattribute.MarshalMap(user{
		Name:       "Peter Ndlovu",
		Registered: true,
	})
	db := &mocks.DynamoDB{
		PutItemFunc: func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
			assert.Equal(t, *input.TableName, "dev_table")
			assert.Equal(t, *input.Item["name"].S, "Peter Ndlovu")
			assert.Equal(t, *input.Item["registered"].BOOL, true)
			return &dynamodb.PutItemOutput{}, nil
		},
	}
	testRecorder := spectest.NewTestRecorder()
	recordingDB := recorder.NewDynamoDB(db, testRecorder)

	_, err := recordingDB.PutItem(&dynamodb.PutItemInput{
		Item:      attr,
		TableName: aws.String("dev_table"),
	})

	assert.NoError(t, err)
	assert.Len(t, testRecorder.Events, 2)

	request := testRecorder.Events[0].(spectest.MessageRequest)
	assert.Equal(t, request.Source, spectest.SystemUnderTestDefaultName)
	assert.Equal(t, request.Target, "DynamoDB")

	response := testRecorder.Events[1].(spectest.MessageResponse)
	assert.Equal(t, response.Source, "DynamoDB")
	assert.Equal(t, response.Target, spectest.SystemUnderTestDefaultName)
}
