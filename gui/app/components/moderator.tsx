"use client";

import "./moderator.css";
import React, { useContext, useEffect, useState, useRef } from "react";
import { AppContext } from "@/context/app-context";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import ModeratorMessageBox from "./moderator_messagebox";
import ClientUtil from "./client_util";
import { Message, ModRequest } from "./types";
import { green } from "@mui/material/colors";
import { v4 as uuid } from "uuid";

const Moderator = () => {
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [inputValue, setInputValue] = useState("");
    const [messages, setMessages] = useState<ModRequest[]>([]);
    const messagePaneRef = useRef<HTMLDivElement>(null);
    const taRef = useRef(null);
    const { loginInfo, login, logout } = useContext(AppContext);

    useEffect(() => {
        // TODO: make it env var
        const ws = new WebSocket("ws://127.0.0.1:80/moderator");

        ws.onopen = () => {
            console.log("useEffect: onopen");
            ClientUtil.sendMessage(ws, "", "", "", "system", "", false, false);
        };

        ws.onmessage = (e) => {
            let msg = JSON.parse(e.data) as ModRequest;
            if (msg.message.kind === "ping") return;
            console.log("useEffect: onmessage: %O", msg);
            if (msg.moderated === true) {
                console.log("removing moderated message");
                setMessages((prevMessages) =>
                    prevMessages.filter((m) => m.id !== msg.id)
                );
            } else {
                setMessages((prevMessages) => [...prevMessages, msg]);
            }
        };

        ws.onerror = (error) => {
            console.error("WebSocket error: %O", error);
        };

        ws.onclose = () => {
            console.log("WebSocket connection closed");
        };

        setSocket(ws);

        return () => {
            console.log("useEffect: return");
            ws.close();
        };
    }, []);

    useEffect(() => {
        if (messagePaneRef.current) {
            messagePaneRef.current.scrollTop =
                messagePaneRef.current.scrollHeight;
        }
    }, [messages]);

    const handleApprove = (msgid) => {
        console.log("handleApprove: %O", msgid);
        ClientUtil.sendMessage(
            socket,
            msgid,
            "",
            loginInfo.email,
            "system",
            "approve",
            false,
            false
        );
    };

    const handleDeny = (msgid) => {
        console.log("handleDeny: %O", msgid);
        ClientUtil.sendMessage(
            socket,
            msgid,
            "",
            loginInfo.email,
            "system",
            "deny",
            false,
            false
        );
    };

    return (
        <div>
            <Box
                sx={{
                    display: "flex",
                    flexDirection: "column",
                    height: "95vh",
                }}
            >
                <Box
                    ref={messagePaneRef}
                    sx={{
                        height: "100%",
                        padding: 2,
                        overflowY: "auto",
                        overflowX: "auto",
                    }}
                >
                    {messages.map(function (msg, i) {
                        return (
                            <ModeratorMessageBox
                                msg={msg}
                                onApprove={handleApprove}
                                onDeny={handleDeny}
                                key={uuid()}
                            />
                        );
                    })}
                </Box>
            </Box>
        </div>
    );
};

export default Moderator;
