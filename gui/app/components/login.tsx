"use client";

import React, { useState, useContext, ChangeEvent } from "react";
import "./login.css";
import { AppContext } from "@/context/app-context";
import Button from "@mui/material/Button";

const Login = () => {
    const context = useContext(AppContext);
    if (context === null) {
        console.error("context not available");
        return;
    }

    const { loginInfo, login, logout } = context;
    const [email, setEmail] = useState("");

    const handleEmail = (e: ChangeEvent<HTMLInputElement>) => {
        setEmail(e.target.value);
    };

    const handleLogin = () => {
        if (email === "") return;
        login(email);
        if (typeof window !== "undefined") window.location.href = "/";
    };

    const handleLogout = () => {
        setEmail("");
        logout();
    };

    return (
        <div className="wrapper">
            <div className="content">
                {loginInfo?.loggedIn ? (
                    <>
                        <p>Employee Email: {loginInfo.email}</p>
                        <Button variant="contained" onClick={handleLogout}>
                            Logout
                        </Button>
                    </>
                ) : (
                    <>
                        <p>
                            Note: The <b>Employee Personal Email</b>(case
                            sensitive) will be used as the chat conversation
                            key.
                        </p>
                        <hr />
                        <div className="input">
                            Employee Personal Email (case sensitive):
                            <input
                                type="text"
                                name="email"
                                id="email"
                                className="username"
                                onChange={handleEmail}
                            />
                        </div>
                        <Button variant="contained" onClick={handleLogin}>
                            Login
                        </Button>
                    </>
                )}
            </div>
        </div>
    );
};

export default Login;
