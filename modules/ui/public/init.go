/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package public

import (
	"github.com/emirpasic/gods/sets/hashset"
	"golang.org/x/oauth2"
	"github.com/huminghe/infini-framework/core/api"
	"github.com/huminghe/infini-framework/core/ui"
)

const (
	githubAuthorizeUrl = "https://github.com/login/oauth/authorize"
	githubTokenUrl     = "https://github.com/login/oauth/access_token"
	redirectUrl        = ""
)

var (
	oauthCfg *oauth2.Config

	// scopes
	scopes = []string{"repo"}
)

func InitUI(cfg ui.AuthConfig) {

	public := PublicUI{}
	ui.HandleUIMethod(api.GET, "/redirect/", public.RedirectHandler)

	if !cfg.Enabled {
		return
	}

	oauthCfg = &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  githubAuthorizeUrl,
			TokenURL: githubTokenUrl,
		},
		RedirectURL: redirectUrl,
		Scopes:      scopes,
	}

	public.admins = hashset.New()
	for _, v := range cfg.AuthorizedAdmins {
		if v != "" {
			public.admins.Add(v)
		}
	}

	ui.HandleUIMethod(api.GET, "/auth/github/", public.AuthHandler)
	ui.HandleUIMethod(api.GET, "/auth/callback/", public.CallbackHandler)
	ui.HandleUIMethod(api.GET, "/auth/login/", public.Login)
	ui.HandleUIMethod(api.GET, "/auth/logout/", public.Logout)
	ui.HandleUIMethod(api.GET, "/auth/success/", public.LoginSuccess)
	ui.HandleUIMethod(api.GET, "/auth/fail/", public.LoginFail)

}
