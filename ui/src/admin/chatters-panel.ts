import { createElement, Component } from '../lib';

@Component
export class ChattersPanel extends HTMLElement {
    static get tag(): string { return 'chatters-panel'; }

    #list: HTMLUListElement;
    connectedCallback() {
        this.#list = createElement<HTMLUListElement>('ul');
        this.append(this.#list);
    }
}