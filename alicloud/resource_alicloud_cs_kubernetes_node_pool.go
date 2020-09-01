package alicloud

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCSKubernetesNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSKubernetesNodePoolCreate,
		Read:   resourceAlicloudCSNodePoolRead,
		Update: resourceAlicloudCSNodePoolUpdate,
		Delete: resourceAlicloudCSNodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems: 1,
			},
			"instance_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 10,
			},
			"password": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ConflictsWith:    []string{"key_name", "kms_encrypted_password"},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"key_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"password", "kms_encrypted_password"},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"kms_encrypted_password": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "key_name"},
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"system_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
				ValidateFunc: validation.StringInSlice([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}, false),
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ValidateFunc: validation.IntBetween(20, 32768),
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"data_disks": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"encrypted": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"tags": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"labels": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"taints": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"node_name_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^customized,[a-z0-9]([-a-z0-9\.])*,([5-9]|[1][0-2]),([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`), "Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test."),
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudCSKubernetesNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	var requestInfo *cs.Client
	var raw interface{}

	// prepare args and set default value
	args, err := buildNodePoolArgs(d, meta)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", "PrepareKubernetesNodePoolArgs", err)
	}

	if err = invoker.Run(func() error {
		raw, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.CreateNodePool(args, d.Get("cluster_id").(string), )
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", "CreateKubernetesNodePool", raw)
	}

	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Params"] = args
		addDebug("CreateKubernetesNodePool", raw, requestInfo, requestMap)
	}

	nodePool, ok := raw.(*cs.CreateNodePoolResponse)
	if ok != true {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", "ParseKubernetesNodePoolResponse", raw)
	}

	d.SetId(nodePool.NodePoolID)

	// reset interval to 10s
	stateConf := BuildStateConf([]string{"initial"}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Get("cluster_id").(string), d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCSNodePoolUpdate(d, meta)
}

func resourceAlicloudCSNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	vpcService := VpcService{client}
	d.Partial(true)
	update := false
	invoker := NewInvoker()
	args := &cs.UpdateNodePoolRequest{
		RegionId:         common.Region(client.RegionId),
		NodePoolInfo:     cs.NodePoolInfo{},
		ScalingGroup:     cs.ScalingGroup{},
		KubernetesConfig: cs.KubernetesConfig{},
		AutoScaling:      cs.AutoScaling{},
	}
	if d.HasChange("node_count") && !d.IsNewResource() {
		update = true
		oldV, newV := d.GetChange("node_count")

		oldValue, ok := oldV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("node_count old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := newV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("node_count new value can not be parsed"), "parseError %d", newValue)
		}

		if newValue < oldValue {
			return WrapErrorf(fmt.Errorf("node_count can not be less than before"), "scaleOutFailed %d:%d", newValue, oldValue)
		}
		args.Count = int64(newValue) - int64(oldValue)
		d.SetPartial("node_count")
	}
	args.NodePoolInfo.Name = d.Get("name").(string)
	if d.HasChange("name") && !d.IsNewResource() {
		update = true
		args.NodePoolInfo.Name = d.Get("name").(string)
		d.SetPartial("name")
	}
	if d.HasChange("vswitch_ids") && !d.IsNewResource() {
		update = true
		var vswitchID string
		if list := expandStringList(d.Get("vswitch_ids").([]interface{})); len(list) > 0 {
			vswitchID = list[0]
		} else {
			vswitchID = ""
		}

		var vpcId string
		if vswitchID != "" {
			vsw, err := vpcService.DescribeVSwitch(vswitchID)
			if err != nil {
				return err
			}
			vpcId = vsw.VpcId
		}
		args.ScalingGroup.VpcId = vpcId
		args.ScalingGroup.VswitchIds = expandStringList(d.Get("vswitch_ids").([]interface{}))
		d.SetPartial("vpc_id")
		d.SetPartial("vswitch_ids")
	}
	if d.HasChange("instance_types") && !d.IsNewResource() {
		update = true
		args.ScalingGroup.InstanceTypes = expandStringList(d.Get("instance_types").([]interface{}))
		d.SetPartial("instance_types")
	}

	// password is required by update method
	args.ScalingGroup.LoginPassword = d.Get("password").(string)
	if d.HasChange("password") && !d.IsNewResource() {
		update = true
		args.ScalingGroup.LoginPassword = d.Get("password").(string)
		d.SetPartial("password")
	}

	args.ScalingGroup.KeyPair = d.Get("key_name").(string)
	if d.HasChange("key_name") && !d.IsNewResource() {
		update = true
		args.ScalingGroup.KeyPair = d.Get("key_name").(string)
		d.SetPartial("key_name")
	}

	if d.HasChange("security_group_id") && !d.IsNewResource() {
		update = true
		args.ScalingGroup.SecurityGroupId = d.Get("security_group_id").(string)
		d.SetPartial("security_group_id")
	}

	if d.HasChange("system_disk_category") && !d.IsNewResource() {
		update = true
		args.ScalingGroup.SystemDiskCategory = d.Get("system_disk_category").(string)
		d.SetPartial("system_disk_category")
	}

	if d.HasChange("system_disk_size") && !d.IsNewResource() {
		update = true
		args.ScalingGroup.SystemDiskSize = int64(d.Get("system_disk_size").(int))
		d.SetPartial("system_disk_size")
	}

	if d.HasChange("image_id") && !d.IsNewResource() {
		update = true
		args.ScalingGroup.ImageId = d.Get("image_id").(string)
		d.SetPartial("image_id")
	}

	if d.HasChange("data_disks") && !d.IsNewResource() {
		update = true
		setNodePoolDataDisks(&args.ScalingGroup, d)
		d.SetPartial("data_disks")
	}

	if d.HasChange("tags") && !d.IsNewResource() {
		update = true
		setNodePoolTags(&args.ScalingGroup, d)
		d.SetPartial("tags")
	}

	if d.HasChange("labels") && !d.IsNewResource() {
		update = true
		setNodePoolLabels(&args.KubernetesConfig, d)
		d.SetPartial("labels")
	}

	if d.HasChange("taints") && !d.IsNewResource() {
		update = true
		setNodePoolTaints(&args.KubernetesConfig, d)
		d.SetPartial("taints")
	}

	if d.HasChange("node_name_mode") && !d.IsNewResource() {
		update = true
		args.KubernetesConfig.NodeNameMode = d.Get("node_name_mode").(string)
		d.SetPartial("node_name_mode")
	}

	if d.HasChange("user_data") && !d.IsNewResource() {
		update = true
		if v := d.Get("user_data").(string); v != "" {
			_, base64DecodeError := base64.StdEncoding.DecodeString(v)
			if base64DecodeError == nil {
				args.KubernetesConfig.UserData = v
			} else {
				args.KubernetesConfig.UserData = base64.StdEncoding.EncodeToString([]byte(v))
			}
			d.SetPartial("user_data")
		}
	}

	if d.HasChange("enable_auto_scaling") && !d.IsNewResource() {
		update = true
		args.AutoScaling.EnableAutoScaling = d.Get("enable_auto_scaling").(bool)
		d.SetPartial("enable_auto_scaling")
	}

	if d.HasChange("max") && !d.IsNewResource() {
		update = true
		args.AutoScaling.MaxInstance = d.Get("max").(int64)
		d.SetPartial("max")
	}

	if d.HasChange("min") && !d.IsNewResource() {
		update = true
		args.AutoScaling.MinInstance = d.Get("min").(int64)
		d.SetPartial("min")
	}

	if update && !d.IsNewResource() {

		var resoponse interface{}
		if err := invoker.Run(func() error {
			var err error
			resoponse, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				resp, err := csClient.UpdateNodePool(d.Get("cluster_id").(string), d.Id(), args)
				return resp, err
			})
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateKubernetesNodePool", DenverdinoAliyungo)
		}
		if debugOn() {
			resizeRequestMap := make(map[string]interface{})
			resizeRequestMap["ClusterId"] = d.Get("cluster_id").(string)
			resizeRequestMap["NodePoolId"] = d.Id()
			resizeRequestMap["Args"] = args
			addDebug("UpdateKubernetesNodePool", resoponse, resizeRequestMap)
		}

		stateConf := BuildStateConf([]string{"scaling"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Get("cluster_id").(string), d.Id(), []string{"deleting", "failed"}))

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	update = false
	d.Partial(false)
	return resourceAlicloudCSNodePoolRead(d, meta)
}

func resourceAlicloudCSNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	object, err := csService.DescribeCsKubernetesNodePool(d.Get("cluster_id").(string), d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("node_count", object.TotalNodes)
	d.Set("cluster_id", d.Get("cluster_id").(string))
	d.Set("name", object.Name)
	d.Set("vpc_id", object.VpcId)
	d.Set("vswitch_ids", object.VswitchIds)
	d.Set("instance_types", object.InstanceTypes)
	//d.Set("password", object.LoginPassword)
	d.Set("key_name", object.KeyPair)
	d.Set("security_group_id", object.SecurityGroupId)
	d.Set("system_disk_category", object.SystemDiskCategory)
	d.Set("system_disk_size", object.SystemDiskSize)
	d.Set("image_id", object.ImageId)
	d.Set("data_disks", object.DataDisks)
	d.Set("node_name_mode", object.NodeNameMode)
	d.Set("user_data", object.UserData)
	if sg, ok := d.GetOk("max"); ok && sg.(string) != "" {
		d.Set("max", object.MaxInstance)
	}
	if _, ok := d.GetOk("enable_auto_scaling"); ok {
		d.Set("enable_auto_scaling", object.EnableAutoScaling)
	}
	if sg, ok := d.GetOk("min"); ok && sg.(string) != "" {
		d.Set("min", object.MinInstance)
	}

	return nil
}

func resourceAlicloudCSNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	var response interface{}
	err := resource.Retry(30*time.Minute, func() *resource.RetryError {
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return nil, csClient.DeleteNodePool(d.Get("cluster_id").(string), d.Id())
			})
			response = raw
			return err
		}); err != nil {
			return resource.RetryableError(err)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Get("cluster_id").(string)
			requestMap["NodePoolId"] = d.Id()
			addDebug("DeleteClusterNodePool", response, d.Id(), requestMap)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorClusterNodePoolNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteClusterNodePool", DenverdinoAliyungo)
	}

	stateConf := BuildStateConf([]string{"active", "deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Get("cluster_id").(string), d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildNodePoolArgs(d *schema.ResourceData, meta interface{}) (*cs.CreateNodePoolRequest, error) {
	client := meta.(*connectivity.AliyunClient)

	vpcService := VpcService{client}

	var vswitchID string
	if list := expandStringList(d.Get("vswitch_ids").([]interface{})); len(list) > 0 {
		vswitchID = list[0]
	} else {
		vswitchID = ""
	}

	var vpcId string
	if vswitchID != "" {
		vsw, err := vpcService.DescribeVSwitch(vswitchID)
		if err != nil {
			return nil, err
		}
		vpcId = vsw.VpcId
	}

	password := d.Get("password").(string)
	if password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return nil, WrapError(err)
			}
			password = decryptResp.Plaintext
		}
	}

	creationArgs := &cs.CreateNodePoolRequest{
		RegionId: common.Region(client.RegionId),
		Count:    int64(d.Get("node_count").(int)),
		NodePoolInfo: cs.NodePoolInfo{
			Name:         d.Get("name").(string),
			NodePoolType: "ess", // hard code the type
		},
		ScalingGroup: cs.ScalingGroup{
			VpcId:              vpcId,
			VswitchIds:         expandStringList(d.Get("vswitch_ids").([]interface{})),
			InstanceTypes:      expandStringList(d.Get("instance_types").([]interface{})),
			LoginPassword:      password,
			KeyPair:            d.Get("key_name").(string),
			SystemDiskCategory: d.Get("system_disk_category").(string),
			SystemDiskSize:     int64(d.Get("system_disk_size").(int)),
			SecurityGroupId:    d.Get("security_group_id").(string),
			ImageId:            d.Get("image_id").(string),
		},
		KubernetesConfig: cs.KubernetesConfig{
			NodeNameMode: d.Get("node_name_mode").(string),
		},
	}

	setNodePoolDataDisks(&creationArgs.ScalingGroup, d)
	setNodePoolTags(&creationArgs.ScalingGroup, d)
	setNodePoolTaints(&creationArgs.KubernetesConfig, d)
	setNodePoolLabels(&creationArgs.KubernetesConfig, d)

	if v := d.Get("user_data").(string); v != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(v)
		if base64DecodeError == nil {
			creationArgs.KubernetesConfig.UserData = v
		} else {
			creationArgs.KubernetesConfig.UserData = base64.StdEncoding.EncodeToString([]byte(v))
		}
	}

	if v, ok := d.GetOk("enable_auto_scaling"); ok {
		if enable, ok := v.(bool); ok && enable {
			creationArgs.AutoScaling = cs.AutoScaling{
				EnableAutoScaling: true,
				MaxInstance:       int64(d.Get("max").(int)),
				MinInstance:       int64(d.Get("min").(int)),
				Type:              d.Get("type").(string),
			}
		}
	}

	return creationArgs, nil
}

func setNodePoolTags(scalingGroup *cs.ScalingGroup, d *schema.ResourceData) error {
	if v, ok := d.GetOk("tags"); ok && len(v.([]interface{})) > 0 {
		tags, ok := v.([]cs.Tag)
		if !ok {
			return fmt.Errorf("failed to convert nodepool tags: %++v", v)
		}
		scalingGroup.Tags = tags
	}

	return nil
}

func setNodePoolLabels(config *cs.KubernetesConfig, d *schema.ResourceData) error {
	if v, ok := d.GetOk("labels"); ok && len(v.([]interface{})) > 0 {
		labels, ok := v.([]cs.Label)
		if !ok {
			return fmt.Errorf("failed to convert nodepool labels: %++v", v)
		}
		config.Labels = labels
	}

	return nil
}

func setNodePoolDataDisks(scalingGroup *cs.ScalingGroup, d *schema.ResourceData) error {
	if dds, ok := d.GetOk("data_disks"); ok {
		disks := dds.([]interface{})
		createDataDisks := make([]cs.NodePoolDataDisk, 0, len(disks))
		for _, e := range disks {
			pack := e.(map[string]interface{})
			dataDisk := cs.NodePoolDataDisk{
				Size:                 pack["size"].(int),
				DiskName:             pack["name"].(string),
				Category:             pack["category"].(string),
				Device:               pack["device"].(string),
				AutoSnapshotPolicyId: pack["auto_snapshot_policy_id"].(string),
				KMSKeyId:             pack["kms_key_id"].(string),
				Encrypted:            pack["encrypted"].(string),
			}
			createDataDisks = append(createDataDisks, dataDisk)
		}
		scalingGroup.DataDisks = createDataDisks
	}

	return nil
}

func setNodePoolTaints(config *cs.KubernetesConfig, d *schema.ResourceData) error {
	if v, ok := d.GetOk("taints"); ok && len(v.([]interface{})) > 0 {
		taints, ok := v.([]cs.Taint)
		if !ok {
			return fmt.Errorf("failed to convert nodepool taints: %++v", v)
		}
		config.Taints = taints
	}

	return nil
}
