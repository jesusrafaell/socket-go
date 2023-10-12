package data

import (
	"sync"
)

type IncidentList struct {
	incidents []*Incident
	// 	incidents map[string]Incident
	sync.RWMutex
}

type Incident struct {
	ID             string  `json:"_id"`
	UserId         string  `json:"user_id"`
	Description    string  `json:"description"`
	Latitude       float32 `json:"lat"`
	Longitude      float32 `json:"lon"`
	Status         int     `json:"status"`
	IncidentTypeId byte    `json:"incident_type_id"`
}

func NewIncidents() *IncidentList {
	return &IncidentList{
		incidents: make([]*Incident, 0),
	}
}

// GetIncidents (ALL)
func (il *IncidentList) GetIncidents() []*Incident {
	il.RLock()
	defer il.RUnlock()

	// create copy list for not current
	incidentsCopy := make([]*Incident, len(il.incidents))
	copy(incidentsCopy, il.incidents)

	return incidentsCopy
}

func (il *IncidentList) GetIncidentById(incidentID string) *Incident {
	il.RLock()
	defer il.RUnlock()

	for _, incident := range il.incidents {
		if incident.ID == incidentID {
			// Devolver una copia del incidente para evitar la modificación no segura en concurrencia
			return &Incident{
				ID:             incident.ID,
				UserId:         incident.UserId,
				Description:    incident.Description,
				Latitude:       incident.Latitude,
				Longitude:      incident.Longitude,
				Status:         incident.Status,
				IncidentTypeId: incident.IncidentTypeId,
			}
		}
	}

	return nil
}

func (il *IncidentList) AddIncident(incident Incident) {
	il.Lock()
	defer il.Unlock()

	il.incidents = append(il.incidents, &incident)
}

func (il *IncidentList) UpdateIncident(incident Incident) {
	il.Lock()
	defer il.Unlock()

	for _, existingIncident := range il.incidents {
		if existingIncident.ID == incident.ID {
			existingIncident.Description = incident.Description
			existingIncident.Latitude = incident.Latitude
			existingIncident.Longitude = incident.Longitude
			existingIncident.Status = incident.Status
			existingIncident.IncidentTypeId = incident.IncidentTypeId

			break
		}
	}
}

func (il *IncidentList) RemoveIncident(incidentID string) {
	il.Lock()
	defer il.Unlock()
	var indexToRemove int = -1

	// index from element delete
	for i, incident := range il.incidents {
		if incident.ID == incidentID {
			indexToRemove = i
			break
		}
	}

	// delete
	if indexToRemove != -1 {
		il.incidents = append(il.incidents[:indexToRemove], il.incidents[indexToRemove+1:]...)
	}
}

func (il *IncidentList) ClearIncidents() {
	il.Lock()
	defer il.Unlock()

	il.incidents = nil
}
