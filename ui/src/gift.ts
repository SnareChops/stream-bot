import { Component, createComponent, createElement } from '.';
import { Signal } from './socket';
import { FakeTerminal } from './components/fake-terminal';

interface GiftSignal extends Signal {
    username: string;
    tier: number;
    amount: number;
    total: number;
    anon: boolean;
}

@Component
export class GiftElement extends HTMLElement {
    static get tag() { return 'twitch-gift'; }
    static create(signal: Signal): GiftElement {
        const gift = createComponent(GiftElement);
        gift.#signal = signal as GiftSignal;
        return gift;
    }

    #signal: GiftSignal | undefined;
    complete = () => { };

    connectedCallback() {
        this.append(FakeTerminal.create(
            `giftsubs -u ${this.#signal?.anon === true ? 'Anonymous' : this.#signal?.username} -t ${this.#signal?.tier} -n ${this.#signal?.amount}`,
            !!this.#signal?.total
                ? `They've gifted ${this.#signal?.total}`
                : `Thanks for the gifted subs!`,
            this.complete.bind(this),
        ));
    }
}

function omit<T extends object, K extends keyof T>(obj: T, prop: K): Omit<T, K> {
    const { [prop]: _, ...result } = obj;
    return result;
}