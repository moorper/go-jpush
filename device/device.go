package device

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	DevicesURL  = "https://device.jpush.cn/v3/devices"
	AliasURL    = "https://device.jpush.cn/v3/aliases"
	TagsURL     = "https://device.jpush.cn/v3/tags"
	Charset     = "UTF-8"
	ContentType = "application/json"
	Accept      = "application/json"
)

type Client struct {
	AppKey       string
	MasterSecret string
}
type DevicesInfo struct {
	Tags   []string `json:"tags"`
	Alias  string   `json:"alias"`
	Mobile string   `json:"mobile"`
}
type UpdateDevicesInfoRequest struct {
	Tags   interface{} `json:"tags,omitempty"`
	Alias  string      `json:"alias,omitempty"`
	Mobile string      `json:"mobile,omitempty"`
}
type AliasQueryResponse struct {
	RegistrationIDS []string `json:"registration_ids"`
}
type FetchTagsResponse struct {
	Tags []string `json:"tags"`
}
type UpdateTagsRequest struct {
	RegistrationIDS map[string][]string `json:"registration_ids"`
}

// GetDeviceInfo 获取指定设备信息
func (client Client) GetDevicesInfo(registrationID string) DevicesInfo {
	url, _ := url.Parse(DevicesURL + "/" + registrationID)

	body, err := client.send(url, http.MethodGet, nil)
	if err != nil {
		fmt.Println(err)
	}

	var devicesInfo DevicesInfo
	json.Unmarshal(body, &devicesInfo)
	return devicesInfo
}

// UpdateDeviceInfo 更新指定设备信息
func (client Client) UpdateDevicesInfo(registrationID string, tags interface{}, alias string, mobile string) {
	url, _ := url.Parse(DevicesURL + "/" + registrationID)
	var req UpdateDevicesInfoRequest
	req.Alias = alias
	req.Mobile = mobile
	if reflect.ValueOf(tags).Kind() == reflect.String {
		req.Tags = ""
	} else if reflect.ValueOf(tags).Kind() == reflect.Map && len(tags.(map[string][]string)) > 0 {
		req.Tags = tags
	}
	dataByte, _ := json.Marshal(req)
	fmt.Println(string(dataByte))
	body, err := client.send(url, http.MethodPost, dataByte)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

// AliasQuery 查询指定别名下的所有设备
func (client Client) FetchDevicesByAlias(alias string, platform []string) []string {
	url, _ := url.Parse(AliasURL + "/" + alias)
	if len(platform) > 0 {
		url.Query().Set("platform", strings.Join(platform, ","))
	}
	body, err := client.send(url, http.MethodGet, nil)
	if err != nil {
		panic(err)
	}
	var resp AliasQueryResponse
	json.Unmarshal(body, &resp)
	return resp.RegistrationIDS
}

// AliasDelete 删除指定 alias 可指定平台
func (client Client) AliasDelete(alias string, platform []string) {
	url, _ := url.Parse(AliasURL + "/" + alias)
	if len(platform) > 0 {
		url.Query().Set("platform", strings.Join(platform, ","))
	}
	_, err := client.send(url, http.MethodDelete, nil)
	if err != nil {
		panic(err)
	}
}

// FetchTags 获取当前账号下的所有 tags
func (client Client) FetchTags() []string {
	url, _ := url.Parse(TagsURL)
	body, err := client.send(url, http.MethodGet, nil)
	if err != nil {
		panic(err)
	}
	var resp FetchTagsResponse
	json.Unmarshal(body, &resp)
	return resp.Tags
}

// IsBoundTags 判断指定设备是否和指定 tags 绑定
func (client Client) IsBoundTags(registrationID string, tags string) bool {
	url, _ := url.Parse(TagsURL + "/" + tags + "/" + "registration_ids" + "/" + registrationID)
	body, err := client.send(url, http.MethodGet, nil)
	if err != nil {
		panic(err)
	}
	return strings.Contains(string(body), "true")
}

// UpdateTags 更新指定 tags 下设备
func (client Client) UpdateTags(tags string, add []string, remove []string) {
	url, _ := url.Parse(TagsURL + "/" + tags)
	var req UpdateTagsRequest
	req.RegistrationIDS = make(map[string][]string)
	if len(add) > 0 {
		req.RegistrationIDS["add"] = add
	}
	if len(remove) > 0 {
		req.RegistrationIDS["remove"] = remove
	}
	data, _ := json.Marshal(req)
	client.send(url, http.MethodPost, data)
}

// DeleteTags 删除指定 tags ，可指定平台
func (client Client) DeleteTags(tags string, platform []string) {
	url, _ := url.Parse(TagsURL + "/" + tags)
	if len(platform) > 0 {
		url.Query().Set("platform", strings.Join(platform, ","))
	}
	client.send(url, http.MethodDelete, nil)
}

// send 发送 http 请求
func (client Client) send(url *url.URL, method string, data []byte) ([]byte, error) {

	httpClient := &http.Client{}
	httpRequest, err := http.NewRequest(method, url.String(), bytes.NewBuffer(data))
	httpRequest.SetBasicAuth(client.AppKey, client.MasterSecret)
	httpRequest.Header.Add("Charset", Charset)
	httpRequest.Header.Add("Content-Type", ContentType)
	httpRequest.Header.Add("Accept", Accept)
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	httpResponseBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode != 200 {
		return httpResponseBody, errors.New(string(httpResponseBody))

	}
	return httpResponseBody, nil
}
