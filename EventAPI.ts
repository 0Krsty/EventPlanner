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

class Cache<T> {
    private cache: Map<string, T> = new Map<string, T>();

    get(key: string): T | undefined {
        return this.cache.get(key);
    }

    set(key: string, value: T): void {
        this.cache.set(key, value);
    }

    clear(): void {
        this.cache.clear();
    }
}

class EventsModule {
    private readonly baseUrl: string;
    private eventsCache = new Cache<Event[]>();

    constructor() {
        this.baseUrl = process.env.API_BASE_URL ?? 'http://localhost:3000';
    }

    private async makeRequest<T>(url: string, method: 'get' | 'post', data?: any): Promise<T> {
        try {
            const response = await axios({ url, method, data });
            if (url.includes("/events")) {
                this.eventsCache.clear(); // Invalidate cache as data might have changed
            }
            console.log(`Operation successful for: ${url}`, response.data);

            // Cache if a GET request to /events
            if (method === 'get' && url === `${this.baseUrl}/events`) {
                this.eventsCache.set(url, response.data);
            }

            return response.data;
        } catch (error) {
            console.error(`Error during ${method.toUpperCase()} request to ${url}:`, error);
            throw error;
        }
    }

    async createEvent(event: Event): Promise<Event> {
        return this.makeRequest<Event>(`${this.baseUrl}/events`, 'post', event);
    }

    async registerParticipant(eventId: string, participant: Participant): Promise<Participant> {
        return this.makeRequest<Participant>(`${this.baseUrl}/events/${eventId}/participants`, 'post', participant);
    }

    async addVendor(eventId: string, vendor: Vendor): Promise<Vendor> {
        return this.makeRequest<Vendor>(`${this.baseUrl}/events/${eventId}/vendors`, 'post', vendor);
    }

    async listEvents(): Promise<Event[]> {
        const cacheKey = `${this.baseUrl}/events`;
        const cachedEvents = this.eventsCache.get(cacheKey);
        if (cachedEvents) {
            return Promise.resolve(cachedEvents);
        }
        return this.makeRequest<Event[]>(cacheKey, 'get');
    }
}

export const eventsModule = new EventsModule();