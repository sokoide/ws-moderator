"use client";

import "./messagebox.css";
import React, { useEffect, useState, useRef } from "react";
import { Message, ModRequest } from "./types";
import Box from "@mui/material/Box";
import RobotIcon from "@mui/icons-material/SmartToy";
import UserIcon from "@mui/icons-material/Person";

interface MessageProps {
    msg: ModRequest;
    moderator: boolean;
}

const MessageBox: React.FC<MessageProps> = ({ msg, moderator }) => {
    console.info("message: %O", msg);

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
                {msg.client_id !== "bot" ? <UserIcon /> : <RobotIcon />}
                {moderator ? <p>{msg.message.data}</p> : <p>{msg.message.data}</p>}
            </Box>
            <br />
        </>
    );
};

export default MessageBox;
