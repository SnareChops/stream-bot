import { Component, createComponent, createElement } from '.';
import { Signal } from './socket';

interface CheerSignal extends Signal {
    username: string;
    amount: number;
    message: string;
}

@Component
export class CheerElement extends HTMLElement {
    static get tag() { return 'twitch-cheer'; }
    static create(signal: Signal): CheerElement {
        const cheer = signal as CheerSignal;
        return createComponent(CheerElement, { username: cheer.username, amount: cheer.amount.toString(), message: cheer.message })
    }

    complete = () => { };

    connectedCallback() {

    }
}