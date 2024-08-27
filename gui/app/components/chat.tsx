"use client";

import "./chat.css";
import React, {
    useContext,
    useEffect,
    useState,
    useRef,
    ChangeEvent,
    KeyboardEvent,
} from "react";
import { AppContext } from "@/context/app-context";
import { TextField, Box, Button } from "@mui/material";

import MessageBox from "./messagebox";
import { ModRequest } from "./types";
import ClientUtil from "./client_util";
import { v4 as uuid } from "uuid";
import { useRouter } from "next/navigation";

const Chat = () => {
    const router = useRouter();
    const context = useContext(AppContext);
    if (context === null) {
        console.error("context not available");
        return;
    }

    const { loginInfo } = context;
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [inputValue, setInputValue] = useState("");
    const [suggestion, setSuggestion] = useState<string | null>(null);
    const [showSuggestion, setShowSuggestion] = useState(false);
    const [isImagineMode, setIsImagineMode] = useState(false); // New state to track imagine mode
    const [messages, setMessages] = useState<ModRequest[]>([]);
    const messagesEndRef = useRef<HTMLDivElement | null>(null);
    const inputContainerRef = useRef<HTMLDivElement | null>(null);
    const eopRef = useRef<HTMLDivElement | null>(null);
    const titleRef = useRef<HTMLInputElement | null>(null);
    const userRef = useRef<HTMLInputElement | null>(null);
    const employeeRef = useRef<HTMLInputElement | null>(null);

    if (loginInfo != null && loginInfo.loggedIn === false) {
        if (typeof window !== "undefined") window.location.href = "/login";
    }

    const scrollToBottom = () => {
        if (messagesEndRef.current) {
            messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
        }
    };

    useEffect(() => {
        const wsUrl = process.env.NEXT_PUBLIC_CHAT_WS ?? "undefined";
        const ws = new WebSocket(wsUrl);

        ws.onopen = () => {
            ClientUtil.sendMessage(
                ws,
                "",
                "",
                loginInfo?.email ?? "",
                "system",
                "",
                false,
                false
            );
        };

        ws.onmessage = (e) => {
            let msg = JSON.parse(e.data) as ModRequest;
            if (
                (msg.client_id === "bot" && msg.id === "") ||
                msg.user_email === (loginInfo?.email ?? "")
            ) {
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
            ws.close();
        };
    }, []);

    useEffect(() => {
        setCheckboxStates(
            messages.map((message) => ({
                msg: message,
                checked: false,
            }))
        );

        const timer = setTimeout(() => {
            scrollToBottom();
        }, 10);

        return () => clearTimeout(timer);
    }, [messages]);

    useEffect(() => {
        if (inputValue.startsWith("/imagine")) {
            setSuggestion("Imagine Mode");
            setShowSuggestion(true);
            setIsImagineMode(true); // Set imagine mode when exact match
        } else if (inputValue.startsWith("/")) {
            setSuggestion("Press Enter to imagine mode");
            setShowSuggestion(true);
            setIsImagineMode(false);
        } else {
            setShowSuggestion(false);
            setIsImagineMode(false); // Reset imagine mode when not in the command
        }
    }, [inputValue]);

    const handleInputChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
        setInputValue(e.target.value);
    };

    const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
        if (showSuggestion && e.key === "Enter") {
            e.preventDefault();
            if (!isImagineMode) {
                setInputValue("/imagine");
                setShowSuggestion(false);
                setIsImagineMode(true); // Ensure imagine mode is set
            }
        }
    };

    const handleSuggestionClick = () => {
        if (!isImagineMode) {
            setInputValue("/imagine");
            setShowSuggestion(false);
            setIsImagineMode(true); // Ensure imagine mode is set
        }
    };

    const onSend = () => {
        if (isImagineMode) {
            console.log("Imagine mode active, handling differently if needed");
            // Handle imagine mode specific behavior
        }

        let msg = ClientUtil.sendMessage(
            socket,
            "",
            "",
            loginInfo?.email ?? "",
            "txt",
            inputValue,
            false,
            false
        );
        if (msg != null) {
            setInputValue("");
            setShowSuggestion(false);
            setIsImagineMode(false); // Reset imagine mode after sending
        } else {
            console.error("WebSocket is not open");
        }
    };

    const onComplete = () => {
        let id_txt_count: number = 0;
        let id_url_count: number = 0;
        let id_txt: string = "";
        let id_url: string = "";
        let email: string = "";

        if (
            !userRef.current ||
            !employeeRef.current ||
            !titleRef.current ||
            userRef.current.value.trim() === "" ||
            employeeRef.current.value.trim() === "" ||
            titleRef.current.value.trim() === ""
        ) {
            alert("please input Title, User and Employee");
            return;
        }

        checkboxStates.map(({ msg, checked }) => {
            if (checked) {
                if (msg.message.kind === "txt") {
                    id_txt_count++;
                    id_txt = msg.id;
                    email = msg.user_email;
                } else if (msg.message.kind === "url") {
                    id_url_count++;
                    id_url = msg.id;
                }
            }
        });
        if (id_txt_count !== 1) {
            alert("please select only one text");
            return;
        }
        if (id_url_count !== 1) {
            alert("please select only one image");
            return;
        }

        let target_page: string = "/confirm?";
        let params: string =
            "msgid_txt=" +
            encodeURIComponent(id_txt) +
            "&msgid_url=" +
            encodeURIComponent(id_url) +
            "&email=" +
            encodeURIComponent(email) +
            "&title=" +
            encodeURIComponent(titleRef.current.value) +
            "&user=" +
            encodeURIComponent(userRef.current.value) +
            "&employee=" +
            encodeURIComponent(employeeRef.current.value);
        const encoded_page = target_page + params;
        router.push(encoded_page);
    };

    const onUncheckAll = () => {
        setCheckboxesState(false);
    };

    const [checkboxStates, setCheckboxStates] = useState<
        { msg: ModRequest; checked: boolean }[]
    >([]);

    const handleCheckboxChange =
        (id: string) => (event: React.ChangeEvent<HTMLInputElement>) => {
            setCheckboxStates((prevStates) =>
                prevStates.map((checkbox) =>
                    checkbox.msg.id === id
                        ? { ...checkbox, checked: event.target.checked }
                        : checkbox
                )
            );
        };

    const setCheckboxesState = (checked: boolean) => {
        setCheckboxStates((prevStates) =>
            prevStates.map((checkbox) => ({ ...checkbox, checked }))
        );
    };

    return (
        <>
            <div ref={eopRef} className="container" style={{ flex: 1, display: "flex", flexDirection: "column", overflow: "hidden" }}>
                <Box p={2} sx={{ flexGrow: 1, overflowY: "auto", paddingBottom: "10px" }}>
                    <Box>
                        {checkboxStates.map(({ msg, checked }) => (
                            <MessageBox
                                msg={msg}
                                checked={checked}
                                onCheckboxChange={handleCheckboxChange(
                                    msg.id
                                )}
                                key={uuid()}
                            />
                        ))}
                    </Box>
                    <div ref={messagesEndRef} />
                </Box>
                <div ref={inputContainerRef} style={{ padding: "10px", borderTop: "1px solid #ccc", background: "#fff", flexShrink: 0 }}>
                    <Box display="flex" flexDirection="column" gap={1}>
                        <Box display="flex" flexDirection="row" gap={1}>
                        <textarea
                            value={inputValue}
                            onChange={handleInputChange}
                            onKeyDown={handleKeyDown}
                            placeholder="Type here..."
                            rows={2}
                            className="bordered-input"
                            style={{ width: "100%", resize: "none" }}
                        />
                            <Button variant="contained" onClick={onSend}>
                                Ask
                            </Button>
                        </Box>
                        {showSuggestion && suggestion && (
                            <Box
                                p={1}
                                className={`suggestion-box ${isImagineMode ? "rainbow-suggestion" : ""}`} // Apply rainbow suggestion class conditionally
                                onClick={handleSuggestionClick}
                                sx={{
                                    cursor: "pointer",
                                    border: "1px solid #ccc",
                                    borderRadius: "4px",
                                    backgroundColor: "#f9f9f9",
                                }}
                            >
                                <p style={{ margin: 0 }}>{suggestion}</p>
                            </Box>
                        )}
                        <Box display="flex" justifyContent="space-between" alignItems="center" flexWrap="wrap">
                            <Box display="flex" gap={2} flexGrow={1}>
                                <TextField
                                    id="title"
                                    inputRef={titleRef}
                                    label="Title of a book"
                                    variant="outlined"
                                    fullWidth
                                />
                                <TextField
                                    id="user"
                                    inputRef={userRef}
                                    label="User who used AI today"
                                    variant="outlined"
                                    fullWidth
                                />
                                <TextField
                                    id="employee"
                                    inputRef={employeeRef}
                                    label="Employee Name"
                                    variant="outlined"
                                    fullWidth
                                />
                            </Box>
                            &nbsp;
                            <Box display="flex" gap={2}>
                                <Button variant="outlined" onClick={onUncheckAll}>
                                    Uncheck All
                                </Button>
                                <Button variant="outlined" onClick={onComplete}>
                                    Generate PDF
                                </Button>
                            </Box>
                        </Box>
                    </Box>
                </div>
            </div>
        </>
    );
};

export default Chat;
