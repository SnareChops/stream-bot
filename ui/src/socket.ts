
export interface Signal {
    kind: string;
}

export type SignalHandler = (message: Signal) => any;

export class Socket extends WebSocket {
    #onSignal: SignalHandler;
    constructor(url: string, onSignal: SignalHandler) {
        super(url);
        this.#onSignal = onSignal;
        this.onopen = () => console.log('Socket connection open');
        this.onclose = () => console.log('Socket connection closed');
        this.onerror = (err) => console.error('Socket connection error', err);
        this.onmessage = this.read.bind(this);
    }

    read(event: MessageEvent<string>) {
        console.log('Socket message received', event);
        try {
            const message = JSON.parse(event.data) as Signal;
            this.#onSignal(message);
        } catch (err) {
            console.log('Failed to parse message data', err);
        }
    }
}