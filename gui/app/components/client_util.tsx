"use client";

import { ModRequest } from "./types";

export const sendMessage = (
    ws: WebSocket,
    id: string,
    client_id: string,
    user_email: string,
    message_kind: string,
    message_data: string,
    approved: boolean,
    moderated: boolean
) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
        const msg: ModRequest = {
            id: id,
            client_id: client_id,
            user_email: user_email,
            message: {
                kind: message_kind,
                data: message_data,
            },
            approved: approved,
            moderated: moderated,
        };
        ws.send(JSON.stringify(msg));
        return msg;
    } else {
        console.error("WebSocket is not open");
        return null;
    }
};

const ClientUtil = {
    sendMessage,
};

export default ClientUtil;
