import { Component, createComponent, createElement } from '../lib';
import { Signal } from '../lib/socket';
import { Message } from './message';
import { FakeTerminal } from '../components/fake-terminal';


interface ResubSignal extends Signal {
    username: string;
    tier: number;
    months: number;
    message: Message;
}

@Component
export class ResubElement extends HTMLElement {
    static get tag() { return 'twitch-resub'; }
    static create(signal: Signal): ResubElement {
        const resub = createComponent(ResubElement);
        resub.#signal = signal as ResubSignal;
        return resub;
    }

    complete = () => { };
    #signal: ResubSignal | undefined;

    connectedCallback() {
        this.append(FakeTerminal.create(
            `resub -u ${this.#signal?.username} -t ${this.#signal?.tier} -m ${this.#signal?.months}`,
            this.#signal?.message || '',
            this.complete.bind(this),
        ));
    }

    formatMessage(message: Message) {

    }
}