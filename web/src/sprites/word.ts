import { RawBounds } from '../../core/bounds';
import { createCanvas } from '../../core/util';
import { Font10, imagesForString, drawString } from '../../core/font';

export class Word extends RawBounds {
    image: CanvasImageSource;

    constructor() {
        super(100, 100)
        imagesForString("A", Font10, [255, 0, 0, 255]).then(lines => {
            const [canvas, context] = createCanvas(this.dx(), this.dy());
            for (let i = 0; i < lines.length; i++) {
                drawString(context, 0, i * (4 + lines[i][0].image.height), 2, lines[i]);
            }
            this.image = canvas;
        });
    }

    Image(): CanvasImageSource {
        return this.image;
    }
}