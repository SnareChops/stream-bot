import { Sprite } from '../core/renderer';
import { Signal } from './socket';
import { AlertBox } from './sprites/alert';

interface FollowSignal extends Signal {
    username: string;
}

export function FollowAlert(signal: Signal) {
    const follow = signal as FollowSignal;
    return new AlertBox(`${follow.username} followed!`, 50, 5000);
}

interface CheerSignal extends Signal {
    username: string;
    amount: number;
    message: string;
}

export function CheerAlert(signal: Signal) {
    const cheer = signal as CheerSignal;
    return new AlertBox(`${cheer.username} wasted ${cheer.amount} bits on this dumb little box${!!cheer.message ? ` to say:\n${cheer.message}` : '.'}`, 50, 5000);
}

interface SubSignal extends Signal {
    username: string;
    tier: number;
    gift: boolean;
}

export function SubAlert(signal: Signal) {
    const sub = signal as SubSignal;
    return new AlertBox(`${sub.username} made a tier ${sub.tier} mistake wasting money on emotes.`, 50, 5000);
}

interface Message {
    text: string;
    emotes: {
        id: string;
        begin: number;
        end: number;
    }[];
}

interface ResubSignal extends Signal {
    username: string;
    tier: number;
    months: number;
    message: Message;
}

export function ResubAlert(signal: Signal) {
    const resub = signal as ResubSignal;
    return new AlertBox(`${resub.username} has wasted tier ${resub.tier} money for ${resub.months} months.`, 50, 5000);
}

interface GiftSignal extends Signal {
    username: string;
    tier: number;
    amount: number;
    total: number;
    anon: boolean;
}

export function GiftAlert(signal: Signal) {
    const gift = signal as GiftSignal;
    return new AlertBox(`${gift.anon ? 'Anonymous' : gift.username} gifted a tier ${gift.tier} subscription to `, 50, 5000);
}