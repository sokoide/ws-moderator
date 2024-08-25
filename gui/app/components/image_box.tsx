"use client";

import "./messagebox.css";
import React, { useEffect, useState, useRef } from "react";
import { Message, ModRequest } from "./types";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import RobotIcon from "@mui/icons-material/SmartToy";
import UserIcon from "@mui/icons-material/Person";

interface MessageProps {
    msg: ModRequest;
    cn: string;
}

const ImageBox: React.FC<MessageProps> = ({ msg, cn}) => {
    const imageUrl = msg.message.data;

    return (
        <>
            <div className={cn}>
                <div className={cn}>
                    <img
                        src={imageUrl}
                        alt="{msg.message.data}"
                        style={{
                            objectFit: "contain",
                            transition: "all 0.3s ease",
                        }}
                    />
                </div>
            </div>
            <br />
        </>
    );
};

export default ImageBox;
