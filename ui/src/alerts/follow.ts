import { Component, createComponent } from '../lib';
import { Signal } from '../lib/socket';
import { FakeTerminal } from '../components/fake-terminal';

interface FollowSignal extends Signal {
    username: string;
}

@Component
export class FollowElement extends HTMLElement {
    static get tag() { return 'twitch-follow'; }
    static create(signal: Signal): FollowElement {
        const follow = createComponent(FollowElement);
        follow.#signal = signal as FollowSignal;
        return follow;
    }

    complete = () => { };
    #signal: FollowSignal | undefined;

    connectedCallback() {
        this.append(FakeTerminal.create(
            `follow -u ${this.#signal?.username}`,
            'Thanks for following!',
            this.complete.bind(this)
        ));
    }
}