package chserver

import (
	"net/http"

	rportplus "github.com/cloudradar-monitoring/rport/plus"
	"github.com/cloudradar-monitoring/rport/plus/capabilities/oauth"
	"github.com/cloudradar-monitoring/rport/server/api"
	"github.com/cloudradar-monitoring/rport/server/routes"
)

const BuiltInAuthProviderName = "built-in"

// AuthProviderInfo contains the provider name and the uris to be used
// for either regular or device flow based authorization
type AuthProviderInfo struct {
	AuthProvider      string `json:"auth_provider"`
	SettingsURI       string `json:"settings_uri"`
	DeviceSettingsURI string `json:"device_settings_uri"`
	MaxTokenLifetime  int    `json:"max_token_lifetime"`
}

// AuthSettings contains the auth info to be used by a regular web app
// type authorization
type AuthSettings struct {
	AuthProvider string           `json:"auth_provider"`
	LoginInfo    *oauth.LoginInfo `json:"details"`
}

// DeviceAuthSettings contains the auth info to be used by a CLI or
// similarly constrained app
type DeviceAuthSettings struct {
	AuthProvider string                 `json:"auth_provider"`
	LoginInfo    *oauth.DeviceLoginInfo `json:"details"`
}

func (al *APIListener) handleGetAuthProvider(w http.ResponseWriter, req *http.Request) {
	var response api.SuccessPayload

	maxTokenLifetime := al.config.API.MaxTokenLifeTimeHours

	if al.config.PlusOAuthEnabled() {
		OAuthProvider := AuthProviderInfo{
			AuthProvider:      al.config.PlusConfig.OAuthConfig.Provider,
			SettingsURI:       routes.AllRoutesPrefix + routes.AuthRoutesPrefix + routes.AuthSettingsRoute,
			DeviceSettingsURI: routes.AllRoutesPrefix + routes.AuthRoutesPrefix + routes.AuthDeviceSettingsRoute,
			MaxTokenLifetime:  maxTokenLifetime,
		}
		response = api.NewSuccessPayload(OAuthProvider)
	} else {
		builtInAuthProvider := AuthProviderInfo{
			AuthProvider:     BuiltInAuthProviderName,
			SettingsURI:      "",
			MaxTokenLifetime: maxTokenLifetime,
		}
		response = api.NewSuccessPayload(builtInAuthProvider)
	}
	al.writeJSONResponse(w, http.StatusOK, response)
}

func (al *APIListener) handleGetAuthSettings(w http.ResponseWriter, req *http.Request) {
	if !al.config.PlusOAuthEnabled() {
		al.jsonErrorResponse(w, http.StatusForbidden, rportplus.ErrPlusNotAvailable)
		return
	}

	plus := al.Server.plusManager
	capEx := plus.GetOAuthCapabilityEx()
	if capEx == nil {
		al.jsonErrorResponse(w, http.StatusForbidden, rportplus.ErrCapabilityNotAvailable(rportplus.PlusOAuthCapability))
		return
	}

	loginInfo, err := capEx.GetLoginInfo()
	if err != nil {
		al.jsonErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	settings := AuthSettings{
		AuthProvider: al.config.PlusConfig.OAuthConfig.Provider,
		LoginInfo:    loginInfo,
	}
	response := api.NewSuccessPayload(settings)
	al.writeJSONResponse(w, http.StatusOK, response)
}

func (al *APIListener) handleGetAuthDeviceSettings(w http.ResponseWriter, req *http.Request) {
	if !al.config.PlusOAuthEnabled() {
		al.jsonErrorResponse(w, http.StatusForbidden, rportplus.ErrPlusNotAvailable)
		return
	}

	plus := al.Server.plusManager
	capEx := plus.GetOAuthCapabilityEx()
	if capEx == nil {
		al.jsonErrorResponse(w, http.StatusForbidden, rportplus.ErrCapabilityNotAvailable(rportplus.PlusOAuthCapability))
		return
	}

	loginInfo, err := capEx.GetLoginInfoForDevice(req)
	if err != nil {
		al.jsonErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	settings := DeviceAuthSettings{
		AuthProvider: al.config.PlusConfig.OAuthConfig.Provider,
		LoginInfo:    loginInfo,
	}

	response := api.NewSuccessPayload(settings)
	al.writeJSONResponse(w, http.StatusOK, response)
}
