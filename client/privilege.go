// This file permit to manage privilege in Ambari cluster
// Ambari documentation: https://community.hortonworks.com/questions/90797/manage-ambari-user-roles.html?childToView=90801

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Privilege object
type Privilege struct {
	PrivilegeInfo *PrivilegeInfo `json:"PrivilegeInfo"`
}
type PrivilegesResponse struct {
	Response
	Items []Privilege `json:"items"`
}
type PrivilegeInfo struct {
	PrivilegeId     int64  `json:"privilege_id,omitempty"`
	PermissionLabel string `json:"permission_label,omitempty"`
	PermissionName  string `json:"permission_name"`
	PrincipalName   string `json:"principal_name"`
	PrincipalType   string `json:"principal_type"`
}

// String return privilege object as Json string
func (p *Privilege) String() string {
	json, _ := json.Marshal(p)
	return string(json)
}

// Privilege return existing privilege on cluster
// It return the pivilege if is found
// It return nil if not found
// It return error if something wrong when it call the API
func (c *AmbariClient) Privilege(clusterName string, id int64) (*Privilege, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("Id: ", id)

	path := fmt.Sprintf("/clusters/%s/privileges/%d", clusterName, id)
	resp, err := c.Client().R().Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Result : ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		} else {
			return nil, NewAmbariError(resp.StatusCode(), resp.Status())
		}
	}
	privilege := &Privilege{}
	err = json.Unmarshal(resp.Body(), privilege)
	if err != nil {
		return nil, err
	}

	log.Debug("Privilege: ", privilege)

	return privilege, nil

}

// CreatePrivilege permit to create new privilege on cluster
// It return the privilege if all work fine
// It return error if something wrong when it call the API
func (c *AmbariClient) CreatePrivilege(clusterName string, privilege *Privilege) (*Privilege, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if privilege == nil {
		panic("Privilege can't be nil")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("Privilege :", privilege)

	// Create the privilege
	path := fmt.Sprintf("/clusters/%s/privileges", clusterName)
	jsonData, err := json.Marshal(privilege)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client().R().SetBody(jsonData).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to create: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	// Get the privilege
	privilege, err = c.SearchPrivilege(clusterName, privilege.PrivilegeInfo.PermissionName, privilege.PrivilegeInfo.PrincipalName, privilege.PrivilegeInfo.PrincipalType)
	if err != nil {
		return nil, err
	}
	if privilege == nil {
		return nil, NewAmbariError(500, "Can't get privilege that just created")
	}

	return privilege, err

}

// DeletePrivilege permit to delete existing privielege on cluster
// It return error if something wrong when it call the API
func (c *AmbariClient) DeletePrivilege(clusterName string, id int64) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)

	path := fmt.Sprintf("/clusters/%s/privileges/%d", clusterName, id)
	resp, err := c.Client().R().Delete(path)
	if err != nil {
		return err
	}
	log.Debug("Response to delete privilege: ", resp)
	if resp.StatusCode() >= 300 {
		return NewAmbariError(resp.StatusCode(), resp.Status())
	}

	return nil

}

// UpdatePrivilege permit to update existing privielege
func (c *AmbariClient) UpdatePrivilege(clusterName string, privilege *Privilege) (*Privilege, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if privilege == nil {
		panic("Privilege can't be nil")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("Privilege: ", privilege)

	// Update the privilege
	path := fmt.Sprintf("/clusters/%s/privileges/%d", clusterName, privilege.PrivilegeInfo.PrivilegeId)
	jsonData, err := json.Marshal(privilege)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client().R().SetBody(jsonData).Put(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to update: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	// Get the privilege because id and permission label change after update
	privilege, err = c.SearchPrivilege(clusterName, privilege.PrivilegeInfo.PermissionName, privilege.PrivilegeInfo.PrincipalName, privilege.PrivilegeInfo.PrincipalType)
	if err != nil {
		return nil, err
	}
	if privilege == nil {
		return nil, NewAmbariError(500, "Can't get privilege that just created")
	}

	return privilege, err

}

// SearchPrivilege permit to get privilege by is name
// It return privielege if is found
// It return nil if is not found
// It return error if something wrong when it call the API
func (c *AmbariClient) SearchPrivilege(clusterName string, permissionName string, principalName string, principalType string) (*Privilege, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if permissionName == "" {
		panic("PermissionName can't be empty")
	}
	if principalName == "" {
		panic("PrincipalName can't be empty")
	}
	if principalType == "" {
		panic("PrincipalType can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("PermissionName: ", permissionName)
	log.Debug("PrincipalName: ", principalName)
	log.Debug("PrincipalType: ", principalType)

	path := fmt.Sprintf("/clusters/%s/privileges", clusterName)
	resp, err := c.Client().R().SetQueryParams(map[string]string{
		"PrivilegeInfo/permission_name": permissionName,
		"PrivilegeInfo/principal_name":  principalName,
		"PrivilegeInfo/principal_type":  principalType,
	}).Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to get: ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		} else {
			return nil, NewAmbariError(resp.StatusCode(), resp.Status())
		}
	}
	privilegeResponses := &PrivilegesResponse{}
	err = json.Unmarshal(resp.Body(), privilegeResponses)
	if err != nil {
		return nil, err
	}
	log.Debug("PrivilegesResponse: ", privilegeResponses)

	if len(privilegeResponses.Items) > 0 {
		log.Debug("Privilege: ", privilegeResponses.Items[0])
		return &privilegeResponses.Items[0], nil
	} else {
		return nil, nil
	}
}
