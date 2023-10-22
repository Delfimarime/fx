package audit

import "github.com/delfimarime/fx/time"

type CreatedAuditInfo struct {
	At        *time.ISODatetime `json:"created_at,omitempty"`
	Principal *Principal        `json:"created_by,omitempty"`
}

type LastModifiedAuditInfo struct {
	At        *time.ISODatetime `json:"last_modified_at,omitempty"`
	Principal *Principal        `json:"last_modified_by,omitempty"`
}

type Principal struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
}
