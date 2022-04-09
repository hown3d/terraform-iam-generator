package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hown3d/terraform-iam-generator/internal/aws"
	"github.com/hown3d/terraform-iam-generator/internal/metrics"
	"github.com/hown3d/terraform-iam-generator/internal/terraform"
)

var directory *string = flag.String("dir", "", "terraform directory to use")

func main() {
	flag.Parse()
	messageChan := make(chan metrics.CsmMessage)
	svc, err := metrics.NewServerAndListen(messageChan)
	defer svc.Stop()
	if err != nil {
		log.Println(err)
		return
	}
	go svc.Read()

	// collect messages
	var msgs []metrics.CsmMessage
	go func() {
		for {
			msg := <-messageChan
			msgs = append(msgs, msg)
		}
	}()

	err = terraform.Apply(terraform.Options{
		Directory: *directory,
	})
	if err != nil {
		log.Println(err)
	}

	err = terraform.Destroy(terraform.Options{
		Directory: *directory,
	})
	if err != nil {
		log.Println(err)
	}

	policy, err := aws.GenerateIamPolicy(msgs)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Result:")
	fmt.Println("-----------------------------------------------")
	fmt.Println("Your terraform code needs the following iam policy:")
	fmt.Println(string(policy))
}
