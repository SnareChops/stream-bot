import { Signal } from '../../core/socket';

export interface AudioSignal extends Signal {
    url: string;
}

export async function playAudio(signal: AudioSignal) {
    const audio = document.createElement('audio');
    audio.setAttribute('src', signal.url);
    audio.setAttribute('autoplay', 'true');
    audio.setAttribute('preload', 'auto');
    audio.addEventListener('ended', function () { this.remove(); }, true);
    document.body.append(audio);
    try {
        await audio.play();
    } catch (err: any) {
        console.error('Audio failed to play', signal, err);
    }
}