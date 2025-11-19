package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (api *APIClient) GetPath(pathName string) (*APIPath, error) {
	reqCtx, reqC := context.WithTimeout(api.context, api.httpTimeout)
	defer reqC()

	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", fmt.Sprintf("%s/v3/paths/get/%s", api.baseURL, pathName),
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

	var successResponse *APIPath
	var errorResponse   *APIError

	switch resp.StatusCode {
	case http.StatusOK:
		err := json.Unmarshal(body, &successResponse)
		if err != nil {
			return nil, err
		}
	case http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError:
		err := json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, err
		} else {
			return nil, fmt.Errorf("code: %s, message: %s", resp.StatusCode, errorResponse.Error)
		}
	default:
		return nil, fmt.Errorf("unhandled status code: %s", resp.StatusCode)
	}

	return successResponse, nil
}
