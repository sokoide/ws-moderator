"use client";

import React, { useContext } from "react";
import Link from "next/link";

import "./navbar.css";
import { AppContext } from "@/context/app-context";

import { AppBar, Toolbar, Typography, IconButton } from "@mui/material";

const Navbar = () => {
    const context = useContext(AppContext);
    if (context === null) {
        console.error("context not available");
        return;
    }
    const { loginInfo } = context;

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
                    <Link href="/">Family Day Event 2024</Link>
                </Typography>
                <div className="links">
                    <div className="user">
                        {loginInfo?.loggedIn ? (
                            <Link href="/login">{loginInfo.email}</Link>
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
