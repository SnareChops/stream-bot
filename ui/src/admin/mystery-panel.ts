import { createElement, createComponent, Component } from '../lib';
import { Signal, Socket } from '../lib/socket';

interface MysteryListSignal {
    kind: string;
    items: MysteryItem[];
}

interface MysteryItem {
    selected: boolean;
    name: string;
    type: string;
    start: number;
    end: number;
    speed: number;
}

@Component
export class MysteryPanel extends HTMLElement {
    static get tag(): string { return 'mystery-panel'; }
    static create(socket: Socket): MysteryPanel {
        const el = createComponent(MysteryPanel);
        el.#socket = socket;
        return el;
    }

    #socket: Socket;
    connectedCallback() {
        this.#socket.addHandler(this.#onSignal.bind(this));
        this.#socket.send({ kind: 'admin.mystery.request' })
    }

    select(name: string) {
        this.#socket.send({ kind: 'admin.mystery.select', args: [name] });
    }

    #onSignal(signal: Signal) {
        switch (signal.kind) {
            case 'admin.mystery.load':
                return this.#loadItems(signal as MysteryListSignal);
            default:
                return
        }
    }

    #loadItems(signal: MysteryListSignal) {
        for (const child of this.children) {
            child.remove();
        }
        for (const item of signal.items) {
            this.append(MysteryListItem.create(this, item))
        }
    }
}

@Component
export class MysteryListItem extends HTMLElement {
    static get tag(): string { return 'mystery-list-item'; }
    static create(panel: MysteryPanel, item: MysteryItem): MysteryListItem {
        const el = createComponent(MysteryListItem);
        el.#panel = panel;
        el.#item = item;
        return el;
    }

    #panel: MysteryPanel;
    #checkbox: HTMLInputElement;
    #item: MysteryItem;
    connectedCallback() {
        this.#checkbox = createElement<HTMLInputElement>('input', { type: 'checkbox' });
        this.#checkbox.checked = this.#item.selected;
        this.#checkbox.onchange = this.#changed.bind(this);
        const nameSpan = createElement('span');
        nameSpan.innerText = this.#item.name;
        const rangeSpan = createElement('span');
        rangeSpan.innerText = this.#item.start + '-' + this.#item.end;
        const speedSpan = createElement('span');
        speedSpan.innerText = '' + this.#item.speed;

        this.append(
            this.#checkbox,
            nameSpan,
            rangeSpan,
            speedSpan,
        );
    }

    #changed() {
        if (this.#checkbox.checked !== this.#item.selected) {
            this.#item.selected = this.#checkbox.checked;
            this.#panel.select(this.#item.name);
        }
    }
}