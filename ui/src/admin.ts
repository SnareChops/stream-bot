import { ChattersPanel } from './admin/chatters-panel';
import { QuickPanel } from './admin/quick-panel';
import { MysteryPanel } from './admin/mystery-panel';
import { createComponent, Component } from './lib';
import { Socket } from './lib/socket';

const socketUrl = `ws://${location.host}/adminws`;

@Component
export class AdminPanel extends HTMLElement {
    static get tag(): string { return 'admin-panel'; }

    #socket: Socket;
    connectedCallback() {
        this.#socket = new Socket(socketUrl)
        this.append(
            createComponent(QuickPanel),
            createComponent(ChattersPanel),
            MysteryPanel.create(this.#socket),
        );
    }
}

document.querySelector('body')?.appendChild(createComponent(AdminPanel))