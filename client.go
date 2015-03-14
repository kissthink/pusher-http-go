package pusher

import (
	"encoding/json"
	"net/url"
	"strings"
)

type Client struct {
	AppId, Key, Secret string
}

func (c *Client) Trigger(channels []string, event string, _data interface{}) (error, string) {
	data, _ := json.Marshal(_data)

	payload, _ := json.Marshal(&EventBody{
		Name:     event,
		Channels: channels,
		Data:     string(data)})

	path := "/apps/" + c.AppId + "/" + "events"

	u := Url{"POST", path, c.Key, c.Secret, payload, nil}

	err, response := Request("POST", u.generate(), payload)

	return err, string(response)
}

func (c *Client) Channels(additional_queries map[string]string) (error, string) {

	path := "/apps/" + c.AppId + "/channels"
	u := Url{"GET", path, c.Key, c.Secret, nil, additional_queries}
	err, response := Request("GET", u.generate(), nil)
	return err, string(response)
}

func (c *Client) Channel(name string, additional_queries map[string]string) (error, *Channel) {

	path := "/apps/" + c.AppId + "/channels/" + name

	u := Url{"GET", path, c.Key, c.Secret, nil, additional_queries}

	err, raw_channel_data := Request("GET", u.generate(), nil)

	channel := &Channel{Name: name, Client: *c}
	json.Unmarshal(raw_channel_data, &channel)
	return err, channel

}

func (c *Client) AuthenticateChannel(_params []byte, presence_data map[string]string) string {
	params, _ := url.ParseQuery(string(_params))

	channel_name := params["channel_name"][0]
	socket_id := params["socket_id"][0]

	string_to_sign := socket_id + ":" + channel_name

	is_presence_channel := strings.HasPrefix(channel_name, "presence-")

	var json_user_data string
	_response := make(map[string]string)

	if is_presence_channel && presence_data != nil {
		_json_user_data, _ := json.Marshal(presence_data)
		json_user_data = string(_json_user_data)
		string_to_sign += ":" + json_user_data

		_response["channel_data"] = json_user_data
	}

	auth_signature := HMACSignature(string_to_sign, c.Secret)

	_response["auth"] = c.Key + ":" + auth_signature

	response, _ := json.Marshal(_response)

	return string(response)
}
