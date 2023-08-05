package user_config

import (
	"buggybox/modules/planner"
	"buggybox/modules/process"
	"buggybox/modules/utils"
	"buggybox/modules/web_server"
	"fmt"
	"strings"
)

var Prepared PreparedConfigType

type PreparedConfigType struct {
	SchemaVersion string
	Process       *process.Process
	Plans         []*planner.Plan
	WebServers    []*web_server.WebServer
}

func (pc *PreparedConfigType) Start() {
	for _, plan := range pc.Plans {
		go plan.ExecuteAll()
	}
}

func (u *PreparedConfigType) preparePlannable(plannable planner.Plannable) error {
	desiredPlans := plannable.GetDesiredPlanNames()

	if plannable.HasCustomPlan() {
		planName := plannable.GetUid()

		customPlan := plannable.MakeCustomPlan()
		customPlan.Name = &planName
		customPlan.MakePrivate()

		u.Plans = append(u.Plans, customPlan)
		desiredPlans = append(desiredPlans, planName)
	}

	for _, planName := range desiredPlans {
		desiredPlan := u.findPlan(planName)

		if desiredPlan == nil {
			return fmt.Errorf("plan %s not found for sub-app %s", planName, plannable.GetUid())
		}

		desiredPlan.Assign(plannable)
	}

	return nil
}

func (u *PreparedConfigType) findPlan(name string) *planner.Plan {
	for _, plan := range u.Plans {
		if plan.Name != nil && *plan.Name == name {
			return plan
		}
	}

	return nil
}

func (pc *PreparedConfigType) Validate() error {
	dups := pc.findDuplicateApps()
	if len(dups) > 0 {
		return fmt.Errorf("there are duplicate sub-apps: %s", strings.Join(dups, ", "))
	}

	dups = pc.findDuplicatePlans()
	if len(dups) > 0 {
		return fmt.Errorf("there are duplicate plans: %s", strings.Join(dups, ", "))
	}

	for _, plan := range pc.Plans {
		err := plan.Validate()
		if err != nil {
			return fmt.Errorf("plan %s is invalid: %v", *plan.Name, err)
		}
	}

	if err := pc.Process.Validate(); err != nil {
		return fmt.Errorf("process manager is invalid: %v", err)
	}

	for _, webServer := range pc.WebServers {
		err := webServer.Validate()
		if err != nil {
			return fmt.Errorf("webserver %s is invalid: %v", webServer.GetUid(), err)
		}

		for _, route := range webServer.Routes {
			err := route.Validate()

			if err != nil {
				return fmt.Errorf("route %s is invalid for webserver %s: %v", route.Path, webServer.GetUid(), err)
			}
		}
	}

	return nil
}

func (u *PreparedConfigType) findDuplicateApps() []string {
	apps := []string{
		u.Process.GetUid(),
	}

	for _, v := range u.WebServers {
		apps = append(apps, v.GetUid())
	}

	return utils.GetDuplicates(apps)
}

func (u *PreparedConfigType) findDuplicatePlans() []string {
	plans := []string{}

	for _, v := range u.Plans {
		plans = append(plans, *v.Name)
	}

	return utils.GetDuplicates(plans)
}