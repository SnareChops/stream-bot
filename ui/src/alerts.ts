import { createComponent, Component } from './lib';
import { Socket, Signal } from './lib/socket';
import { FollowElement } from './alerts/follow';
import { CheerElement } from './alerts/cheer';
import { SubElement } from './alerts/sub';
import { ResubElement } from './alerts/resub';
import { GiftElement } from './alerts/gift';
import { AudioSignal, AudioPlayer } from './components/audio-player';
import { VideoSignal, VideoPlayer } from './components/video-player';

const socketUrl = `ws://${location.host}/ui/ws`;

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
        this.#socket = new Socket(socketUrl);
        this.#socket.addHandler(this.#onSignal.bind(this));
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
                return this.append(AudioPlayer.create((signal as unknown as AudioSignal).url));
            case "video.small":
                return this.#enqueue(VideoPlayer.create(signal as unknown as VideoSignal));
        }
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

document.querySelector('body')?.appendChild(createComponent(Alerts));