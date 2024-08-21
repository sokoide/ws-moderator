"use client";

export interface ChatItem {
    direction: "received" | "sent"; // direction of the message
    userEmail: string; // email of the user who sent this message
    data: string; // message to show in chat window
    moderated: boolean;
    approved: boolean;
}
