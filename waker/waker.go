package waker

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Pair struct {
	key   string
	value string
}

func gen_auth(user_name string, user_password string) string {
	input := []byte(user_name + ":" + user_password)
	encodeString := base64.StdEncoding.EncodeToString(input)
	return encodeString
}

func gen_form_body(values []Pair) string {
	var result string
	count := 0
	for _, pair := range values { //取map中的值
		escaped := url.QueryEscape(pair.value)
		if count > 0 {
			result += fmt.Sprintf("&%s=%s", pair.key, escaped)
		} else {
			result += fmt.Sprintf("%s=%s", pair.key, escaped)
		}
		count++
	}
	return result
}

type Waker struct {
	url    string
	client *http.Client
	token  string
	auth   string
}

func New(url, user, password string) Waker {
	return Waker{
		url,
		&http.Client{},
		"",
		gen_auth(user, password),
	}
}

func (w *Waker) Login() error {
	current_page := "Main_Login.asp"
	next_page := "index.asp"
	action_page := "login.cgi"

	value_package := make([]Pair, 0, 10)
	value_package = append(value_package, Pair{"group_id", ""})
	value_package = append(value_package, Pair{"action_mode", ""})
	value_package = append(value_package, Pair{"action_script", ""})
	value_package = append(value_package, Pair{"action_wait", "5"})
	value_package = append(value_package, Pair{"current_page", current_page})
	value_package = append(value_package, Pair{"next_page", next_page})
	value_package = append(value_package, Pair{"login_authorization", w.auth})

	body := gen_form_body(value_package)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", w.url, action_page), strings.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", fmt.Sprintf("%s/%s", w.url, current_page))

	resp, err := w.client.Do(req)

	if err != nil {
		return err
	} else {
		cookie := resp.Header.Get("Set-Cookie")
		parts := strings.Split(cookie, ";")
		if len(parts) > 1 {
			value := strings.Split(parts[0], "=")
			if len(value) == 2 {
				//key := value[0]
				w.token = value[1]
			} else {
				return errors.New("failed get token")
			}
		} else {
			return errors.New("unexpected response")
		}

		defer func() {
			err := resp.Body.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Waker) ExecuteCommand(cmd string) error {
	current_page := "Main_WOL_Content.asp"
	next_page := "Main_WOL_Content.asp"
	action_page := "apply.cgi"

	value_package := make([]Pair, 0, 10)
	value_package = append(value_package, Pair{"group_id", ""})
	value_package = append(value_package, Pair{"action_mode", " Refresh "})
	value_package = append(value_package, Pair{"action_script", ""})
	value_package = append(value_package, Pair{"action_wait", ""})
	value_package = append(value_package, Pair{"firmver", "3.0.0.4"})
	value_package = append(value_package, Pair{"first_time", ""})

	value_package = append(value_package, Pair{"preferred_lang", "CN"})
	value_package = append(value_package, Pair{"destIP", "C4:09:38:F5:73:92"})
	value_package = append(value_package, Pair{"next_page", next_page})
	value_package = append(value_package, Pair{"current_page", current_page})
	value_package = append(value_package, Pair{"SystemCmd", cmd})
	value_package = append(value_package, Pair{"wollist_macAddr", ""})

	body := gen_form_body(value_package)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", w.url, action_page), strings.NewReader(body))
	if err != nil {
		return err
	}

	cookie := fmt.Sprintf("asus_token=%s", w.token)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", fmt.Sprintf("%s/%s", w.url, current_page))
	req.Header.Set("Cookie", cookie)

	rsp, err := w.client.Do(req)
	if rsp.StatusCode != http.StatusOK {
		return errors.New(rsp.Status)
	}

	if err != nil {
		return err
	}
	return nil
}

func (w *Waker) Logout() error {

	current_page := "Main_WOL_Content.asp"
	//	next_page := "Main_Login.asp"
	action_page := "Logout.asp"
	/*
		value_package := make([]Pair, 0, 10)
		value_package = append(value_package, Pair{"group_id", ""})
		value_package = append(value_package, Pair{"action_mode", " Refresh "})
		value_package = append(value_package, Pair{"action_script", ""})
		value_package = append(value_package, Pair{"action_wait", ""})
		value_package = append(value_package, Pair{"firmver", "3.0.0.4"})
		value_package = append(value_package, Pair{"first_time", ""})

		value_package = append(value_package, Pair{"preferred_lang", "CN"})
		value_package = append(value_package, Pair{"next_page", next_page})
		value_package = append(value_package, Pair{"current_page", current_page})
	*/
	body := "" // gen_form_body(value_package)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", w.url, action_page), strings.NewReader(body))
	if err != nil {
		return err
	}

	cookie := fmt.Sprintf("asus_token=%s", w.token)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", fmt.Sprintf("%s/%s", w.url, current_page))
	req.Header.Set("Cookie", cookie)

	rsp, err := w.client.Do(req)
	if rsp.StatusCode != http.StatusOK {
		return errors.New(rsp.Status)
	}

	if err != nil {
		return err
	}
	return nil
}
