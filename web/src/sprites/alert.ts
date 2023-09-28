import { Bounds, RawBounds } from '../../core/bounds';
import { drawStringBlock, Font15, imagesForString, Letter, wrap } from '../../core/font';
import { Sprite } from '../../core/renderer';
import { createCanvas } from '../../core/util';

enum AlertStatus {
    PENDING,
    ANIMATING,
    WAITING,
    DONE,
}

export interface Alert extends Sprite {
    start(done: () => void);
    update(delta: number);
}

export class AlertBox extends RawBounds {
    box: HTMLImageElement;
    canvas: HTMLCanvasElement;
    context: CanvasRenderingContext2D;
    cooldown: number = 0;
    status: AlertStatus = AlertStatus.PENDING;
    index: number = 0;
    lines: Letter[][];
    kerning: number = 2;
    leading: number = 4;
    textBounds: Bounds = new RawBounds(377, 119);
    done: () => void;

    constructor(private text: string, private rate: number, private wait: number) {
        super(0, 0);
        this.textBounds.vector.setVec2(22, 15);
        this.box = new Image();
        this.box.src = 'assets/AlertBox.png';
        this.box.addEventListener('load', (x => {
            const target = x.target as HTMLImageElement;
            this.setSize(target.naturalWidth, target.naturalHeight);
            [this.canvas, this.context] = createCanvas(...this.size());
            this.render();
        }).bind(this));
        imagesForString(this.text, Font15, [1, 249, 255, 255]).then(lines => {
            this.lines = wrap(this.textBounds, this.kerning, this.leading, lines);
            this.render();
        });
    }

    start(done: () => void) {
        this.cooldown = this.rate;
        this.index = 0;
        this.status = AlertStatus.ANIMATING;
        this.done = done;
    }

    update(delta: number) {
        if (this.status === AlertStatus.ANIMATING) {
            this.cooldown -= delta;
            if (this.cooldown < 1) {
                this.cooldown = this.rate + this.cooldown;
                this.index += 1;
                if (this.index === this.text.length) {
                    this.status = AlertStatus.WAITING;
                }
                this.render();
                return;
            }
        }
        if (this.status === AlertStatus.WAITING) {
            this.wait -= delta;
            if (this.wait < 1) {
                this.status = AlertStatus.DONE;
                this.done();
            }
        }
    }

    Image(): CanvasImageSource {
        return this.canvas;
    }

    render() {
        if (!this.context || !this.lines) return;
        this.context.reset();
        this.context.drawImage(this.box, 0, 0);
        if (this.index === 0) return;
        let count = 0;
        const linesToDraw: Letter[][] = [];
        for (const line of this.lines) {
            if (line.length + count < this.index) {
                count += line.length;
                linesToDraw.push(line);
            } else {
                const end = this.index - count;
                linesToDraw.push(line.slice(0, end));
                break;
            }
        }
        drawStringBlock(this.context, this.textBounds, this.kerning, this.leading, linesToDraw);
    }
}