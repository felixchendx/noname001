package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (api *APIClient) GetGlobalConfiguration() (*Conf, error) {
	reqCtx, reqC := context.WithTimeout(api.context, api.httpTimeout)
	defer reqC()

	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", api.baseURL + "/v3/config/global/get",
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var successResponse *Conf
	var errorResponse   *APIError

	switch resp.StatusCode {
	case http.StatusOK:
		err := json.Unmarshal(body, &successResponse)
		if err != nil {
			return nil, err
		}
	case http.StatusBadRequest, http.StatusInternalServerError:
		err := json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, err
		} else {
			return nil, fmt.Errorf("code: %s, message: %s", resp.StatusCode, errorResponse.Error)
		}
	case http.StatusNotFound:
		return nil, fmt.Errorf("code: 404, data not found")
	default:
		return nil, fmt.Errorf("unhandled status code: %s", resp.StatusCode)
	}

	return successResponse, nil
}

func (api *APIClient) AddPathConfiguration(pathName string, pathConfig *SimplePathConfiguration) (error) {
	reqCtx, reqC := context.WithTimeout(api.context, api.httpTimeout)
	defer reqC()

	reqBytes, _ := json.Marshal(pathConfig)
	req, err := http.NewRequestWithContext(
		reqCtx,
		"POST", fmt.Sprintf("%s/v3/config/paths/add/%s", api.baseURL, pathName),
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var errorResponse *APIError

	switch resp.StatusCode {
	case http.StatusOK:
		// no response body, only 200 indicating success
	case http.StatusBadRequest, http.StatusInternalServerError:
		err := json.Unmarshal(body, &errorResponse)
		if err != nil {
			return err
		} else {
			return fmt.Errorf("code: %s, message: %s", resp.StatusCode, errorResponse.Error)
		}
	case http.StatusNotFound:
		return fmt.Errorf("code: 404, data not found")
	default:
		return fmt.Errorf("unhandled status code: %s", resp.StatusCode)
	}

	return nil
}

// not implemented, if necessary, do implement
// func (api *APIClient) PatchPathConfiguration

func (api *APIClient) ReplacePathConfiguration(pathName string, pathConfig *SimplePathConfiguration) (error) {
	reqCtx, reqC := context.WithTimeout(api.context, api.httpTimeout)
	defer reqC()

	reqBytes, _ := json.Marshal(pathConfig)
	req, err := http.NewRequestWithContext(
		reqCtx,
		"POST", fmt.Sprintf("%s/v3/config/paths/replace/%s", api.baseURL, pathName),
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var errorResponse *APIError

	switch resp.StatusCode {
	case http.StatusOK:
		// no response body, only 200 indicating success
	case http.StatusBadRequest, http.StatusInternalServerError:
		err := json.Unmarshal(body, &errorResponse)
		if err != nil {
			return err
		} else {
			return fmt.Errorf("code: %s, message: %s", resp.StatusCode, errorResponse.Error)
		}
	case http.StatusNotFound:
		return fmt.Errorf("code: 404, data not found")
	default:
		return fmt.Errorf("unhandled status code: %s", resp.StatusCode)
	}

	return nil
}

func (api *APIClient) DeletePathConfiguration(pathName string) (error) {
	reqCtx, reqC := context.WithTimeout(api.context, api.httpTimeout)
	defer reqC()

	req, err := http.NewRequestWithContext(
		reqCtx,
		"DELETE", fmt.Sprintf("%s/v3/config/paths/delete/%s", api.baseURL, pathName),
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var errorResponse *APIError

	switch resp.StatusCode {
	case http.StatusOK:
		// no response body, only 200 indicating success
	case http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError:
		err := json.Unmarshal(body, &errorResponse)
		if err != nil {
			return err
		} else {
			return fmt.Errorf("code: %s, message: %s", resp.StatusCode, errorResponse.Error)
		}
	default:
		return fmt.Errorf("unhandled status code: %s", resp.StatusCode)
	}

	return nil
}
