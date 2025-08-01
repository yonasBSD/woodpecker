// Copyright 2021 Woodpecker Authors
// Copyright 2018 Drone.IO Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"fmt"
	"strings"
)

type ApprovalMode string

const (
	RequireApprovalNone         ApprovalMode = "none"          // require approval for no events
	RequireApprovalForks        ApprovalMode = "forks"         // require approval for PRs from forks (default)
	RequireApprovalPullRequests ApprovalMode = "pull_requests" // require approval for all PRs
	RequireApprovalAllEvents    ApprovalMode = "all_events"    // require approval for all external events
)

func (mode ApprovalMode) Valid() bool {
	switch mode {
	case RequireApprovalNone,
		RequireApprovalForks,
		RequireApprovalPullRequests,
		RequireApprovalAllEvents:
		return true
	default:
		return false
	}
}

// Repo represents a repository.
type Repo struct {
	ID      int64 `json:"id,omitempty"                    xorm:"pk autoincr 'id'"`
	UserID  int64 `json:"-"                               xorm:"INDEX 'user_id'"`
	ForgeID int64 `json:"forge_id,omitempty"              xorm:"forge_id"`
	// ForgeRemoteID is the unique identifier for the repository on the forge.
	ForgeRemoteID                ForgeRemoteID        `json:"forge_remote_id"                 xorm:"forge_remote_id"`
	OrgID                        int64                `json:"org_id"                          xorm:"INDEX 'org_id'"`
	Owner                        string               `json:"owner"                           xorm:"UNIQUE(name) 'owner'"`
	Name                         string               `json:"name"                            xorm:"UNIQUE(name) 'name'"`
	FullName                     string               `json:"full_name"                       xorm:"UNIQUE 'full_name'"`
	Avatar                       string               `json:"avatar_url,omitempty"            xorm:"varchar(500) 'avatar'"`
	ForgeURL                     string               `json:"forge_url,omitempty"             xorm:"varchar(1000) 'forge_url'"`
	Clone                        string               `json:"clone_url,omitempty"             xorm:"varchar(1000) 'clone'"`
	CloneSSH                     string               `json:"clone_url_ssh"                   xorm:"varchar(1000) 'clone_ssh'"`
	Branch                       string               `json:"default_branch,omitempty"        xorm:"varchar(500) 'branch'"`
	PREnabled                    bool                 `json:"pr_enabled"                      xorm:"DEFAULT TRUE 'pr_enabled'"`
	Timeout                      int64                `json:"timeout,omitempty"               xorm:"timeout"`
	Visibility                   RepoVisibility       `json:"visibility"                      xorm:"varchar(10) 'visibility'"`
	IsSCMPrivate                 bool                 `json:"private"                         xorm:"private"`
	Trusted                      TrustedConfiguration `json:"trusted"                         xorm:"json 'trusted'"`
	RequireApproval              ApprovalMode         `json:"require_approval"                xorm:"varchar(50) require_approval"`
	ApprovalAllowedUsers         []string             `json:"approval_allowed_users"          xorm:"json approval_allowed_users"`
	IsActive                     bool                 `json:"active"                          xorm:"active"`
	AllowPull                    bool                 `json:"allow_pr"                        xorm:"allow_pr"`
	AllowDeploy                  bool                 `json:"allow_deploy"                    xorm:"allow_deploy"`
	Config                       string               `json:"config_file"                     xorm:"varchar(500) 'config_path'"`
	Hash                         string               `json:"-"                               xorm:"varchar(500) 'hash'"`
	Perm                         *Perm                `json:"-"                               xorm:"-"`
	CancelPreviousPipelineEvents []WebhookEvent       `json:"cancel_previous_pipeline_events" xorm:"json 'cancel_previous_pipeline_events'"`
	NetrcTrustedPlugins          []string             `json:"netrc_trusted"                   xorm:"json 'netrc_trusted'"`
} //	@name	Repo

// TableName return database table name for xorm.
func (Repo) TableName() string {
	return "repos"
}

func (r *Repo) ResetVisibility() {
	r.Visibility = VisibilityPublic
	if r.IsSCMPrivate {
		r.Visibility = VisibilityPrivate
	}
}

// ParseRepo parses the repository owner and name from a string.
func ParseRepo(str string) (user, repo string, err error) {
	before, after, _ := strings.Cut(str, "/")
	if before == "" || after == "" {
		err = fmt.Errorf("invalid or missing repository (e.g. octocat/hello-world)")
		return
	}
	user = before
	repo = after
	return
}

// Update updates the repository with values from the given Repo.
func (r *Repo) Update(from *Repo) {
	if from.ForgeRemoteID.IsValid() {
		r.ForgeRemoteID = from.ForgeRemoteID
	}
	r.Owner = from.Owner
	r.Name = from.Name
	r.FullName = from.FullName
	r.Avatar = from.Avatar
	r.ForgeURL = from.ForgeURL
	r.PREnabled = from.PREnabled
	if len(from.Clone) > 0 {
		r.Clone = from.Clone
	}
	if len(from.CloneSSH) > 0 {
		r.CloneSSH = from.CloneSSH
	}
	r.Branch = from.Branch
	if from.IsSCMPrivate != r.IsSCMPrivate {
		if from.IsSCMPrivate {
			r.Visibility = VisibilityPrivate
		} else {
			r.Visibility = VisibilityPublic
		}
	}
	r.IsSCMPrivate = from.IsSCMPrivate
}

// RepoPatch represents a repository patch object.
type RepoPatch struct {
	Config                       *string                    `json:"config_file,omitempty"`
	RequireApproval              *string                    `json:"require_approval,omitempty"`
	ApprovalAllowedUsers         *[]string                  `json:"approval_allowed_users,omitempty"`
	Timeout                      *int64                     `json:"timeout,omitempty"`
	Visibility                   *string                    `json:"visibility,omitempty"`
	AllowPull                    *bool                      `json:"allow_pr,omitempty"`
	AllowDeploy                  *bool                      `json:"allow_deploy,omitempty"`
	CancelPreviousPipelineEvents *[]WebhookEvent            `json:"cancel_previous_pipeline_events"`
	NetrcTrusted                 *[]string                  `json:"netrc_trusted"`
	Trusted                      *TrustedConfigurationPatch `json:"trusted"`
} //	@name	RepoPatch

type ForgeRemoteID string

func (r ForgeRemoteID) IsValid() bool {
	return r != "" && r != "0"
}

type TrustedConfiguration struct {
	Network  bool `json:"network"`
	Volumes  bool `json:"volumes"`
	Security bool `json:"security"`
}

type TrustedConfigurationPatch struct {
	Network  *bool `json:"network"`
	Volumes  *bool `json:"volumes"`
	Security *bool `json:"security"`
}

// RepoLastPipeline represents a repository with last pipeline execution information.
type RepoLastPipeline struct {
	*Repo
	LastPipeline *Pipeline `json:"last_pipeline,omitempty"`
} //	@name	RepoLastPipeline
