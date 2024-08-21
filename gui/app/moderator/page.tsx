"use client";

import React, { useCallback, useEffect, useState } from "react";

const Moderator = () => {
    const socket = new WebSocket("ws://127.0.0.1:80/chat");
    const [message, setMessage] = useState("");

    useEffect(() => {
        socket.onopen = () => {
            console.log("useEffect: onopen");
            setMessage("Connected");
        };

        socket.onmessage = (e) => {
            console.log("useEffect: onmessage");
            setMessage("Get message from server: " + e.data);
        };

        return () => {
            console.log("useEffect: return");
            // socket.close();
        };
    }, []);

    return <>Moderator</>;
};

export default Moderator;
