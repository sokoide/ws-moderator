import "./moderator.css";
import React, {
    useContext,
    useEffect,
    useState,
    useRef,
    ChangeEvent,
} from "react";
import { AppContext } from "@/context/app-context";
import { Box, Button, Typography, List, ListItem } from "@mui/material";
import ModeratorMessageBox from "./moderator_messagebox";
import ClientUtil from "./client_util";
import { Message, ModRequest } from "./types";
import { green } from "@mui/material/colors";
import { v4 as uuid } from "uuid";

const Moderator = () => {
    const context = useContext(AppContext);
    if (context === null) {
        console.error("context not available");
        return;
    }
    const { loginInfo } = context;

    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [inputValue, setInputValue] = useState("");
    const [messages, setMessages] = useState<ModRequest[]>([]);
    const messagePaneRef = useRef<HTMLDivElement>(null);
    const taRef = useRef(null);
    const [logs, setLogs] = useState<string[]>([]);
    const MAX_LOG_MESSAGES = 1000;
    const logEndRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const wsUrl = process.env.NEXT_PUBLIC_MODERATOR_WS ?? "undefined";
        console.log("wsUrl: %O", wsUrl);
        const ws = new WebSocket(wsUrl);

        ws.onopen = () => {
            console.log("useEffect: onopen");
            ClientUtil.sendMessage(ws, "", "", "", "system", "", false, false);
        };

        ws.onmessage = (e) => {
            let msg = JSON.parse(e.data) as ModRequest;

            if (msg.message.kind === "ping") return;
            console.log("useEffect: onmessage: %O", msg);

            addLogMessage(formatMessage(msg), true);
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

    // scroll to the bottom
    useEffect(() => {
        logEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [logs]);

    const formatMessage = (msg: ModRequest) : string => {
        return `${msg.user_email} | ${msg.message.kind} | aprv:${msg.approved} | mdrt:${msg.moderated} | ${msg.message.data.slice(0, 256)}`
    }

    const addLogMessage = (message: string, add_timestamp: boolean) => {
        var msg: string
        if (add_timestamp === true) {
            msg = `${new Date().toLocaleTimeString()}: ${message}`;
        } else {
            msg = message;
        }
        setLogs((prevLogs) => {
            // Ensure we keep only the last MAX_LOG_MESSAGES messages
            const newLogs = [...prevLogs, msg];
            if (newLogs.length > MAX_LOG_MESSAGES) {
                newLogs.shift(); // Remove the oldest message
            }
            return newLogs;
        });
    };

    const handleApprove = (msgid: string) => {
        console.log("handleApprove: %O", msgid);
        ClientUtil.sendMessage(
            socket,
            msgid,
            "",
            loginInfo?.email ?? "",
            "system",
            "approve",
            false,
            false
        );
    };

    const handleDeny = (msgid: string) => {
        console.log("handleDeny: %O", msgid);
        ClientUtil.sendMessage(
            socket,
            msgid,
            "",
            loginInfo?.email ?? "",
            "system",
            "deny",
            false,
            false
        );
    };

    return (
        <div style={{ width: '100%' }}> {/* Added full width */}
            <Box
                sx={{
                    display: "flex",
                    flexDirection: "column",
                    height: "60vh",
                    width: "100%",  // Ensure the Box uses the full width
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
            <Box
                sx={{
                    maxHeight: "33vh",
                    height: "33vh",
                    overflowY: "auto",
                    padding: 2,
                    width: '100%',  // Ensure the log monitor uses the full width
                    backgroundColor: "black",
                    color: "green",
                    borderRadius: "4px",
                    fontFamily: "monospace",
                }}
            >
                <Typography variant="h6" sx={{ marginBottom: 1 }}>
                    Log Monitor
                </Typography>
                <List>
                    {logs.map((log, index) => (
                        <ListItem
                            key={index}
                            sx={{ padding: 0, wordWrap: "break-word" }}
                        >
                            <Typography variant="body2">{log}</Typography>
                        </ListItem>
                    ))}
                    <div ref={logEndRef} />
                </List>
            </Box>
        </div>
    );
};

export default Moderator;
