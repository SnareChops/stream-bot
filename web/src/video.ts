import { Signal } from './socket';

export interface VideoSignal extends Signal {
    url: string;
}

export async function playVideo(signal: VideoSignal) {
    const video = document.createElement('video');
    video.setAttribute('width', '320');
    video.setAttribute('height', '240');
    const source = document.createElement('source');
    source.setAttribute('src', signal.url);
    source.setAttribute('type', 'video/mp4');
    video.append(source);
    document.body.append(video);
    video.addEventListener('ended', function () { this.remove(); });
    try {
        await video.play();
    } catch (err: any) {
        console.error('Video failed to play', signal, err);
    }
}