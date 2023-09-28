import { PhysicsBounds } from '../../core/bounds';
import { CollisionType } from '../collision';
import { createCanvas } from '../../core/util';

export class TestSprite extends PhysicsBounds {
    image: CanvasImageSource;

    constructor() {
        super(50, 50);
        const [canvas, context] = createCanvas(50, 50);
        this.image = canvas;
        context.fillStyle = 'blue';
        context.fillRect(0, 0, 50, 50);
    }

    Image(): CanvasImageSource {
        return this.image;
    }

    type(): number {
        return CollisionType.ENTITY;
    }

    collision(): number {
        return CollisionType.WALL;
    }
}