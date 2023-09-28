import { Scene } from '../core/scene';
export class Game {
    constructor(public scene: Scene) { }

    update(delta: number) {
        this.scene.update(delta);
    }

    draw(context: CanvasRenderingContext2D) {
        this.scene.draw(context);
    }

    loadScene(scene: Scene) {
        this.scene = scene;
    }
}