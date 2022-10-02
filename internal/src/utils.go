package src

import (
	"encoding/json"
	"os"
	"strings"

	uuid "github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func LogInfo(location string, event string, message ...string) {
	mess := ""
	for _, msg := range message {
		mess = mess + msg
	}
	Logger.WithFields(logrus.Fields{
		"location": location,
		"event":    event,
	}).Info(mess)
}

func LogError(location string, event string, message ...string) {
	errorId := uuid.New()
	mess := ""
	for _, msg := range message {
		mess = mess + msg
	}
	Logger.WithFields(logrus.Fields{
		"id":       errorId,
		"location": location,
		"event":    event,
	}).Error(mess)
}

func LogFatal(location string, event string, message ...string) {
	errorId := uuid.New()
	mess := ""
	for _, msg := range message {
		mess = mess + msg
	}
	Logger.WithFields(logrus.Fields{
		"id":       errorId,
		"location": location,
		"event":    event,
	}).Fatal(mess)
}

func GetEnv(key, fallback string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		return fallback
	}
	return value
}

func CheckError(err error, location string, event string, message ...string) {
	if err != nil {
		message = append(message, ":", err.Error())
		LogError(location, event, message...)
	}
}

func Contains(a string, b string) bool {
	return strings.Contains(
		strings.ToLower(a),
		strings.ToLower(b),
	)
}

func FilterInputMap(refObj interface{}, inputMap map[string]interface{}) (map[string]any, error) {
	var refMap map[string]interface{}
	ref, err := json.Marshal(refObj)
	if err != nil {
		CheckError(err, "Input filtering", "Error during marshalling input reference", err.Error())
		return nil, err
	}
	json.Unmarshal(ref, &refMap)
	resMap := map[string]any{}
	for k, v := range refMap {

		if inputMap[k] != v && inputMap[k] != nil && k != "id" {
			resMap[k] = inputMap[k]
		}
	}
	return resMap, nil
}
