package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/urfave/cli/v2"
)

func NewPutMetricDataInput(namespace *string, metricName *string, clusterName *string, serviceName *string, value *float64, unit types.StandardUnit) *cloudwatch.PutMetricDataInput {
	return &cloudwatch.PutMetricDataInput{
		Namespace: namespace,
		MetricData: []types.MetricDatum{
			{
				MetricName: metricName,
				Dimensions: []types.Dimension{
					{
						Name:  aws.String("ClusterName"),
						Value: clusterName,
					},
					{
						Name:  aws.String("ServiceName"),
						Value: serviceName,
					},
				},
				Value: value,
				Unit:  unit,
			},
		},
	}
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.Float64Flag{
				Name:  "min",
				Value: 1.0,
				Usage: "Minimum value of the metric",
			},
			&cli.Float64Flag{
				Name:  "max",
				Value: 100.0,
				Usage: "Maximum value of the metric",
			},
			&cli.StringFlag{
				Name:  "metric",
				Value: "JvdbcCustomMetric2",
				Usage: "Name of the metric",
			},
			&cli.StringFlag{
				Name:  "service",
				Value: "pv-11-api-service",
				Usage: "Name of the ECS service",
			},
			&cli.StringFlag{
				Name:  "cluster",
				Value: "pv-dev-ecs-cluster",
				Usage: "Name of the ECS cluster",
			},
			&cli.StringFlag{
				Name:  "namespace",
				Value: "ECS/CloudWatch/Custom",
				Usage: "Cloudwatch namespace",
			},
		},
		Name:  "gmetric-aws",
		Usage: "Send a custom metric to AWS CloudWatch",
		Action: func(cliCtx *cli.Context) error {
			min := cliCtx.Float64("min")
			max := cliCtx.Float64("max")
			metricName := cliCtx.String("metric")
			serviceName := cliCtx.String("service")
			clusterName := cliCtx.String("cluster")
			namespace := cliCtx.String("namespace")

			if min >= max {
				return cli.Exit("Minimum value must be less than maximum value", 1)
			}

			bckCtx := context.Background()
			sdkConfig, err := config.LoadDefaultConfig(bckCtx)
			if err != nil {
				return cli.Exit(fmt.Sprintf("Couldn't load default configuration. Have you set up your AWS account? %v", err), 1)
			}

			// Create a CloudWatch client
			cwClient := cloudwatch.NewFromConfig(sdkConfig)

			// Trigger the code every second
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				value := min + rand.Float64()*(max-min)
				metric := NewPutMetricDataInput(&namespace, &metricName, &clusterName, &serviceName, &value, types.StandardUnitPercent)
				_, err := cwClient.PutMetricData(bckCtx, metric)
				if err != nil {
					log.Printf("Error sending metric data: %v", err)
					continue
				}

				log.Printf("Successfully sent metric data: %f", value)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
