package services

import (
	"crashsaver/websocket/data"
	"encoding/json"
	"log"
)

type Incident struct {
	// incidents *data.IncidentList  // No necesitas tener la referencia directa a IncidentList aquí
	manager *Manager
	// sync.RWMutex
}

func NewIncidents(manager *Manager) *Incident {
	return &Incident{
		manager: manager,
	}
}

// HandleWebSocketMessage maneja los mensajes WebSocket relacionados con incidentes.
func (is *Incident) HandleWebSocketMessage(client *Client, payload data.WebSocketMessage) {
	// Obtiene la lista de incidentes desde el Manager
	incidentList := is.manager.incidents

	switch payload.Type {
	case "create":
		var incident data.Incident
		err := json.Unmarshal(payload.Data, &incident)

		if err != nil {
			log.Printf("Error al decodificar el incidente: %v\n", err)
			return
		}

		//valid not null or format?
		//
		incidentList.AddIncident(incident)
		// Enviar notificación a todos los clientes

		is.emitIncidentToClients("created", incident)

	case "update":
		var incident data.Incident
		err := json.Unmarshal(payload.Data, &incident)
		if err != nil {
			log.Printf("Error al decodificar el incidente: %v\n", err)
			return
		}
		incidentList.UpdateIncident(incident)
		//send notification
		// Enviar el incidente creado a todos los clientes
		is.emitIncidentToClients("updated", incident)

	case "delete":
		log.Printf("delete incident %s\n", payload.Data)
		var incident data.Incident
		err := json.Unmarshal(payload.Data, &incident)
		if err != nil {
			log.Printf("Error al decodificar el incidente: %v\n", err)
			return
		}
		incidentList.RemoveIncident(incident.ID)

		//send deleted incidented
		// is.emitIncidentToClients("delete", incident)

		//send list completed incidented
		strIncidents, err := is.handleGetIncidents(client, incidentList)
		if err != nil {
			// Manejar el error aquí
			log.Printf("Error al obtener la lista de incidentes: %v\n", err)
			return
		}
		//send a one client all
		is.manager.sendMessageToAllClients(strIncidents)
	case "get":
		//get list incidents

		strIncidents, err := is.handleGetIncidents(client, incidentList)
		if err != nil {
			// Manejar el error aquí
			log.Printf("Error al obtener la lista de incidentes: %v\n", err)
			return
		}
		//send a one client all
		is.manager.sendMessageToClient(client, strIncidents)
	}
}

func (is *Incident) handleGetIncidents(client *Client, incidentList *data.IncidentList) (string, error) {
	// Obtener la lista de incidentes
	allIncidents := incidentList.GetIncidents()

	// Convertir la lista de incidentes a JSON
	incidentsJSON, err := json.Marshal(allIncidents)
	if err != nil {
		log.Printf("Error al convertir la lista de incidentes a JSON: %v\n", err)
		return "", err
	}

	// Crear un mensaje WebSocket con el tipo de acción y la lista de incidentes
	message := data.WebSocketMessage{
		Type: "listIncidents",
		Data: json.RawMessage(incidentsJSON),
	}

	// Convertir el mensaje a JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al convertir el mensaje a JSON: %v\n", err)
		return "", err
	}

	return string(messageJSON), nil
	// Enviar el mensaje al cliente actual

}

// emitIncidentToClients envía un incidente a todos los clientes.
func (is *Incident) emitIncidentToClients(event string, incident data.Incident) {
	// Convertir el incidente a JSON
	incidentJSON, err := json.Marshal(incident)
	if err != nil {
		log.Printf("Error al convertir el incidente a JSON: %v\n", err)
		return
	}

	// Crear un mensaje WebSocket con el tipo de acción y los datos del incidente
	message := data.WebSocketMessage{
		Type: event,
		Data: json.RawMessage(incidentJSON),
	}

	// Convertir el mensaje a JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al convertir el mensaje a JSON: %v\n", err)
		return
	}

	// Enviar el mensaje a todos los clientes
	is.manager.sendMessageToAllClients(string(messageJSON))
}
