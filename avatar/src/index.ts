import { Engine } from '../../core/engine';
import { createCanvas } from '../../core/util';
import { AvatarScene } from './scene';
import { Avatar } from './avatar';

const body = document.querySelector('body') as HTMLBodyElement;

Avatar.load().then(avatar => {
    const [w, h] = avatar.size();
    const [canvas, context] = createCanvas(w, h);
    body.style.width = w + 'px';
    body.style.height = h + 'px';
    body.append(canvas);
    new Engine(context).runGame(new AvatarScene(avatar));
});