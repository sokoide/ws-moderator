"use client";

export interface ChatItem {
    direction: "received" | "sent"; // direction of the message
    data: string; // message to show in chat window
    moderated: boolean;
    approved: boolean;
}
