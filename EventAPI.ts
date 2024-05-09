import axios from 'axios';

interface Event {
    id: string;
    title: string;
    description: string;
    date: Date;
}

interface Participant {
    id: string;
    name: string;
    email: string;
}

interface Vendor {
    id: string;
    name: string;
    serviceProvided: string;
}

class EventsModule {
    private readonly baseUrl: string;

    constructor() {
        this.baseUrl = process.env.API_BASE_URL ?? 'http://localhost:3000';
    }

    async createEvent(event: Event): Promise<Event> {
        try {
            const response = await axios.post(`${this.baseUrl}/events`, event);
            console.log('Event created successfully:', response.data);
            return response.data;
        } catch (error) {
            console.error('Error creating event:', error);
            throw error;
        }
    }

    async registerParticipant(eventId: string, participant: Participant): Promise<Participant> {
        try {
            const response = await axios.post(`${this.baseUrl}/events/${eventId}/participants`, participant);
            console.log('Participant registered successfully:', response.data);
            return response.data;
        } catch (error) {
            console.error('Error registering participant:', error);
            throw error;
        }
    }

    async addVendor(eventId: string, vendor: Vendor): Promise<Vendor> {
        try {
            const response = await axios.post(`${this.baseUrl}/events/${eventId}/vendors`, vendor);
            console.log('Vendor added successfully:', response.data);
            return response.data;
        } catch (error) {
            console.error('Error adding vendor:', error);
            throw error;
        }
    }

    async listEvents(): Promise<Event[]> {
        try {
            const response = await axios.get(`${this.baseUrl}/events`);
            return response.data;
        } catch (error) {
            console.error('Error listing events:', error);
            throw error;
        }
    }
}

export const eventsModule = new EventsModule();