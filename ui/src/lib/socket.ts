export interface Signal {
    kind: string;
    args?: string[];
}

export type SignalHandler = (message: Signal) => any;

export class Socket {
    #url: string;
    #socket: WebSocket | undefined;
    #handlers: SignalHandler[] = [];
    #interval: number = 1000;
    #maxInterval: number = 30000;
    #decay: number = 2;

    constructor(url: string) {
        this.#url = url;
        this.#connect();
    }

    addHandler(handler: SignalHandler) {
        this.#handlers.push(handler);
    }

    send(signal: Signal) {
        this.#socket?.send(JSON.stringify(signal));
    }

    #connect() {
        this.#socket = new WebSocket(this.#url);
        this.#socket.onopen = () => {
            console.log('Socket connection open');
            this.#interval = 1000;
        }
        this.#socket.onclose = event => {
            console.log('Socket connection closed', event);
            this.#reconnect();
        }
        this.#socket.onerror = err => console.error('Socket connection error', err);
        this.#socket.onmessage = event => {
            console.log('Socket message received', event);
            try {
                const message = JSON.parse(event.data) as Signal;
                for (const handler of this.#handlers) {
                    handler(message);
                }
            } catch (err) {
                console.log('Failed to parse message data', err);
            }
        }
    }

    #reconnect() {
        setTimeout(() => {
            console.log(`Socket reconnecting... (delay: ${this.#interval}ms)`);
            this.#connect();
            this.#interval = Math.min(
                this.#interval * this.#decay,
                this.#maxInterval
            );
        }, this.#interval);
    }
}