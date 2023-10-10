import { createElement, Component } from '../lib';

@Component
export class QuickPanel extends HTMLElement {
    static get tag(): string { return 'quick-panel'; }

    connectedCallback() {
        const button1 = createElement('button');
        button1.innerText = 'Button 1';
        const button2 = createElement('button');
        button2.innerText = 'Button 2';
        const button3 = createElement('button');
        button3.innerText = 'Button 3';
        const button4 = createElement('button');
        button4.innerText = 'Button 4';
        this.append(
            button1,
            button2,
            button3,
            button4,
        );
    }
}