package awsclient

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/yukels/util/context"
)

type CloudWatchClient struct {
	client *cloudwatch.CloudWatch
}

func NewCloudWatchClient(ctx context.Context) (*CloudWatchClient, error) {
	if err := createSession(ctx); err != nil {
		return nil, err
	}

	return &CloudWatchClient{
		client: cloudwatch.New(awsSession),
	}, nil
}

func (c *CloudWatchClient) GetMetricLast(ctx context.Context, namespace, metric, dimensionName, dimensionValue string) (float64, error) {
	stats, err := c.client.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		MetricName: aws.String(metric),
		Namespace:  aws.String(namespace),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String(dimensionName),
				Value: aws.String(dimensionValue),
			},
		},
		Statistics: []*string{aws.String("Average")},
		StartTime:  aws.Time(time.Now().Add(-time.Duration(10 * time.Minute))),
		EndTime:    aws.Time(time.Now()),
		Period:     aws.Int64(60),
	})
	if err != nil {
		return 0, err
	}

	// find last dp
	var minTimestamp *time.Time
	var value float64
	for _, dp := range stats.Datapoints {
		if minTimestamp == nil || minTimestamp.Before(*dp.Timestamp) {
			value = *dp.Average
			minTimestamp = dp.Timestamp
		}
	}

	return value, nil
}
