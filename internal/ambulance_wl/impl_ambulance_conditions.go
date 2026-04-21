package ambulance_wl

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

type implAmbulanceConditionsAPI struct {
}

func NewAmbulanceConditionsApi() AmbulanceConditionsAPI {
	return &implAmbulanceConditionsAPI{}
}

// GetConditions - Provides the list of conditions associated with ambulance
func (o implAmbulanceConditionsAPI) GetConditions(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		result := ambulance.PredefinedConditions
		if result == nil {
			result = []Condition{}
		}
		return nil, result, http.StatusOK
	})
}

// GetCondition - Provides details about a specific predefined condition
func (o implAmbulanceConditionsAPI) GetCondition(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		conditionCode := c.Param("conditionCode")
		idx := slices.IndexFunc(ambulance.PredefinedConditions, func(cond Condition) bool {
			return cond.Code == conditionCode
		})
		if idx < 0 {
			return nil, gin.H{"status": "Not Found", "message": "Condition not found"}, http.StatusNotFound
		}
		return nil, ambulance.PredefinedConditions[idx], http.StatusOK
	})
}

// CreateCondition - Saves new condition into the predefined conditions list
func (o implAmbulanceConditionsAPI) CreateCondition(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var condition Condition
		if err := c.ShouldBindJSON(&condition); err != nil {
			return nil, gin.H{"status": "Bad Request", "message": "Invalid request body", "error": err.Error()}, http.StatusBadRequest
		}
		if condition.Value == "" {
			return nil, gin.H{"status": "Bad Request", "message": "Condition value is required"}, http.StatusBadRequest
		}
		if condition.Code != "" {
			conflictIdx := slices.IndexFunc(ambulance.PredefinedConditions, func(cond Condition) bool {
				return cond.Code == condition.Code
			})
			if conflictIdx >= 0 {
				return nil, gin.H{"status": "Conflict", "message": "Condition with this code already exists"}, http.StatusConflict
			}
		}
		ambulance.PredefinedConditions = append(ambulance.PredefinedConditions, condition)
		return ambulance, condition, http.StatusOK
	})
}

// UpdateCondition - Updates specific predefined condition
func (o implAmbulanceConditionsAPI) UpdateCondition(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		conditionCode := c.Param("conditionCode")
		var condition Condition
		if err := c.ShouldBindJSON(&condition); err != nil {
			return nil, gin.H{"status": "Bad Request", "message": "Invalid request body", "error": err.Error()}, http.StatusBadRequest
		}
		idx := slices.IndexFunc(ambulance.PredefinedConditions, func(cond Condition) bool {
			return cond.Code == conditionCode
		})
		if idx < 0 {
			return nil, gin.H{"status": "Not Found", "message": "Condition not found"}, http.StatusNotFound
		}
		if condition.Value != "" {
			ambulance.PredefinedConditions[idx].Value = condition.Value
		}
		if condition.Code != "" {
			ambulance.PredefinedConditions[idx].Code = condition.Code
		}
		if condition.Reference != "" {
			ambulance.PredefinedConditions[idx].Reference = condition.Reference
		}
		if condition.TypicalDurationMinutes > 0 {
			ambulance.PredefinedConditions[idx].TypicalDurationMinutes = condition.TypicalDurationMinutes
		}
		return ambulance, ambulance.PredefinedConditions[idx], http.StatusOK
	})
}

// DeleteCondition - Deletes specific predefined condition
func (o implAmbulanceConditionsAPI) DeleteCondition(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		conditionCode := c.Param("conditionCode")
		idx := slices.IndexFunc(ambulance.PredefinedConditions, func(cond Condition) bool {
			return cond.Code == conditionCode
		})
		if idx < 0 {
			return nil, gin.H{"status": "Not Found", "message": "Condition not found"}, http.StatusNotFound
		}
		ambulance.PredefinedConditions = append(ambulance.PredefinedConditions[:idx], ambulance.PredefinedConditions[idx+1:]...)
		return ambulance, nil, http.StatusNoContent
	})
}
