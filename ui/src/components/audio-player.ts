import { Component, createComponent, createElement } from '../lib';

export interface AudioSignal {
    url: string;
}

@Component
export class AudioPlayer extends HTMLElement {
    static get tag() { return 'audio-player'; }
    static create(url: string): AudioPlayer {
        const result = createComponent(AudioPlayer);
        result.#url = url;
        return result;
    }

    #url: string = '';

    connectedCallback() {
        const audio = createElement<HTMLAudioElement>('audio');
        audio.setAttribute('src', this.#url);
        audio.setAttribute('autoplay', 'true');
        audio.setAttribute('preload', 'auto');
        audio.addEventListener('ended', this.remove.bind(this), true);
        this.appendChild(audio);
        try {
            audio.play()
        } catch (err: any) {
            console.log('Audio failed to play', err);
        }
    }
}