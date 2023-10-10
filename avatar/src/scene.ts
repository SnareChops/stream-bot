import { Signal, Socket } from '../../core/socket';
import { Avatar, AvatarSignal } from './avatar';

export class AvatarScene {
    socket: Socket = new Socket(`ws://${location.host}/avatar/ws`);
    constructor(public avatar: Avatar) {
        this.socket.addHandler(this.signal.bind(this));
    }

    signal(signal: Signal) {
        switch (signal.kind) {
            case "avatar":
                return this.avatar.signal(signal as AvatarSignal);
        }
    }

    update(delta: number) {
        this.avatar.update(delta);
    }

    draw(screen: CanvasRenderingContext2D) {
        screen.drawImage(this.avatar.Image()!!, ...this.avatar.vector.vec2())
    }
}