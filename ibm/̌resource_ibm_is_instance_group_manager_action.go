// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIBMISInstanceGroupManagerAction() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISInstanceGroupManagerActionCreate,
		Read:     resourceIBMISInstanceGroupManagerActionRead,
		Update:   resourceIBMISInstanceGroupManagerActionUpdate,
		Delete:   resourceIBMISInstanceGroupManagerActionDelete,
		Exists:   resourceIBMISInstanceGroupManagerActionExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_is_instance_group_manager", "name"),
				Description:  "instance group manager action name",
			},

			"instance_group_manager_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance group manager action ID",
			},

			"instance_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance group ID",
			},

			"instance_group_manager_scheduled": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance group manager ID of type scheduled",
			},

			"run_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The date and time the scheduled action will run.",
			},

			"membership_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of members the instance group should have at the scheduled time.",
			},

			"max_membership_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_is_instance_group_manager_action", "max_membership_count"),
				Description:  "The maximum number of members in a managed instance group",
			},

			"min_membership_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: InvokeValidator("ibm_is_instance_group_manager_action", "min_membership_count"),
				Description:  "The minimum number of members in a managed instance group",
			},

			"instance_group_manager_autoscale": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier for this instance group manager of type autoscale.",
			},

			"instance_group_manager_autoscale_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance group manager name of type autoscale.",
			},

			"cron_spec": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cron specification for a recurring scheduled action. Actions can be applied a maximum of one time within a 5 min period.",
			},

			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the instance group action- `active`: Action is ready to be run- `completed`: Action was completed successfully- `failed`: Action could not be completed successfully- `incompatible`: Action parameters are not compatible with the group or manager- `omitted`: Action was not applied because this action's manager was disabled.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the instance group manager action was modified.",
			},
			"action_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of action for the instance group.",
			},

			"last_applied_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time the scheduled action was last applied. If empty the action has never been applied.",
			},
			"next_run_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time the scheduled action will next run. If empty the system is currently calculating the next run time.",
			},
		},
	}
}

func resourceIBMISInstanceGroupManagerActionValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "max_membership_count",
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			MinValue:                   "1",
			MaxValue:                   "1000"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "min_membership_count",
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			MinValue:                   "1",
			MaxValue:                   "1000"})

	ibmISInstanceGroupManagerResourceValidator := ResourceValidator{ResourceName: "ibm_is_instance_group_manager_action", Schema: validateSchema}
	return &ibmISInstanceGroupManagerResourceValidator
}

func resourceIBMISInstanceGroupManagerActionCreate(d *schema.ResourceData, meta interface{}) error {
	// CreateInstanceGroupManagerAction
	instanceGroupID := d.Get("instance_group").(string)
	instancegroupmanagerscheduledID := d.Get("instance_group_manager_scheduled").(string)

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	instanceGroupManagerActionOptions := vpcv1.CreateInstanceGroupManagerActionOptions{}
	instanceGroupManagerActionOptions.InstanceGroupID = &instanceGroupID
	instanceGroupManagerActionOptions.InstanceGroupManagerID = &instancegroupmanagerscheduledID

	instanceGroupManagerActionPrototype := vpcv1.InstanceGroupManagerActionPrototype{}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		instanceGroupManagerActionPrototype.Name = &name
	}

	if v, ok := d.GetOk("run_at"); ok {
		runat := v.(string)
		instanceGroupManagerActionPrototype.RunAt = &runat
	}

	if v, ok := d.GetOk("cron_spec"); ok {
		runat := v.(string)
		instanceGroupManagerActionPrototype.RunAt = &runat
	}

	if v, ok := d.GetOk("membership_count"); ok {
		membershipCount := int64(v.(int))
		instanceGroupManagerScheduledActionGroupPrototype := vpcv1.InstanceGroupManagerScheduledActionGroupPrototype{}
		instanceGroupManagerScheduledActionGroupPrototype.MembershipCount = &membershipCount
		instanceGroupManagerActionPrototype.Group = &instanceGroupManagerScheduledActionGroupPrototype
	}

	instanceGroupManagerScheduledActionByManagerManager := vpcv1.InstanceGroupManagerScheduledActionByManagerManager{}
	if v, ok := d.GetOk("min_membership_count"); ok {
		minmembershipCount := int64(v.(int))
		instanceGroupManagerScheduledActionByManagerManager.MinMembershipCount = &minmembershipCount
	}

	if v, ok := d.GetOk("max_membership_count"); ok {
		maxmembershipCount := int64(v.(int))
		instanceGroupManagerScheduledActionByManagerManager.MaxMembershipCount = &maxmembershipCount
	}

	if v, ok := d.GetOk("instance_group_manager_autoscale"); ok {
		instanceGroupManagerAutoScale := v.(string)
		instanceGroupManagerScheduledActionByManagerManager.ID = &instanceGroupManagerAutoScale
		instanceGroupManagerActionPrototype.Manager = &instanceGroupManagerScheduledActionByManagerManager
	}

	instanceGroupManagerActionOptions.InstanceGroupManagerActionPrototype = &instanceGroupManagerActionPrototype

	instanceGroupManagerActionIntf, response, err := sess.CreateInstanceGroupManagerAction(&instanceGroupManagerActionOptions)
	instanceGroupManagerAction := instanceGroupManagerActionIntf.(*vpcv1.InstanceGroupManagerAction)
	if err != nil || instanceGroupManagerAction == nil {
		return fmt.Errorf("Error creating InstanceGroup manager Action: %s\n%s", err, response)
	}
	d.SetId(fmt.Sprintf("%s/%s%s", instanceGroupID, instancegroupmanagerscheduledID, *instanceGroupManagerAction.ID))

	return resourceIBMISInstanceGroupManagerActionRead(d, meta)

}

func resourceIBMISInstanceGroupManagerActionUpdate(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	var changed bool
	instanceGroupManagerActionPatchModel := &vpcv1.InstanceGroupManagerActionPatch{}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		instanceGroupManagerActionPatchModel.Name = &name
		changed = true
	}

	if d.HasChange("cron_spec") {
		cronspec := d.Get("cron_spec").(string)
		instanceGroupManagerActionPatchModel.CronSpec = &cronspec
		changed = true
	}

	if d.HasChange("run_at") {
		runat := d.Get("run_at").(string)
		instanceGroupManagerActionPatchModel.RunAt = &runat
		changed = true
	}

	// if d.HasChange("max_membership_count") {
	// 	maxMembershipCount := int64(d.Get("max_membership_count").(int))
	// 	instanceGroupManagerPatchModel.MaxMembershipCount = &maxMembershipCount
	// 	changed = true
	// }

	// if d.HasChange("min_membership_count") {
	// 	minMembershipCount := int64(d.Get("min_membership_count").(int))
	// 	instanceGroupManagerPatchModel.MinMembershipCount = &minMembershipCount
	// 	changed = true
	// }

	// if d.HasChange("enable_manager") {
	// 	enableManager := d.Get("enable_manager").(bool)
	// 	instanceGroupManagerPatchModel.ManagementEnabled = &enableManager
	// 	changed = true
	// }

	if changed {

		parts, err := idParts(d.Id())
		if err != nil {
			return err
		}

		instanceGroupID := parts[0]
		instancegroupmanagerscheduledID := parts[1]
		instanceGroupManagerActionID := parts[2]

		updateInstanceGroupManagerActionOptions := &vpcv1.UpdateInstanceGroupManagerActionOptions{}
		updateInstanceGroupManagerActionOptions.InstanceGroupID = &instanceGroupID
		updateInstanceGroupManagerActionOptions.InstanceGroupManagerID = &instancegroupmanagerscheduledID
		updateInstanceGroupManagerActionOptions.ID = &instanceGroupManagerActionID

		instanceGroupManagerActionPatch, err := instanceGroupManagerActionPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for instanceGroupManagerActionPatch: %s", err)
		}
		updateInstanceGroupManagerActionOptions.InstanceGroupManagerActionPatch = instanceGroupManagerActionPatch
		_, response, err := sess.UpdateInstanceGroupManagerAction(updateInstanceGroupManagerActionOptions)
		if err != nil {
			return fmt.Errorf("Error updating InstanceGroup manager action: %s\n%s", err, response)
		}
	}
	return resourceIBMISInstanceGroupManagerRead(d, meta)
}

func resourceIBMISInstanceGroupManagerActionRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	instanceGroupID := parts[0]
	instancegroupmanagerscheduledID := parts[1]
	instanceGroupManagerActionID := parts[2]

	getInstanceGroupManagerActionOptions := &vpcv1.GetInstanceGroupManagerActionOptions{
		InstanceGroupID:        &instanceGroupID,
		InstanceGroupManagerID: &instancegroupmanagerscheduledID,
		ID:                     &instanceGroupManagerActionID,
	}

	instanceGroupManagerActionIntf, response, err := sess.GetInstanceGroupManagerAction(getInstanceGroupManagerActionOptions)
	instanceGroupManagerAction := instanceGroupManagerActionIntf.(*vpcv1.InstanceGroupManagerAction)
	if err != nil || instanceGroupManagerAction == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting InstanceGroup Manager Action: %s\n%s", err, response)
	}

	if err = d.Set("auto_delete", instanceGroupManagerAction.AutoDelete); err != nil {
		return fmt.Errorf("Error setting auto_delete: %s", err)
	}

	if err = d.Set("auto_delete_timeout", intValue(instanceGroupManagerAction.AutoDeleteTimeout)); err != nil {
		return fmt.Errorf("Error setting auto_delete_timeout: %s", err)
	}
	if err = d.Set("created_at", instanceGroupManagerAction.CreatedAt.String()); err != nil {
		return fmt.Errorf("Error setting created_at: %s", err)
	}

	if err = d.Set("instance_group_manager_action", instanceGroupManagerActionID); err != nil {
		return fmt.Errorf("Error setting instance_group_manager_action : %s", err)
	}

	if err = d.Set("name", instanceGroupManagerAction.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err = d.Set("resource_type", instanceGroupManagerAction.ResourceType); err != nil {
		return fmt.Errorf("Error setting resource_type: %s", err)
	}
	if err = d.Set("status", instanceGroupManagerAction.Status); err != nil {
		return fmt.Errorf("Error setting status: %s", err)
	}
	if err = d.Set("updated_at", instanceGroupManagerAction.UpdatedAt.String()); err != nil {
		return fmt.Errorf("Error setting updated_at: %s", err)
	}
	if err = d.Set("action_type", instanceGroupManagerAction.ActionType); err != nil {
		return fmt.Errorf("Error setting action_type: %s", err)
	}

	if err = d.Set("cron_spec", instanceGroupManagerAction.CronSpec); err != nil {
		return fmt.Errorf("Error setting cron_spec: %s", err)
	}
	if err = d.Set("last_applied_at", instanceGroupManagerAction.LastAppliedAt.String()); err != nil {
		return fmt.Errorf("Error setting last_applied_at: %s", err)
	}
	if err = d.Set("next_run_at", instanceGroupManagerAction.NextRunAt.String()); err != nil {
		return fmt.Errorf("Error setting next_run_at: %s", err)
	}

	instanceGroupManagerScheduledActionGroupGroup := instanceGroupManagerAction.Group
	if instanceGroupManagerScheduledActionGroupGroup != nil && instanceGroupManagerScheduledActionGroupGroup.MembershipCount != nil {
		d.Set("membership_count", intValue(instanceGroupManagerScheduledActionGroupGroup.MembershipCount))
	}
	instanceGroupManagerScheduledActionManagerManagerInt := instanceGroupManagerAction.Manager
	if instanceGroupManagerScheduledActionManagerManagerInt != nil {
		instanceGroupManagerScheduledActionManagerManager := instanceGroupManagerScheduledActionManagerManagerInt.(*vpcv1.InstanceGroupManagerScheduledActionManagerManager)
		if instanceGroupManagerScheduledActionManagerManager != nil && instanceGroupManagerScheduledActionManagerManager.ID != nil {

			if instanceGroupManagerScheduledActionManagerManager.MaxMembershipCount != nil {
				d.Set("max_membership_count", intValue(instanceGroupManagerScheduledActionManagerManager.MaxMembershipCount))
			}
			d.Set("min_membership_count", intValue(instanceGroupManagerScheduledActionManagerManager.MinMembershipCount))
			d.Set("instance_group_manager_autoscale_name", instanceGroupManagerScheduledActionManagerManager.Name)
			d.Set("instance_group_manager_autoscale", instanceGroupManagerScheduledActionManagerManager.ID)
		}
	}

	return nil
}

func resourceIBMISInstanceGroupManagerActionDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	instanceGroupID := parts[0]
	instancegroupmanagerscheduledID := parts[1]
	instanceGroupManagerActionID := parts[2]

	deleteInstanceGroupManagerActionOptions := &vpcv1.DeleteInstanceGroupManagerActionOptions{}
	deleteInstanceGroupManagerActionOptions.InstanceGroupID = &instanceGroupID
	deleteInstanceGroupManagerActionOptions.InstanceGroupManagerID = &instancegroupmanagerscheduledID
	deleteInstanceGroupManagerActionOptions.ID = &instanceGroupManagerActionID

	response, err := sess.DeleteInstanceGroupManagerAction(deleteInstanceGroupManagerActionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Deleting the InstanceGroup Manager Action: %s\n%s", err, response)
	}
	return nil
}

func resourceIBMISInstanceGroupManagerActionExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	instanceGroupID := parts[0]
	instancegroupmanagerscheduledID := parts[1]
	instanceGroupManagerActionID := parts[2]

	getInstanceGroupManagerActionOptions := &vpcv1.GetInstanceGroupManagerActionOptions{
		InstanceGroupID:        &instanceGroupID,
		InstanceGroupManagerID: &instancegroupmanagerscheduledID,
		ID:                     &instanceGroupManagerActionID,
	}

	_, response, err := sess.GetInstanceGroupManagerAction(getInstanceGroupManagerActionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error Getting InstanceGroup Manager Action: %s\n%s", err, response)
	}

	return true, nil
}
