package grafana

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/machines"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/utils"
)

type RegisterAction struct{}

func (ra *RegisterAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (ra *RegisterAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	sysID := uuid.Nil
	userUUID := uuid.Nil

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.")

	connPayload, err := utils.Cast[models.ConnectionPayload](data)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	grafanaConn, err := utils.Cast[GrafanaConn](connPayload.MetaData)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	grafanaCred, err := utils.Cast[GrafanaCred](connPayload.CredentialSecret)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	grafanaClient := models.NewGrafanaClient()
	err = grafanaClient.Validate(ctx, grafanaConn.URL, grafanaCred.APIKey)
	if err != nil {
		return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": models.ErrGrafanaScan(err)}).Build(), nil
	}
	return machines.NoOp, nil, nil
}

func (ra *RegisterAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
