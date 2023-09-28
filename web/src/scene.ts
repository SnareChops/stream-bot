import { BOTTOM, CENTER, RIGHT } from '../core/bounds';
import { PhysicsSimulator } from '../core/physics';
import { ParticleEmitter } from '../core/particle';
import { createCanvas } from '../core/util';
import { Renderer } from '../core/renderer';
import { Alert } from './sprites/alert';
import { Cursor } from './cursor';
import { Signal, Socket } from './socket';
import { WallCollider } from './collision';
import { FollowAlert, CheerAlert, SubAlert, ResubAlert, GiftAlert } from './alert';
import { AudioSignal, playAudio } from './audio';
import { Avatar, AvatarSignal } from './sprites/avatar';
import { VideoSignal, playVideo } from './video';

export class AlertScene {
    socket: Socket = new Socket('ws://localhost:3000/ws');
    renderer: Renderer = new Renderer();
    physics: PhysicsSimulator = new PhysicsSimulator();
    emitter: ParticleEmitter = new ParticleEmitter();
    cursor: Cursor = new Cursor();

    avatar: Avatar | undefined;
    walls: WallCollider[] = [];
    queue: Alert[] = [];
    active: Alert[] = [];

    init(): this {
        const [canvas, ctx] = createCanvas(4, 4);
        ctx.fillStyle = 'red';
        ctx.fillRect(0, 0, 4, 4);

        this.emitter = new ParticleEmitter().init(1000, [canvas])
        this.emitter.setLife(100, 200);
        this.emitter.setVelocity(.2, .5);
        this.emitter.setDensity(30);
        this.emitter.setAngle(0, 2 * Math.PI);
        this.renderer.addEffect(this.emitter);

        this.walls.push(new WallCollider(-10, 0, 10, 1080));
        this.walls.push(new WallCollider(0, -10, 1920, 10));
        this.walls.push(new WallCollider(1920, 0, 10, 1080));
        this.walls.push(new WallCollider(0, 1080, 1920, 10));
        for (const collider of this.walls) {
            this.physics.addCollider(collider);
        }
        Avatar.load().then(avatar => {
            this.avatar = avatar;
            this.avatar.setAnchor(CENTER, CENTER);
            const [w, h] = this.avatar.size();
            this.avatar.vector.setVec2(1920 - (w / 2), 1080 - (h / 2));
            console.log(this.avatar.vector.vec2());
            this.renderer.addToScreen(avatar);
        });
        this.socket.addHandler(this.onSignal.bind(this));
        return this;
    }

    onSignal(signal: Signal) {
        switch (signal.kind) {
            case "twitch.follow":
                return this.queue.push(FollowAlert(signal));
            case "twitch.cheer":
                return this.queue.push(CheerAlert(signal));
            case "twitch.sub":
                return this.queue.push(SubAlert(signal));
            case "twitch.resub":
                return this.queue.push(ResubAlert(signal));
            case "twitch.gift":
                return this.queue.push(GiftAlert(signal));
            case "audio":
                return playAudio(signal as AudioSignal);
            case "avatar":
                return this.avatar?.signal(signal as AvatarSignal);
            case "video.small":
                return playVideo(signal as VideoSignal);
        }
    }

    update(delta: number) {
        while (true) {
            // If there is space on the screen, pop the next item off the queue and start
            if (this.queue.length > 0) {
                const pos = findSpace(this.active, this.queue[0], 10, 10);
                if (!pos) break
                const alert = this.queue.shift() as Alert;
                alert.vector.setVec2(...pos);
                alert.start(() => this.active.splice(this.active.indexOf(alert), 1));
                this.active.push(alert);
            } else break;
        }
        if (this.cursor.clicked) {
            this.emitter.bounds.vector.setVec2(this.cursor.x, this.cursor.y);
            this.emitter.start(100);
        }
        this.physics.update(delta);
        this.emitter.update(delta);
        this.avatar?.update(delta);
        for (const alert of this.active) {
            alert.update(delta);
        }
        this.cursor.update();
    }

    draw(ctx: CanvasRenderingContext2D) {
        this.renderer.draw(ctx);
    }
}


function findSpace(active: Alert[], next: Alert, paddingX: number, paddingY: number): [number, number] | undefined {
    const [lastX, lastY] = active.reduce<[number, number]>((last: [number, number], alert: Alert) => {
        const [x, y] = alert.vecOf(RIGHT, BOTTOM);
        // Check for alerts to the right
        if (x > last[0]) return [x, y];
        // Check for alerts below
        if (y > last[1]) return [last[0], y];
        return last;
    }, [0, 0]);
    // Attempt to add to the bottom of the current list
    const [w, h] = next.size()
    // Try to place below last alert
    if (lastY + h + paddingY <= 1080)
        return [lastX + w, lastY + h + paddingY];
    // Try to place the alert at the top to the right of the last alert
    if (lastX + w + paddingX <= 1920)
        return [lastX + w + paddingX, paddingY];
    return void 0;
}