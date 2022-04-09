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
	start := time.Now().UTC()
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
	end := time.Now().UTC()

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
	events, err := aws.NewCloudtrailService()
	if err != nil {
		log.Println(err)
		return
	}

	// cut arn to get only the name
	userName := userArn[strings.LastIndex(userArn, "/")+1:]
	log.Printf("Using %s as userName", userName)

	log.Printf("Lookup events between %s and %s", start, end)
	iamActions, err := events.GetIamActions(ctx, userName, &start, &end)
	if err != nil {
		log.Println(err)
		return
	}

	for _, e := range iamActions {
		fmt.Printf("APICall:%s\n", e.APICall)
		fmt.Printf("Resources:%v\n", e.Resources)
		fmt.Printf("Service:%v\n", e.Service)
		fmt.Println("-------")
	}
}