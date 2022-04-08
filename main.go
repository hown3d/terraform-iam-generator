package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hown3d/terraform-iam-generator/internal/aws"
	"github.com/hown3d/terraform-iam-generator/internal/terraform"
)

var directory *string = flag.String("dir", "", "terraform directory to use")

func main() {
	flag.Parse()
	start := time.Now()
	err := terraform.Apply(terraform.Options{
		Directory: *directory,
	})
	if err != nil {
		log.Println(err)
		return
	}
	terraform.Destroy(terraform.Options{
		Directory: *directory,
	})
	end := time.Now()

	sts, err := aws.NewStsService()
	if err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()

	userArn, err := sts.GetCallerIdentity(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	events, err := aws.NewEventService()
	if err != nil {
		log.Println(err)
		return
	}

	// cut arn to get only the name
	userName := userArn[strings.LastIndex(userArn, ",")+1:]

	collectedEvents, err := events.GetEventsOfUser(ctx, userName, &start, &end)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(collectedEvents)
}
