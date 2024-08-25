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
}

const ImageBox: React.FC<MessageProps> = ({ msg }) => {
    // src="http://localhost/images/alice_gmail.com/2e632d95-96db-4e81-a22c-7c7d39db5eb6.png"
    // src={msg.message.data}
    const imageUrl = msg.message.data;
    console.log("* %O", imageUrl)

    return (
        <>
            <div className="zoom">
                <p>{msg.message.data}</p>
                <div className="zoom">
                    <img
                        src={imageUrl}
                        alt="image"
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
