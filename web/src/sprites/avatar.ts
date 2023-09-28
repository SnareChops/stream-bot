import { RawBounds } from '../../core/bounds';
import { createCanvas, isSet } from '../../core/util';
import { Signal } from '../socket';

export enum AvatarState {
    MOUTH_CLOSED = 1 << 0,
    MOUTH_OPEN = 1 << 1,
    EYES_CLOSED = 1 << 2,
    EYES_OPEN = 1 << 3,
    WHEEZE = 1 << 4,
    POG = 1 << 5,
}

export interface AvatarSignal extends Signal {
    type: 'talking';
    talking: boolean;
}

export class Avatar extends RawBounds {
    state: number = AvatarState.MOUTH_CLOSED | AvatarState.EYES_OPEN;
    blink: number = 0;
    static async load(): Promise<Avatar> {
        const images = await Promise.all([
            loadImage("assets/gremlin/idle.png"),
            loadImage("assets/gremlin/talk.png"),
            loadImage("assets/gremlin/xd.png"),
            loadImage("assets/gremlin/xdTalk.png"),
            loadImage("assets/gremlin/wheeze.png"),
            loadImage("assets/gremlin/pog.png"),
        ]);
        const scaled = images.map(x => scaleImage(x, 250, 335));
        console.log(scaled);
        const map = new Map<number, CanvasImageSource>([
            [AvatarState.MOUTH_CLOSED | AvatarState.EYES_OPEN, scaled[0]],
            [AvatarState.MOUTH_OPEN | AvatarState.EYES_OPEN, scaled[1]],
            [AvatarState.MOUTH_CLOSED | AvatarState.EYES_CLOSED, scaled[2]],
            [AvatarState.MOUTH_OPEN | AvatarState.EYES_CLOSED, scaled[3]],
            [AvatarState.WHEEZE, scaled[4]],
            [AvatarState.POG, scaled[5]],
        ]);
        const { width, height } = scaled.values().next().value;
        return new Avatar(width, height, map);
    }

    constructor(width: number, height: number, public images: Map<number, CanvasImageSource>) {
        super(width, height);
    }

    signal(signal: AvatarSignal) {
        switch (signal.type) {
            case 'talking':
                if (signal.talking) {
                    this.state &= ~AvatarState.MOUTH_CLOSED;
                    this.state |= AvatarState.MOUTH_OPEN;
                } else {
                    this.state &= ~AvatarState.MOUTH_OPEN;
                    this.state |= AvatarState.MOUTH_CLOSED;
                }
                return
        }
    }

    update(delta: number) {
        this.blink -= delta;
        if (this.blink <= 0) {
            if (isSet(this.state, AvatarState.EYES_CLOSED)) {
                this.state &= ~AvatarState.EYES_CLOSED;
                this.state |= AvatarState.EYES_OPEN;
                this.blink = Math.random() * 19000 + 100;
            } else {
                this.state &= ~AvatarState.EYES_OPEN;
                this.state |= AvatarState.EYES_CLOSED;
                this.blink = Math.random() * 750 + 50;
            }
        }
    }

    Image(): CanvasImageSource | undefined {
        if (isSet(this.state, AvatarState.WHEEZE)) return this.images.get(AvatarState.WHEEZE);
        if (isSet(this.state, AvatarState.POG)) return this.images.get(AvatarState.POG);
        return this.images.get(this.state);
    }
}

async function loadImage(url: string): Promise<HTMLImageElement> {
    return new Promise<HTMLImageElement>(resolve => {
        const image = new Image();
        image.src = url;
        image.addEventListener('load', function () { resolve(this); });
    });
}

function scaleImage(image: HTMLImageElement, maxWidth: number, maxHeight: number): CanvasImageSource {
    const scaleWidth = maxWidth / image.width;
    const scaleHeight = maxHeight / image.height;
    const scale = Math.min(scaleWidth, scaleHeight);
    const width = image.width * scale;
    const height = image.height * scale;
    const [canvas, context] = createCanvas(width, height);
    context.drawImage(image, 0, 0, width, height);
    return canvas;
}