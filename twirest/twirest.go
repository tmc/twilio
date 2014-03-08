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

// httpRequest creates a http REST request from the supplied request struct
// and the account Sid
func httpRequest(reqStruct interface{}, accSid string) (
	httpReq *http.Request, err error) {

	url := "https://api.twilio.com/" + ApiVer + "/Accounts"

	switch reqStruct := reqStruct.(type) {
	default:
		err = fmt.Errorf("invalid type %T in Request", reqStruct)
	case SendMessage, Messages:
		url = url + "/" + accSid + "/Messages"
	case Message:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Messages/" + reqStruct.Sid
		if reqStruct.Media == true {
			url = url + "/Media"
			if reqStruct.MediaSid != "" {
				url = url + "/" + reqStruct.MediaSid
			}
		}
	case MakeCall, Calls:
		url = url + "/" + accSid + "/Calls"
	case Call:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Calls/" + reqStruct.Sid
		if reqStruct.Recordings == true {
			url = url + "/Recordings"
		} else if reqStruct.Notifications == true {
			url = url + "/Notifications"
		}
	case ModifyCall:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Calls/" + reqStruct.Sid
	case Accounts:
		url = url
	case Account:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid
	case Notifications:
		url = url + "/" + accSid + "/Notifications"
	case Notification:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Notifications/" + reqStruct.Sid
	case DeleteNotification:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Notifications/" + reqStruct.Sid
	case Recordings:
		url = url + "/" + accSid + "/Recordings"
	case Recording:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Recordings/" + reqStruct.Sid
	case DeleteRecording:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Recordings/" + reqStruct.Sid
	case UsageRecords:
		url = url + "/" + accSid + "/Usage/Records/" +
			reqStruct.SubResource
	case Queues, CreateQueue:
		url = url + "/" + accSid + "/Queues"
	case Queue:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Queues/" + reqStruct.Sid
	case ChangeQueue:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Queues/" + reqStruct.Sid
	case DeleteQueue:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Queues/" + reqStruct.Sid
	case QueueMembers:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Queues/" + reqStruct.Sid +
			"/Members"
	case QueueMember:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Queues/" + reqStruct.Sid
		if reqStruct.Front {
			url = url + "/Members/Front"
		} else {
			err = required(reqStruct.Sid, reqStruct.CallSid)
			url = url + "/Members/" + reqStruct.CallSid
		}
	case DeQueue:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Queues/" + reqStruct.Sid +
			"/Members"
		if reqStruct.Front {
			url = url + "/Front"
		} else {
			err = required(reqStruct.Sid, reqStruct.CallSid)
			url = url + "/" + reqStruct.CallSid
		}
	case Conferences:
		url = url + "/" + accSid + "/Conferences"
	case Conference:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Conferences/" + reqStruct.Sid
	case Participants:
		err = required(reqStruct.Sid)
		url = url + "/" + accSid + "/Conferences/" + reqStruct.Sid +
			"/Participants"
	case Participant:
		err = required(reqStruct.Sid, reqStruct.CallSid)
		url = url + "/" + accSid + "/Conferences/" + reqStruct.Sid +
			"/Participants/" + reqStruct.CallSid
	case UpdateParticipant:
		err = required(reqStruct.Sid, reqStruct.CallSid)
		url = url + "/" + accSid + "/Conferences/" + reqStruct.Sid +
			"/Participants/" + reqStruct.CallSid
	case DeleteParticipant:
		err = required(reqStruct.Sid, reqStruct.CallSid)
		url = url + "/" + accSid + "/Conferences/" + reqStruct.Sid +
			"/Participants/" + reqStruct.CallSid
	}
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
	case DeleteNotification, DeleteRecording, DeleteParticipant,
		DeleteQueue:
		if queryStr != "" {
			url = url + "?" + queryStr
		}
		httpReq, err = http.NewRequest("DELETE", url, nil)
	// POST query method
	case SendMessage, MakeCall, ModifyCall, CreateQueue, ChangeQueue,
		DeQueue, UpdateParticipant:
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
		Notifications, Recordings, Accounts, UsageRecords, CreateQueue,
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
