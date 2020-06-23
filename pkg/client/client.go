package client

import (
	"fmt"
	"strings"

	"github.com/abtasty/flagship-go-sdk/pkg/cache"
	"github.com/abtasty/flagship-go-sdk/pkg/model"

	"github.com/abtasty/flagship-go-sdk/pkg/bucketing"
	"github.com/abtasty/flagship-go-sdk/pkg/utils"

	"github.com/abtasty/flagship-go-sdk/pkg/decision"
	"github.com/abtasty/flagship-go-sdk/pkg/logging"
	"github.com/abtasty/flagship-go-sdk/pkg/tracking"
)

// DecisionMode represents the decision mode of the Client engine
type DecisionMode string

// The different decision modes
const (
	API       DecisionMode = "API"
	Bucketing DecisionMode = "Bucketing"
)

// Client represent the Flagship SDK client object
type Client struct {
	envID             string
	decisionMode      DecisionMode
	decisionClient    decision.ClientInterface
	trackingAPIClient tracking.APIClientInterface
	cacheManager      cache.Manager
}

var clientLogger = logging.CreateLogger("FS Client")

// Create creates a Client from options
func Create(f *Options) (*Client, error) {
	clientLogger.Info(fmt.Sprintf("Creating FS Client with Decision Mode : %s", f.decisionMode))
	client := &Client{
		envID: f.EnvID,
	}

	var err error

	if len(f.cacheManagerOptions) > 0 {
		cacheManager, err := cache.InitManager(f.cacheManagerOptions...)
		if err != nil {
			clientLogger.Error("Got error when initializing cache", err)
		}
		client.cacheManager = cacheManager
	}

	if client.trackingAPIClient == nil {
		client.trackingAPIClient, err = tracking.NewAPIClient(client.envID, f.decisionAPIOptions...)
	}

	if client.decisionClient == nil {
		if f.decisionMode == Bucketing {
			client.decisionClient, err = bucketing.NewEngine(client.envID, client.cacheManager, f.bucketingOptions...)
			if err != nil {
				clientLogger.Error("Got error when creating bucketing engine", err)
			}
		} else {
			client.decisionClient, err = decision.NewAPIClient(client.envID, f.decisionAPIOptions...)
			if err != nil {
				clientLogger.Error("Got error when creating Decision API engine", err)
			}
		}
	}

	return client, err
}

// NewVisitor returns a new Visitor from ID and context
func (c *Client) NewVisitor(visitorID string, context model.Context) (visitor *Visitor, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = utils.HandleRecovered(r, clientLogger)
		}
	}()

	clientLogger.Info(fmt.Sprintf("Creating new visitor with id : %s", visitorID))

	errs := context.Validate()
	if len(errs) > 0 {
		errorStrings := []string{}
		for _, e := range errs {
			clientLogger.Error("Context error", e)
			errorStrings = append(errorStrings, e.Error())
		}
		return nil, fmt.Errorf("Invalid context : %s", strings.Join(errorStrings, ", "))
	}

	return &Visitor{
		ID:                visitorID,
		Context:           context,
		decisionClient:    c.decisionClient,
		decisionMode:      c.decisionMode,
		trackingAPIClient: c.trackingAPIClient,
		cacheManager:      c.cacheManager,
	}, nil
}

// SendHit sends a tracking hit to the Data Collect API
func (c *Client) SendHit(visitorID string, hit model.HitInterface) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = utils.HandleRecovered(r, clientLogger)
		}
	}()

	clientLogger.Info(fmt.Sprintf("Sending hit for visitor with id : %s", visitorID))
	err = c.trackingAPIClient.SendHit(visitorID, hit)

	if err != nil {
		err = fmt.Errorf("Error when sending hit: %s", err.Error())
	}
	return err
}

// Dispose disposes the Client and close all connections
func (c *Client) Dispose() (err error) {
	return err
}

// GetEnvID returns the current set env id
func (c *Client) GetEnvID() string {
	return c.envID
}
