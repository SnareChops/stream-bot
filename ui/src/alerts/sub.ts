import { Component, createComponent } from '../lib';
import { Signal } from '../lib/socket';
import { FakeTerminal } from '../components/fake-terminal';

interface SubSignal extends Signal {
    username: string;
    tier: number;
    gift: boolean;
}

@Component
export class SubElement extends HTMLElement {
    static get tag() { return 'twitch-sub'; }
    static create(signal: Signal): SubElement {
        const sub = createComponent(SubElement);
        sub.#signal = signal as SubSignal;
        return sub;
    }

    complete = () => { };
    #signal: SubSignal | undefined;

    connectedCallback() {
        this.append(FakeTerminal.create(
            `subscribe -u ${this.#signal?.username} -t ${this.#signal?.tier}`,
            'Thanks for subscribing!',
            this.complete.bind(this)
        ));
    }
}