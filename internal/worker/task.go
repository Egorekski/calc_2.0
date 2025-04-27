package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"distributed_calculator/internal/constants"

	"distributed_calculator/internal/app/models"

	"go.uber.org/zap"
)

func (a *Agent) getTask() (*models.Task, error) {
	resp, err := a.httpClient.Get(fmt.Sprintf(constants.PathInternalTask, a.config.OrchestratorURL))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			a.logger.Error(constants.ErrFailedCloseRespBody, zap.Error(err))
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(constants.ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var taskResp models.TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return nil, err
	}

	return &taskResp.Task, nil
}

func (a *Agent) sendResult(taskID string, result float64) error {
	taskResult := models.TaskResult{
		ID:     taskID,
		Result: result,
	}

	body, err := json.Marshal(taskResult)
	if err != nil {
		return err
	}

	resp, err := a.httpClient.Post(
		fmt.Sprintf(constants.PathInternalTask, a.config.OrchestratorURL),
		constants.ContentTypeJSON,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			a.logger.Error(constants.ErrFailedCloseRespBody, zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(constants.ErrUnexpectedStatusCode, resp.StatusCode)
	}

	return nil
}
