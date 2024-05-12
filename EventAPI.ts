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

    private async postRequest<T>(url: string, data: T): Promise<T> {
        const response = await axios.post(url, data);
        console.log(`Operation successfully for: ${url}`, response.data);
        return response.data;
    }

    private async getRequest<T>(url: string): Promise<T> {
        const response = await axios.get(url);
        console.log(`Operation successfully for: ${url}`, response.data);
        return response.data;
    }

    async createEvent(event: Event): Promise<Event> {
        return this.postRequest(`${this.baseUrl}/events`, event)
            .catch(error => {
                console.error('Error creating event:', error);
                throw error;
            });
    }

    async registerParticipant(eventId: string, participant: Participant): Promise<Participant> {
        return this.postRequest(`${this.baseUrl}/events/${eventId}/participants`, participant)
            .catch(error => {
                console.error('Error registering participant:', error);
                throw error;
            });
    }

    async addVendor(eventId: string, vendor: Vendor): Promise<Vendor> {
        return this.postRequest(`${this.baseUrl}/events/${eventId}/vendors`, vendor)
            .catch(error => {
                console.error('Error adding vendor:', error);
                throw error;
            });
    }

    async listEvents(): Promise<Event[]> {
        return this.getRequest<Event[]>(`${this.baseUrl}/events`)
            .catch(error => {
                console.error('Error listing events:', error);
                throw error;
            });
    }
}

export const eventsModule = new EventsModule();