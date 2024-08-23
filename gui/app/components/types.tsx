"use client";

export interface Message {
    kind: string;
    data: string;
}

export interface ModRequest {
    id: string,
    client_id: string,
    user_email: string,
    message: Message,
    approved: boolean,
    moderated: boolean,
}