package main

import (
	"log"
	"time"

	"github.com/abtasty/flagship-go-sdk"
	"github.com/abtasty/flagship-go-sdk/pkg/bucketing"
	"github.com/abtasty/flagship-go-sdk/pkg/client"
)

var testEnvId = "blvo2kijq6pg023l8edg"
var modifKey = "testCache"
var modifDefaultValue = "default"

func main() {

	fsClient, err := flagship.Start(testEnvId)

	if err != nil {
		log.Printf("Flagship client error %v", err)
	}

	testGetValue(fsClient)

	fsClient, err = flagship.Start(testEnvId, client.WithBucketing(bucketing.PollingInterval(2*time.Second)))

	if err != nil {
		log.Printf("Flagship client error %v", err)
	}

	testGetValue(fsClient)
}

func testGetValue(client *client.Client) {
	visitor, err := client.NewVisitor("test1", nil)
	if err != nil {
		log.Printf("Flagship client error %v", err)
	}

	err = visitor.SynchronizeModifications()

	if err != nil {
		log.Printf("Flagship visitor error %v", err)
	}

	modifValue, err := visitor.GetModificationString(modifKey, modifDefaultValue, true)

	if err != nil {
		log.Printf("Flagship modification error %v", err)
	}

	log.Printf("Got modification %v with value : %v", modifKey, modifValue)
}
