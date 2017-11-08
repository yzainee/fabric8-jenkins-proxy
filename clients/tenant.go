package clients

import (
	"errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	"net/http"
)

func NewTenant(tenantServiceURL string, authToken string) Tenant {
	return Tenant{
		authToken: authToken,
		tenantServiceURL: tenantServiceURL,
	}
}

type TenantInfo struct {
	Data TenantInfoData
}

type TenantInfoData struct {
	Attributes Attributes
	Id string
	Type string
}

type Attributes struct {
	CreatedAt time.Time `json:"created-at"`
	Email string
	Namespaces []Namespace
}

type Namespace struct {
	ClusterURL string `json:"cluster-url"`
	Name string
	State string
	Type string
}

type Tenant struct {
	tenantServiceURL string
	authToken string
}

func (t Tenant) GetTenantInfo(tenantId string) (*TenantInfo, error) {
	
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/tenants/%s", t.tenantServiceURL, tenantId), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", t.authToken)
	
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
	
		tenantInfo := &TenantInfo{}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, tenantInfo)
		if err != nil {
			return nil, err
		}
	
		return tenantInfo, nil
	}

	func (t Tenant) GetNamespaceByType(ti *TenantInfo, typ string) (r *Namespace, err error) {
		if ti == nil {
			err = errors.New("Cannot find namepsace - no info about tenant passed in.")
			return
		}

		for i:=0;i<len(ti.Data.Attributes.Namespaces);i++ {
			n := ti.Data.Attributes.Namespaces[i]
			if n.Type == typ {
				r = &n
				return
			}		
		}
	
		err = errors.New(fmt.Sprintf("Could not find tenant %s Jenkins namespace.", ti.Data.Attributes.Email))
		return
	}