import React, { useState, useEffect } from 'react';
import './EventPlanner.css';

interface Event {
  id: string;
  name: string;
  description: string;
  date: string;
  participants: Participant[];
  vendors: Vendor[];
}

interface Participant {
  id: string;
  name: string;
}

interface Vendor {
  id: string;
  name: string;
  services: string[];
}

const EventPlanner: React.FC = () => {
  const [events, setEvents] = useState<Event[]>([]);
  const [selectedEvent, setSelectedEvent] = useState<Event | null>(null);

  useEffect(() => {
    fetchEvents();
  }, []);

  const fetchEvents = async () => {
    const response = await fetch(`${process.env.REACT_APP_API_URL}/events`);
    const data = await response.json();
    setEvents(data);
  };

  const selectEvent = (eventId: string) => {
    const event = events.find(e => e.id === eventId);
    setSelectedEvent(event);
  };

  return (
    <div className="event-planner-dashboard">
      <div className="events-sidebar">
        <h2>Events</h2>
        <ul>
          {events.map(event => (
            <li key={event.id} onClick={() => selectEvent(event.id)}>
              {event.name}
            </li>
          ))}
        </ul>
      </div>
      <div className="event-details">
        {selectedEvent ? (
          <div>
            <h2>{selectedEvent.name}</h2>
            <p>{selectedEvent.description}</p>
            <p><strong>Date:</strong> {selectedEvent.date}</p>
            <h3>Participants</h3>
            <ul>
              {selectedEvent.participants.map(participant => (
                <li key={participant.id}>{participant.name}</li>
              ))}
            </ul>
            <h3>Vendors</h3>
            <ul>
              {selectedEvent.vendors.map(vendor => (
                <li key={vendor.id}>{vendor.name} - Services: {vendor.services.join(', ')}</li>
              ))}
            </ul>
          </div>
        ) : (
          <p>Please select an event to see the details.</p>
        )}
      </div>
    </div>
  );
};

export default EventPlanner;