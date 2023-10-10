import { Avatar } from '../../avatar/src/avatar';
import { Engine } from '../../core/engine';
import { initFonts } from '../../core/font';
import { Game } from './game';
import { AlertScene } from './scene';

const $body = document.querySelector<HTMLBodyElement>('body') as HTMLBodyElement;
const $canvas = document.createElement('canvas') as HTMLCanvasElement;
const ctx = $canvas.getContext('2d') as CanvasRenderingContext2D;

$canvas.width = 1920;
$canvas.height = 1080;

$body.append($canvas);
initFonts().then(() => {
    new Engine(ctx).runGame(new Game(new AlertScene()));
});