"use client";

import "./messagebox.css";
import React, { useEffect, useState, useRef } from "react";
import { ChatItem } from "./types";
import Box from "@mui/material/Box";
import RobotIcon from "@mui/icons-material/SmartToy";
import UserIcon from "@mui/icons-material/Person";

interface MessageProps {
    message: ChatItem;
}

const MessageBox: React.FC<MessageProps> = ({ message }) => {
    console.info("message.data: %O", message.data);

    if (!message.userEmail.endsWith("@_system")){
        return (
            <>
                <Box
                    whiteSpace="pre-line"
                    my={0}
                    display="flex"
                    gap={2}
                    p={2}
                    sx={{
                        width: "95%",
                        border: "2px solid grey",
                        borderRadius: 1,
                    }}
                >
                    {message.direction === "sent" ? (
                        <UserIcon />
                    ) : (
                        <RobotIcon />
                    )}
                    <p>{message.data}</p>
                </Box>
                <br />
            </>
        );
    } else {
        return <>TODO:</>;
    }
};

export default MessageBox;
