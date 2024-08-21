"use client";

export interface Message {
    kind: "received" | "sent";
    data: string;
    moderated: boolean;
    approved: boolean;
}
