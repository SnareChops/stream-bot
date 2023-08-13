import { Component } from '.';
import { Socket, Signal } from './socket';
import { FollowElement } from './follow';
import { CheerElement } from './cheer';
import { SubElement } from './sub';
import { ResubElement } from './resub';
import { GiftElement } from './gift';
import { AudioSignal, AudioPlayer } from './components/audio-player'

const socketUrl = `ws://${location.host}/ws`;

interface AlertElement extends HTMLElement {
    complete: () => void;
}

@Component
export class Alerts extends HTMLElement {
    static get tag() { return 'custom-alerts' }
    #socket: Socket | undefined;
    #queue: AlertElement[] = [];
    #active: AlertElement[] = [];
    connectedCallback() {
        this.#socket = new Socket(socketUrl, this.#onSignal.bind(this));
    }

    #onSignal(signal: Signal) {
        console.log('Processing signal', signal);
        switch (signal.kind) {
            case "twitch.follow":
                return this.#enqueue(FollowElement.create(signal));
            case "twitch.cheer":
                return this.#enqueue(CheerElement.create(signal));
            case "twitch.sub":
                return this.#enqueue(SubElement.create(signal));
            case "twitch.resub":
                return this.#enqueue(ResubElement.create(signal));
            case "twitch.gift":
                return this.#enqueue(GiftElement.create(signal));
            case "audio":
                return this.#playAudio(signal as unknown as AudioSignal)
        }
    }

    #playAudio(signal: AudioSignal) {
        this.appendChild(AudioPlayer.create(signal.url));
    }

    #enqueue(component: AlertElement) {
        this.#queue.push(component)
        this.#tick();
    }

    #activate(component: AlertElement) {
        component.complete = this.#complete.bind(this, component);
        this.#queue.splice(this.#queue.indexOf(component), 1);
        this.#active.push(component);
        this.appendChild(component);
    }

    #complete(component: AlertElement) {
        console.log('Removing component', component);
        this.removeChild(component);
        this.#active.splice(this.#active.indexOf(component), 1);
        this.#tick();
    }

    #tick() {
        while (this.#active.length < 20) {
            if (this.#queue.length <= 0) break;
            const next = this.#queue.shift();
            if (!!next) {
                this.#activate(next);
            }
        }
    }
}