package bitbucket

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func hashValue(v string) string {
	b := md5.Sum([]byte(v))
	return fmt.Sprintf("%x", b)
}

// RepositoryVariable structure for handling key info
type RepositoryVariable struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	UUID    string `json:"uuid,omitempty"`
	Secured bool   `json:"secured"`
}

func resourceRepositoryVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryVariableCreate,
		Update: resourceRepositoryVariableUpdate,
		Read:   resourceRepositoryVariableRead,
		Delete: resourceRepositoryVariableDelete,

		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"hash": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"secured": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func newRepositoryVariableFromResource(d *schema.ResourceData) *RepositoryVariable {
	dk := &RepositoryVariable{
		Key:     d.Get("key").(string),
		Value:   d.Get("value").(string),
		Secured: d.Get("secured").(bool),
	}
	return dk
}

func resourceRepositoryVariableCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*Client)
	rvcr := newRepositoryVariableFromResource(d)
	bytedata, err := json.Marshal(rvcr)

	if err != nil {
		return err
	}
	req, err := client.Post(fmt.Sprintf("2.0/repositories/%s/pipelines_config/variables/",
		d.Get("repository").(string),
	), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	var rv RepositoryVariable

	body, readerr := ioutil.ReadAll(req.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &rv)
	if decodeerr != nil {
		return decodeerr
	}
	d.Set("uuid", rv.UUID)
	d.Set("hash", hashValue(rvcr.Value))
	d.SetId(rv.Key)

	return resourceRepositoryVariableRead(d, m)
}

func resourceRepositoryVariableRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*Client)
	rvReq, _ := client.Get(fmt.Sprintf("2.0/repositories/%s/pipelines_config/variables/%s",
		d.Get("repository").(string),
		d.Get("uuid").(string),
	))

	log.Printf("ID: %s", url.PathEscape(d.Id()))

	if rvReq.StatusCode == 200 {
		var rv RepositoryVariable
		body, readerr := ioutil.ReadAll(rvReq.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &rv)
		if decodeerr != nil {
			return decodeerr
		}

		d.Set("uuid", rv.UUID)
		d.Set("key", rv.Key)
		d.Set("secured", rv.Secured)
		switch {
		case rv.Secured:
			v := d.Get("value").(string)
			h := d.Get("hash").(string)
			if hashValue(v) != h {
				d.Set("value", "")
				d.Set("hash", "")
			}
		default:
			d.Set("value", rv.Value)
		}
	}

	if rvReq.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceRepositoryVariableUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	rvcr := newRepositoryVariableFromResource(d)
	bytedata, err := json.Marshal(rvcr)

	if err != nil {
		return err
	}
	req, err := client.Put(fmt.Sprintf("2.0/repositories/%s/pipelines_config/variables/%s",
		d.Get("repository").(string),
		d.Get("uuid").(string),
	), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	if req.StatusCode != 200 {
		return nil
	}

	d.Set("hash", hashValue(rvcr.Value))
	return resourceRepositoryVariableRead(d, m)
}

func resourceRepositoryVariableDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, err := client.Delete(fmt.Sprintf(fmt.Sprintf("2.0/repositories/%s/pipelines_config/variables/%s",
		d.Get("repository").(string),
		d.Get("uuid").(string),
	)))
	return err
}
