package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

func checkWebhookURL() bool {
	return config.WebhookURL != ""
}

func webhookGetInfo(id string) (Entry, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	var entry Entry

	if !checkWebhookURL() {
		return entry, errors.New("Missing webhook url")
	}

	resp, err := client.Get(config.WebhookURL + "?id=" + id)
	if err != nil {
		return entry, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return entry, nil
	}

	err = json.Unmarshal(body, &entry)
	if err != nil {
		return entry, err
	}

	err = entry.Validate()
	if err != nil {
		return entry, err
	}

	return entry, nil
}

func webhookPutInfo(entry *Entry) error {
	if !checkWebhookURL() {
		return errors.New("Missing webhook url")
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	json, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	resp, err := client.Post(config.WebhookURL, "application/json", bytes.NewReader(json))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	println(string(body))

	return nil
}

func webhookDelete(id string) error {
	if !checkWebhookURL() {
		return errors.New("Missing webhook url")
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodDelete, config.WebhookURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	println(string(body))

	return nil
}
