package core

import (
	"fmt"
	"log"
	"sync"

	"github.com/hown3d/terraform-iam-generator/internal/aws"
	"github.com/hown3d/terraform-iam-generator/internal/metrics"
	"github.com/hown3d/terraform-iam-generator/internal/terraform"
)

func Run(dir string) {
	messageChan := make(chan metrics.CsmMessage)
	svc, err := metrics.NewServerAndListen(messageChan)
	if err != nil {
		log.Println(err)
		return
	}
	go svc.Read()

	var wg sync.WaitGroup
	// collect messages
	var msgs []metrics.CsmMessage
	go func() {
		wg.Add(1)
		for msg := range messageChan {
			msgs = append(msgs, msg)
		}
		wg.Done()
	}()

	err = terraform.Apply(terraform.Options{
		Directory: dir,
	})
	if err != nil {
		log.Println(err)
	}

	err = terraform.Destroy(terraform.Options{
		Directory: dir,
	})
	if err != nil {
		log.Println(err)
	}
	svc.Stop()
	wg.Wait()

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
