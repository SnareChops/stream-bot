export class Cursor {
    clicked: boolean = false;
    x: number = 0;
    y: number = 0;

    constructor() {
        document.addEventListener('click', event => {
            this.clicked = true;
            this.x = event.clientX;
            this.y = event.clientY;
        });
    }

    update() {
        this.clicked = false;
    }
}