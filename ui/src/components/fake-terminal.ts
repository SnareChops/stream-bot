import { Component, createComponent, createElement } from '../lib';
import { Message, formatMessage } from '../alerts/message';

@Component
export class FakeTerminal extends HTMLElement {
    static get tag() { return 'fake-terminal'; }
    static create(command: string, response: string | Message, done: () => void): FakeTerminal {
        const result = createComponent(FakeTerminal);
        result.#command = command;
        result.#response = response;
        result.#done = done;
        return result;
    }

    #command: string = '';
    #response: string | Message = '';
    #done = () => { };

    connectedCallback() {
        this.append(
            createComponent(FakeTerminalHeader),
            FakeTerminalBody.create(this.#command || '', this.#response, this.#done),
        )
    }
}

@Component
export class FakeTerminalBody extends HTMLElement {
    static get tag() { return 'fake-terminal-body'; }
    static create(command: string, response: string | Message, done: () => void): FakeTerminalBody {
        const result = createComponent(FakeTerminalBody);
        result.#command = command;
        result.#response = response;
        result.#done = done;
        return result;
    }

    #$command: HTMLSpanElement = createElement('span');
    #command: string = '';
    #response: string | Message = '';
    #cooldown: number = 0;
    #index: number = 0;
    #done = () => { };

    connectedCallback() {
        const p = createElement('p')
        p.append(
            createComponent(FakeTerminalLeader),
            this.#$command,
        )
        this.append(p);
        requestAnimationFrame(this.tick.bind(this));
    }

    tick() {
        this.#cooldown--;
        if (this.#cooldown < 1) {
            this.#cooldown = 4;
            this.#index += 1;
            if (this.#index > this.#command.length) {
                let p1;
                if (typeof this.#response === 'string') {
                    p1 = createElement('p');
                    p1.innerText = this.#response;
                } else {
                    p1 = formatMessage(this.#response as Message)
                }
                const p2 = createElement('p');
                p2.append(createComponent(FakeTerminalLeader), createComponent(FakeCursor));
                this.append(p1, p2);
                setTimeout(this.#done.bind(this), 5000);
                return // Escape the tick loop
            } else {
                this.#$command.innerText = this.#command.substring(0, this.#index);
            }
        }
        requestAnimationFrame(this.tick.bind(this));
    }
}

@Component
export class FakeTerminalLeader extends HTMLElement {
    static get tag() { return 'fake-terminal-leader'; }

    connectedCallback() {
        const user = createElement('span');
        user.innerText = 'snarechops';
        const tilde = createElement('span');
        tilde.innerText = '~';
        const dollar = createElement('span');
        dollar.innerText = '$';
        this.append(user, tilde, dollar, ' ');
    }
}

@Component
export class FakeTerminalHeader extends HTMLElement {
    static get tag() { return 'fake-terminal-header'; }

    connectedCallback() {
        this.append(createComponent(FakeButtons))
    }
}

@Component
export class FakeButtons extends HTMLElement {
    static get tag() { return 'fake-buttons'; }

    connectedCallback() {
        this.append(
            createElement('fake-close'),
            createElement('fake-minimize'),
            createElement('fake-zoom'),
        );
    }
}

@Component
export class FakeCursor extends HTMLElement {
    static get tag() { return 'fake-cursor'; }
}