"use client";

import "./messagebox.css";
import React, { useEffect, useState, useRef } from "react";
import { Message, ModRequest } from "./types";
import Box from "@mui/material/Box";
import RobotIcon from "@mui/icons-material/SmartToy";
import UserIcon from "@mui/icons-material/Person";

interface MessageProps {
    msg: ModRequest;
}

const MessageBox: React.FC<MessageProps> = ({ msg }) => {
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
                {msg.user_email !== "system@system" ? <UserIcon /> : <RobotIcon />}
                <p>{msg.message.data}</p>
            </Box>
            <br />
        </>
    );
};

export default MessageBox;
