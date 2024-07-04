package tfid

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildAutomationJobSchedule(b *client.ClientBuilder, id armid.ResourceId, _ string) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewAutomationJobScheduleClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.JobSchedule.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	runBook := props.Runbook
	if runBook == nil {
		return "", fmt.Errorf("unexpected nil properties.runBook in response")
	}
	runBookName := runBook.Name
	if runBookName == nil {
		return "", fmt.Errorf("unexpected nil properties.runBook.name in response")
	}
	schedule := props.Schedule
	if schedule == nil {
		return "", fmt.Errorf("unexpected nil properties.schedule in response")
	}
	scheduleName := schedule.Name
	if scheduleName == nil {
		return "", fmt.Errorf("unexpected nil properties.schedule.name in response")
	}

	scheduleId := id.Parent().Clone().(*armid.ScopedResourceId)
	scheduleId.AttrTypes = append(scheduleId.AttrTypes, "schedules")
	scheduleId.AttrNames = append(scheduleId.AttrNames, *scheduleName)

	runBookId := id.Parent().Clone().(*armid.ScopedResourceId)
	runBookId.AttrTypes = append(runBookId.AttrTypes, "runbooks")
	runBookId.AttrNames = append(runBookId.AttrNames, *runBookName)

	return scheduleId.String() + "|" + runBookId.String(), nil
}
