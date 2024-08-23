"use client";

import "./messagebox.css";
import React, { useEffect, useState, useRef } from "react";
import { Message, ModRequest } from "./types";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Divider from "@mui/material/Divider";
import RobotIcon from "@mui/icons-material/SmartToy";
import UserIcon from "@mui/icons-material/Person";

interface MessageProps {
    msg: ModRequest;
}

interface ModerateButtonsProps {
    msgid: string;
}

const ModerateButtons: React.FC<MessageProps> = ({ msg }) => {
    return (
        <>
            <Box flexDirection={"column"}>
                <Divider/>
                <Box>
                    <p>email: {msg.user_email}</p>
                    <p>msgid: {msg.id}</p>
                </Box>
                <Box flexDirection={"row"}>
                    <Button>Approve</Button>
                    <Button>Deny</Button>
                </Box>
            </Box>
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
                    <Box flexDirection={"column"}>
                        <Box>
                            <p>{msg.message.data}</p>
                        </Box>
                        {msg.moderated || msg.approved ? (
                            ""
                        ) : (
                            <ModerateButtons msg={msg} />
                        )}
                    </Box>
                </Box>
                <br />
            </>
        );
    }
};

export default ModeratorMessageBox;
