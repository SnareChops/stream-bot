import { createElement } from '.';
export interface Message {
    text: string;
    emotes: {
        id: string;
        begin: number;
        end: number;
    }[];
}

export function formatMessage(message: Message): HTMLParagraphElement {
    const p = createElement<HTMLParagraphElement>('p');
    let cursor = 0;
    for (const emote of message.emotes) {
        const span = createElement('span');
        span.innerText = message.text.slice(cursor, emote.begin);
        const img = createElement('img', {
            href: `https://static-cdn.jtvnw.net/emoticons/v2/emotesv2_${emote.id}/default/dark/3.0`
        });
        p.append(span, img);
        cursor = emote.end
    }
    const span = createElement('span');
    span.innerText = message.text.slice(cursor, message.text.length);
    return p;
}
