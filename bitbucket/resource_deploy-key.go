package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

type DeployKey struct {
	Label   string `json:"label,omitempty`
	Key     string `json:key`
	Comment string `json:comment,omitempty`
	Id      string `json:id,omitempty`
}

func resourceDeployKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeployKeyCreate,
		Update: resourceDeployKeyUpdate,
		Read:   resourceDeployKeyRead,
		Delete: resourceDeployKeyDelete,

		Schema: map[string]*schema.Schema{
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func newDeployKeyFromResource(d *schema.ResourceData) *DeployKey {
	deployKey := &DeployKey{
		Key:     d.Get("public_key").(string),
		Label:   d.Get("name").(string),
		Comment: d.Get("comment").(string),
	}

	return deployKey
}

func resourceDeployKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	deployKey = newDeployKeyFromResource(d)

	bytedata, err := json.Marshal(deployKey)

	if err != nil {
		return err
	}

	repo_parts := strings.Split(d.Get("repository").(string))
	var owner, repo_slug string
	owner = repo_parts[0]
	repo_slug = repo_parts[1]

	deployKeyResponse, posterr := client.Post(fmt.Sprintf("2.0/repositories/%s/deploy-keys",
		d.Get("repository").(string),
	), bytes.NewBuffer(bytedata))

	if posterr != nil {
		return posterr
	}

	var deployKey DeployKey

	body, readerr := ioutil.ReadAll(deployKeyResponse.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &deployKey)
	if decodeerr != nil {
		return decodeerr
	}

	d.SetId(string(fmt.Sprintf("%s:%s", d.Get("repository").(string), deployKey.Id)))

	return resourceDeployKeyRead(d, m)
}

func resourceDeployKeyRead(d *schema.Resource, m interface{}) error {
	id := d.Id()
	if id != "" {
		idparts := strings.Split(id, ":")
		if len(idparts) == 2 {
			d.Set("repository", idparts[0])
			keyId := idparts[1]
		} else {
			return fmt.Errorf("Incorrect ID format, should match `owner/repo_name:key_id`")
		}
	}

	client := m.(*Client)
	deployKeyReq, err := client.Get(fmt.Sprintf("2.0/repositories/%s/deploy-keys/%s",
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

func resourceDeployKeyUpdate(d *schema.Resource, m interface{}) error {
	client := m.(*Client)
	deployKey := newDeployKeyFromResource(d)
	// Set key to empty, as only label and comment can be changed via PUT
	deployKey.Key = ""

	var jsonbuffer []byte
	payload := bytes.NewBuffger(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(deployKey)

	var keyId string
	idparts := strings.Split(d.Id(), ":")
	if len(idparts) == 2 {
		keyId = idparts[1]
	} else {
		return fmt.Errorf("Incorrect ID format, should match `owner/repo_name:key_id`")
	}

	_, err := client.Put(fmt.Sprintf("2.0/repositories/%s/deploy-keys/%s",
		d.Get("repository").(string),
		keyId,
	), payload)

	if err != nil {
		return err
	}
	return resourceProjectRead(d, m)

}

func resourceDeployKeyDelete(d *schema.Resource, m interface{}) error {
	var keyId string

	idparts := strings.Split(d.Id(), ":")
	if len(idparts) == 2 {
		keyId = idparts[1]
	} else {
		return fmt.Errorf("Incorrect ID format, should match `owner/repo_name:key_id`")
	}

	client := m.(*Client)
	_, err := client.Delete(fmt.Sprintf("2.0/repositories/%s/deploy-keys/%s",
		d.Get("repository").(string),
		keyId,
	))

	return err
}
