package jpush

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

/**
 * HostPushSsl push url
 * HostScheduleSsl schedule url
 *
 */
const (
	HostPushSsl     = "https://api.jpush.cn/v3/push"
	HostScheduleSsl = "https://api.jpush.cn/v3/schedules"
	CHARSET         = "UTF-8"
	CONTENTTYPE     = "application/json"
)

type Client struct {
	AppKey       string
	MasterSecret string
	url          *url.URL
	method       string
}
type PushSendSuccessResponse struct {
	Sendno string `json:"sendno"`
	MsgID  string `json:"msg_id"`
}
type ScheduleSendSuccessResponse struct {
	ScheduleId string `json:"schedule_id"`
	Name       string `json:"name"`
}
type ScheduleListResponse struct {
	TotalCount int        `json:"total_count"`
	TotalPages int        `json:"total_pages"`
	Page       int        `json:"page"`
	Schedules  []Schedule `json:"schedules"`
}

func (this *Client) send(data []byte) ([]byte, error) {

	httpClient := &http.Client{}
	httpRequest, err := http.NewRequest(this.method, this.url.String(), bytes.NewBuffer(data))
	httpRequest.SetBasicAuth(this.AppKey, this.MasterSecret)
	httpRequest.Header.Add("Charset", CHARSET)
	httpRequest.Header.Add("Content-Type", CONTENTTYPE)
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

func (client *Client) PushSend(push *PayLoad) (PushSendSuccessResponse, error) {
	client.method = "POST"
	client.url, _ = url.Parse(HostPushSsl)
	bytes, _ := push.ToBytes()
	result, err := client.send(bytes)
	var pushSendSuccessResponse PushSendSuccessResponse
	var errorResponse ErrorResponse
	if err != nil {
		json.Unmarshal(result, &errorResponse)
		return pushSendSuccessResponse, errorResponse
	}
	json.Unmarshal(result, &pushSendSuccessResponse)
	return pushSendSuccessResponse, nil
}
func (client *Client) ScheduleSend(schedule *Schedule) (ScheduleSendSuccessResponse, error) {
	client.method = "POST"
	client.url, _ = url.Parse(HostScheduleSsl)
	bytes, _ := schedule.ToBytes()
	result, err := client.send(bytes)
	var scheduleSendSuccessResponse ScheduleSendSuccessResponse
	var errorReponse ErrorResponse
	if err != nil {
		json.Unmarshal(result, &errorReponse)
		return scheduleSendSuccessResponse, errorReponse
	}
	json.Unmarshal(result, &scheduleSendSuccessResponse)
	return scheduleSendSuccessResponse, nil
}
func (client *Client) ScheduleList(page int) (ScheduleListResponse, error) {
	client.method = "GET"
	client.url, _ = url.Parse(HostScheduleSsl)
	client.url.Query().Add("page", strconv.Itoa(page))
	result, err := client.send(nil)
	var scheduleListResponse ScheduleListResponse
	var errorReponse ErrorResponse
	if err != nil {
		json.Unmarshal(result, &errorReponse)
		return scheduleListResponse, errorReponse
	}
	json.Unmarshal(result, &scheduleListResponse)
	return scheduleListResponse, nil
}
func (client *Client) ScheduleShow(id string) (Schedule, error) {
	client.method = "GET"
	client.url, _ = url.Parse(HostScheduleSsl + "/" + id)
	result, err := client.send(nil)
	var schedule Schedule
	var errResponse ErrorResponse
	if err != nil {
		json.Unmarshal(result, &errResponse)
		return schedule, errResponse
	}
	json.Unmarshal(result, &schedule)
	return schedule, nil
}
func (client *Client) ScheduleDelete(id string) error {
	client.method = "DELETE"
	client.url, _ = url.Parse(HostScheduleSsl + "/" + id)
	result, err := client.send(nil)

	if err != nil {
		var errorResult ErrorResponse
		json.Unmarshal(result, &errorResult)
		return errorResult
	}
	return nil
}
