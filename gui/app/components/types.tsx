"use client";

export interface ChatItem {
    from: "user" | "moderator" | "ai" | "system"; // who sent it?
    userEmail: string; // email of the user who sent this message
    data: string; // message to show in chat window
    moderated: boolean;
    approved: boolean;
}

export interface Message {
    kind: string;
    data: string;
}

export interface ModRequest {
    id: string,
    client_id: string,
    user: string,
    user_email: string,
    message: Message,
    approved: boolean,
    moderated: boolean,
}