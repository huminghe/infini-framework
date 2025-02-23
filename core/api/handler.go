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

package api

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/jmoiron/jsonq"
	"github.com/huminghe/infini-framework/core/api/router"
	"github.com/huminghe/infini-framework/core/errors"
	"github.com/huminghe/infini-framework/core/util"
	"io/ioutil"
	"net/http"
	"strings"
)

// Method is object of http method
type Method string

const (
	// GET is http get method
	GET Method = "GET"
	// POST is http post method
	POST Method = "POST"
	// PUT is http put method
	PUT Method = "PUT"
	// DELETE is http delete method
	DELETE Method = "DELETE"
	// HEAD is http head method
	HEAD Method = "HEAD"

	OPTIONS Method = "OPTIONS"
)

// String return http method as string
func (method Method) String() string {
	switch method {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case HEAD:
		return "HEAD"
	}
	return "N/A"
}

// Handler is the object of http handler
type Handler struct {
	wroteHeader bool

	//w http.ResponseWriter
	//req *http.Request
	//
	formParsed bool
}

// WriteHeader write status code to http header
func (handler Handler) WriteHeader(w http.ResponseWriter, code int) {

	if apiConfig.TLSConfig.TLSEnabled {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	}

	w.WriteHeader(code)
	handler.wroteHeader = true
}

// Get http parameter or return default value
func (handler Handler) Get(req *http.Request, key string, defaultValue string) string {
	if !handler.formParsed {
		req.ParseForm()
	}
	if len(req.Form) > 0 {
		return req.Form.Get(key)
	}
	return defaultValue
}

// GetHeader return specify http header or return default value if not set
func (handler Handler) GetHeader(req *http.Request, key string, defaultValue string) string {
	v := req.Header.Get(key)
	if strings.TrimSpace(v) == "" {
		return defaultValue
	}
	return v
}

// EncodeJSON encode the object to json string
func (handler Handler) EncodeJSON(v interface{}) (b []byte, err error) {

	//if(w.Get("pretty","false")=="true"){
	b, err = json.MarshalIndent(v, "", "  ")
	//}else{
	//	b, err = json.Marshal(v)
	//}

	if err != nil {
		return nil, err
	}
	return b, nil
}

// WriteJSONHeader will write standard json header
func (handler Handler) WriteJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	handler.wroteHeader = true
}

// Result is a general json result
type Result struct {
	Total  int         `json:"total"`
	Result interface{} `json:"result"`
}

// WriteJSONListResult output result list to json format
func (handler Handler) WriteJSONListResult(w http.ResponseWriter, total int, v interface{}, statusCode int) error {
	result := Result{}
	result.Total = total
	result.Result = v
	return handler.WriteJSON(w, result, statusCode)
}

// WriteJSON output signal result with json format
func (handler Handler) WriteJSON(w http.ResponseWriter, v interface{}, statusCode int) error {
	if !handler.wroteHeader {
		handler.WriteJSONHeader(w)
		w.WriteHeader(statusCode)
	}

	b, err := handler.EncodeJSON(v)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (handler Handler) WriteAckJSON(w http.ResponseWriter, ack bool) error {
	if !handler.wroteHeader {
		handler.WriteJSONHeader(w)
		w.WriteHeader(200)
	}

	v := map[string]bool{}
	v["ok"] = ack

	b, err := handler.EncodeJSON(v)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// GetParameter return query parameter with argument name
func (handler Handler) GetParameter(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// GetParameterOrDefault return query parameter or return default value
func (handler Handler) GetParameterOrDefault(r *http.Request, key string, defaultValue string) string {
	v := r.URL.Query().Get(key)
	if len(v) > 0 {
		return v
	}
	return defaultValue
}

// GetIntOrDefault return parameter or default, data type is int
func (handler Handler) GetIntOrDefault(r *http.Request, key string, defaultValue int) int {

	v := handler.GetParameter(r, key)
	s, ok := util.ToInt(v)
	if ok != nil {
		return defaultValue
	}
	return s

}

// GetJSON return json input
func (handler Handler) GetJSON(r *http.Request) (*jsonq.JsonQuery, error) {

	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, errors.NewWithCode(err, errors.JSONIsEmpty, r.URL.String())
	}

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(string(content)))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	return jq, nil
}

func (handler Handler) DecodeJSON(r *http.Request, o interface{}) error {

	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return errors.NewWithCode(err, errors.JSONIsEmpty, r.URL.String())
	}

	return json.Unmarshal(content, o)
}

// GetRawBody return raw http request body
func (handler Handler) GetRawBody(r *http.Request) ([]byte, error) {

	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, errors.NewWithCode(err, errors.BodyEmpty, r.URL.String())
	}
	return content, nil
}

// Write response to client
func (handler Handler) Write(w http.ResponseWriter, b []byte) (int, error) {
	return w.Write(b)
}

// Error404 output 404 response
func (handler Handler) Error404(w http.ResponseWriter) {
	handler.WriteJSON(w, map[string]interface{}{"error": 404}, http.StatusNotFound)
}

// Error500 output 500 response
func (handler Handler) Error500(w http.ResponseWriter, msg string) {
	handler.WriteJSON(w, map[string]interface{}{"error": msg}, http.StatusInternalServerError)
}

// Error output custom error
func (handler Handler) Error(w http.ResponseWriter, err error) {
	handler.WriteJSON(w, map[string]interface{}{"error": err.Error()}, http.StatusInternalServerError)
}

// Flush flush response message
func (handler Handler) Flush(w http.ResponseWriter) {
	flusher := w.(http.Flusher)
	flusher.Flush()
}

// BasicAuth register api with basic auth
func BasicAuth(h httprouter.Handle, requiredUser, requiredPassword string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

var authEnabled = false

func NeedPermission(permission string, h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if !authEnabled || CheckPermission(w, r, permission) {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			//TODO redirect url configurable
			http.Redirect(w, r, "/auth/login/?redirect_url="+util.UrlEncode(r.URL.String()), 302)
		}
	}
}

func EnableAuth(enable bool) {
	authEnabled = enable
}

func IsAuthEnable() bool {
	return authEnabled
}

func Login(w http.ResponseWriter, r *http.Request, user, role string) {
	SetSession(w, r, "user", user)
	SetSession(w, r, "role", role)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	SetSession(w, r, "user", "")
	SetSession(w, r, "role", "")
}

func GetLoginInfo(w http.ResponseWriter, r *http.Request) (user, role string) {
	ok1, u := GetSession(w, r, "user")
	ok2, v := GetSession(w, r, "role")
	if !(ok1 && ok2) {
		return "", ""
	}
	return u.(string), v.(string)
}

func CheckPermission(w http.ResponseWriter, r *http.Request, requiredPermission string) bool {
	permissions := []string{}
	permissions = append(permissions, requiredPermission)
	return CheckPermissions(w, r, permissions)
}

func CheckPermissions(w http.ResponseWriter, r *http.Request, requiredPermissions []string) bool {
	user, role := GetLoginInfo(w, r)
	log.Trace("check user, ", user, ",", role, ",", requiredPermissions)
	if user != "" && role != "" {
		//TODO remove hard-coded permission check
		if role == ROLE_ADMIN {
			return true
		}

		perms, err := GetPermissionsByRole(role)
		if err != nil {
			log.Error(err)
			return false
		}

		for _, v := range requiredPermissions {
			if v != "" && !perms.Contains(v) {
				log.Tracef("user %s with role: %s do not have permission: %s", user, role, v)
				return false
			}
		}

		log.Trace("user logged in, ", user, ",", role, ",", requiredPermissions)
		return true
	}

	log.Trace("user not logged in, ", user, ",", role, ",", requiredPermissions)
	return false
}
