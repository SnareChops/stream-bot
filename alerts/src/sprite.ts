import { RawBounds } from '../../core/bounds';

export class ImageSprite extends RawBounds {
    image: CanvasImageSource;

    constructor(url: string) {
        super(0, 0);
        const self = this;
        this.image = new Image();
        this.image.src = url;
        this.image.addEventListener('load', x => {
            const target = x.target as HTMLImageElement;
            self.setSize(target.naturalWidth, target.naturalWidth);
        });
    }

    Image(): CanvasImageSource {
        return this.image;
    }
}
