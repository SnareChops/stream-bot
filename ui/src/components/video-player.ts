import { Component, createComponent, createElement } from '../lib';

export interface VideoSignal {
    url: string;
}

@Component
export class VideoPlayer extends HTMLElement {
    static get tag() { return 'video-player'; }
    static create(signal: VideoSignal): VideoPlayer {
        const result = createComponent(VideoPlayer);
        result.#signal = signal;
        return result;
    }

    complete = () => { };
    #signal: VideoSignal | undefined;

    connectedCallback() {
        if (!this.#signal) return this.remove();
        const video = createElement<HTMLVideoElement>('video');
        video.setAttribute('width', '320');
        video.setAttribute('height', '240');
        const source = createElement<HTMLSourceElement>('source');
        source.setAttribute('src', this.#signal.url);
        source.setAttribute('type', 'video/mp4');
        video.append(source);
        this.append(video);
        video.addEventListener('ended', this.remove.bind(this))
        video.play();
    }
}