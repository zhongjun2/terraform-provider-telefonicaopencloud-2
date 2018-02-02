package telefonicaopencloud

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/cloudeyeservice/alarmrule"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestCESAlarmRule_basic(t *testing.T) {
	var ar alarmrule.AlarmRule

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCESAlarmRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testCESAlarmRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testCESAlarmRuleExists("telefonicaopencloud_ces_alarmrule.alarmrule_1", &ar),
				),
			},
			resource.TestStep{
				Config: testCESAlarmRule_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"telefonicaopencloud_ces_alarmrule.alarmrule_1", "alarm_enabled", "false"),
				),
			},
		},
	})
}

func testCESAlarmRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.loadCESClient(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating TelefonicaOpenCloud ces client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "telefonicaopencloud_ces_alarmrule" {
			continue
		}

		id := rs.Primary.ID
		_, err := alarmrule.Get(networkingClient, id).Extract()
		if err == nil {
			return fmt.Errorf("Alarm rule still exists")
		}
	}

	return nil
}

func testCESAlarmRuleExists(n string, ar *alarmrule.AlarmRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.loadCESClient(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating TelefonicaOpenCloud ces client: %s", err)
		}

		id := rs.Primary.ID
		found, err := alarmrule.Get(networkingClient, id).Extract()
		if err != nil {
			return err
		}

		*ar = *found

		return nil
	}
}

var testCESAlarmRule_basic = fmt.Sprintf(`
resource "telefonicaopencloud_compute_instance_v2" "vm_1" {
  name = "instance_1"
  network {
    uuid = "%s"
  }
}

resource "telefonicaopencloud_ces_alarmrule" "alarmrule_1" {
  "alarm_name" = "alarm_rule1"

  "metric" {
    "namespace" = "SYS.ECS"
    "metric_name" = "network_outgoing_bytes_rate_inband"
    "dimensions" {
        "name" = "instance_id"
        "value" = "${telefonicaopencloud_compute_instance_v2.vm_1.id}"
    }
  }
  "condition"  {
    "period" = 300
    "filter" = "average"
    "comparison_operator" = ">"
    "value" = 6
    "unit" = "B/s"
    "count" = 1
  }
  "alarm_action_enabled" = false
}
`, OS_NETWORK_ID)

var testCESAlarmRule_update = fmt.Sprintf(`
resource "telefonicaopencloud_compute_instance_v2" "vm_1" {
  name = "instance_1"
  network {
    uuid = "%s"
  }
}

resource "telefonicaopencloud_ces_alarmrule" "alarmrule_1" {
  "alarm_name" = "alarm_rule1"

  "metric" {
    "namespace" = "SYS.ECS"
    "metric_name" = "network_outgoing_bytes_rate_inband"
    "dimensions" {
        "name" = "instance_id"
        "value" = "${telefonicaopencloud_compute_instance_v2.vm_1.id}"
    }
  }
  "condition"  {
    "period" = 300
    "filter" = "average"
    "comparison_operator" = ">"
    "value" = 6
    "unit" = "B/s"
    "count" = 1
  }
  "alarm_action_enabled" = false
  "alarm_enabled" = false
}
`, OS_NETWORK_ID)