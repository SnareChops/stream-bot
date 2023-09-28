import { createElement } from '../lib';

export class CanvasControl {
    static create(): [CanvasControl, HTMLCanvasElement] {
        const control = new CanvasControl();
        control.#canvas = createElement<HTMLCanvasElement>('canvas', { width: '1920', height: '1080' });
        const ctx = control.#canvas.getContext('2d');
        if (!ctx) {
            console.error("2d context failed to load");
            return [control, control.#canvas];
        }
        control.#ctx = ctx;
        return [control, control.#canvas];
    }

    #canvas: HTMLCanvasElement;
    #ctx: CanvasRenderingContext2D;


}