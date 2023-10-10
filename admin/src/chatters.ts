import { Signal } from '../../core/socket';

export interface Chatter {
    name: string;
    time: number;
    element: HTMLLIElement;
}

let element: Element | undefined;
const chatters: Chatter[] = [];
const list = document.createElement('ul');
let interval: Number;

export function hydrateChatters() {
    element = document.querySelector('chatters-panel') || void 0;
    while (element?.firstChild) {
        element.firstChild.remove();
    }
    element?.append(list);
    interval = setInterval(refresh, 800);
}

export interface JoinSignal extends Signal {
    username: string;
    time: number;
}

export function chatterJoin(signal: JoinSignal) {
    console.log('chatter join', signal);
    for (const chatter of chatters) {
        if (chatter.name == signal.username) return;
    }
    const li = document.createElement('li');
    const span = document.createElement('span');
    li.innerText = signal.username;
    li.append(span);
    list?.append(li);
    chatters.push({
        name: signal.username,
        time: signal.time,
        element: li,
    });
}

export interface LeaveSignal extends Signal {
    username: string;
}

export function chatterLeave(signal: LeaveSignal) {
    for (const chatter of chatters) {
        if (chatter.name == signal.username) {
            chatter.element.remove();
            chatters.splice(chatters.indexOf(chatter), 1);
            return
        }
    }
}

export interface MessageSignal extends Signal {
    username: string;
}

export function chatterMessage(signal: MessageSignal) {
    for (const chatter of chatters) {
        if (chatter.name == signal.username) {
            const span = chatter.element.querySelector('span');
            span!!.innerText = timeSince(chatter.time);
            return;
        }
    }
    chatterJoin({ ...signal, time: new Date().getTime() });
}

function refresh() {
    for (const chatter of chatters) {
        const span = chatter.element.querySelector('span');
        span!!.innerText = timeSince(chatter.time);
    }
}

const SECOND = 1000;
const MINUTE = 60 * SECOND;
const HOUR = 60 * MINUTE;

function timeSince(time: number): string {
    const since = (new Date()).getTime() - time;
    if (since > HOUR) return `${Math.floor(since / HOUR)}h`;
    if (since > MINUTE) return `${Math.floor(since / MINUTE)}m`;
    return `${Math.floor(since / SECOND)}s`;
}