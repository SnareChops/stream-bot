import { chatterJoin, hydrateChatters, chatterLeave, chatterMessage, MessageSignal, JoinSignal, LeaveSignal } from './chatters';
import { Socket } from '../../core/socket';

const socket = new Socket(`ws://${location.host}/admin/ws`);
socket.addHandler(signal => {
    switch (signal.kind) {
        case "twitch.join":
            chatterJoin(signal as JoinSignal);
        case "twitch.leave":
            chatterLeave(signal as LeaveSignal);
        case "twitch.message":
            chatterMessage(signal as MessageSignal)
    }
});

hydrateChatters();