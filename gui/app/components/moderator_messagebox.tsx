"use client";

import "./messagebox.css";
import React, { useEffect, useState, useRef } from "react";
import { Message, ModRequest } from "./types";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import RobotIcon from "@mui/icons-material/SmartToy";
import UserIcon from "@mui/icons-material/Person";

interface MessageProps {
    msg: ModRequest;
}

const ModerateButtons = () => {
    return (
        <>
            <br />
            <Button>Approve</Button>
            <Button>Deny</Button>
        </>
    );
};

const ModeratorMessageBox: React.FC<MessageProps> = ({ msg }) => {
    console.info("message: %O", msg);

    if (msg.moderated) {
        return <></>;
    } else {
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
                    <p>{msg.message.data}</p>
                    {msg.moderated || msg.approved ? "" : <ModerateButtons />}
                </Box>
                <br />
            </>
        );
    }
};

export default ModeratorMessageBox;
