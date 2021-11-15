package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DeployKey struct {
	Label   string `json:"label,omitempty"`
	Key     string `json:"key,omitempty"`
	Comment string `json:"comment,omitempty"`
	Id      int    `json:"id,omitempty"`
}

func resourceDeployKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeployKeyCreate,
		Read:   resourceDeployKeyRead,
		Delete: resourceDeployKeyDelete,

		Schema: map[string]*schema.Schema{
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.TrimSpace(old) == strings.TrimSpace(new) {
						return true
					}
					return false
				},
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func newDeployKeyFromResource(d *schema.ResourceData) *DeployKey {
	deployKey := &DeployKey{
		Key:     strings.TrimSpace(d.Get("public_key").(string)),
		Label:   d.Get("name").(string),
		Comment: d.Get("comment").(string),
	}

	return deployKey
}

func resourceDeployKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	deployKey := newDeployKeyFromResource(d)

	bytedata, err := json.Marshal(deployKey)

	if err != nil {
		return err
	}

	deployKeyResponse, posterr := client.Post(fmt.Sprintf("2.0/repositories/%s/deploy-keys",
		d.Get("repository").(string),
	), bytes.NewBuffer(bytedata))

	if posterr != nil {
		return posterr
	}

	body, readerr := ioutil.ReadAll(deployKeyResponse.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &deployKey)
	if decodeerr != nil {
		return decodeerr
	}

	d.SetId(string(fmt.Sprintf("%s:%d", d.Get("repository").(string), deployKey.Id)))

	return resourceDeployKeyRead(d, m)
}

func resourceDeployKeyRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	var keyId int
	var iderr error
	if id != "" {
		idparts := strings.Split(id, ":")
		if len(idparts) == 2 {
			d.Set("repository", idparts[0])
			keyId, iderr = strconv.Atoi(idparts[1])
			if iderr != nil {
				return iderr
			}
		} else {
			return fmt.Errorf("Incorrect ID format, should match `owner/repo_name:key_id`")
		}
	}

	client := m.(*Client)
	deployKeyReq, err := client.Get(fmt.Sprintf("2.0/repositories/%s/deploy-keys/%d",
		d.Get("repository").(string),
		keyId,
	))

	if err != nil {
		return err
	}

	var deployKey DeployKey

	if err != nil {
		return err
	}

	body, readerr := ioutil.ReadAll(deployKeyReq.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &deployKey)
	if decodeerr != nil {
		return decodeerr
	}

	d.Set("public_key", deployKey.Key)
	d.Set("name", deployKey.Label)
	d.Set("comment", deployKey.Comment)

	return nil

}

func resourceDeployKeyDelete(d *schema.ResourceData, m interface{}) error {
	var keyId int
	var iderr error

	idparts := strings.Split(d.Id(), ":")
	if len(idparts) == 2 {
		keyId, iderr = strconv.Atoi(idparts[1])
		if iderr != nil {
			return iderr
		}
	} else {
		return fmt.Errorf("Incorrect ID format, should match `owner/repo_name:key_id`")
	}

	client := m.(*Client)
	_, err := client.Delete(fmt.Sprintf("2.0/repositories/%s/deploy-keys/%d",
		d.Get("repository").(string),
		keyId,
	))

	return err
}
