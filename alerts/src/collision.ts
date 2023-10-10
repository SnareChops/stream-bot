import { RawBounds } from '../../core/bounds';

export enum CollisionType {
    NONE = 0,
    WALL = 1 << 0,
    ENTITY = 1 << 1,
}

export class WallCollider extends RawBounds {
    constructor(x: number, y: number, w: number, h: number) {
        super(w, h);
        this.vector.setVec2(x, y);
    }

    type(): number {
        return CollisionType.WALL;
    }
}