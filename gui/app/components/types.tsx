"use client";

export interface Message {
    kind: "received" | "sent";
    data: string;
}
