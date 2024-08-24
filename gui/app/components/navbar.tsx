"use client";

import React, { useContext } from "react";
import Link from "next/link";

import "./navbar.css";
import { AppContext } from "@/context/app-context";

import {
    AppBar,
    Toolbar,
    Typography,
    IconButton,
} from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";

const Navbar = () => {
    const { loginInfo } = useContext(AppContext);

    return (
        <AppBar position="static">
            <Toolbar className="toolbar">
                <IconButton
                    edge="start"
                    color="inherit"
                    aria-label="menu"
                    sx={{ mr: 2 }}
                >
                    {/* <MenuIcon /> */}
                </IconButton>
                <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                    Family Day Event 2024
                </Typography>
                <div className="links">
                    <Link href="/">Chat</Link>
                    &nbsp;
                    <div className="user">
                        {loginInfo.loggedIn ? (
                            <Link href="/login">{loginInfo.username}</Link>
                        ) : (
                            <Link href="/login">Login</Link>
                        )}
                    </div>
                </div>
            </Toolbar>
        </AppBar>
    );
};

export default Navbar;
