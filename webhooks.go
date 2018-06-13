package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	validator "gopkg.in/go-playground/validator.v8"
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

	if err = validate.Struct(Entry{}); err != nil {

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.NameNamespace)
			fmt.Println(err.Field)
			fmt.Println(err.Tag)
			fmt.Println(err.ActualTag)
			fmt.Println(err.Kind)
			fmt.Println(err.Type)
			fmt.Println(err.Value)
			fmt.Println(err.Param)
			fmt.Println()
		}

	}

	fmt.Println("got webhook response", entry)

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
