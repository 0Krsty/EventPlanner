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

// Separate EventListItem into its own component 
const EventListItem = ({ event, onSelect }: { event: Event; onSelect: (id: string) => void }) => (
  <li onClick={() => onSelect(event.id)}>
    {event.name}
  </li>
);

// Separate ParticipantDetails into its own component
const ParticipantDetails = ({ participants }: { participants: Participant[] }) => (
  <ul>
    {participants.map(participant => (
      <li key={participant.id}>{participant.name}</li>
    ))}
  </ul>
);

// Separate VendorDetails into its own component
const VendorDetails = ({ vendors }: { vendors: Vendor[] }) => (
  <ul>
    {vendors.map(vendor => (
      <li key={vendor.id}>{vendor.name} - Services: {vendor.services.join(', ')}</li>
    ))}
  </ul>
);

const EventPlanner: React.FC = () => {
  const [events, setEvents] = useState<Event[]>([]);
  const [selectedEvent, setSelectedEvent] = useState<Event | null>(null);

  // Fetch events once component is mounted
  useEffect(() => {
    const fetchEvents = async () => {
      const response = await fetch(`${process.env.REACT_APP_API_URL}/events`);
      const data = await response.json();
      setEvents(data);
    };
    fetchEvents();
  }, []);

  // Event selection handler
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
            <EventListItem key={event.id} event={event} onSelect={selectEvent} />
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
            <ParticipantDetails participants={selectedEvent.participants} />
            <h3>Vendors</h3>
            <VendorDetails vendors={selectedEvent.vendors} />
          </div>
        ) : (
          <p>Please select an event to see the details.</p>
        )}
      </div>
    </div>
  );
};

export default EventPlanner;