package services

import (
	"crashsaver/websocket/data"
	"crashsaver/websocket/pkg/socket"
	"encoding/json"
	"log"
	"sync"
)

type IncidentService struct {
	list    *data.IncidentList
	manager *socket.Manager
	sync.RWMutex
}

func NewIncidentService(incidentList *data.IncidentList, manager *socket.Manager) *IncidentService {
	return &IncidentService{
		list:    incidentList,
		manager: manager,
	}
}

func (i *IncidentService) HandleWebSocketMessage(payload data.WebSocketMessage) {
	switch payload.Type {
	case "create":
		i.createIncident(payload)
		break
	case "update":
		i.updateIncident(payload)
		break
	case "delete":
		i.deleteIncident(payload)
		break
	case "get":
		//get list incidents
		incidents := i.GetIncidents()

		strIncidents, err := i.listIncidentsToString(incidents)

		if err != nil {
			// Manejar el error aquí
			log.Printf("Error al obtener la lista de incidentes: %v\n", err)
			return
		}

		log.Print(strIncidents)

		// client.WriteMessage(strIncidents)
		break
	}
}

func (i *IncidentService) createIncident(payload data.WebSocketMessage) *data.Incident {
	var incident data.Incident
	err := json.Unmarshal(payload.Data, &incident)

	if err != nil {
		log.Printf("Error al decodificar el incidente: %v\n", err)
		return nil
	}

	i.list.AddIncident(incident)

	i.emitIncidentToClients("created", &incident)
	return &incident
}

func (i *IncidentService) GetIncidents() []*data.Incident {
	allIncidents := i.list.GetIncidents()

	return allIncidents
}

func (i *IncidentService) updateIncident(payload data.WebSocketMessage) *data.Incident {
	var incident data.Incident
	err := json.Unmarshal(payload.Data, &incident)
	if err != nil {
		log.Printf("Error al decodificar el incidente: %v\n", err)
		return nil
	}
	i.list.UpdateIncident(incident)

	i.emitIncidentToClients("updated", &incident)
	return &incident
}

func (i *IncidentService) deleteIncident(payload data.WebSocketMessage) *data.Incident {
	log.Printf("delete incident %s\n", payload.Data)
	var incident data.Incident
	err := json.Unmarshal(payload.Data, &incident)
	if err != nil {
		log.Printf("Error al decodificar el incidente: %v\n", err)
		return nil
	}

	i.list.RemoveIncident(incident.ID)

	i.emitIncidentToClients("deleted", &incident)
	return &incident
}

// parse list incident to msg
func (i *IncidentService) listIncidentsToString(incidents []*data.Incident) (string, error) {

	// convert to json
	incidentsJSON, err := json.Marshal(incidents)
	if err != nil {
		log.Printf("Error al convertir la lista de incidentes a JSON: %v\n", err)
		return "", err
	}

	// create format msg
	message := data.WebSocketMessage{
		Type: "listIncidents",
		Data: json.RawMessage(incidentsJSON),
	}

	// conver to json
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al convertir el mensaje a JSON: %v\n", err)
		return "", err
	}

	return string(messageJSON), nil
}

// emitIncidentToClients send to all clients
func (i *IncidentService) emitIncidentToClients(event string, incident *data.Incident) {
	incidentJSON, err := json.Marshal(incident)
	if err != nil {
		log.Printf("Error al convertir el incidente a JSON: %v\n", err)
		return
	}

	message := data.WebSocketMessage{
		Type: event,
		Data: json.RawMessage(incidentJSON),
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al convertir el mensaje a JSON: %v\n", err)
		return
	}

	log.Println(messageJSON)

	i.manager.MessageToAllClients(string(messageJSON))
}
