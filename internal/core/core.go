package core

import (
	"fmt"
	"log"
	"sync"

	"github.com/hown3d/terraform-iam-generator/internal/aws"
	"github.com/hown3d/terraform-iam-generator/internal/metrics"
	"github.com/hown3d/terraform-iam-generator/internal/terraform"
)

func Run(dir string, tfVars []terraform.Variable, tfVarFiles []terraform.VariableFile) {
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

	tfOpts := terraform.Options{
		Directory: dir,
		VarsFiles: tfVarFiles,
		Vars:      tfVars,
	}
	err = terraform.Apply(tfOpts)
	if err != nil {
		log.Println(err)
	}

	err = terraform.Destroy(tfOpts)
	if err != nil {
		log.Println(err)
		return
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
