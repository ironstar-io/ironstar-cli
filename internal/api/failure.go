package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
)

func (err *APIError) Error(output string) {
	if strings.ToLower(output) == "json" {
		utils.PrintInterfaceAsJSON(err)
		return
	}

	fmt.Println()

	switch err.StatusCode {
	case 400:
		color.Red("Ironstar API call failed! (400: Bad Request)")
	case 401:
		color.Red("Ironstar API call failed! (401: Unauthorized)")
	case 403:
		color.Red("Ironstar API call failed! (403: Forbidden)")
	case 404:
		color.Red("Ironstar API call failed! (404: Not Found)")
	case 405:
		color.Red("Ironstar API call failed! (405: Method Not Allowed)")
	case 500:
		color.Red("Ironstar API call failed! (500: Server Error)")
		fmt.Println()
		color.Yellow("Please contact Ironstar Support - support@ironstar.io")
	default:
		color.Red("Ironstar API call failed!")
	}

	fmt.Println()
	if err.StatusCode != 0 {
		fmt.Printf("Status Code: %+v\n", err.StatusCode)
	}

	if err.IronstarCode != "" {
		fmt.Println("Ironstar Code: " + err.IronstarCode)
	}

	if err.CorrelationId != "" {
		fmt.Println("Correlation ID: " + err.CorrelationId)
	}

	if err.Message != "" {
		fmt.Println("Message: " + err.Message)
	}
}

var ErrIronstarAPICall = errors.New("Ironstar API call was unsuccessful!")

func (res *RawResponse) HandleFailure(output string) error {
	var apiErr APIError
	var correlationId string

	if res.Header["X-Correlation-Id"] != nil && res.Header["X-Correlation-Id"][0] != "" {
		correlationId = res.Header["X-Correlation-Id"][0]
	}

	if res.StatusCode == 500 {
		apiErr = APIError{
			StatusCode:    res.StatusCode,
			IronstarCode:  "INTERNAL_SERVER_ERROR",
			CorrelationId: correlationId,
			Message:       "An unexpected error occurred in the Ironstar API server.",
		}
	} else {
		f := &FailureBody{}
		err := json.Unmarshal(res.Body, f)
		if err != nil {
			if strings.ToLower(output) != "json" {
				fmt.Println()
				fmt.Println("An unexpected error occurred")
				fmt.Println()
				fmt.Println("The server responded with status code " + strconv.Itoa(res.StatusCode))
				if res.Body != nil {
					fmt.Println(string(res.Body))
				}
			}
			return errors.New("Unable to read server response body")
		}

		apiErr = APIError{
			StatusCode:    res.StatusCode,
			IronstarCode:  f.Code,
			CorrelationId: correlationId,
			Message:       f.Message,
		}
	}

	apiErr.Result = "error"
	apiErr.Error(output)

	return ErrIronstarAPICall
}

func (err *APIError) ExternalError() {
	fmt.Println()

	switch err.StatusCode {
	case 400:
		color.Red("External API call failed! (Bad Request)")
	case 401:
		color.Red("External API call failed! (Unauthorized)")
	case 403:
		color.Red("External API call failed! (Forbidden)")
	case 404:
		color.Red("External API call failed! (Not Found)")
	case 500:
		color.Red("External API call failed! (Server Error)")
	default:
		color.Red("External API call failed!")
	}

	fmt.Println()
	fmt.Printf("Status Code: %+v\n", err.StatusCode)
	fmt.Printf("Message: %+v\n", err.Message)
	fmt.Printf("Method: %+v\n", err.CallMethod)
	fmt.Printf("URL: %+v\n", err.CallURL)
}

var ErrExternalAPICall = errors.New("External API call was unsuccessful!")

func (res *RawResponse) HandleExternalFailure() error {
	var apiErr APIError
	var correlationId string

	if res.Header["X-Correlation-Id"] != nil && res.Header["X-Correlation-Id"][0] != "" {
		correlationId = res.Header["X-Correlation-Id"][0]
	}

	if res.StatusCode == 500 {
		apiErr = APIError{
			StatusCode:    res.StatusCode,
			IronstarCode:  "INTERNAL_SERVER_ERROR",
			CorrelationId: correlationId,
			Message:       "An unexpected error occurred reaching an external API server.",
		}
	} else {
		f := &FailureBody{}
		err := json.Unmarshal(res.Body, f)
		if err != nil {
			return err
		}

		apiErr = APIError{
			StatusCode:    res.StatusCode,
			CallURL:       res.CallURL,
			CallMethod:    res.CallMethod,
			CorrelationId: correlationId,
			Message:       f.Message,
		}
	}

	apiErr.ExternalError()

	return ErrExternalAPICall
}
