"use client";

import React, { useState, useContext } from "react";
import "./login.css";
import { AppContext } from "@/context/app-context";
import Button from "@mui/material/Button";

const Login = () => {
    const { loginInfo, login, logout } = useContext(AppContext);

    const [username, setUsername] = useState("");
    const [employee, setEmployee] = useState("");
    const [email, setEmail] = useState("");

    const handleUsername = (e) => {
        setUsername(e.target.value);
    };

    const handleEmployee = (e) => {
        setEmployee(e.target.value);
    };

    const handleEmail = (e) => {
        setEmail(e.target.value);
    };

    const handleLogin = () => {
        if (username === "") return;
        if (employee === "") return;
        if (email === "") return;
        login(username, employee, email);
        window.location.href = "/";
    };

    const handleLogout = () => {
        setUsername("");
        logout();
    };

    return (
        <div className="wrapper">
            <div className="content">
                {loginInfo.loggedIn ? (
                    <>
                        <p>Author: {loginInfo.username}</p>
                        <p>Employee Name: {loginInfo.employee}</p>
                        <p>Employee Email: {loginInfo.email}</p>
                        <Button variant="contained" onClick={handleLogout}>
                            Logout
                        </Button>
                    </>
                ) : (
                    <>
                        <div className="input">
                            Author Name:
                            <input
                                type="text"
                                name="username"
                                id="username"
                                className="username"
                                onChange={handleUsername}
                            />
                        </div>
                        <div className="input">
                            Employee Name:
                            <input
                                type="text"
                                name="employee"
                                id="employee"
                                className="username"
                                onChange={handleEmployee}
                            />
                        </div>
                        <div className="input">
                            Employee Personal Email:
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
