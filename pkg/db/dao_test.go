package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type dynamodbContainer struct {
	testcontainers.Container
	URL string
}

func setupDynamoDBContainer(ctx context.Context) (*dynamodbContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "amazon/dynamodb-local",
		ExposedPorts: []string{"8000/tcp"},
		Cmd:          []string{"-jar DynamoDBLocal.jar -inMemory -sharedDb"},
		WaitingFor:   wait.ForListeningPort("8000/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "8000")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &dynamodbContainer{Container: container, URL: uri}, nil
}

// TODO fix this
func _TestGetWebsites(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	dynamoContainer, err := setupDynamoDBContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	defer dynamoContainer.Terminate(ctx)

	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{
		Endpoint: aws.String(dynamoContainer.URL),
	})
	dao := NewProbesDao(db)

	result := dao.GetWebsites()

	assert.NotNil(t, result)
}
