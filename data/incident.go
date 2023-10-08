package data

import (
	"sync"
)

// type IncidentList struct {
// 	incidents map[string]Incident
// 	sync.RWMutex
// }

// Cambia a un slice para mayor coherencia
type IncidentList struct {
	incidents []*Incident
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

// // Singleton para IncidentList
// var singletonIncidentList *IncidentList
// var once sync.Once

func NewIncidentsList() *IncidentList {
	return &IncidentList{
		incidents: make([]*Incident, 0),
	}
}

// GetIncidents devuelve todos los incidentes en la lista
func (il *IncidentList) GetIncidents() []*Incident {
	il.RLock()
	defer il.RUnlock()

	// Crear una copia de la lista para evitar la modificación no segura en concurrencia
	incidentsCopy := make([]*Incident, len(il.incidents))
	copy(incidentsCopy, il.incidents)

	return incidentsCopy
}

// GetIncident devuelve un incidente específico según el ID
func (il *IncidentList) GetIncident(incidentID string) *Incident {
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

// AddIncident agrega un nuevo incidente a la lista
func (il *IncidentList) AddIncident(incident Incident) {
	il.Lock()
	defer il.Unlock()

	// Crear una copia del incidente para evitar la modificación no segura en concurrencia
	incidentCopy := &Incident{
		ID:             incident.ID,
		UserId:         incident.UserId,
		Description:    incident.Description,
		Latitude:       incident.Latitude,
		Longitude:      incident.Longitude,
		Status:         incident.Status,
		IncidentTypeId: incident.IncidentTypeId,
	}

	il.incidents = append(il.incidents, incidentCopy)
}

// RemoveIncident elimina un incidente de la lista según el ID
func (il *IncidentList) RemoveIncident(incidentID string) {
	il.Lock()
	defer il.Unlock()

	var updatedIncidents []*Incident

	for _, incident := range il.incidents {
		if incident.ID != incidentID {
			// Conservar solo los incidentes que no coinciden con el ID proporcionado
			incidentCopy := &Incident{
				ID:             incident.ID,
				UserId:         incident.UserId,
				Description:    incident.Description,
				Latitude:       incident.Latitude,
				Longitude:      incident.Longitude,
				Status:         incident.Status,
				IncidentTypeId: incident.IncidentTypeId,
			}
			updatedIncidents = append(updatedIncidents, incidentCopy)
		}
	}

	il.incidents = updatedIncidents
}

func (il *IncidentList) UpdateIncident(incident Incident) {
	il.Lock()
	defer il.Unlock()

	for _, existingIncident := range il.incidents {
		if existingIncident.ID == incident.ID {
			// Actualiza los campos necesarios del incidente existente
			existingIncident.Description = incident.Description
			existingIncident.Latitude = incident.Latitude
			existingIncident.Longitude = incident.Longitude
			existingIncident.Status = incident.Status
			existingIncident.IncidentTypeId = incident.IncidentTypeId

			// Rompe el bucle después de encontrar el incidente
			break
		}
	}
}

// ClearIncidents borra todos los incidentes de la lista
func (il *IncidentList) ClearIncidents() {
	il.Lock()
	defer il.Unlock()

	il.incidents = nil
}
