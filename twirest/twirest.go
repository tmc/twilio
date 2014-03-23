// Copyright (C) 2014 Cristoffer Kvist. All rights reserved.
// This project is licensed under the terms of the MIT license in LICENSE.

// Package twirest provides a interface to Twilio REST API allowing the user to
// query meta-data from their account and, to initiate calls and send SMS.
package twirest

import (
	"crypto/tls"
	//"crypto/x509"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const ApiVer string = "2010-04-01"

// TwilioClient struct for holding a http client and user credentials
type TwilioClient struct {
	httpclient            *http.Client
	accountSid, authToken string
}

// Create a new client
func NewClient(accountSid, authToken string) *TwilioClient {
	// certPool := x509.NewCertPool()
	// pemFile, err := os.Open("cacert.pem")
	// if err != nil {
	// 	err = fmt.Errorf("Using host's root CA\n\t%s", err)
	// 	certPool = nil
	// } else {
	// 	defer pemFile.Close()
	// 	bytes, _ := ioutil.ReadAll(pemFile)
	// 	certPool.AppendCertsFromPEM(bytes)
	// }
	tr := &http.Transport{TLSClientConfig: &tls.Config{RootCAs: nil},
		DisableCompression: true}
	client := &http.Client{Transport: tr}

	return &TwilioClient{client, accountSid, authToken}
}

// Request makes a REST resource or action request from twilio servers and
// returns the response. The type of request is determined by the request
// struct supplied.
func (twiClient *TwilioClient) Request(reqStruct interface{}) (
	TwilioResponse, error) {

	twiResp := TwilioResponse{}

	// setup a POST/GET/DELETE http request from request struct
	httpReq, err := httpRequest(reqStruct, twiClient.accountSid)
	if err != nil {
		return twiResp, err
	}
	// add authentication and headers to the http request
	httpReq.SetBasicAuth(twiClient.accountSid, twiClient.authToken)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Accept", "*/*")

	response, err := twiClient.httpclient.Do(httpReq)
	if err != nil {
		return twiResp, err
	}

	// Save http status code to response struct
	twiResp.Status.Http = response.StatusCode

	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	// parse xml response into twilioResponse struct
	xml.Unmarshal(body, &twiResp)

	twiResp.Status.Twilio, err = exceptionToErr(twiResp)
	return twiResp, err
}

// exceptiontToErr converts a Twilio response exception (if any) to a go error
func exceptionToErr(twir TwilioResponse) (code int, err error) {
	if twir.Exception != nil {
		return twir.Exception.Code, fmt.Errorf("%s",
			twir.Exception.Message)
	}
	return
}

func urlString(reqStruct interface{}, accSid string) (url string, err error) {

	url = "https://api.twilio.com/" + ApiVer + "/Accounts"

	switch reqSt := reqStruct.(type) {
	default:
		for i := 0; i < reflect.ValueOf(reqSt).NumField(); i++ {
			fldType := reflect.ValueOf(reqSt).Type().Field(i).Type
			fldTag := reflect.ValueOf(reqSt).Type().Field(i).Tag
			fldName := reflect.ValueOf(reqSt).Type().Field(i).Name
			fldValue := reflect.ValueOf(reqSt).Field(i).String()

			if fldType.Name() == "uri" {
				url = url + "/" + accSid + string(fldTag)
			}
			if fldName == "Sid" {
				err = required(fldValue)
				url = url + "/" + fldValue
			}
		}
	}

	switch reqSt := reqStruct.(type) {
	default:
	case Message:
		if reqSt.Media == true {
			url = url + "/Media"
			if reqSt.MediaSid != "" {
				url = url + "/" + reqSt.MediaSid
			}
		}
	case Call:
		if reqSt.Recordings == true {
			url = url + "/Recordings"
		} else if reqSt.Notifications == true {
			url = url + "/Notifications"
		}
	case UsageRecords:
		url = url + "/" + reqSt.SubResource
	case QueueMember:
		if reqSt.Front {
			url = url + "/Members/Front"
		} else {
			err = required(reqSt.CallSid)
			url = url + "/Members/" + reqSt.CallSid
		}
	case DeQueue:
		if reqSt.Front {
			url = url + "/Members/Front"
		} else {
			err = required(reqSt.CallSid)
			url = url + "/Members/" + reqSt.CallSid
		}
	case Participants:
		url = url + "/Participants"
	case Participant, UpdateParticipant:
		err = required(reqSt.CallSid)
		url = url + "/Participants/" + reqSt.CallSid
	case UpdateParticipant:
		err = required(reqSt.CallSid)
		url = url + "/Participants/" + reqSt.CallSid
	case DeleteParticipant:
		err = required(reqSt.CallSid)
		url = url + "/Participants/" + reqSt.CallSid
	}
	return url, err
}

// httpRequest creates a http REST request from the supplied request struct
// and the account Sid
func httpRequest(reqStruct interface{}, accountSid string) (
	httpReq *http.Request, err error) {

	url, err := urlString(reqStruct, accountSid)

	if err != nil {
		return httpReq, err
	}

	queryStr := queryString(reqStruct)

	switch reqStruct.(type) {
	// GET query method
	default:
		if queryStr != "" {
			url = url + "?" + queryStr
		}
		httpReq, err = http.NewRequest("GET", url, nil)
	// DELETE query method
	case DeleteNotification, DeleteOutgoingCallerId,
		DeleteRecording, DeleteParticipant, DeleteQueue:
		if queryStr != "" {
			url = url + "?" + queryStr
		}
		httpReq, err = http.NewRequest("DELETE", url, nil)
	// POST query method
	case SendMessage, MakeCall, ModifyCall, CreateQueue, ChangeQueue,
		DeQueue, UpdateParticipant, UpdateOutgoingCallerId,
		AddOutgoingCallerId:
		requestBody := strings.NewReader(queryStr)
		httpReq, err = http.NewRequest("POST", url, requestBody)

	}

	return httpReq, err
}

// queryString constructs the request string by combining struct tags and
// elements from the request struct. Each element string is being url
// encoded/escaped before included.
func queryString(reqSt interface{}) (qryStr string) {
	switch reqSt := reqSt.(type) {
	default:
	case SendMessage, Messages, MakeCall, Calls, ModifyCall,
		Notifications, OutgoingCallerIds, Recordings, Accounts,
		UsageRecords, CreateQueue,
		ChangeQueue, DeQueue, Conferences, Participants:
		for i := 0; i < reflect.ValueOf(reqSt).NumField(); i++ {
			fld := reflect.ValueOf(reqSt).Field(i)
			//fldName := reflect.ValueOf(reqSt).Type().Field(i).Name
			fldType := reflect.ValueOf(reqSt).Type().Field(i).Type
			fldTag := reflect.ValueOf(reqSt).Type().Field(i).Tag

			if fldType.Kind() == reflect.String &&
				string(fldTag) != "" && fld.String() != "" {

				qryStr += string(fldTag) +
					url.QueryEscape(fld.String()) + "&"
			}
		}
		// remove the last '&' if we created a query string
		if len(qryStr) > 0 {
			qryStr = qryStr[:len(qryStr)-1]
		}
	}
	return qryStr
}

// check that string(s) is(are) not empty, return error otherwise
func required(rs ...string) (err error) {
	for _, s := range rs {
		if s == "" {
			return fmt.Errorf("required field missing")
		}
	}
	return
}
